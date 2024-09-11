// Package main contains the main function for k6registry CLI tool.
package main

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/google/shlex"
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

func main() {
	log.SetFlags(0)
	log.Writer()

	runCmd(newCmd(getArgs(), initLogging()))
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

	if isGitHubAction() {
		if err := emitOutput(); err != nil {
			log.Fatal(formatError(err))
		}
	}
}

//nolint:forbidigo
func isGitHubAction() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

//nolint:forbidigo
func getArgs() []string {
	if !isGitHubAction() {
		return os.Args[1:]
	}

	var args []string

	if getenv("INPUT_QUIET", "false") == "true" {
		args = append(args, "--quiet")
	}

	if getenv("INPUT_VERBOSE", "false") == "true" {
		args = append(args, "--verbose")
	}

	if getenv("INPUT_LOOSE", "false") == "true" {
		args = append(args, "--loose")
	}

	if getenv("INPUT_LINT", "false") == "true" {
		args = append(args, "--lint")
	}

	if getenv("INPUT_COMPACT", "false") == "true" {
		args = append(args, "--compact")
	}

	if getenv("INPUT_CATALOG", "false") == "true" {
		args = append(args, "--catalog")
	}

	if api := getenv("INPUT_API", ""); len(api) != 0 {
		args = append(args, "--api", api)
	}

	if out := getenv("INPUT_OUT", ""); len(out) != 0 {
		args = append(args, "--out", out)
	}

	if paths := getenv("INPUT_TEST", ""); len(paths) != 0 {
		parts, err := shlex.Split(paths)
		if err == nil {
			paths = strings.Join(parts, ",")
		}

		args = append(args, "--test", paths)
	}

	args = append(args, getenv("INPUT_IN", ""))

	return args
}

//nolint:forbidigo
func getenv(name string, defval string) string {
	value, found := os.LookupEnv(name)
	if found {
		return value
	}

	return defval
}
