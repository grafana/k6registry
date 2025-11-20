package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v62/github"
	"github.com/grafana/k6registry"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"gopkg.in/yaml.v3"
)

var errCompliance = errors.New("compliance check failed")

type loadOptions struct {
	lint             bool
	ignoreLintErrors bool
	lintChecks       []string
}

func k6AsExtension() k6registry.Extension {
	return k6registry.Extension{
		Module:      k6Module,
		Description: k6Description,
		Tier:        k6registry.TierOfficial,
		Imports:     []string{k6ImportPath},
	}
}

func loadSource(in io.Reader) (k6registry.Registry, error) {
	var (
		raw []byte
		err error
	)

	slog.Debug("Validate source")

	raw, err = validateWithSchema(in)
	if err != nil {
		return nil, err
	}

	var registry k6registry.Registry

	slog.Debug("Unmarshal source")

	if err := yaml.Unmarshal(raw, &registry); err != nil {
		return nil, err
	}

	k6 := false

	for idx := range registry {
		if registry[idx].Module == k6Module {
			k6 = true

			break
		}
	}

	if !k6 {
		registry = append(registry, k6AsExtension())
	}

	return registry, nil
}

func loadOne(ctx context.Context, ext *k6registry.Extension, lint bool, checks []string) error {
	if len(ext.Tier) == 0 {
		ext.Tier = k6registry.TierCommunity
	}

	repo, tags, err := loadRepository(ctx, ext)
	if err != nil {
		return err
	}

	ext.Repo = repo

	if len(ext.Versions) == 0 {
		ext.Versions = tagsToVersions(tags)
	}

	if !lint || ext.Module == k6Module {
		return nil
	}

	if ext.Compliance == nil {
		ext.Compliance = make(k6registry.ExtensionCompliance)
	}

	complianceErrors := []error{}
	for _, version := range ext.Versions {
		official := ext.Tier == k6registry.TierOfficial

		compliance, err := checkCompliance(
			ctx,
			ext.Module,
			version,
			official,
			checks,
			repo.CloneURL,
			int64(repo.Timestamp),
		)
		if err != nil {
			return err
		}

		var issues []string

		for _, check := range compliance.Checks {
			if !check.Passed {
				issues = append(issues, string(check.ID))
			}
		}

		if len(issues) > 0 {
			complianceErrors = append(complianceErrors, fmt.Errorf("%w %s@%s", errCompliance, ext.Module, version))
		}

		ext.Compliance[version] = k6registry.Compliance{
			Issues: issues,
		}
	}

	return errors.Join(complianceErrors...)
}

func load(
	ctx context.Context,
	in io.Reader,
	opts loadOptions,
) (k6registry.Registry, error) {
	registry, err := loadSource(in)
	if err != nil {
		return nil, err
	}

	compliancedErrors := []error{}
	for idx := range registry {
		ext := &registry[idx]

		slog.Debug("Process extension", "module", ext.Module)

		err := loadOne(ctx, ext, opts.lint, opts.lintChecks)
		if err != nil {
			if !errors.Is(err, errCompliance) {
				return nil, err
			}

			compliancedErrors = append(compliancedErrors, err)
		}

		if len(ext.Constraints) > 0 {
			constraints, err := semver.NewConstraint(ext.Constraints)
			if err != nil {
				return nil, err
			}

			ext.Versions = filterVersions(ext.Versions, constraints)
		}

		if err := sortVersions(ext.Versions); err != nil {
			return nil, err
		}
	}

	if len(compliancedErrors) == 0 {
		return registry, nil
	}

	slog.Warn(errors.Join(compliancedErrors...).Error())

	if opts.ignoreLintErrors {
		return registry, nil
	}

	return registry, errCompliance
}

func loadRepository(ctx context.Context, ext *k6registry.Extension) (*k6registry.Repository, []string, error) {
	module := ext.Module

	if ext.Repo != nil && len(ext.Repo.CloneURL) > 0 {
		versions, err := loadGit(ctx, module, ext.Repo.CloneURL)
		if err != nil {
			return nil, nil, err
		}

		return ext.Repo, versions, nil
	}

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

func moduleToOwnerAndName(module string) (string, string) {
	if module == k6Module {
		return "grafana", "k6"
	}

	const maxParts = 4
	parts := strings.SplitN(module, "/", maxParts)

	return parts[1], parts[2]
}

func loadGitHub(ctx context.Context, module string) (*k6registry.Repository, []string, error) {
	slog.Debug("Loading GitHub repository", "module", module)

	client, err := contextGitHubClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	owner, name := moduleToOwnerAndName(module)

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

	const maxTags = 100

	repoTags, _, err := client.Repositories.ListTags(ctx, owner, name, &github.ListOptions{PerPage: maxTags})
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
	slog.Debug("Loading GitLab repository", "module", module)

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

	const maxTags = 50

	rels, _, err := client.Releases.ListReleases(pid,
		&gitlab.ListReleasesOptions{
			ListOptions: gitlab.ListOptions{PerPage: maxTags},
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

func loadGit(ctx context.Context, module string, cloneURL string) ([]string, error) {
	base, err := modulesDir(ctx)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(base, module)

	repo, err := openRepo(ctx, dir, cloneURL)
	if err != nil {
		return nil, err
	}

	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	const tagPrefix = "refs/tags/"

	versions := make([]string, 0)

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag := strings.TrimPrefix(ref.Name().String(), tagPrefix)

		if _, err := semver.NewVersion(tag); err == nil {
			versions = append(versions, tag)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func openRepo(ctx context.Context, dir string, cloneURL string) (*git.Repository, error) {
	_, err := os.Stat(dir)
	notfound := (err != nil && errors.Is(err, os.ErrNotExist))

	if err != nil && !notfound {
		return nil, err
	}

	if notfound {
		slog.Debug("Clone", "url", cloneURL)

		return git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{URL: cloneURL})
	}

	return git.PlainOpen(dir)
}

func checkoutModVersion(ctx context.Context, dir string, cloneURL string, version string) error {
	repo, err := openRepo(ctx, dir, cloneURL)
	if err != nil {
		return err
	}

	wtree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// If a version is specified, fetch and checkout that tag
	if version != "" {
		slog.Debug("Fetch tags", "url", cloneURL)

		// Fetch all tags to ensure we have the requested one
		err = repo.FetchContext(ctx, &git.FetchOptions{
			RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
			Force:    true,
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}

		// Checkout the specific tag
		tagRef := plumbing.NewTagReferenceName(version)
		slog.Debug("Checkout tag", "tag", version)

		err = wtree.Checkout(&git.CheckoutOptions{Force: true, Branch: tagRef})

		return err
	}

	// No version specified, use default branch
	head, err := repo.Head()
	if err != nil {
		return err
	}

	err = wtree.Checkout(&git.CheckoutOptions{Force: true, Branch: head.Name()})
	if err != nil {
		return err
	}

	slog.Debug("Pull", "url", cloneURL)

	err = wtree.Pull(&git.PullOptions{Force: true})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}

const (
	ghModulePrefix = "github.com/"
	glModulePrefix = "gitlab.com/"
	k6Module       = "go.k6.io/k6"
	k6ImportPath   = "k6"
	k6Description  = "A modern load testing tool, using Go and JavaScript"
)

var errUnsupportedModule = errors.New("unsupported module")
