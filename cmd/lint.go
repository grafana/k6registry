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
	"strings"
	"time"
)

const (
	// TTL for compliance cache (1 week).
	complianceCacheTTL = 60 * 60 * 24 * 7

	xk6Binary = "xk6"

	// xk6 lint returns 2 if some check failed. Other codes mean the lint failed.
	lintFailedRC = 2
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

	// Compliance check timestamp in Unix time
	Timestamp int64 `json:"timestamp" mapstructure:"timestamp" yaml:"timestamp"`
}

func loadCompliance(ctx context.Context, module string, version string, timestamp int64) (*Compliance, bool, error) {
	base, err := checksDir(ctx)
	if err != nil {
		return nil, false, err
	}

	filename := filepath.Join(base, module, version) + ".json"

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

func saveCompliance(ctx context.Context, module string, version string, comp *Compliance) error {
	base, err := checksDir(ctx)
	if err != nil {
		return err
	}

	filename := filepath.Join(base, module, version) + ".json"

	if err := os.MkdirAll(filepath.Dir(filename), permDir); err != nil {
		return err
	}

	data, err := json.Marshal(comp)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, permFile)
}

func checkCompliance(
	ctx context.Context,
	module string,
	version string,
	official bool,
	checks []string,
	cloneURL string,
	tstamp int64,
) (*Compliance, error) {
	com, found, err := loadCompliance(ctx, module, version, tstamp)
	if found {
		slog.Debug("Compliance from cache", "module", module, "version", version)

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

	if err := checkoutModVersion(ctx, dir, cloneURL, version); err != nil {
		return nil, err
	}

	slog.Debug("Check compliance", "module", module)

	_, err = exec.LookPath(xk6Binary)
	if err != nil {
		return nil, fmt.Errorf("searching xk6 path %w", err)
	}

	lintOut := &bytes.Buffer{}
	lintErr := &bytes.Buffer{}
	lintArgs := []string{"lint", "--json", "-v"}

	if len(checks) > 0 {
		lintArgs = append(lintArgs, "--enable-only", strings.Join(checks, ","))
	}

	lintCmd := exec.Command(xk6Binary, lintArgs...)

	lintCmd.Stdout = lintOut
	lintCmd.Stderr = lintErr
	lintCmd.Dir = dir

	err = lintCmd.Run()
	if err != nil {
		rc := lintCmd.ProcessState.ExitCode()
		slog.Debug("xk6 execution failed", "rc", rc, "stderr", lintErr.String())

		if rc != lintFailedRC {
			return nil, fmt.Errorf("xk6 lint failed %w", err)
		}
	}

	compliance := &Compliance{}
	err = json.Unmarshal(lintOut.Bytes(), compliance)
	if err != nil {
		return nil, err
	}

	for idx := range compliance.Checks {
		compliance.Checks[idx].Details = ""
	}

	err = saveCompliance(ctx, module, version, compliance)
	if err != nil {
		return nil, err
	}

	return compliance, nil
}
