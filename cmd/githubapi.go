package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/k6registry"
)

// githubAPIBaseURL is a var (not const) so tests can point it at an httptest server.
var githubAPIBaseURL = "https://api.github.com" //nolint:gochecknoglobals // test seam

const githubUserAgent = "k6registry"

var errGitHubAPI = errors.New("github API request failed")

type githubRepoResponse struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	License *struct {
		SPDXID string `json:"spdx_id"`
	} `json:"license"`
	PushedAt        *time.Time `json:"pushed_at"`
	HTMLURL         string     `json:"html_url"`
	Name            string     `json:"name"`
	Homepage        string     `json:"homepage"`
	Description     string     `json:"description"`
	Visibility      string     `json:"visibility"`
	CloneURL        string     `json:"clone_url"`
	Topics          []string   `json:"topics"`
	StargazersCount int        `json:"stargazers_count"`
	Archived        bool       `json:"archived"`
}

type githubTagResponse struct {
	Name string `json:"name"`
}

func githubGet(ctx context.Context, token string, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", githubUserAgent)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint:errcheck // response body close error is not actionable

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("%w: %s: %s", errGitHubAPI, resp.Status, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func loadGitHub(ctx context.Context, module string) (*k6registry.Repository, []string, error) {
	slog.Debug("Loading GitHub repository", "module", module)

	token, err := contextGitHubToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	owner, name := moduleToOwnerAndName(module)

	var rep githubRepoResponse

	repoURL := fmt.Sprintf("%s/repos/%s/%s", githubAPIBaseURL, owner, name)

	if err := githubGet(ctx, token, repoURL, &rep); err != nil {
		return nil, nil, err
	}

	repo := new(k6registry.Repository)

	repo.Topics = rep.Topics
	repo.URL = rep.HTMLURL
	repo.Name = rep.Name
	repo.Owner = rep.Owner.Login

	repo.Homepage = rep.Homepage
	if len(repo.Homepage) == 0 {
		repo.Homepage = repo.URL
	}

	repo.Archived = rep.Archived
	repo.Description = rep.Description
	repo.Stars = rep.StargazersCount

	if rep.License != nil {
		repo.License = rep.License.SPDXID
	}

	repo.Public = rep.Visibility == "public"

	if rep.PushedAt != nil && !rep.PushedAt.IsZero() {
		repo.Timestamp = float64(rep.PushedAt.Unix())
	}

	repo.CloneURL = rep.CloneURL

	const maxTags = 100

	var repoTags []githubTagResponse

	tagsURL := fmt.Sprintf("%s/repos/%s/%s/tags?per_page=%d", githubAPIBaseURL, owner, name, maxTags)

	if err := githubGet(ctx, token, tagsURL, &repoTags); err != nil {
		return nil, nil, err
	}

	tags := make([]string, 0, len(repoTags))

	for _, tag := range repoTags {
		tags = append(tags, tag.Name)
	}

	return repo, tags, nil
}
