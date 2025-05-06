package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/config"
	"github.com/google/go-github/v62/github"
)

type githubClientKey struct{}

var errInvalidContext = errors.New("invalid context")

// contextGitHubClient returns a *github.Client from context.
func contextGitHubClient(ctx context.Context) (*github.Client, error) {
	value := ctx.Value(githubClientKey{})
	if value != nil {
		if client, ok := value.(*github.Client); ok {
			return client, nil
		}
	}

	return nil, fmt.Errorf("%w: missing github.Client", errInvalidContext)
}

// newContext prepares GitHub CLI extension context with http.Client and github.Client values.
// You can use ContextHTTPClient and ContextGitHubClient later to get client instances from the context.
func newContext(ctx context.Context, appname string) (context.Context, error) {
	htc, err := newHTTPClient()
	if err != nil {
		return nil, err
	}

	cacheDir, err := xdg.CacheFile(appname)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cacheDir, permDir)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, cacheDirKey{}, cacheDir)

	return context.WithValue(ctx, githubClientKey{}, github.NewClient(htc)), nil
}

const cacheTTL = 2 * time.Hour

var errMissingAuthToken = errors.New("missing authentication token")

func newHTTPClient() (*http.Client, error) {
	var opts api.ClientOptions

	opts.Host, _ = auth.DefaultHost()

	opts.AuthToken, _ = auth.TokenForHost(opts.Host)
	if opts.AuthToken == "" {
		return nil, fmt.Errorf("%w: host %s", errMissingAuthToken, opts.Host)
	}

	if cfg, _ := config.Read(nil); cfg != nil {
		opts.UnixDomainSocket, _ = cfg.Get([]string{"http_unix_socket"})
	}

	opts.EnableCache = true
	opts.CacheTTL = cacheTTL

	return api.NewHTTPClient(opts)
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
	if err := os.MkdirAll(dir, permDir); err != nil {
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
