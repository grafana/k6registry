// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package k6registry

type Category string

const CategoryAuthentication Category = "authentication"
const CategoryBrowser Category = "browser"
const CategoryData Category = "data"
const CategoryKubernetes Category = "kubernetes"
const CategoryMessaging Category = "messaging"
const CategoryMisc Category = "misc"
const CategoryObservability Category = "observability"
const CategoryProtocol Category = "protocol"
const CategoryReporting Category = "reporting"

// The result of the extension's k6 compliance checks.
type Compliance struct {
	// Compliance expressed as a grade.
	//
	// The `grade` property contains a grade (A-F) of how well the extension complies
	// with best practices.
	// The value of the `grade` can be `A`,`B`,`C`,`D`,`E`,`F` and is calculated from
	// the `level` property.
	//
	Grade Grade `json:"grade" yaml:"grade" mapstructure:"grade"`

	// Compliance expressed as a percentage.
	//
	// The `level` property contains a percentage of how well the extension complies
	// with best practices.
	// The value of the `level` can be between `0-100` and is determined by the
	// weighted and normalized sum of the scores of the compliance checks.
	//
	Level interface{} `json:"level" yaml:"level" mapstructure:"level"`
}

// Properties of the registered k6 extension.
//
// Only those properties of the extensions are registered, which either cannot be
// detected automatically, or delegation to the extension is not allowed.
//
// Properties that are available using the repository manager API are intentionally
// not registered.
//
// The string like properties that are included in the generated Grafana
// documentation are intentionally not accessed via the API of the repository
// manager. It is not allowed to inject arbitrary text into the Grafana
// documentation site without approval. Therefore, these properties are registered
// (eg `description`)
type Extension struct {
	// The categories to which the extension belongs.
	//
	// If the property is missing or empty in the registry source, the default value
	// is "misc".
	//
	// Possible values:
	//
	//   - authentication
	//   - browser
	//   - data
	//   - kubernetes
	//   - messaging
	//   - misc
	//   - observability
	//   - protocol
	//   - reporting
	//
	Categories []Category `json:"categories,omitempty" yaml:"categories,omitempty" mapstructure:"categories,omitempty"`

	// The result of the extension's k6 compliance checks.
	//
	Compliance *Compliance `json:"compliance,omitempty" yaml:"compliance,omitempty" mapstructure:"compliance,omitempty"`

	// Brief description of the extension.
	//
	Description string `json:"description" yaml:"description" mapstructure:"description"`

	// List of JavaScript import paths registered by the extension.
	//
	// Currently, paths must start with the prefix `k6/x/`.
	//
	// The extensions used by k6 scripts are automatically detected based on the
	// values specified here, therefore it is important that the values used here are
	// consistent with the values registered by the extension at runtime.
	//
	Imports []string `json:"imports,omitempty" yaml:"imports,omitempty" mapstructure:"imports,omitempty"`

	// The extension's go module path.
	//
	// This is the unique identifier of the extension.
	// More info about module paths: https://go.dev/ref/mod#module-path
	//
	// The extension has no name property, the module path or part of it can be used
	// as the extension name. For example, using the first two elements of the module
	// path after the host name, the name `grafana/xk6-dashboard` can be formed from
	// the module path `github.com/grafana/xk6-dashboard`. This is typically the
	// repository owner name and the repository name in the repository manager.
	//
	// The extension has no URL property, a URL can be created from the module path
	// that refers to the extension within the repository manager.
	//
	Module string `json:"module" yaml:"module" mapstructure:"module"`

	// List of output names registered by the extension.
	//
	// The extensions used by k6 scripts are automatically detected based on the
	// values specified here, therefore it is important that the values used here are
	// consistent with the values registered by the extension at runtime.
	//
	Outputs []string `json:"outputs,omitempty" yaml:"outputs,omitempty" mapstructure:"outputs,omitempty"`

	// Products in which the extension can be used.
	//
	// Some extensions are not available in all k6 products.
	// This may be for a technological or business reason, or the functionality of the
	// extension may not make sense in the given product.
	//
	// Possible values:
	//
	//   - oss: Extensions are available in k6 OSS
	//   - cloud: Extensions are available in Grafana Cloud k6
	//
	// If the property is missing or empty in the source of the registry, it means
	// that the extension is only available in the k6 OSS product.
	// In this case, the registry will be filled in accordingly during generation.
	//
	//
	Products []Product `json:"products,omitempty" yaml:"products,omitempty" mapstructure:"products,omitempty"`

	// Repository metadata.
	//
	// Metadata provided by the extension's git repository manager. Repository
	// metadata are not registered, they are queried at runtime using the repository
	// manager API.
	//
	Repo *Repository `json:"repo,omitempty" yaml:"repo,omitempty" mapstructure:"repo,omitempty"`

	// Maintainer of the extension.
	//
	// Possible values:
	//
	//   - official: Extensions owned, maintained, and designated by Grafana as
	// "official"
	//   - partner: Extensions written, maintained, validated, and published by
	// third-party companies against their own projects.
	//   - community: Extensions are listed on the Registry by individual maintainers,
	// groups of maintainers, or other members of the k6 community.
	//
	// Extensions owned by the `grafana` GitHub organization are not officially
	// supported by Grafana by default.
	// There are several k6 extensions owned by the `grafana` GitHub organization,
	// which were created for experimental or example purposes only.
	// The `official` tier value is needed so that officially supported extensions can
	// be distinguished from them.
	//
	// If it is missing from the registry source, it will be set with the default
	// "community" value during generation.
	//
	//
	Tier Tier `json:"tier,omitempty" yaml:"tier,omitempty" mapstructure:"tier,omitempty"`

	// List of supported versions.
	//
	// Versions are tags whose format meets the requirements of semantic versioning.
	// Version tags often start with the letter `v`, which is not part of the semantic
	// version.
	//
	Versions []string `json:"versions,omitempty" yaml:"versions,omitempty" mapstructure:"versions,omitempty"`
}

