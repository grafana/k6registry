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
```

<details>
<summary><b>Example registry</b></summary>

Registry generated from the source above.

```json file=docs/example.json
[
  {
    "categories": [
      "reporting",
      "observability"
    ],
    "description": "Web-based metrics dashboard for k6",
    "module": "github.com/grafana/xk6-dashboard",
    "outputs": [
      "dashboard"
    ],
    "products": [
      "oss"
    ],
    "repo": {
      "description": "A k6 extension that makes k6 metrics available on a web-based dashboard.",
      "homepage": "https://github.com/grafana/xk6-dashboard",
      "license": "AGPL-3.0",
      "name": "xk6-dashboard",
      "owner": "grafana",
      "public": true,
      "stars": 323,
      "topics": [
        "xk6",
        "xk6-official",
        "xk6-output-dashboard"
      ],
      "url": "https://github.com/grafana/xk6-dashboard",
      "versions": [
        "v0.7.5",
        "v0.7.4",
        "v0.7.3",
        "v0.7.3-alpha.1",
        "v0.7.2",
        "v0.7.1",
        "v0.7.0",
        "v0.7.0-apha.3",
        "v0.7.0-alpha.5",
        "v0.7.0-alpha.4",
        "v0.7.0-alpha.3",
        "v0.7.0-alpha.2",
        "v0.7.0-alpha.1",
        "v0.6.1",
        "v0.6.0",
        "v0.5.5",
        "v0.5.4",
        "v0.5.3",
        "v0.5.2",
        "v0.5.1",
        "v0.5.0",
        "v0.4.4",
        "v0.4.3",
        "v0.4.2",
        "v0.4.1",
        "v0.4.0",
        "v0.3.2",
        "v0.3.1",
        "v0.3.0",
        "v0.2.0",
        "v0.1.3",
        "v0.1.2",
        "v0.1.1",
        "v0.1.0"
      ]
    },
    "tier": "official"
  },
  {
    "categories": [
      "data"
    ],
    "description": "Load test SQL Servers",
    "imports": [
      "k6/x/sql"
    ],
    "module": "github.com/grafana/xk6-sql",
    "products": [
      "cloud",
      "oss"
    ],
    "repo": {
      "description": "k6 extension to load test RDBMSs (PostgreSQL, MySQL, MS SQL and SQLite3)",
      "homepage": "https://github.com/grafana/xk6-sql",
      "license": "Apache-2.0",
      "name": "xk6-sql",
      "owner": "grafana",
      "public": true,
      "stars": 104,
      "topics": [
        "k6",
        "sql",
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-sql",
      "versions": [
        "v0.4.0",
        "v0.3.0",
        "v0.2.1",
        "v0.2.0",
        "v0.1.1",
        "v0.1.0",
        "v0.0.1"
      ]
    },
    "tier": "official"
  },
  {
    "categories": [
      "kubernetes"
    ],
    "description": "Inject faults to test",
    "imports": [
      "k6/x/disruptor"
    ],
    "module": "github.com/grafana/xk6-disruptor",
    "products": [
      "oss"
    ],
    "repo": {
      "description": "Extension for injecting faults into k6 tests",
      "homepage": "https://k6.io/docs/javascript-api/xk6-disruptor/",
      "license": "AGPL-3.0",
      "name": "xk6-disruptor",
      "owner": "grafana",
      "public": true,
      "stars": 87,
      "topics": [
        "chaos-engineering",
        "fault-injection",
        "k6",
        "testing",
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-disruptor",
      "versions": [
        "v0.3.11",
        "v0.3.10",
        "v0.3.9",
        "v0.3.8",
        "v0.3.7",
        "v0.3.6",
        "v0.3.5",
        "v0.3.5-rc2",
        "v0.3.5-rc1",
        "v0.3.4",
        "v0.3.3",
        "v0.3.2",
        "v0.3.1",
        "v0.3.0",
        "v0.2.1",
        "v0.2.0",
        "v0.1.3",
        "v0.1.2",
        "v0.1.1",
        "v0.1.0"
      ]
    },
    "tier": "official"
  },
  {
    "categories": [
      "data"
    ],
    "description": "Generate random fake data",
    "imports": [
      "k6/x/faker"
    ],
    "module": "github.com/szkiba/xk6-faker",
    "products": [
      "oss"
    ],
    "repo": {
      "description": "Random fake data generator for k6.",
      "homepage": "http://ivan.szkiba.hu/xk6-faker/",
      "license": "AGPL-3.0",
      "name": "xk6-faker",
      "owner": "szkiba",
      "public": true,
      "stars": 49,
      "topics": [
        "xk6",
        "xk6-javascript-k6-x-faker"
      ],
      "url": "https://github.com/szkiba/xk6-faker",
      "versions": [
        "v0.3.0",
        "v0.3.0-alpha.1",
        "v0.2.2",
        "v0.2.1",
        "v0.2.0",
        "v0.1.0"
      ]
    },
    "tier": "community"
  },
  {
    "categories": [
      "misc"
    ],
    "description": "A modern load testing tool, using Go and JavaScript",
    "module": "go.k6.io/k6",
    "products": [
      "cloud",
      "oss"
    ],
    "repo": {
      "description": "A modern load testing tool, using Go and JavaScript - https://k6.io",
      "homepage": "https://github.com/grafana/k6",
      "license": "AGPL-3.0",
      "name": "k6",
      "owner": "grafana",
      "public": true,
      "stars": 24285,
      "topics": [
        "es6",
        "go",
        "golang",
        "hacktoberfest",
        "javascript",
        "load-generator",
        "load-testing",
        "performance"
      ],
      "url": "https://github.com/grafana/k6",
      "versions": [
        "v0.53.0",
        "v0.52.0",
        "v0.51.0",
        "v0.50.0",
        "v0.49.0",
        "v0.48.0",
        "v0.47.0",
        "v0.46.0",
        "v0.45.1",
        "v0.45.0",
        "v0.44.1",
        "v0.44.0",
        "v0.43.1",
        "v0.43.0",
        "v0.42.0",
        "v0.41.0",
        "v0.40.0",
        "v0.39.0",
        "v0.38.3",
        "v0.38.2",
        "v0.38.1",
        "v0.38.0",
        "v0.37.0",
        "v0.36.0",
        "v0.35.0",
        "v0.34.1",
        "v0.34.0",
        "v0.33.0",
        "v0.32.0",
        "v0.31.1",
        "v0.31.0",
        "v0.30.0",
        "v0.29.0",
        "v0.28.0",
        "v0.27.1",
        "v0.27.0",
        "v0.26.2",
        "v0.26.1",
        "v0.26.0",
        "v0.25.1",
        "v0.25.0",
        "v0.24.0",
        "v0.23.1",
        "v0.23.0",
        "v0.22.1",
        "v0.22.0",
        "v0.21.1",
        "v0.21.0",
        "v0.20.0",
        "v0.19.0",
        "v0.18.2",
        "v0.18.1",
        "v0.18.0",
        "v0.17.2",
        "v0.17.1",
        "v0.17.0",
        "v0.16.0",
        "v0.15.0",
        "v0.14.0",
        "v0.13.0",
        "v0.12.2",
        "v0.12.1",
        "v0.11.0",
        "v0.10.0",
        "v0.9.3",
        "v0.9.2",
        "v0.9.1",
        "v0.9.0",
        "v0.8.5",
        "v0.8.4",
        "v0.8.3",
        "v0.8.2",
        "v0.8.1",
        "v0.8.0",
        "v0.7.0",
        "v0.6.0",
        "v0.5.2",
        "v0.5.1",
        "v0.5.0",
        "v0.4.5",
        "v0.4.4",
        "v0.4.3",
        "v0.4.2",
        "v0.4.1",
        "v0.4.0",
        "v0.3.0",
        "v0.2.1",
        "v0.2.0",
        "v0.0.2",
        "v0.0.1"
      ]
    },
    "tier": "official"
  }
]
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
mute   |    no   | `false` | no output, only validation
loose  |    no   | `false` | skip JSON schema validation
lint   |    no   | `false` | enable built-in linter
compact|    no   | `false` | compact instead of pretty-printed output
ref    |    no   |         | reference output URL for change detection

In GitHub action mode, the change can be indicated by comparing the output to a reference output. The reference output URL can be passed in the `ref` action parameter. The `changed` output variable will be `true` or `false` depending on whether the output has changed or not compared to the reference output.

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

The source of the extension registry is read from the YAML format file specified as command line argument. If it is missing, the source is read from the standard input.

Repository metadata is collected using the API of the extensions' git repository managers. Currently only the GitHub API is supported.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.


```
k6registry [flags] [file]
```

### Flags

```
  -o, --out string   write output to file instead of stdout
  -m, --mute         no output, only validation
      --loose        skip JSON schema validation
      --lint         enable built-in linter
  -c, --compact      compact instead of pretty-printed output
  -V, --version      print version
  -h, --help         help for k6registry
```

<!-- #endregion cli -->

## Contribure 

If you want to contribute, start by reading [CONTRIBUTING.md](CONTRIBUTING.md).
