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
		k6registry.Extension{Module: k6Module, Description: k6Description, Cloud: true, Official: true})

	for idx, ext := range registry {
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

const (
	ghModulePrefix = "github.com/"
	k6Module       = "go.k6.io/k6"
	k6Description  = "A modern load testing tool, using Go and JavaScript"
)

var errUnsupportedModule = errors.New("unsupported module")
