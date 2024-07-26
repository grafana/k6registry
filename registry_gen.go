// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package k6registry

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
	// Cloud-enabled extension flag.
	//
	// The `true` value of the `cloud` flag indicates that the extension is also
	// available in the Grafana k6 cloud.
	//
	// The use of certain extensions is not supported in a cloud environment. There
	// may be a technological reason for this, or the extension's functionality is
	// meaningless in the cloud.
	//
	Cloud bool `json:"cloud,omitempty" yaml:"cloud,omitempty" mapstructure:"cloud,omitempty"`

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

	// Officially supported extension flag.
	//
	// The `true` value of the `official` flag indicates that the extension is
	// officially supported by Grafana.
	//
	// Extensions owned by the `grafana` GitHub organization are not officially
	// supported by Grafana by default. There are several k6 extensions owned by the
	// `grafana` GitHub organization, which were created for experimental or example
	// purposes only. The `official` flag is needed so that officially supported
	// extensions can be distinguished from them.
	//
	Official bool `json:"official,omitempty" yaml:"official,omitempty" mapstructure:"official,omitempty"`

	// List of output names registered by the extension.
	//
	// The extensions used by k6 scripts are automatically detected based on the
	// values specified here, therefore it is important that the values used here are
	// consistent with the values registered by the extension at runtime.
	//
	Outputs []string `json:"outputs,omitempty" yaml:"outputs,omitempty" mapstructure:"outputs,omitempty"`

	// Repository metadata.
	//
	// Metadata provided by the extension's git repository manager. Repository
	// metadata are not registered, they are queried at runtime using the repository
	// manager API.
	//
	Repo *Repository `json:"repo,omitempty" yaml:"repo,omitempty" mapstructure:"repo,omitempty"`
}

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
	Url string `json:"url" yaml:"url" mapstructure:"url"`

	// List of supported versions.
	//
	// Versions are tags whose format meets the requirements of semantic versioning.
	// Version tags often start with the letter `v`, which is not part of the semantic
	// version.
	//
	Versions []string `json:"versions,omitempty" yaml:"versions,omitempty" mapstructure:"versions,omitempty"`
}
