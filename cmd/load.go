package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v62/github"
	"github.com/grafana/k6registry"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

func load(ctx context.Context, in io.Reader, loose bool, lint bool) (interface{}, error) {
	var (
		raw []byte
		err error
	)

	if loose {
		raw, err = io.ReadAll(in)
	} else {
		raw, err = validateWithSchema(in)
	}

	if err != nil {
		return nil, err
	}

	var registry k6registry.Registry

	if err := yaml.Unmarshal(raw, &registry); err != nil {
		return nil, err
	}

	registry = append(registry,
		k6registry.Extension{
			Module:      k6Module,
			Description: k6Description,
			Tier:        k6registry.TierOfficial,
			Products: []k6registry.Product{
				k6registry.ProductCloud,
				k6registry.ProductOss,
			},
		})

	for idx, ext := range registry {
		if len(ext.Tier) == 0 {
			registry[idx].Tier = k6registry.TierCommunity
		}

		if len(ext.Products) == 0 {
			registry[idx].Products = append(registry[idx].Products, k6registry.ProductOss)
		}

		if len(ext.Categories) == 0 {
			registry[idx].Categories = append(registry[idx].Categories, k6registry.CategoryMisc)
		}

		if ext.Repo != nil {
			continue
		}

		repo, err := loadRepository(ctx, ext.Module)
		if err != nil {
			return nil, err
		}

		registry[idx].Repo = repo
	}

	if lint {
		if err := validateWithLinter(registry); err != nil {
			return nil, err
		}
	}

	bin, err := json.Marshal(registry)
	if err != nil {
		return nil, err
	}

	var result []interface{}

	if err := json.Unmarshal(bin, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func loadRepository(ctx context.Context, module string) (*k6registry.Repository, error) {
	if strings.HasPrefix(module, k6Module) || strings.HasPrefix(module, ghModulePrefix) {
		repo, err := loadGitHub(ctx, module)
		if err != nil {
			return nil, err
		}

		// Some unused metadata in the k6 repository changes too often
		if strings.HasPrefix(module, k6Module) {
			repo.Stars = 0
			repo.Timestamp = 0
		}

		return repo, nil
	}

	if strings.HasPrefix(module, glModulePrefix) {
		repo, err := loadGitLab(ctx, module)
		if err != nil {
			return nil, err
		}

		return repo, nil
	}

	return nil, fmt.Errorf("%w: %s", errUnsupportedModule, module)
}

func loadGitHub(ctx context.Context, module string) (*k6registry.Repository, error) {
	client, err := contextGitHubClient(ctx)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	repo.Topics = rep.Topics

	repo.Url = rep.GetHTMLURL()
	repo.Name = rep.GetName()
	repo.Owner = rep.GetOwner().GetLogin()

	repo.Homepage = rep.GetHomepage()
	if len(repo.Homepage) == 0 {
		repo.Homepage = repo.Url
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

	repo.CloneUrl = rep.GetCloneURL()

	tags, _, err := client.Repositories.ListTags(ctx, owner, name, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		name := tag.GetName()

		if _, err := semver.NewVersion(name); err != nil {
			continue
		}

		repo.Versions = append(repo.Versions, name)
	}

	return repo, nil
}

func loadGitLab(ctx context.Context, module string) (*k6registry.Repository, error) {
	client, err := gitlab.NewClient("")
	if err != nil {
		return nil, err
	}

	pid := strings.TrimPrefix(module, glModulePrefix)

	lic := true

	proj, _, err := client.Projects.GetProject(pid, &gitlab.GetProjectOptions{License: &lic}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	repo := new(k6registry.Repository)

	repo.Owner = proj.Namespace.FullPath
	repo.Name = proj.Name
	repo.Description = proj.Description
	repo.Stars = proj.StarCount
	repo.Archived = proj.Archived
	repo.Url = proj.WebURL
	repo.Homepage = proj.WebURL
	repo.Topics = proj.Topics
	repo.Public = len(proj.Visibility) == 0 || proj.Visibility == gitlab.PublicVisibility

	repo.CloneUrl = proj.HTTPURLToRepo

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
		return nil, err
	}

	for _, rel := range rels {
		if _, err := semver.NewVersion(rel.TagName); err != nil {
			continue
		}

		repo.Versions = append(repo.Versions, rel.TagName)
	}

	return repo, nil
}

const (
	ghModulePrefix = "github.com/"
	glModulePrefix = "gitlab.com/"
	k6Module       = "go.k6.io/k6"
	k6Description  = "A modern load testing tool, using Go and JavaScript"
)

var errUnsupportedModule = errors.New("unsupported module")
