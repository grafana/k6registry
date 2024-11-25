// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/grafana/k6registry"
	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	out     string
	compact bool
	catalog string
	quiet   bool
	verbose bool
	loose   bool
	lint    bool
	api     string
	test    []string
	origin  string
	ref     string
}

// New creates new cobra command for exec command.
func New(levelVar *slog.LevelVar) (*cobra.Command, error) {
	opts := new(options)

	legacy := false

	root := &cobra.Command{
		Use:               "k6registry [flags] [source-file]",
		Short:             "k6 Extension Registry/Catalog Generator",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Args:              cobra.MaximumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.verbose && levelVar != nil {
				levelVar.Set(slog.LevelDebug)
			}

			if legacy {
				return legacyConvert(cmd.Context())
			}

			return run(cmd.Context(), args, opts)
		},
	}

	ctx, err := newContext(context.TODO(), root.Root().Name())
	if err != nil {
		return nil, err
	}

	root.SetContext(ctx)

	flags := root.Flags()

	flags.SortFlags = false

	flags.StringVarP(&opts.out, "out", "o", "", "write output to file instead of stdout")
	flags.StringVar(&opts.api, "api", "", "write outputs to directory instead of stdout")
	flags.StringVar(&opts.origin, "origin", "", "external registry URL for default values")
	flags.StringVar(&opts.ref, "ref", "", "reference output URL for change detection")
	flags.StringSliceVar(&opts.test, "test", []string{}, "test api path(s) (example: /registry.json,/catalog.json)")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "no output, only validation")
	flags.BoolVar(&opts.loose, "loose", false, "skip JSON schema validation")
	flags.BoolVar(&opts.lint, "lint", false, "enable built-in linter")
	flags.BoolVarP(&opts.compact, "compact", "c", false, "compact instead of pretty-printed output")
	flags.StringVar(&opts.catalog, "catalog", "", "generate catalog to the specified file")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "verbose logging")
	root.MarkFlagsMutuallyExclusive("compact", "quiet")
	root.MarkFlagsMutuallyExclusive("api", "out")
	root.MarkFlagsMutuallyExclusive("api", "catalog")

	flags.BoolP("version", "V", false, "print version")

	flags.BoolVar(&legacy, "legacy", false, "convert legacy registry")

	cobra.CheckErr(flags.MarkHidden("legacy"))

	return root, nil
}

//nolint:forbidigo
func run(ctx context.Context, args []string, opts *options) (result error) {
	if len(opts.api) != 0 {
		if err := os.MkdirAll(opts.api, permDir); err != nil {
			return err
		}

		if len(opts.out) == 0 {
			opts.quiet = true
		}
	}

	input := os.Stdin

	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		defer func() {
			err := file.Close()
			if result == nil && err != nil {
				result = err
			}
		}()

		input = file
	}

	output := os.Stdout

	if len(opts.out) > 0 {
		file, err := os.Create(opts.out)
		if err != nil {
			return err
		}

		defer func() {
			err := file.Close()
			if result == nil && err != nil {
				result = err
			}
		}()

		output = file
	}

	registry, err := load(ctx, input, opts.loose, opts.lint, opts.origin)
	if err != nil {
		return err
	}

	return postRun(ctx, registry, output, opts)
}

func postRun(ctx context.Context, registry k6registry.Registry, output io.Writer, opts *options) error {
	if len(opts.api) != 0 {
		if err := writeAPI(registry, opts.api); err != nil {
			return err
		}

		return testAPI(opts.test, opts.api)
	}

	if len(opts.catalog) > 0 {
		if err := writeCatalog(registry, opts.catalog, opts.compact); err != nil {
			return err
		}
	}

	if !opts.quiet {
		if err := writeOutput(registry, output, opts.compact, false); err != nil {
			return err
		}
	}

	if isGitHubAction() && len(opts.out) > 0 && len(opts.ref) > 0 {
		return emitOutput(ctx, opts.out, opts.ref)
	}

	return nil
}

//nolint:forbidigo
func writeCatalog(registry k6registry.Registry, filename string, compact bool) (result error) {
	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if result == nil && err != nil {
			result = err
		}
	}()

	if err := writeOutput(registry, file, compact, true); err != nil {
		result = err
	}

	return
}

func writeOutput(registry k6registry.Registry, output io.Writer, compact, catalog bool) error {
	encoder := json.NewEncoder(output)

	if !compact {
		encoder.SetIndent("", "  ")
	}

	encoder.SetEscapeHTML(false)

	var source interface{} = registry

	if catalog {
		source = k6registry.RegistryToCatalog(registry)
	}

	return encoder.Encode(source)
}

const (
	permFile fs.FileMode = 0o644
	permDir  fs.FileMode = 0o755
)
