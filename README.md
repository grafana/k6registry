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
  official: true

- module: github.com/grafana/xk6-sql
  description: Load test SQL Servers
  imports:
    - k6/x/sql
  cloud: true
  official: true

- module: github.com/grafana/xk6-disruptor
  description: Inject faults to test
  imports:
    - k6/x/disruptor
  official: true

- module: github.com/szkiba/xk6-faker
  description: Generate random fake data
  imports:
    - k6/x/faker
```

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
  uses: grafana/k6registry@v0.1.4
  with:
    in: "registry.yaml"
    out: "registry.json"
    filter: "."
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
