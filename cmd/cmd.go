// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	out     string
	compact bool
	quiet   bool
	loose   bool
	lint    bool
	api     string
}

// New creates new cobra command for exec command.
func New() (*cobra.Command, error) {
	opts := new(options)

	legacy := false

	root := &cobra.Command{
		Use:               "k6registry [flags] [source-file]",
		Short:             "k6 extension registry generator",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Args:              cobra.MaximumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
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
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "no output, only validation")
	flags.BoolVar(&opts.loose, "loose", false, "skip JSON schema validation")
	flags.BoolVar(&opts.lint, "lint", false, "enable built-in linter")
	flags.BoolVarP(&opts.compact, "compact", "c", false, "compact instead of pretty-printed output")
	root.MarkFlagsMutuallyExclusive("compact", "quiet")

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

	registry, err := load(ctx, input, opts.loose, opts.lint)
	if err != nil {
		return err
	}

	if len(opts.api) != 0 {
		return writeAPI(registry, opts.api)
	}

	if opts.quiet {
		return nil
	}

	encoder := json.NewEncoder(output)

	if !opts.compact {
		encoder.SetIndent("", "  ")
	}

	err = encoder.Encode(registry)
	if err != nil {
		return err
	}

	return nil
}

const (
	permFile fs.FileMode = 0o644
	permDir  fs.FileMode = 0o755
)
