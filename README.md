# k6registry

**Data model for the k6 extension registry**

This repository contains the [JSON schema](registry.schema.json) of the k6 extension registry and the golang data model generated from it.

## Concept

The k6 extension registry is a JSON file that contains the most important properties of extensions.

### Registered Properties

Only those properties of the extensions are registered, which either cannot be detected automatically, or delegation to the extension is not allowed.

Properties that are available using the repository manager API are intentionally not registered.

The string like properties that are included in the generated Grafana documentation are intentionally not accessed via the API of the repository manager. It is not allowed to inject arbitrary text into the Grafana documentation site without approval. Therefore, these properties are registered (eg `description`)

### Extension Identification

The primary identifier of an extension is the extension's [go module path](https://go.dev/ref/mod#module-path).

The extension has no `name` property, the module path or part of it can be used as the extension name. For example, using the first two elements of the module path after the host name, the name `grafana/xk6-dashboard` can be formed from the module path `github.com/grafana/xk6-dashboard`. This is typically the repository owner name (`grafana`) and the repository name in the repository manager (`xk6-dashboard`).

The extension has no URL property, a URL can be created from the module path that refers to the extension within the repository manager.

### JavaScript Modules

The JavaScript module names implemented by the extension can be specified in the `imports` property. An extension can register multiple JavaScript module names, so this is an array property.

### Output Names

The output names implemented by the extension can be specified in the `outputs` property. An extension can register multiple output names, so this is an array property.

### Cloud Enabled

The `true` value of the `cloud` flag indicates that the extension is also available in the Grafana k6 cloud.

### Officially Supported

The `true` value of the `official` flag indicates that the extension is officially supported by Grafana.

## Example

```json file=example.json
{
  "extensions": [
    {
      "module": "github.com/grafana/xk6-dashboard",
      "description": "Web-based metrics dashboard for k6",
      "outputs": ["dashboard"],
      "official": true
    },
    {
      "module": "github.com/grafana/xk6-sql",
      "description": "Load test SQL Servers",
      "imports": ["k6/x/sql"],
      "official": true
    },
    {
      "module": "github.com/grafana/xk6-distruptor",
      "description": "Inject faults to test",
      "imports": ["k6/x/distruptor"],
      "official": true
    },
    {
      "module": "github.com/szkiba/xk6-faker",
      "description": "Generate random fake data",
      "imports": ["k6/x/faker"]
    }
  ]
}
```

## Development

The source of the schema is contained in the [registry.schema.yaml](registry.schema.yaml) file. This file must be edited if required. The `go generate` command generates the [registry_gen.go](registry_gen.go) and [registry.schema.json](registry.schema.json) files and updates the example code block in the README.md from the [example.json](example.json) file.