// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

//go:embed help.md
var help string

type options struct {
	out     string
	compact bool
	raw     bool
	yaml    bool
	mute    bool
	loose   bool
	lint    bool
}

// New creates new cobra command for exec command.
func New() (*cobra.Command, error) {
	opts := new(options)

	legacy := false

	root := &cobra.Command{
		Use:               "k6registry [flags] <jq filter> [file]",
		Short:             "k6 extension registry processor",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Args:              cobra.RangeArgs(1, 2),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if legacy {
				return legacyConvert(cmd.Context())
			}

			return run(cmd.Context(), args, opts)
		},
	}

	ctx, err := newContext(context.TODO())
	if err != nil {
		return nil, err
	}

	root.SetContext(ctx)

	flags := root.Flags()

	flags.SortFlags = false

	flags.StringVarP(&opts.out, "out", "o", "", "write output to file instead of stdout")
	flags.BoolVarP(&opts.mute, "mute", "m", false, "no output, only validation")
	flags.BoolVar(&opts.loose, "loose", false, "skip JSON schema validation")
	flags.BoolVar(&opts.lint, "lint", false, "enable built-in linter")
	flags.BoolVarP(&opts.compact, "compact", "c", false, "compact instead of pretty-printed output")
	flags.BoolVarP(&opts.raw, "raw", "r", false, "output raw strings, not JSON texts")
	flags.BoolVarP(&opts.yaml, "yaml", "y", false, "output YAML instead of JSON")
	root.MarkFlagsMutuallyExclusive("raw", "compact", "yaml", "mute")

	flags.BoolP("version", "V", false, "print version")

	flags.BoolVar(&legacy, "legacy", false, "convert legacy registry")

	cobra.CheckErr(flags.MarkHidden("legacy"))

	return root, nil
}

//nolint:forbidigo
func run(ctx context.Context, args []string, opts *options) (result error) {
	input := os.Stdin

	if len(args) > 1 {
		file, err := os.Open(args[1])
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

	if err := jq(registry, args[0], printer(output, opts)); err != nil {
		return err
	}

	return result
}

func printer(output io.Writer, opts *options) func(interface{}) error {
	if opts.raw {
		return func(v interface{}) error {
			_, err := fmt.Fprintln(output, v)

			return err
		}
	}

	if opts.yaml {
		encoder := yaml.NewEncoder(output)

		return encoder.Encode
	}

	if opts.mute {
		return func(_ interface{}) error {
			return nil
		}
	}

	encoder := json.NewEncoder(output)

	if !opts.compact {
		encoder.SetIndent("", "  ")
	}

	return encoder.Encode
}
