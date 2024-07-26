package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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
func newContext(ctx context.Context) (context.Context, error) {
	htc, err := newHTTPClient()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, githubClientKey{}, github.NewClient(htc)), nil
}

func newHTTPClient() (*http.Client, error) {
	var opts api.ClientOptions

	opts.Host, _ = auth.DefaultHost()

	opts.AuthToken, _ = auth.TokenForHost(opts.Host)
	if opts.AuthToken == "" {
		return nil, fmt.Errorf("authentication token not found for host %s", opts.Host)
	}

	if cfg, _ := config.Read(nil); cfg != nil {
		opts.UnixDomainSocket, _ = cfg.Get([]string{"http_unix_socket"})
	}

	opts.EnableCache = true
	opts.CacheTTL = 2 * time.Hour

	return api.NewHTTPClient(opts)
}
