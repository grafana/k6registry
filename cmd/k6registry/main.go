// Package main contains the main function for k6registry CLI tool.
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/grafana/k6registry/cmd"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "dev"

func initLogging() *slog.LevelVar {
	levelVar := new(slog.LevelVar)

	logrus.SetLevel(logrus.DebugLevel)

	logger := slog.New(sloglogrus.Option{Level: levelVar}.NewLogrusHandler())

	slog.SetDefault(logger)

	return levelVar
}

//nolint:forbidigo
func main() {
	log.SetFlags(0)
	log.Writer()

	runCmd(newCmd(os.Args[1:], initLogging()))
}

func newCmd(args []string, levelVar *slog.LevelVar) *cobra.Command {
	cmd, err := cmd.New(levelVar)
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
