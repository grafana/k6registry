package cmd //nolint:testpackage

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var errUnexpectedContent = errors.New("unexpected VERSION file content")

func requireGit(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath(gitBinary); err != nil {
		t.Skip("git executable not found")
	}
}

func runGitT(t *testing.T, dir string, args ...string) []byte {
	t.Helper()

	cmd := exec.CommandContext(context.Background(), gitBinary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), //nolint:forbidigo // test helper
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test.com",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}

	return out
}

func writeFileT(t *testing.T, dir string, name string, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), permFile); err != nil { //nolint:forbidigo // test fixture
		t.Fatal(err)
	}
}

// newTestRemote creates a non-bare repo, usable as a local clone source, with:
//   - tag v1.0.0 -> VERSION file containing "v1.0.0"
//   - tag v1.1.0 -> VERSION file containing "v1.1.0"
//   - tag not-a-version (non-semver, should be filtered out downstream)
//   - an untagged tip commit on the default branch -> VERSION file containing "tip".
func newTestRemote(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	runGitT(t, dir, "init", "-b", "main")

	writeFileT(t, dir, "VERSION", "v1.0.0")
	runGitT(t, dir, "add", ".")
	runGitT(t, dir, "commit", "-m", "v1.0.0")
	runGitT(t, dir, "tag", "v1.0.0")

	writeFileT(t, dir, "VERSION", "v1.1.0")
	runGitT(t, dir, "add", ".")
	runGitT(t, dir, "commit", "-m", "v1.1.0")
	runGitT(t, dir, "tag", "v1.1.0")
	runGitT(t, dir, "tag", "not-a-version")

	writeFileT(t, dir, "VERSION", "tip")
	runGitT(t, dir, "add", ".")
	runGitT(t, dir, "commit", "-m", "tip")

	return dir
}

func TestOpenOrCloneBareRepo(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	remote := newTestRemote(t)
	dest := filepath.Join(t.TempDir(), "repo")

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatalf("clone: %v", err)
	}

	out, err := runGit(ctx, dest, "rev-parse", "--is-bare-repository")
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}

	if strings.TrimSpace(string(out)) != "true" {
		t.Fatalf("expected a bare repository, got %q", out)
	}

	before, err := os.Stat(dest) //nolint:forbidigo // test
	if err != nil {
		t.Fatal(err)
	}

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatalf("second call: %v", err)
	}

	after, err := os.Stat(dest) //nolint:forbidigo // test
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf("expected second call to be a no-op, directory was modified")
	}
}

func TestOpenOrCloneBareRepo_BadURL(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	dest := filepath.Join(t.TempDir(), "repo")

	err := openOrCloneBareRepo(ctx, dest, filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected an error cloning a nonexistent remote")
	}
}

func TestListTags(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	remote := newTestRemote(t)
	dest := filepath.Join(t.TempDir(), "repo")

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatal(err)
	}

	tags, err := listTags(ctx, dest)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{"v1.0.0": true, "v1.1.0": true, "not-a-version": true}

	if len(tags) != len(want) {
		t.Fatalf("got tags %v, want %v distinct entries", tags, len(want))
	}

	for _, tag := range tags {
		if !want[tag] {
			t.Fatalf("unexpected tag %q in %v", tag, tags)
		}
	}
}

func TestLoadGit(t *testing.T) {
	requireGit(t)
	t.Parallel()

	remote := newTestRemote(t)
	ctx := context.WithValue(context.Background(), cacheDirKey{}, t.TempDir())

	versions, err := loadGit(ctx, "example.com/mod", remote)
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(versions)

	want := []string{"v1.0.0", "v1.1.0"}

	if len(versions) != len(want) {
		t.Fatalf("got %v, want %v", versions, want)
	}

	for i, v := range want {
		if versions[i] != v {
			t.Fatalf("got %v, want %v", versions, want)
		}
	}
}

func TestCheckoutWorktree_Version(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	remote := newTestRemote(t)
	dest := filepath.Join(t.TempDir(), "repo")

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatal(err)
	}

	worktreeDir, cleanup, err := checkoutWorktree(ctx, dest, "v1.0.0")
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(worktreeDir, "VERSION")) //nolint:forbidigo // test file in temp dir
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "v1.0.0" {
		t.Fatalf("got %q, want %q", content, "v1.0.0")
	}

	if err := cleanup(); err != nil {
		t.Fatalf("cleanup: %v", err)
	}

	_, err = os.Stat(worktreeDir) //nolint:forbidigo // test
	if !os.IsNotExist(err) {      //nolint:forbidigo // test
		t.Fatalf("expected worktree dir to be removed, stat err=%v", err)
	}

	out, err := runGit(ctx, dest, "worktree", "list")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(out), worktreeDir) {
		t.Fatalf("dangling worktree entry left behind: %s", out)
	}
}

func TestCheckoutWorktree_DefaultBranch(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	remote := newTestRemote(t)
	dest := filepath.Join(t.TempDir(), "repo")

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatal(err)
	}

	worktreeDir, cleanup, err := checkoutWorktree(ctx, dest, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = cleanup() })

	content, err := os.ReadFile(filepath.Join(worktreeDir, "VERSION")) //nolint:forbidigo // test file in temp dir
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "tip" {
		t.Fatalf("got %q, want %q", content, "tip")
	}
}

func TestCheckoutWorktree_UnknownVersion(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	remote := newTestRemote(t)
	dest := filepath.Join(t.TempDir(), "repo")

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatal(err)
	}

	if _, _, err := checkoutWorktree(ctx, dest, "v9.9.9"); err == nil {
		t.Fatal("expected an error checking out a nonexistent tag")
	}
}

func TestCheckoutWorktree_Concurrent(t *testing.T) {
	requireGit(t)
	t.Parallel()

	ctx := context.Background()
	remote := newTestRemote(t)
	dest := filepath.Join(t.TempDir(), "repo")

	if err := openOrCloneBareRepo(ctx, dest, remote); err != nil {
		t.Fatal(err)
	}

	versions := []string{"v1.0.0", "v1.1.0"}
	results := make(chan error, len(versions))

	for _, version := range versions {
		go func(version string) {
			worktreeDir, cleanup, err := checkoutWorktree(ctx, dest, version)
			if err != nil {
				results <- err

				return
			}

			defer func() { _ = cleanup() }()

			content, err := os.ReadFile(filepath.Join(worktreeDir, "VERSION")) //nolint:forbidigo // test file in temp dir
			if err != nil {
				results <- err

				return
			}

			if string(content) != version {
				results <- errUnexpectedContent

				return
			}

			results <- nil
		}(version)
	}

	for range versions {
		if err := <-results; err != nil {
			t.Fatalf("concurrent checkout failed: %v", err)
		}
	}
}

func TestCheckGitAvailable_MissingBinary(t *testing.T) {
	t.Setenv("PATH", "")

	if err := checkGitAvailable(); err == nil {
		t.Fatal("expected an error when git is not on PATH")
	}
}
