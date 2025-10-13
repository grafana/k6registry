package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
)

const (
	// TTL for compliance cache (1 week).
	complianceCacheTTL = 60 * 60 * 24 * 7

	xk6Binary = "xk6"
)

// Check is the result of a particular inspection.
type Check struct {
	// Textual explanation of the check result.
	Details string `json:"details,omitempty" mapstructure:"details,omitempty" yaml:"details,omitempty"`

	// The ID of the checker.
	// It identifies the method of check, not the execution of the check.
	ID string `json:"id" mapstructure:"id" yaml:"id"`

	// The result of the check.
	// A true value of the passed property indicates a successful check, while a false
	// value indicates a failure.
	Passed bool `json:"passed" mapstructure:"passed" yaml:"passed"`
}

// Compliance is the result of the extension's k6 compliance checks.
type Compliance struct {
	// Results of individual checks.
	Checks []Check `json:"checks,omitempty" mapstructure:"checks,omitempty" yaml:"checks,omitempty"`

	// The results of the checks are in the form of a grade.
	Grade Grade `json:"grade" mapstructure:"grade" yaml:"grade"`

	// Compliance expressed as a percentage.
	Level int `json:"level" mapstructure:"level" yaml:"level"`

	// Compliance check timestamp in Unix time
	Timestamp int64 `json:"timestamp" mapstructure:"timestamp" yaml:"timestamp"`
}

// Grade defines the extension grading according to the compliance level.
type Grade string

const (
	GradeA Grade = "A" //nolint:revive
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeE Grade = "E"
	GradeF Grade = "F"
	GradeG Grade = "G"
)

func loadCompliance(ctx context.Context, module string, timestamp int64) (*Compliance, bool, error) {
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

	var comp Compliance

	if err := json.Unmarshal(data, &comp); err != nil {
		return nil, false, err
	}

	age := time.Now().Unix() - comp.Timestamp

	if comp.Timestamp >= timestamp && age <= complianceCacheTTL {
		return &comp, true, nil
	}

	return nil, false, nil
}

func saveCompliance(ctx context.Context, module string, comp *Compliance) error {
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

func checkCompliance(
	ctx context.Context,
	module string,
	official bool,
	cloneURL string,
	tstamp int64,
) (*Compliance, error) {
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

	_, err = exec.LookPath(xk6Binary)
	if err != nil {
		return nil, fmt.Errorf("searching xk6 path %w", err)
	}

	lintOut := &bytes.Buffer{}
	lintErr := &bytes.Buffer{}
	lintCmd := exec.Command(xk6Binary, "lint", "--json")
	lintCmd.Stdout = lintOut
	lintCmd.Stderr = lintErr

	err = lintCmd.Run()
	if err != nil {
		slog.Debug("xk6 execution failed", "rc", lintCmd.ProcessState.ExitCode(), "stderr", lintErr.String())

		return nil, fmt.Errorf("xk6 lint failed %w", err)
	}

	compliance := &Compliance{}
	err = json.Unmarshal(lintOut.Bytes(), compliance)
	if err != nil {
		return nil, err
	}

	for idx := range compliance.Checks {
		compliance.Checks[idx].Details = ""
	}

	err = saveCompliance(ctx, module, compliance)
	if err != nil {
		return nil, err
	}

	return compliance, nil
}
