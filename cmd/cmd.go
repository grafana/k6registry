// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	out     string
	compact bool
	raw     bool
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

	flags.BoolP("version", "V", false, "print version")
	flags.StringVarP(&opts.out, "out", "o", "", "write output to file instead of stdout")
	flags.BoolVarP(&opts.compact, "compact", "c", false, "compact instead of pretty-printed output")
	flags.BoolVarP(&opts.raw, "raw", "r", false, "output raw strings, not JSON texts")

	flags.BoolVar(&legacy, "legacy", false, "convert legacy registry")

	cobra.CheckErr(flags.MarkHidden("legacy"))

	return root, nil
}

//nolint:forbidigo
func run(ctx context.Context, args []string, opts *options) error {
	var result error

	input := os.Stdin

	if len(args) > 1 {
		file, err := os.Open(args[1])
		if err != nil {
			return err
		}

		defer func() {
			result = file.Close()
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
			result = file.Close()
		}()

		output = file
	}

	registry, err := load(ctx, input)
	if err != nil {
		return err
	}

	var fn func(interface{}) error

	if opts.raw {
		fn = func(v interface{}) error {
			_, err := fmt.Fprintln(output, v)

			return err
		}
	} else {
		encoder := json.NewEncoder(output)

		if !opts.compact {
			encoder.SetIndent("", "  ")
		}

		fn = encoder.Encode
	}

	if err := jq(registry, args[0], fn); err != nil {
		return err
	}

	return result
}
