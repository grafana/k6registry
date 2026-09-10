package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/grafana/k6registry"
)

// gitlabAPIBaseURL is a var (not const) so tests can point it at an httptest server.
var gitlabAPIBaseURL = "https://gitlab.com/api/v4" //nolint:gochecknoglobals // test seam

var errGitLabAPI = errors.New("gitlab API request failed")

type gitlabProjectResponse struct {
	Namespace struct {
		FullPath string `json:"full_path"`
	} `json:"namespace"`
	License *struct {
		Key string `json:"key"`
	} `json:"license"`
	LastActivityAt *time.Time `json:"last_activity_at"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	WebURL         string     `json:"web_url"`
	Visibility     string     `json:"visibility"`
	HTTPURLToRepo  string     `json:"http_url_to_repo"`
	Topics         []string   `json:"topics"`
	StarCount      int        `json:"star_count"`
	Archived       bool       `json:"archived"`
}

type gitlabReleaseResponse struct {
	TagName string `json:"tag_name"`
}

func gitlabGet(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint:errcheck // response body close error is not actionable

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("%w: %s: %s", errGitLabAPI, resp.Status, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func loadGitLab(ctx context.Context, module string) (*k6registry.Repository, []string, error) {
	slog.Debug("Loading GitLab repository", "module", module)

	pid := url.QueryEscape(strings.TrimPrefix(module, glModulePrefix))

	var proj gitlabProjectResponse

	projectURL := fmt.Sprintf("%s/projects/%s?license=true", gitlabAPIBaseURL, pid)

	if err := gitlabGet(ctx, projectURL, &proj); err != nil {
		return nil, nil, err
	}

	repo := new(k6registry.Repository)

	repo.Owner = proj.Namespace.FullPath
	repo.Name = proj.Name
	repo.Description = proj.Description
	repo.Stars = proj.StarCount
	repo.Archived = proj.Archived
	repo.URL = proj.WebURL
	repo.Homepage = proj.WebURL
	repo.Topics = proj.Topics
	repo.Public = len(proj.Visibility) == 0 || proj.Visibility == "public"

	repo.CloneURL = proj.HTTPURLToRepo

	if proj.LastActivityAt != nil {
		repo.Timestamp = float64(proj.LastActivityAt.Unix())
	}

	if proj.License != nil {
		for key := range validLicenses {
			if strings.EqualFold(key, proj.License.Key) {
				repo.License = key
			}
		}
	}

	const maxTags = 50

	var releases []gitlabReleaseResponse

	releasesURL := fmt.Sprintf("%s/projects/%s/releases?per_page=%d", gitlabAPIBaseURL, pid, maxTags)

	if err := gitlabGet(ctx, releasesURL, &releases); err != nil {
		return nil, nil, err
	}

	tags := make([]string, 0, len(releases))

	for _, rel := range releases {
		tags = append(tags, rel.TagName)
	}

	return repo, tags, nil
}
