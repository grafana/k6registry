<h1 name="title">k6registry</h1>

**k6 Extension Registry/Catalog Generator**

k6registry is a CLI tool and a GitHub Action that enables the generation of the k6 extension registry. The generation source is a YAML (or JSON) file that contains the most important properties of extensions. The generator generates the missing properties from the repository metadata. Repository metadata is collected using the repository manager APIs. GitHub and GitLab APIs are currently supported.

The generator also performs static analysis of extensions. The result of the analysis is the level of compliance with best practices (0-100%). A compliance grade (A-F) is calculated from the compliance level. The compliance level and grade are stored in the registry for each extension. Based on the compliance grade, an SVG compliance badge is created for each extension. Example badge:

![xk6-sql](https://registry.k6.io/module/github.com/grafana/xk6-sql/badge.svg)

The k6 Extension Catalog is an alternative representation of the k6 Extension Registry. The output of the generation can be in k6 Extension Catalog format. This format is optimized to resolve extensions as dependencies.

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

- module: github.com/grafana/xk6-faker
  description: Generate random fake data
  imports:
    - k6/x/faker
  tier: official
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
[
  {
    "categories": [
      "reporting",
      "observability"
    ],
    "compliance": {
      "grade": "C",
      "level": 80
    },
    "description": "Web-based metrics dashboard for k6",
    "module": "github.com/grafana/xk6-dashboard",
    "outputs": [
      "dashboard"
    ],
    "products": [
      "oss"
    ],
    "repo": {
      "clone_url": "https://github.com/grafana/xk6-dashboard.git",
      "description": "A k6 extension that makes k6 metrics available on a web-based dashboard.",
      "homepage": "https://github.com/grafana/xk6-dashboard",
      "license": "AGPL-3.0",
      "name": "xk6-dashboard",
      "owner": "grafana",
      "public": true,
      "stars": 330,
      "timestamp": 1719907965,
      "topics": [
        "xk6",
        "xk6-official",
        "xk6-output-dashboard"
      ],
      "url": "https://github.com/grafana/xk6-dashboard"
    },
    "tier": "official",
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
  {
    "categories": [
      "data"
    ],
    "compliance": {
      "grade": "A",
      "level": 100
    },
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
      "clone_url": "https://github.com/grafana/xk6-sql.git",
      "description": "k6 extension to load test RDBMSs (PostgreSQL, MySQL, MS SQL and SQLite3)",
      "homepage": "https://github.com/grafana/xk6-sql",
      "license": "Apache-2.0",
      "name": "xk6-sql",
      "owner": "grafana",
      "public": true,
      "stars": 107,
      "timestamp": 1725628398,
      "topics": [
        "k6",
        "sql",
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-sql"
    },
    "tier": "official",
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
  {
    "categories": [
      "kubernetes"
    ],
    "compliance": {
      "grade": "C",
      "level": 80
    },
    "description": "Inject faults to test",
    "imports": [
      "k6/x/disruptor"
    ],
    "module": "github.com/grafana/xk6-disruptor",
    "products": [
      "oss"
    ],
    "repo": {
      "clone_url": "https://github.com/grafana/xk6-disruptor.git",
      "description": "Extension for injecting faults into k6 tests",
      "homepage": "https://k6.io/docs/javascript-api/xk6-disruptor/",
      "license": "AGPL-3.0",
      "name": "xk6-disruptor",
      "owner": "grafana",
      "public": true,
      "stars": 91,
      "timestamp": 1725571356,
      "topics": [
        "chaos-engineering",
        "fault-injection",
        "k6",
        "testing",
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-disruptor"
    },
    "tier": "official",
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
  {
    "categories": [
      "data"
    ],
    "compliance": {
      "grade": "A",
      "level": 100
    },
    "description": "Generate random fake data",
    "imports": [
      "k6/x/faker"
    ],
    "module": "github.com/grafana/xk6-faker",
    "products": [
      "oss"
    ],
    "repo": {
      "clone_url": "https://github.com/grafana/xk6-faker.git",
      "description": "Random fake data generator for k6.",
      "homepage": "https://faker.x.k6.io",
      "license": "AGPL-3.0",
      "name": "xk6-faker",
      "owner": "grafana",
      "public": true,
      "stars": 50,
      "timestamp": 1725533453,
      "topics": [
        "xk6"
      ],
      "url": "https://github.com/grafana/xk6-faker"
    },
    "tier": "official",
    "versions": [
      "v0.4.0",
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
    "categories": [
      "misc"
    ],
    "compliance": {
      "grade": "A",
      "level": 100
    },
    "description": "Print ASCII art banner from k6 test",
    "imports": [
      "k6/x/banner"
    ],
    "module": "gitlab.com/szkiba/xk6-banner",
    "products": [
      "oss"
    ],
    "repo": {
      "clone_url": "https://gitlab.com/szkiba/xk6-banner.git",
      "description": "Print ASCII art banner from k6 test.",
      "homepage": "https://gitlab.com/szkiba/xk6-banner",
      "license": "MIT",
      "name": "xk6-banner",
      "owner": "szkiba",
      "public": true,
      "timestamp": 1724312566,
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
    "categories": [
      "misc"
    ],
    "description": "A modern load testing tool, using Go and JavaScript",
    "imports": [
      "k6"
    ],
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
      "url": "https://github.com/grafana/k6"
    },
    "tier": "official",
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

`grafana/k6registry` is a GitHub Action that enables k6 Extension Registry and Catalog generation.

**Inputs**

name   | reqired | default | description
-------|---------|---------|-------------
in     |   yes   |         | input file name
out    |    no   |  stdout | output file name
api    |    no   |         | output directory name
test   |    no   |         | api path(s) to test after generation
quiet  |    no   | `false` | no output, only validation
loose  |    no   | `false` | skip JSON schema validation
lint   |    no   | `false` | enable built-in linter
compact|    no   | `false` | compact instead of pretty-printed output
catalog|    no   | `false` | generate catalog instead of registry
ref    |    no   |         | reference output URL for change detection
origin |    no   |         | external registry URL for default values

In GitHub action mode, the change can be indicated by comparing the output to a reference output. The reference output URL can be passed in the `ref` action parameter. The `changed` output variable will be `true` or `false` depending on whether the output has changed or not compared to the reference output.

The `api` parameter can be used to specify a directory into which the outputs are written. The `registry.json` file is placed in the root directory. The `extension.json` file and the `badge.svg` file are placed in a directory with the same name as the go module path of the extension (if the `lint` parameter is `true`).

The `test` parameter can be used to test registry and catalog files generated with the `api` parameter. The test is successful if the file is not empty, contains `k6` and at least one extension, and if all extensions meet the minimum requirements (e.g. it has versions).

**Outputs**

name    | description
--------|------------
changed | `true` if the output has changed compared to `ref`, otherwise `false`

**Example usage**

```yaml
- name: Generate extension registry
  uses: grafana/k6registry@v0.1.18
  with:
    in: "registry.yaml"
    out: "registry.json"
    lint: "true"
```

## CLI

<!-- #region cli -->
## k6registry

k6 Extension Registry/Catalog Generator

### Synopsis

Generate k6 extension registry/catalog from source.

The generation source is a YAML (or JSON) file that contains the most important properties of extensions. The generator generates the missing properties from the repository metadata. Repository metadata is collected using the repository manager APIs. GitHub and GitLab APIs are currently supported.

The generator also performs static analysis of extensions. The result of the analysis is the level of compliance with best practices (0-100%). A compliance grade (A-F) is calculated from the compliance level. The compliance level and grade are stored in the registry for each extension. Based on the compliance grade, an SVG compliance badge is created for each extension. 

The source is read from file specified as command line argument. If it is missing, the source is read from the standard input.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.

The `--api` flag can be used to specify a directory to which the outputs will be written. The `registry.json` file is placed in the root directory. The `extension.json` file and the `badge.svg` file (if the `--lint` flag is used) are placed in a directory with the same name as the extension's go module path.

The `--test` flag can be used to test registry and catalog files generated with the `--api` flag. The test is successful if the file is not empty, contains `k6` and at least one extension, and if all extensions meet the minimum requirements (e.g. it has versions).


```
k6registry [flags] [source-file]
```

### Flags

```
  -o, --out string      write output to file instead of stdout
      --api string      write outputs to directory instead of stdout
      --origin string   external registry URL for default values
      --test strings    test api path(s) (example: /registry.json,/catalog.json)
  -q, --quiet           no output, only validation
      --loose           skip JSON schema validation
      --lint            enable built-in linter
  -c, --compact         compact instead of pretty-printed output
      --catalog         generate catalog instead of registry
  -v, --verbose         verbose logging
  -V, --version         print version
  -h, --help            help for k6registry
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
├── catalog.json
├── registry.json
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
│   └── at-least
│       ├── A.json
│       ├── B.json
│       ├── C.json
│       ├── D.json
│       ├── E.json
│       └── F.json
├── module
│   ├── github.com
│   │   └── grafana
│   │       ├── xk6-dashboard
│   │       │   ├── badge.svg
│   │       │   └── extension.json
│   │       ├── xk6-disruptor
│   │       │   ├── badge.svg
│   │       │   └── extension.json
│   │       ├── xk6-faker
│   │       │   ├── badge.svg
│   │       │   └── extension.json
│   │       └── xk6-sql
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
│   ├── cloud-catalog.json
│   ├── cloud.json
│   ├── oss-catalog.json
│   └── oss.json
└── tier
    ├── community-catalog.json
    ├── community.json
    ├── official-catalog.json
    ├── official.json
    ├── partner-catalog.json
    ├── partner.json
    └── at-least
        ├── community-catalog.json
        ├── community.json
        ├── official-catalog.json
        ├── official.json
        ├── partner-catalog.json
        └── partner.json
```

The primary purpose of the `--api` flag is to support a custom *k6 extension registry* instance.

## Contribure 

If you want to contribute, start by reading [CONTRIBUTING.md](CONTRIBUTING.md).
