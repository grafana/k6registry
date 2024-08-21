# k6 Extension Registry Concept

## Registry Source

The k6 extension registry source is a YAML file that contains the most important properties of extensions.

### File format

The k6 extension registry format is YAML, because the registry is edited by humans and the YAML format is more human-friendly than JSON. The files generated from the registry are typically in JSON format, because they are processed by programs and JSON is more widely supported than YAML. A JSON format is also generated from the entire registry, so that it can also be processed by programs.

### Registered Properties

Only those properties of the extensions are registered, which either cannot be detected automatically, or delegation to the extension is not allowed.

Properties that are available using the repository manager API (GitHub API, GitLab API, etc) are intentionally not registered. For example, the number of stars can be queried via the repository manager API, so this property is not registered.

Exceptions are the string-like properties that are embedded in the Grafana documentation. These properties are registered because it is not allowed to inject arbitrary text into the Grafana documentation site without approval. Therefore, these properties are registered (eg `description`)

The properties provided by the repository managers ([Repository Metadata]) are queried during registry processing and can be used to produce the output properties.

### Extension Identification

The primary identifier of an extension is the extension's [go module path].

The extension does not have a `name` property, the [repository metadata] can be used to construct a `name` property. Using the repository owner and the repository name, for example, `grafana/xk6-dashboard` can be generated for the `github.com/grafana/xk6-dashboard` extension.

The extension does not have a `url` property, but there is a `url` property in the [repository metadata].

[go module path]: https://go.dev/ref/mod#module-path
[Repository Metadata]: #repository-metadata

### JavaScript Modules

The JavaScript module names implemented by the extension can be specified in the `imports` property. An extension can register multiple JavaScript module names, so this is an array property.

### Output Names

The output names implemented by the extension can be specified in the `outputs` property. An extension can register multiple output names, so this is an array property.

### Tier

Extensions can be classified according to who maintains the extension. This usually also specifies who the user can get support from.

The `tier` property refers to the maintainer of the extension.

Possible values:

  - **official**: Extensions owned, maintained, and designated by Grafana as "official"
  - **partner**: Extensions written, maintained, validated, and published by third-party companies against their own projects.
  - **community**: Extensions are listed on the Registry by individual maintainers, groups of maintainers, or other members of the k6 community.

Extensions owned by the `grafana` GitHub organization are not officially supported by Grafana by default. There are several k6 extensions owned by the `grafana` GitHub organization, which were created for experimental or example purposes only. The `official` tier value is needed so that officially supported extensions can be distinguished from them.

If it is missing from the registry source, it will be set with the default `community` value during generation.

### Product

The `product` property contains the names of the k6 products in which the extension is available.

Some extensions are not available in all k6 products. This may be for a technological or business reason, or the functionality of the extension may not make sense in the given product.

Possible values:

  - **oss**: Extensions are available in *k6 OSS*
  - **cloud**: Extensions are available in *Grafana Cloud k6*

If the property is missing or empty in the source of the registry, it means that the extension is only available in the *k6 OSS* product. In this case, the registry will be filled in accordingly during generation.

### Example registry

```yaml file=example.yaml
- module: github.com/grafana/xk6-dashboard
  description: Web-based metrics dashboard for k6
  outputs:
    - dashboard
  tier: official
  categories:
    - reporting
    - observability

- module: github.com/grafana/xk6-sql
  description: Load test SQL Servers
  imports:
    - k6/x/sql
  tier: official
  product: ["cloud", "oss"]
  categories:
    - data

- module: github.com/grafana/xk6-disruptor
  description: Inject faults to test
  imports:
    - k6/x/disruptor
  tier: official
  categories:
    - kubernetes

- module: github.com/szkiba/xk6-faker
  description: Generate random fake data
  imports:
    - k6/x/faker
  categories:
    - data
```

### Categories

The `categories` property contains the categories to which the extension belongs.

If the property is missing or empty in the registry source, the default value is "misc".

Possible values:

  - **authentication**
  - **browser**
  - **data**
  - **kubernetes**
  - **messaging**
  - **misc**
  - **observability**
  - **protocol**
  - **reporting**

### Repository Metadata

Repository metadata provided by the extension's git repository manager. Repository metadata are not registered, they are queried at processing time using the repository manager API.

#### Owner

The `owner` property contains the owner of the extension's git repository.

#### Name

The `name` property contains the name of the extension's git repository.

#### License

The `license` property contains the SPDX ID of the extension's license. For more information about SPDX, visit https://spdx.org/licenses/

#### Public

The `true` value of the `public` flag indicates that the repository is public, available to anyone.

#### URL

The `url` property contains the URL of the repository. The `url` is provided by the repository manager and can be displayed in a browser.

#### Homepage

The `homepage` property contains the project homepage URL. If no homepage is set, the value is the same as the `url` property.

#### Stars

The `stars` property contains the number of stars in the extension's repository. The extension's popularity is indicated by how many users have starred the extension's repository.

#### Topics

The `topics` property contains the repository topics. Topics make it easier to find the repository. It is recommended to set the `xk6` topic to the extensions repository.

#### Versions

The `versions` property contains the list of supported versions. Versions are tags whose format meets the requirements of semantic versioning. Version tags often start with the letter `v`, which is not part of the semantic version.

#### Archived

The `true` value of the `archived` flag indicates that the repository is archived, read only.

If a repository is archived, it usually means that the owner has no intention of maintaining it. Such extensions should be removed from the registry.

## Registry Processing

The source of the registry is a YAML file optimized for human use. For programs that use the registry, it is advisable to generate an output in JSON format optimized for the given application (for example, an extension catalog for the k6build service).

```mermaid
---
title: registry processing
---
erDiagram
    "registry processor" ||--|| "extension registry" : input
    "registry processor" ||--|{ "repository manager" : metadata
    "registry processor" ||--|| "jq expression" : filter
    "registry processor" ||--|| "custom JSON" : output
    "custom JSON" }|--|{ "application" : uses
```

The registry is processed based on the popular `jq` expressions.

The input of the processing is the extension registry supplemented with repository metadata (for example, available versions, number of stars, etc). The output of the processing is defined by the `jq` filter.

### Registry Validation

The registry is validated using [JSON schema](https://grafana.github.io/k6registry/registry.schema.json). Requirements that cannot be validated using the JSON schema are validated using custom linter.

Custom linter checks the following for each extension:

  - Is the go module path valid?
  - Is there at least one versioned release?
  - Is a valid license configured?
  - Is the xk6 topic set for the repository?
  - Is the repository not archived?

It is strongly recommended to lint the extension registry after each modification, but at least before approving the change.
