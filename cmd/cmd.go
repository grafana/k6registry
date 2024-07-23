// Package cmd contains run cobra command factory function.
package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

// New creates new cobra command for exec command.
func New() *cobra.Command {
	root := &cobra.Command{
		Use:               "k6registry",
		Short:             "k6 extension registry processor",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return root
}
