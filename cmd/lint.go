package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
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

func fixWorkdirPerm(dir string) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var mode fs.FileMode

		if info.IsDir() {
			mode = permDir
		} else {
			mode = permFile
		}

		return os.Chmod(path, mode) //nolint:forbidigo
	})
}

//nolint:forbidigo
func updateWorkdir(ctx context.Context, dir string, cloneURL string) error {
	_, err := os.Stat(dir)
	notfound := err != nil && errors.Is(err, os.ErrNotExist)

	if err != nil && !notfound {
		return err
	}

	if notfound {
		_, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{URL: cloneURL})
		if err != nil {
			return err
		}

		return fixWorkdirPerm(dir)
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	wtree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wtree.Pull(&git.PullOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return fixWorkdirPerm(dir)
}

func checkCompliance(ctx context.Context, module string, cloneURL string, tstamp float64) (*k6lint.Compliance, error) {
	com, found, err := loadCompliance(ctx, module, tstamp)
	if found {
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
