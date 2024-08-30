<h1 name="title">k6registry</h1>

**Data model and tooling for the k6 extension registry**

This repository contains the [JSON schema](docs/registry.schema.json) of the k6 extension registry and the [`k6registry`](#k6registry) command line tool for generating registry from source. The command line tool can also be used as a [GitHub Action](#github-action).

Check [k6 Extension Registry Concept](docs/registry.md) for information on design considerations.

**Example registry source**

```yaml file=docs/example.yaml
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
  products: ["cloud", "oss"]
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

- module: gitlab.com/szkiba/xk6-banner
  description: Print ASCII art banner from k6 test
  imports:
    - k6/x/banner
  categories:
    - misc
```

<details>
<summary><b>Example registry</b></summary>

Registry generated from the source above.

```json file=docs/example.json
```

</details>

A [legacy extension registry](docs/legacy.yaml) converted to the new format is also a good example.

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6registry/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6registry/cmd/k6registry@latest
```

## GitHub Action

`grafana/k6registry` is a GitHub Action that enables k6 extension registry processing and the generation of customized JSON output for different applications. Processing is based on popular `jq` expressions using an embedded `jq` implementation.

The jq filter expression can be specified in the `filter` parameter.

The extension registry is read from the YAML format file specified in the `in` parameter.

Repository metadata is collected using the repository manager APIs. Currently only the GitHub API is supported.

The output of the processing will be written to the standard output by default. The output can be saved to a file using the `out` parameter.

**Inputs**

name   | reqired | default | description
-------|---------|---------|-------------
in     |   yes   |         | input file name
out    |    no   |  stdout | output file name
api    |    no   |         | output directory name
mute   |    no   | `false` | no output, only validation
loose  |    no   | `false` | skip JSON schema validation
lint   |    no   | `false` | enable built-in linter
compact|    no   | `false` | compact instead of pretty-printed output
ref    |    no   |         | reference output URL for change detection

In GitHub action mode, the change can be indicated by comparing the output to a reference output. The reference output URL can be passed in the `ref` action parameter. The `changed` output variable will be `true` or `false` depending on whether the output has changed or not compared to the reference output.

The `api` parameter can be used to specify a directory into which the outputs are written. The `registry.json` file is placed in the root directory. The `extension.json` file and the `badge.svg` file are placed in a directory with the same name as the go module path of the extension (if the `lint` parameter is `true`).

**Outputs**

name    | description
--------|------------
changed | `true` if the output has changed compared to `ref`, otherwise `false`

**Example usage**

```yaml
- name: Generate registry in JSON format
  uses: grafana/k6registry@v0.1.10
  with:
    in: "registry.yaml"
    out: "registry.json"
    lint: "true"
```

## CLI

<!-- #region cli -->
## k6registry

k6 extension registry generator

### Synopsis

Command line k6 extension registry generator.

The source of the extension registry contains only the most important properties of the extensions. The rest of the properties are collected by k6registry using the API of the extensions' git repository managers.

The source of the extension registry is read from the YAML (or JSON) format file specified as command line argument. If it is missing, the source is read from the standard input.

Repository metadata is collected using the API of the extensions' git repository managers. Currently only the GitHub API is supported.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.

The `--api` flag can be used to specify a directory to which the outputs will be written. The `registry.json` file is placed in the root directory. The `extension.json` file and the `badge.svg` file (if the `--lint` flag is used) are placed in a directory with the same name as the extension's go module path.

```
k6registry [flags] [source-file]
```

### Flags

```
  -o, --out string   write output to file instead of stdout
      --api string   write outputs to directory instead of stdout
  -q, --quiet        no output, only validation
      --loose        skip JSON schema validation
      --lint         enable built-in linter
  -c, --compact      compact instead of pretty-printed output
  -V, --version      print version
  -h, --help         help for k6registry
```

<!-- #endregion cli -->

## API in the filesystem

By using the `--api` flag, files are created with relative paths in a base directory with a kind of REST API logic:
- in the `module` directory, a directory with the same name as the path of the extension module
   - `badge.svg` badge generated based on the compliance grade
   - `extension.json` extension data in a separate file
- the subdirectories of the base directory contain subsets of the registry broken down according to different properties (`tier`, `product`, `category`, `grade`)

```ascii file=docs/example-api.txt
docs/example-api
├── registry.json
├── registry.schema.json
├── category
│   ├── authentication.json
│   ├── browser.json
│   ├── data.json
│   ├── kubernetes.json
│   ├── messaging.json
│   ├── misc.json
│   ├── observability.json
│   ├── protocol.json
│   └── reporting.json
├── grade
│   ├── A.json
│   ├── B.json
│   ├── C.json
│   ├── D.json
│   ├── E.json
│   ├── F.json
│   └── passing
│       ├── A.json
│       ├── B.json
│       ├── C.json
│       ├── D.json
│       ├── E.json
│       └── F.json
├── module
│   ├── github.com
│   │   ├── grafana
│   │   │   ├── xk6-dashboard
│   │   │   │   ├── badge.svg
│   │   │   │   └── extension.json
│   │   │   ├── xk6-disruptor
│   │   │   │   ├── badge.svg
│   │   │   │   └── extension.json
│   │   │   └── xk6-sql
│   │   │       ├── badge.svg
│   │   │       └── extension.json
│   │   └── szkiba
│   │       └── xk6-faker
│   │           ├── badge.svg
│   │           └── extension.json
│   ├── gitlab.com
│   │   └── szkiba
│   │       └── xk6-banner
│   │           ├── badge.svg
│   │           └── extension.json
│   └── go.k6.io
│       └── k6
│           └── extension.json
├── product
│   ├── cloud.json
│   └── oss.json
└── tier
    ├── community.json
    └── official.json
```

The primary purpose of the `--api` flag is to support a custom *k6 extension registry* instance.

## Contribure 

If you want to contribute, start by reading [CONTRIBUTING.md](CONTRIBUTING.md).
