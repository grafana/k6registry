<h1 name="title">k6registry</h1>

**k6 Extension Registry Generator**

k6registry is a CLI tool and a GitHub Action that enables the generation of the k6 extension registry. The generation source is a YAML (or JSON) file that contains the most important properties of extensions. The generator generates the missing properties from the repository metadata. Repository metadata is collected using the repository manager APIs. GitHub and GitLab APIs are currently supported.

The generator also performs static analysis of extensions using [xk6 lint](https://github.com/grafana/xk6?tab=readme-ov-file#xk6-lint) command. The result of the analysis is a list of issues detected.

Check [k6 Extension Registry Concept](docs/registry.md) for information on design considerations.

**Example registry source**

```yaml file=docs/example.yaml
- module: github.com/grafana/xk6-dashboard
  description: Web-based metrics dashboard for k6
  outputs:
    - dashboard
  subcommands:
    - dashboard
  tier: official

- module: github.com/grafana/xk6-sql
  description: Load test SQL Servers
  imports:
    - k6/x/sql
  tier: official

- module: github.com/grafana/xk6-faker
  description: Generate random fake data
  imports:
    - k6/x/faker
  tier: official

- module: gitlab.com/szkiba/xk6-banner
  description: Print ASCII art banner from k6 test
  imports:
    - k6/x/banner
```

<details>
<summary><b>Example registry</b></summary>

Registry generated from the source above.

```json file=docs/example.json
[
  {
    "description": "Web-based metrics dashboard for k6",
    "module": "github.com/grafana/xk6-dashboard",
    "outputs": [
      "dashboard"
    ],
    "repo": {
      "clone_url": "https://github.com/grafana/xk6-dashboard.git",
      "description": "A k6 extension that makes k6 metrics available on a web-based dashboard.",
      "homepage": "http://dashboard.x.k6.io/",
      "license": "AGPL-3.0",
      "name": "xk6-dashboard",
      "owner": "grafana",
      "public": true,
      "stars": 469,
      "timestamp": 1770397660,
      "topics": [
        "k6",
        "xk6",
        "xk6-output-dashboard"
      ],
      "url": "https://github.com/grafana/xk6-dashboard"
    },
    "subcommands": [
      "dashboard"
    ],
    "tier": "official",
    "versions": [
      "v0.8.0",
      "v0.7.14",
      "v0.7.13",
      "v0.7.12",
      "v0.7.11",
      "v0.7.10",
      "v0.7.9",
      "v0.7.8",
      "v0.7.7",
      "v0.7.6",
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
  {
    "description": "Load test SQL Servers",
    "imports": [
      "k6/x/sql"
    ],
    "module": "github.com/grafana/xk6-sql",
    "repo": {
      "clone_url": "https://github.com/grafana/xk6-sql.git",
      "description": "Use SQL databases from k6 tests.",
      "homepage": "http://sql.x.k6.io/",
      "license": "Apache-2.0",
      "name": "xk6-sql",
      "owner": "grafana",
      "public": true,
      "stars": 179,
      "timestamp": 1770735796,
      "topics": [
        "k6",
        "sql",
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-sql"
    },
    "tier": "official",
    "versions": [
      "v1.0.6",
      "v1.0.5",
      "v1.0.4",
      "v1.0.3",
      "v1.0.2",
      "v1.0.1",
      "v1.0.0",
      "v0.4.1",
      "v0.4.0",
      "v0.3.0",
      "v0.2.1",
      "v0.2.0",
      "v0.1.1",
      "v0.1.0",
      "v0.0.1"
    ]
  },
  {
    "description": "Generate random fake data",
    "imports": [
      "k6/x/faker"
    ],
    "module": "github.com/grafana/xk6-faker",
    "repo": {
      "clone_url": "https://github.com/grafana/xk6-faker.git",
      "description": "Random fake data generator for k6.",
      "homepage": "http://faker.x.k6.io/",
      "license": "AGPL-3.0",
      "name": "xk6-faker",
      "owner": "grafana",
      "public": true,
      "stars": 113,
      "timestamp": 1770391184,
      "topics": [
        "k6",
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-faker"
    },
    "tier": "official",
    "versions": [
      "v0.4.4",
      "v0.4.3",
      "v0.4.2",
      "v0.4.1",
      "v0.4.0",
      "v0.3.2",
      "v0.3.1",
      "v0.3.0",
      "v0.3.0-alpha.1",
      "v0.2.2",
      "v0.2.1",
      "v0.2.0",
      "v0.1.0"
    ]
  },
  {
    "description": "Print ASCII art banner from k6 test",
    "imports": [
      "k6/x/banner"
    ],
    "module": "gitlab.com/szkiba/xk6-banner",
    "repo": {
      "clone_url": "https://gitlab.com/szkiba/xk6-banner.git",
      "description": "Print ASCII art banner from k6 test.",
      "homepage": "https://gitlab.com/szkiba/xk6-banner",
      "license": "MIT",
      "name": "xk6-banner",
      "owner": "szkiba",
      "public": true,
      "timestamp": 1725896396,
      "topics": [
        "xk6"
      ],
      "url": "https://gitlab.com/szkiba/xk6-banner"
    },
    "tier": "community",
    "versions": [
      "v0.1.0"
    ]
  },
  {
    "description": "A modern load testing tool, using Go and JavaScript",
    "imports": [
      "k6"
    ],
    "module": "go.k6.io/k6",
    "repo": {
      "description": "A modern load testing tool, using Go and JavaScript",
      "homepage": "https://grafana.com/oss/k6/",
      "license": "AGPL-3.0",
      "name": "k6",
      "owner": "grafana",
      "public": true,
      "topics": [
        "es6",
        "go",
        "golang",
        "hacktoberfest",
        "javascript",
        "k6",
        "load-generator",
        "load-testing",
        "performance"
      ],
      "url": "https://github.com/grafana/k6"
    },
    "tier": "official",
    "versions": [
      "v1.6.0",
      "v1.5.0",
      "v1.4.2",
      "v1.4.1",
      "v1.4.1-release-1",
      "v1.4.1-release",
      "v1.4.0",
      "v1.3.0",
      "v1.2.3",
      "v1.2.2",
      "v1.2.1",
      "v1.2.0",
      "v1.1.0",
      "v1.0.0",
      "v1.0.0-rc2",
      "v1.0.0-rc1",
      "v0.59.0",
      "v0.58.0",
      "v0.57.0",
      "v0.56.0",
      "v0.55.2",
      "v0.55.1",
      "v0.55.0",
      "v0.54.0",
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
      "v0.6.0"
    ]
  }
]
```

</details>


## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6registry/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6registry/cmd/k6registry@latest
```

## GitHub Workflows

When used in GitHub Workflows, the change can be indicated by comparing the output to a reference output. The reference output URL can be passed in the `--ref` flag. The `changed` output variable will be `true` or `false` depending on whether the output has changed or not compared to the reference output.

**Outputs**

name    | description
--------|------------
changed | `true` if the output has changed compared to `ref`, otherwise `false`

## CLI

<!-- #region cli -->
## k6registry

k6 Extension Registry/Catalog Generator

### Synopsis

Generate k6 extension registry from source.

The generation source is a YAML (or JSON) file that contains the most important properties of extensions. The generator generates the missing properties from the repository metadata. Repository metadata is collected using the repository manager APIs. GitHub and GitLab APIs are currently supported.

The generator also performs static analysis of extensions using [xk6 lint](https://github.com/grafana/xk6?tab=readme-ov-file#xk6-lint) command.

The source is read from file specified as command line argument. If it is missing, the source is read from the standard input.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.



```
k6registry [flags] [source-file]
```

### Flags

```
  -o, --out string            write output to file instead of stdout
  -q, --quiet                 no output, only validation
      --lint                  enable built-in linter
      --ignore-lint-errors    don't fail on lint errors
      --lint-checks strings   lint checks to apply. Check xk6 documentation for available options.
  -c, --compact               compact instead of pretty-printed output
  -v, --verbose               verbose logging
  -V, --version               print version
  -h, --help                  help for k6registry
```

### Commands

* [k6registry schema](#k6registry-schema)	 - Output the JSON schema to stdout

---
## k6registry schema

Output the JSON schema to stdout

### Synopsis

Output the JSON schema for the k6 extension registry to stdout

```
k6registry schema [flags]
```

### Flags

```
  -h, --help   help for schema
```

### SEE ALSO

* [k6registry](#k6registry)	 - k6 Extension Registry/Catalog Generator

<!-- #endregion cli -->

## Contribure 

If you want to contribute, start by reading [CONTRIBUTING.md](CONTRIBUTING.md).