type Grade string

const GradeA Grade = "A"
const GradeB Grade = "B"
const GradeC Grade = "C"
const GradeD Grade = "D"
const GradeE Grade = "E"
const GradeF Grade = "F"
const GradeG Grade = "G"

type Product string

const ProductCloud Product = "cloud"
const ProductOSS Product = "oss"

// k6 Extension Registry.
//
// The k6 extension registry contains the most important properties of registered
// extensions.
type Registry []Extension

// Repository metadata.
//
// Metadata provided by the extension's git repository manager. Repository metadata
// are not registered, they are queried at runtime using the repository manager
// API.
type Repository struct {
	// Archived repository flag.
	//
	// A `true` value indicates that the repository is archived, read only.
	//
	// If a repository is archived, it usually means that the owner has no intention
	// of maintaining it. Such extensions should be removed from the registry.
	//
	Archived bool `json:"archived,omitempty" yaml:"archived,omitempty" mapstructure:"archived,omitempty"`

	// URL for the git clone operation.
	//
	// The clone_url property contains a (typically HTTP) URL, which is used to clone
	// the repository.
	//
	CloneURL string `json:"clone_url,omitempty" yaml:"clone_url,omitempty" mapstructure:"clone_url,omitempty"`

	// Repository description.
	//
	Description string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The URL to the project homepage.
	//
	// If no homepage is set, the value is the same as the url property.
	//
	Homepage string `json:"homepage,omitempty" yaml:"homepage,omitempty" mapstructure:"homepage,omitempty"`

	// The SPDX ID of the extension's license.
	//
	// For more information about SPDX, visit https://spdx.org/licenses/
	//
	License string `json:"license,omitempty" yaml:"license,omitempty" mapstructure:"license,omitempty"`

	// The name of the repository.
	//
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The owner of the repository.
	//
	Owner string `json:"owner" yaml:"owner" mapstructure:"owner"`

	// Public repository flag.
	//
	// A `true` value indicates that the repository is public, available to anyone.
	//
	Public bool `json:"public,omitempty" yaml:"public,omitempty" mapstructure:"public,omitempty"`

	// The number of stars in the extension's repository.
	//
	// The extension's popularity is indicated by how many users have starred the
	// extension's repository.
	//
	Stars int `json:"stars,omitempty" yaml:"stars,omitempty" mapstructure:"stars,omitempty"`

	// Last modification timestamp.
	//
	// The timestamp property contains the timestamp of the last modification of the
	// repository in UNIX time format (the number of non-leap seconds that have
	// elapsed since 00:00:00 UTC on 1st January 1970).
	// Its value depends on the repository manager, in the case of GitHub it contains
	// the time of the last push operation, in the case of GitLab the time of the last
	// repository activity.
	//
	Timestamp float64 `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Repository topics.
	//
	// Topics make it easier to find the repository. It is recommended to set the xk6
	// topic to the extensions repository.
	//
	Topics []string `json:"topics,omitempty" yaml:"topics,omitempty" mapstructure:"topics,omitempty"`

	// URL of the repository.
	//
	// The URL is provided by the repository manager and can be displayed in a
	// browser.
	//
	URL string `json:"url" yaml:"url" mapstructure:"url"`
}

type Tier string

const TierCommunity Tier = "community"
const TierOfficial Tier = "official"
const TierPartner Tier = "partner"
