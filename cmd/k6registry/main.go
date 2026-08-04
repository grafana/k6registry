// Package main contains the main function for k6registry CLI tool.
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/grafana/k6registry/cmd"
)

var version = "dev"

func initLogging() *slog.LevelVar {
	levelVar := new(slog.LevelVar)

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: levelVar}) //nolint:forbidigo // CLI tool
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return levelVar
}

func main() {
	log.SetFlags(0)
	log.Writer()

	err := newCmd(os.Args[1:], initLogging()).Execute() //nolint:forbidigo // CLI tool
	if err != nil {
		slog.Error(err.Error()) //nolint:gosec // CLI tool error output
		os.Exit(1)              //nolint:forbidigo // CLI tool
	}
}

func newCmd(args []string, levelVar *slog.LevelVar) *cobra.Command {
	cmd, err := cmd.New(levelVar)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd.Version = version
	cmd.SetArgs(args)

	return cmd
}
