Command line k6 extension registry processor.

k6registry is a command line tool that enables k6 extension registry processing and the generation of customized JSON output for different applications. Processing is based on popular `jq` expressions using an embedded `jq` implementation.

The first argument is the jq filter expression. This is the basis for processing.

The extension registry is read from the YAML format file specified in the second argument. If it is missing, the extension registry is read from the standard input.

Repository metadata is collected using the repository manager APIs. Currently only the GitHub API is supported.

The output of the processing will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.
