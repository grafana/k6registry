package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

const githubHost = "github.com"

var errMissingAuthToken = errors.New("missing authentication token")

// githubToken resolves a GitHub API token for github.com, mirroring the gh
// CLI's own precedence: GH_TOKEN/GITHUB_TOKEN env vars, then gh's hosts.yml
// config file, then (if the gh binary is available) `gh auth token`, which
// can also read tokens stored in the OS keyring.
func githubToken(ctx context.Context) (string, error) {
	if token := os.Getenv("GH_TOKEN"); token != "" { //nolint:forbidigo // gh CLI env var
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" { //nolint:forbidigo // gh CLI env var
		return token, nil
	}

	if token := tokenFromGHConfig(); token != "" {
		return token, nil
	}

	if token := tokenFromGHCli(ctx); token != "" {
		return token, nil
	}

	return "", fmt.Errorf("%w: host %s", errMissingAuthToken, githubHost)
}

type ghHostConfig struct {
	OauthToken string `yaml:"oauth_token"`
}

func tokenFromGHConfig() string {
	filename := filepath.Join(ghConfigDir(), "hosts.yml")

	data, err := os.ReadFile(filepath.Clean(filename)) //nolint:forbidigo // gh CLI config file
	if err != nil {
		return ""
	}

	var hosts map[string]ghHostConfig

	if err := yaml.Unmarshal(data, &hosts); err != nil {
		return ""
	}

	return hosts[githubHost].OauthToken
}

// ghConfigDir mirrors the gh CLI's own config directory precedence:
// GH_CONFIG_DIR, XDG_CONFIG_HOME, AppData (Windows only), then $HOME/.config/gh.
func ghConfigDir() string {
	if dir := os.Getenv("GH_CONFIG_DIR"); dir != "" { //nolint:forbidigo // gh CLI env var
		return dir
	}

	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" { //nolint:forbidigo // gh CLI env var
		return filepath.Join(dir, "gh")
	}

	if runtime.GOOS == "windows" {
		if dir := os.Getenv("AppData"); dir != "" { //nolint:forbidigo // gh CLI env var
			return filepath.Join(dir, "GitHub CLI")
		}
	}

	home, _ := os.UserHomeDir() //nolint:forbidigo // gh CLI config file

	return filepath.Join(home, ".config", "gh")
}

func tokenFromGHCli(ctx context.Context) string {
	ghExe := os.Getenv("GH_PATH") //nolint:forbidigo // gh CLI env var

	if ghExe == "" {
		var err error

		ghExe, err = exec.LookPath("gh")
		if err != nil {
			return ""
		}
	}

	//nolint:gosec // fixed subcommand/flags, ghExe is either GH_PATH or resolved via LookPath
	out, err := exec.CommandContext(ctx, ghExe, "auth", "token", "--secure-storage", "--hostname", githubHost).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
