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

	"github.com/grafana/k6registry"
	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	out     string
	compact bool
	quiet   bool
	verbose bool
	loadOptions
}

// New creates new cobra command for exec command.
func New(levelVar *slog.LevelVar) (*cobra.Command, error) {
	opts := new(options)

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

			return run(cmd.Context(), args, opts)
		},
	}

	root.AddCommand(schemaCmd())

	ctx, err := newContext(context.TODO(), root.Root().Name())
	if err != nil {
		return nil, err
	}

	root.SetContext(ctx)

	flags := root.Flags()

	flags.SortFlags = false

	flags.StringVarP(&opts.out, "out", "o", "", "write output to file instead of stdout")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "no output, only validation")
	flags.BoolVar(&opts.lint, "lint", false, "enable built-in linter")
	flags.BoolVar(&opts.ignoreLintErrors, "ignore-lint-errors", false, "don't fail on lint errors")
	flags.StringSliceVar(
		&opts.lintChecks,
		"lint-checks",
		nil,
		"lint checks to apply. Check xk6 documentation for available options.",
	)
	flags.BoolVarP(&opts.compact, "compact", "c", false, "compact instead of pretty-printed output")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "verbose logging")
	root.MarkFlagsMutuallyExclusive("compact", "quiet")

	flags.BoolP("version", "V", false, "print version")

	return root, nil
}

func schemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema",
		Short: "Output the JSON schema to stdout",
		Long:  "Output the JSON schema for the k6 extension registry to stdout",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := cmd.OutOrStdout().Write(k6registry.Schema)

			return err
		},
	}
}

//nolint:funlen
func run(ctx context.Context, args []string, opts *options) (result error) {
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

	registry, err := load(ctx, input, opts.loadOptions)
	if err != nil {
		return err
	}

	if err := postRun(registry, output, opts); err != nil {
		return err
	}

	return nil
}

func postRun(registry k6registry.Registry, output io.Writer, opts *options) error {
	if opts.quiet {
		return nil
	}

	return writeOutput(registry, output, opts.compact)
}

func writeOutput(registry k6registry.Registry, output io.Writer, compact bool) error {
	encoder := json.NewEncoder(output)

	if !compact {
		encoder.SetIndent("", "  ")
	}

	encoder.SetEscapeHTML(false)

	var source interface{} = registry

	return encoder.Encode(source)
}

const (
	permFile fs.FileMode = 0o644
	permDir  fs.FileMode = 0o755
)
