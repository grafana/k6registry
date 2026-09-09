package cmd //nolint:testpackage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/grafana/k6registry"
)

func withGitHubAPIBaseURL(t *testing.T, url string) {
	t.Helper()

	old := githubAPIBaseURL
	githubAPIBaseURL = url

	t.Cleanup(func() { githubAPIBaseURL = old })
}

func TestLoadGitHub(t *testing.T) { //nolint:paralleltest // mutates the shared githubAPIBaseURL test seam; must run serially
	var gotAuth string
	var gotUserAgent string

	mux := http.NewServeMux()

	mux.HandleFunc("/repos/grafana/xk6-faker", func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUserAgent = r.Header.Get("User-Agent")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"html_url":         "https://github.com/grafana/xk6-faker",
			"name":             "xk6-faker",
			"owner":            map[string]any{"login": "grafana"},
			"homepage":         "",
			"archived":         false,
			"description":      "a faker extension",
			"stargazers_count": 42,
			"license":          map[string]any{"spdx_id": "MIT"},
			"visibility":       "public",
			"pushed_at":        "2024-01-02T03:04:05Z",
			"clone_url":        "https://github.com/grafana/xk6-faker.git",
			"topics":           []string{"k6", "extension"},
		})
	})

	mux.HandleFunc("/repos/grafana/xk6-faker/tags", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]string{{"name": "v1.0.0"}, {"name": "v1.1.0"}})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	withGitHubAPIBaseURL(t, server.URL)

	ctx := context.WithValue(context.Background(), githubTokenKey{}, "test-token")

	repo, tags, err := loadGitHub(ctx, "github.com/grafana/xk6-faker")
	if err != nil {
		t.Fatal(err)
	}

	if gotAuth != "Bearer test-token" {
		t.Fatalf("got Authorization header %q", gotAuth)
	}

	if gotUserAgent != githubUserAgent {
		t.Fatalf("got User-Agent header %q, want %q", gotUserAgent, githubUserAgent)
	}

	want := &k6registry.Repository{
		URL:         "https://github.com/grafana/xk6-faker",
		Name:        "xk6-faker",
		Owner:       "grafana",
		Homepage:    "https://github.com/grafana/xk6-faker",
		Description: "a faker extension",
		Stars:       42,
		License:     "MIT",
		Public:      true,
		CloneURL:    "https://github.com/grafana/xk6-faker.git",
		Topics:      []string{"k6", "extension"},
		Timestamp:   float64(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC).Unix()),
	}

	if !reflect.DeepEqual(repo, want) {
		t.Fatalf("got %+v\nwant %+v", repo, want)
	}

	if !reflect.DeepEqual(tags, []string{"v1.0.0", "v1.1.0"}) {
		t.Fatalf("got tags %v", tags)
	}
}

func TestLoadGitHub_NoHomepage(t *testing.T) { //nolint:paralleltest // mutates the shared githubAPIBaseURL test seam; must run serially
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/grafana/xk6-faker", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"html_url": "https://github.com/grafana/xk6-faker",
			"name":     "xk6-faker",
			"owner":    map[string]any{"login": "grafana"},
		})
	})

	mux.HandleFunc("/repos/grafana/xk6-faker/tags", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]string{})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	withGitHubAPIBaseURL(t, server.URL)

	ctx := context.WithValue(context.Background(), githubTokenKey{}, "test-token")

	repo, _, err := loadGitHub(ctx, "github.com/grafana/xk6-faker")
	if err != nil {
		t.Fatal(err)
	}

	if repo.Homepage != repo.URL {
		t.Fatalf("expected homepage to fall back to repo URL, got %q", repo.Homepage)
	}
}

func TestLoadGitHub_ErrorStatus(t *testing.T) { //nolint:paralleltest // mutates the shared githubAPIBaseURL test seam; must run serially
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	withGitHubAPIBaseURL(t, server.URL)

	ctx := context.WithValue(context.Background(), githubTokenKey{}, "test-token")

	if _, _, err := loadGitHub(ctx, "github.com/grafana/does-not-exist"); err == nil {
		t.Fatal("expected an error for a 404 response")
	}
}

func TestLoadGitHub_MissingToken(t *testing.T) { //nolint:paralleltest // mutates the shared githubAPIBaseURL test seam; must run serially
	if _, _, err := loadGitHub(context.Background(), "github.com/grafana/xk6-faker"); err == nil {
		t.Fatal("expected an error when the context has no github token")
	}
}
