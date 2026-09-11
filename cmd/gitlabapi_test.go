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

func withGitLabAPIBaseURL(t *testing.T, url string) {
	t.Helper()

	old := gitlabAPIBaseURL
	gitlabAPIBaseURL = url

	t.Cleanup(func() { gitlabAPIBaseURL = old })
}

func TestLoadGitLab(t *testing.T) { //nolint:paralleltest // mutates the shared gitlabAPIBaseURL test seam; must run serially
	mux := http.NewServeMux()

	mux.HandleFunc("/projects/grafana%2Fxk6-faker", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"namespace":        map[string]any{"full_path": "grafana"},
			"name":             "xk6-faker",
			"description":      testFakerDescription,
			"star_count":       7,
			"archived":         false,
			"web_url":          "https://gitlab.com/grafana/xk6-faker",
			"topics":           []string{"k6", "extension"},
			"visibility":       "public",
			"http_url_to_repo": "https://gitlab.com/grafana/xk6-faker.git",
			"last_activity_at": "2024-01-02T03:04:05Z",
			"license":          map[string]any{"key": "mit"},
		})
	})

	mux.HandleFunc("/projects/grafana%2Fxk6-faker/releases", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]string{{"tag_name": "v1.0.0"}, {"tag_name": "v1.1.0"}})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	withGitLabAPIBaseURL(t, server.URL)

	repo, tags, err := loadGitLab(context.Background(), "gitlab.com/grafana/xk6-faker")
	if err != nil {
		t.Fatal(err)
	}

	want := &k6registry.Repository{
		Owner:       "grafana",
		Name:        "xk6-faker",
		Description: testFakerDescription,
		Stars:       7,
		URL:         "https://gitlab.com/grafana/xk6-faker",
		Homepage:    "https://gitlab.com/grafana/xk6-faker",
		Topics:      []string{"k6", "extension"},
		Public:      true,
		CloneURL:    "https://gitlab.com/grafana/xk6-faker.git",
		Timestamp:   float64(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC).Unix()),
		License:     "MIT",
	}

	if !reflect.DeepEqual(repo, want) {
		t.Fatalf("got %+v\nwant %+v", repo, want)
	}

	if !reflect.DeepEqual(tags, []string{"v1.0.0", "v1.1.0"}) {
		t.Fatalf("got tags %v", tags)
	}
}

func TestLoadGitLab_ErrorStatus(t *testing.T) { //nolint:paralleltest // mutates the shared gitlabAPIBaseURL test seam; must run serially
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"404 Project Not Found"}`))
	}))
	defer server.Close()

	withGitLabAPIBaseURL(t, server.URL)

	if _, _, err := loadGitLab(context.Background(), "gitlab.com/does/not-exist"); err == nil {
		t.Fatal("expected an error for a 404 response")
	}
}
