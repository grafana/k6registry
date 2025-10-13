Generate k6 extension registry from source.

The generation source is a YAML (or JSON) file that contains the most important properties of extensions. The generator generates the missing properties from the repository metadata. Repository metadata is collected using the repository manager APIs. GitHub and GitLab APIs are currently supported.

The generator also performs static analysis of extensions using [xk6 lint](https://github.com/grafana/xk6?tab=readme-ov-file#xk6-lint) command.

The source is read from file specified as command line argument. If it is missing, the source is read from the standard input.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.

