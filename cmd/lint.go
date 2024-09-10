package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/grafana/k6lint"
)

//nolint:forbidigo
func loadCompliance(ctx context.Context, module string, timestamp float64) (*k6lint.Compliance, bool, error) {
	base, err := checksDir(ctx)
	if err != nil {
		return nil, false, err
	}

	filename := filepath.Join(base, module) + ".json"

	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var comp k6lint.Compliance

	if err := json.Unmarshal(data, &comp); err != nil {
		return nil, false, err
	}

	if comp.Timestamp >= timestamp {
		return &comp, true, nil
	}

	return nil, false, nil
}

//nolint:forbidigo
func saveCompliance(ctx context.Context, module string, comp *k6lint.Compliance) error {
	base, err := checksDir(ctx)
	if err != nil {
		return err
	}

	filename := filepath.Join(base, module) + ".json"

	if err := os.MkdirAll(filepath.Dir(filename), permDir); err != nil {
		return err
	}

	data, err := json.Marshal(comp)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, permFile)
}

//nolint:forbidigo
func updateWorkdir(ctx context.Context, dir string, cloneURL string) error {
	_, err := os.Stat(dir)
	notfound := err != nil && errors.Is(err, os.ErrNotExist)

	if err != nil && !notfound {
		return err
	}

	if notfound {
		slog.Debug("Clone", "url", cloneURL)

		_, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{URL: cloneURL})
		return err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	wtree, err := repo.Worktree()
	if err != nil {
		return err
	}

	slog.Debug("Pull", "url", cloneURL)

	err = wtree.Pull(&git.PullOptions{Force: true})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		if !errors.Is(err, git.ErrWorktreeNotClean) && !errors.Is(err, git.ErrUnstagedChanges) {
			return err
		}

		slog.Debug("Retry pull", "url", cloneURL)

		head, err := repo.Head()
		if err != nil {
			return err
		}

		err = wtree.Checkout(&git.CheckoutOptions{Force: true, Branch: head.Name()})
		if err != nil {
			return err
		}

		err = wtree.Pull(&git.PullOptions{Force: true})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}
	}

	return nil
}

func checkCompliance(ctx context.Context, module string, cloneURL string, tstamp float64) (*k6lint.Compliance, error) {
	com, found, err := loadCompliance(ctx, module, tstamp)
	if found {
		slog.Debug("Compliance from cache", "module", module)

		return com, nil
	}

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	base, err := modulesDir(ctx)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(base, module)

	if err := updateWorkdir(ctx, dir, cloneURL); err != nil {
		return nil, err
	}

	slog.Debug("Check compliance", "module", module)

	compliance, err := k6lint.Lint(ctx, dir, &k6lint.Options{
		Passed: []k6lint.Checker{k6lint.CheckerLicense, k6lint.CheckerVersions, k6lint.CheckerGit},
	})
	if err != nil {
		return nil, err
	}

	compliance.Checks = nil

	err = saveCompliance(ctx, module, compliance)
	if err != nil {
		return nil, err
	}

	return compliance, nil
}
