package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type githubTokenKey struct{}

var errInvalidContext = errors.New("invalid context")

// contextGitHubToken returns the GitHub API token from context.
func contextGitHubToken(ctx context.Context) (string, error) {
	value := ctx.Value(githubTokenKey{})
	if value != nil {
		if token, ok := value.(string); ok {
			return token, nil
		}
	}

	return "", fmt.Errorf("%w: missing github token", errInvalidContext)
}

// newContext resolves the GitHub API token and prepares the on-disk cache directory,
// storing both in the returned context for later use by loadGitHub and the cache helpers.
func newContext(ctx context.Context, appname string) (context.Context, error) {
	token, err := githubToken(ctx)
	if err != nil {
		return nil, err
	}

	cacheDir, err := xdg.CacheFile(appname)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cacheDir, permDir) //nolint:forbidigo // CLI tool
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, cacheDirKey{}, cacheDir)

	return context.WithValue(ctx, githubTokenKey{}, token), nil
}

type cacheDirKey struct{}

func contextCacheDir(ctx context.Context) (string, error) {
	value := ctx.Value(cacheDirKey{})
	if value != nil {
		if client, ok := value.(string); ok {
			return client, nil
		}
	}

	return "", fmt.Errorf("%w: missing cache dir", errInvalidContext)
}

func cacheSubDir(ctx context.Context, subdir string) (string, error) {
	base, err := contextCacheDir(ctx)
	if err != nil {
		return "", err
	}

	dir := filepath.Join(base, subdir)
	if err := os.MkdirAll(dir, permDir); err != nil { //nolint:forbidigo // CLI tool
		return "", err
	}

	return dir, nil
}

func modulesDir(ctx context.Context) (string, error) {
	return cacheSubDir(ctx, "modules")
}

func checksDir(ctx context.Context) (string, error) {
	return cacheSubDir(ctx, "checks")
}
