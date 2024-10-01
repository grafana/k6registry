package cmd

import "github.com/Masterminds/semver/v3"

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
