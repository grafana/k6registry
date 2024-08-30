Command line k6 extension registry generator.

The source of the extension registry contains only the most important properties of the extensions. The rest of the properties are collected by k6registry using the API of the extensions' git repository managers.

The source of the extension registry is read from the YAML (or JSON) format file specified as command line argument. If it is missing, the source is read from the standard input.

Repository metadata is collected using the API of the extensions' git repository managers. Currently only the GitHub API is supported.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.

The `--api` flag can be used to specify a directory to which the outputs will be written. The `registry.json` file is placed in the root directory. The `extension.json` file and the `badge.svg` file (if the `--lint` flag is used) are placed in a directory with the same name as the extension's go module path.