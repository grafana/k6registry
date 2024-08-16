Command line k6 extension registry generator.

The source of the extension registry contains only the most important properties of the extensions. The rest of the properties are collected by k6registry using the API of the extensions' git repository managers.

The source of the extension registry is read from the YAML (or JSON) format file specified as command line argument. If it is missing, the source is read from the standard input.

Repository metadata is collected using the API of the extensions' git repository managers. Currently only the GitHub API is supported.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.
