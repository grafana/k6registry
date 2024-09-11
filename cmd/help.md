Generate k6 extension registry/catalog from source.

The generation source is a YAML (or JSON) file that contains the most important properties of extensions. The generator generates the missing properties from the repository metadata. Repository metadata is collected using the repository manager APIs. GitHub and GitLab APIs are currently supported.

The generator also performs static analysis of extensions. The result of the analysis is the level of compliance with best practices (0-100%). A compliance grade (A-F) is calculated from the compliance level. The compliance level and grade are stored in the registry for each extension. Based on the compliance grade, an SVG compliance badge is created for each extension. 

The source is read from file specified as command line argument. If it is missing, the source is read from the standard input.

The output of the generation will be written to the standard output by default. The output can be saved to a file using the `-o/--out` flag.

The `--api` flag can be used to specify a directory to which the outputs will be written. The `registry.json` file is placed in the root directory. The `extension.json` file and the `badge.svg` file (if the `--lint` flag is used) are placed in a directory with the same name as the extension's go module path.

The `--test` flag can be used to test registry and catalog files generated with the `--api` flag. The test is successful if the file is not empty, contains `k6` and at least one extension, and if all extensions meet the minimum requirements (e.g. it has versions).
