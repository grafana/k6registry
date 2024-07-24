<h1 name="title">k6registry</h1>

**Data model and tooling for the k6 extension registry**

This repository contains the [JSON schema](docs/registry.schema.json) of the k6 extension registry and the [`k6registry`](#k6registry) command line tool for registry processing.

Check [k6 Extension Registry Concept](docs/registry.md) for information on design considerations.

**Example registry**

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

## Usage

<!-- #region cli -->
## k6registry

k6 extension registry processor

### Synopsis

Command line k6 extension registry processor.

k6registry is a command line tool that enables k6 extension registry processing and the generation of customized JSON output for different applications. Processing is based on popular `jq` expressions using an embedded `jq` implementation.

The first argument is the jq filter expression. This is the basis for processing.

The extension registry is read from the YAML format file specified in the second argument. If it is missing, the extension registry is read from the standard input.

Repository metadata is collected using the repository manager APIs. Currently only the GitHub API is supported.

The output of the processing will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.


```
k6registry [flags] <jq filter> [file]
```

### Flags

```
  -c, --compact      compact instead of pretty-printed output
  -h, --help         help for k6registry
  -o, --out string   write output to file instead of stdout
  -r, --raw          output raw strings, not JSON texts
  -V, --version      print version
```

<!-- #endregion cli -->

## Contribure 

If you want to contribute, start by reading [CONTRIBUTING.md](CONTRIBUTING.md).
