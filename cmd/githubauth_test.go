package cmd //nolint:testpackage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// clearGitHubEnv resets everything githubToken consults, so tests are isolated
// from the host's real environment (e.g. a developer machine with gh installed
// and authenticated, or GH_TOKEN set). t.Setenv restores the prior value after
// the test, and can't be combined with t.Parallel, so these tests run serially.
func clearGitHubEnv(t *testing.T) {
	t.Helper()

	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_CONFIG_DIR", t.TempDir())
	t.Setenv("GH_PATH", "")
	t.Setenv("PATH", t.TempDir()) // hide any real `gh` binary that might be on PATH
}

func writeExecutableT(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o755); err != nil { //nolint:forbidigo // test fixture
		t.Fatal(err)
	}
}

func TestGithubToken_GHTokenEnv(t *testing.T) {
	clearGitHubEnv(t)
	t.Setenv("GH_TOKEN", "gh-token-value")

	token, err := githubToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token != "gh-token-value" {
		t.Fatalf("got %q", token)
	}
}

func TestGithubToken_GitHubTokenEnvFallback(t *testing.T) {
	clearGitHubEnv(t)
	t.Setenv("GITHUB_TOKEN", "legacy-token-value")

	token, err := githubToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token != "legacy-token-value" {
		t.Fatalf("got %q", token)
	}
}

func TestGithubToken_GHTokenTakesPrecedence(t *testing.T) {
	clearGitHubEnv(t)
	t.Setenv("GH_TOKEN", "primary")
	t.Setenv("GITHUB_TOKEN", "secondary")

	token, err := githubToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token != "primary" {
		t.Fatalf("got %q, want %q", token, "primary")
	}
}

func TestGithubToken_ConfigFile(t *testing.T) {
	clearGitHubEnv(t)

	configDir := t.TempDir()
	t.Setenv("GH_CONFIG_DIR", configDir)

	hosts := "github.com:\n    oauth_token: config-token-value\n"

	//nolint:forbidigo // test fixture
	if err := os.WriteFile(filepath.Join(configDir, "hosts.yml"), []byte(hosts), 0o600); err != nil {
		t.Fatal(err)
	}

	token, err := githubToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token != "config-token-value" {
		t.Fatalf("got %q", token)
	}
}

func TestGithubToken_GHCli(t *testing.T) {
	clearGitHubEnv(t)

	script := filepath.Join(t.TempDir(), "fake-gh.sh")
	writeExecutableT(t, script, "#!/bin/sh\necho cli-token-value\n")
	t.Setenv("GH_PATH", script)

	token, err := githubToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token != "cli-token-value" {
		t.Fatalf("got %q", token)
	}
}

func TestGithubToken_NotFound(t *testing.T) { //nolint:paralleltest // uses t.Setenv, which is incompatible with t.Parallel
	clearGitHubEnv(t)

	if _, err := githubToken(context.Background()); err == nil {
		t.Fatal("expected an error when no token is available anywhere")
	}
}

func TestGhConfigDir_Precedence(t *testing.T) {
	t.Setenv("GH_CONFIG_DIR", "/explicit")
	t.Setenv("XDG_CONFIG_HOME", "/xdg")

	if got := ghConfigDir(); got != "/explicit" {
		t.Fatalf("got %q, want /explicit", got)
	}

	t.Setenv("GH_CONFIG_DIR", "")

	if got, want := ghConfigDir(), filepath.Join("/xdg", "gh"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	t.Setenv("XDG_CONFIG_HOME", "")

	home := t.TempDir()
	t.Setenv("HOME", home)

	if got, want := ghConfigDir(), filepath.Join(home, ".config", "gh"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
