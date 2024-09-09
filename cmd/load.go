package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v62/github"
	"github.com/grafana/k6registry"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

func k6AsExtension() k6registry.Extension {
	return k6registry.Extension{
		Module:      k6Module,
		Description: k6Description,
		Tier:        k6registry.TierOfficial,
		Products: []k6registry.Product{
			k6registry.ProductCloud,
			k6registry.ProductOSS,
		},
		Imports: []string{k6ImportPath},
	}
}

func load(ctx context.Context, in io.Reader, loose bool, lint bool) (k6registry.Registry, error) {
	var (
		raw []byte
		err error
	)

	if loose {
		slog.Debug("Read source")
		raw, err = io.ReadAll(in)
	} else {
		slog.Debug("Validate source")

		raw, err = validateWithSchema(in)
	}

	if err != nil {
		return nil, err
	}

	var registry k6registry.Registry

	slog.Debug("Unmarshal source")

	if err := yaml.Unmarshal(raw, &registry); err != nil {
		return nil, err
	}

	registry = append(registry, k6AsExtension())

	for idx, ext := range registry {
		slog.Debug("Process extension", "module", ext.Module)
		if len(ext.Tier) == 0 {
			registry[idx].Tier = k6registry.TierCommunity
		}

		if len(ext.Products) == 0 {
			registry[idx].Products = append(registry[idx].Products, k6registry.ProductOSS)
		}

		if len(ext.Categories) == 0 {
			registry[idx].Categories = append(registry[idx].Categories, k6registry.CategoryMisc)
		}

		if ext.Repo != nil {
			continue
		}

		repo, tags, err := loadRepository(ctx, ext.Module)
		if err != nil {
			return nil, err
		}

		registry[idx].Repo = repo

		if len(registry[idx].Versions) == 0 {
			registry[idx].Versions = filterVersions(tags)
		}

		if lint && ext.Module != k6Module {
			compliance, err := checkCompliance(ctx, ext.Module, repo.CloneURL, repo.Timestamp)
			if err != nil {
				return nil, err
			}

			registry[idx].Compliance = &k6registry.Compliance{Grade: k6registry.Grade(compliance.Grade), Level: compliance.Level}
		}
	}

	if lint {
		if err := validateWithLinter(registry); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

func loadRepository(ctx context.Context, module string) (*k6registry.Repository, []string, error) {
	slog.Debug("Loading repository", "module", module)
	if strings.HasPrefix(module, k6Module) || strings.HasPrefix(module, ghModulePrefix) {
		repo, tags, err := loadGitHub(ctx, module)
		if err != nil {
			return nil, nil, err
		}

		// Some unused metadata in the k6 repository changes too often
		if strings.HasPrefix(module, k6Module) {
			repo.Stars = 0
			repo.Timestamp = 0
			repo.CloneURL = ""
		}

		return repo, tags, nil
	}

	if strings.HasPrefix(module, glModulePrefix) {
		repo, tags, err := loadGitLab(ctx, module)
		if err != nil {
			return nil, nil, err
		}

		return repo, tags, nil
	}

	return nil, nil, fmt.Errorf("%w: %s", errUnsupportedModule, module)
}

func filterVersions(tags []string) []string {
	versions := make([]string, 0, len(tags))

	for _, tag := range tags {
		if _, err := semver.NewVersion(tag); err != nil {
			continue
		}

		versions = append(versions, tag)
	}

	return versions
}

func loadGitHub(ctx context.Context, module string) (*k6registry.Repository, []string, error) {
	client, err := contextGitHubClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	var owner, name string

	if module == k6Module {
		owner = "grafana"
		name = "k6"
	} else {
		parts := strings.SplitN(module, "/", 4)

		owner = parts[1]
		name = parts[2]
	}

	repo := new(k6registry.Repository)

	rep, _, err := client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, nil, err
	}

	repo.Topics = rep.Topics

	repo.URL = rep.GetHTMLURL()
	repo.Name = rep.GetName()
	repo.Owner = rep.GetOwner().GetLogin()

	repo.Homepage = rep.GetHomepage()
	if len(repo.Homepage) == 0 {
		repo.Homepage = repo.URL
	}

	repo.Archived = rep.GetArchived()

	repo.Description = rep.GetDescription()
	repo.Stars = rep.GetStargazersCount()

	if lic := rep.GetLicense(); lic != nil {
		repo.License = lic.GetSPDXID()
	}

	repo.Public = rep.GetVisibility() == "public"

	if ts := rep.GetPushedAt(); !ts.IsZero() {
		repo.Timestamp = float64(ts.Unix())
	}

	repo.CloneURL = rep.GetCloneURL()

	repoTags, _, err := client.Repositories.ListTags(ctx, owner, name, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, nil, err
	}

	tags := make([]string, 0, len(repoTags))

	for _, tag := range repoTags {
		tags = append(tags, tag.GetName())
	}

	return repo, tags, nil
}

func loadGitLab(ctx context.Context, module string) (*k6registry.Repository, []string, error) {
	client, err := gitlab.NewClient("")
	if err != nil {
		return nil, nil, err
	}

	pid := strings.TrimPrefix(module, glModulePrefix)

	lic := true

	proj, _, err := client.Projects.GetProject(pid, &gitlab.GetProjectOptions{License: &lic}, gitlab.WithContext(ctx))
	if err != nil {
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
	repo.Public = len(proj.Visibility) == 0 || proj.Visibility == gitlab.PublicVisibility

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

	rels, _, err := client.Releases.ListReleases(pid,
		&gitlab.ListReleasesOptions{
			ListOptions: gitlab.ListOptions{PerPage: 50},
		})
	if err != nil {
		return nil, nil, err
	}

	tags := make([]string, 0, len(rels))

	for _, rel := range rels {
		tags = append(tags, rel.TagName)
	}

	return repo, tags, nil
}

const (
	ghModulePrefix = "github.com/"
	glModulePrefix = "gitlab.com/"
	k6Module       = "go.k6.io/k6"
	k6ImportPath   = "k6"
	k6Description  = "A modern load testing tool, using Go and JavaScript"
)

var errUnsupportedModule = errors.New("unsupported module")
