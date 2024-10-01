package cmd

import (
	"sort"

	"github.com/Masterminds/semver/v3"
)

func tagsToVersions(tags []string) []string {
	versions := make([]string, 0, len(tags))

	for _, tag := range tags {
		_, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}

		versions = append(versions, tag)
	}

	return versions
}

func filterVersions(tags []string, constraints *semver.Constraints) []string {
	versions := make([]string, 0, len(tags))

	for _, tag := range tags {
		version, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}

		if constraints.Check(version) {
			versions = append(versions, tag)
		}
	}

	return versions
}

func sortVersions(versions []string) error {
	type version struct {
		source string
		parsed *semver.Version
	}

	all := make([]*version, 0, len(versions))

	for _, source := range versions {
		parsed, err := semver.NewVersion(source)
		if err != nil {
			return err
		}

		all = append(all, &version{source: source, parsed: parsed})
	}

	sort.Slice(all, func(i, j int) bool { return all[i].parsed.Compare(all[j].parsed) > 0 })

	for idx := range all {
		versions[idx] = all[idx].source
	}

	return nil
}
