package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const gitBinary = "git"

var errGitNotFound = errors.New("git executable not found")

func checkGitAvailable() error {
	if _, err := exec.LookPath(gitBinary); err != nil {
		return fmt.Errorf("%w: %w", errGitNotFound, err)
	}

	return nil
}

// runGit runs git with args, using dir as the working directory (ignored if empty).
func runGit(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, gitBinary, args...) //nolint:gosec // git is a fixed, trusted binary
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}

	return stdout.Bytes(), nil
}

func splitLines(out []byte) []string {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	result := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			result = append(result, line)
		}
	}

	return result
}

// openOrCloneBareRepo ensures a mirror clone of cloneURL exists at dir.
// A mirror is a bare repository whose refs (branches and tags) are kept in
// sync 1:1 with the remote on fetch, so it never materializes a working tree.
func openOrCloneBareRepo(ctx context.Context, dir string, cloneURL string) error {
	if err := checkGitAvailable(); err != nil {
		return err
	}

	_, err := os.Stat(dir) //nolint:gosec,forbidigo // modules cache dir
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) { //nolint:forbidigo // CLI tool
		return err
	}

	_, err = runGit(ctx, "", "clone", "--mirror", cloneURL, dir)

	return err
}

// listTags returns the tag names present in the mirror repo at dir.
func listTags(ctx context.Context, dir string) ([]string, error) {
	out, err := runGit(ctx, dir, "tag", "--list")
	if err != nil {
		return nil, err
	}

	return splitLines(out), nil
}

// defaultBranch returns the branch name that dir's HEAD points to.
func defaultBranch(ctx context.Context, dir string) (string, error) {
	out, err := runGit(ctx, dir, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// checkoutWorktree checks out version (or the default branch if version is
// empty) from the mirror repo at repoDir into a new temporary worktree,
// returning its path and a cleanup function that removes it.
func checkoutWorktree(ctx context.Context, repoDir string, version string) (string, func() error, error) {
	if err := checkGitAvailable(); err != nil {
		return "", nil, err
	}

	if _, err := runGit(ctx, repoDir, "fetch", "--prune", "origin"); err != nil {
		return "", nil, err
	}

	ref := version

	if ref == "" {
		branch, err := defaultBranch(ctx, repoDir)
		if err != nil {
			return "", nil, err
		}

		ref = branch
	}

	worktreeDir, err := os.MkdirTemp("", "k6registry-*") //nolint:forbidigo // ephemeral checkout
	if err != nil {
		return "", nil, err
	}

	if _, err := runGit(ctx, repoDir, "worktree", "add", "--detach", worktreeDir, ref); err != nil {
		_ = os.RemoveAll(worktreeDir) //nolint:forbidigo // cleanup on failure

		return "", nil, err
	}

	cleanup := func() error {
		if _, err := runGit(ctx, repoDir, "worktree", "remove", "--force", worktreeDir); err != nil {
			_ = os.RemoveAll(worktreeDir) //nolint:forbidigo // best-effort cleanup fallback
			_, _ = runGit(ctx, repoDir, "worktree", "prune")

			return err
		}

		return nil
	}

	return worktreeDir, cleanup, nil
}
