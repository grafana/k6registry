// Package main contains the main function for k6registry CLI tool.
package main

import (
	"log"
	"os"

	"github.com/grafana/k6registry/cmd"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	log.SetFlags(0)
	log.Writer()

	runCmd(newCmd(os.Args[1:])) //nolint:forbidigo
}

func newCmd(args []string) *cobra.Command {
	cmd, err := cmd.New()
	if err != nil {
		log.Fatal(formatError(err))
	}

	cmd.Version = version
	cmd.SetArgs(args)

	return cmd
}

func runCmd(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		log.Fatal(formatError(err))
	}
}
