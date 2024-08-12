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

	runCmd(newCmd(getArgs()))
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

	if getenv("INPUT_MUTE", "false") == "true" {
		args = append(args, "--mute")
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

	if getenv("INPUT_RAW", "false") == "true" {
		args = append(args, "--raw")
	}

	if getenv("INPUT_YAML", "false") == "true" {
		args = append(args, "--yaml")
	}

	if out := getenv("INPUT_OUT", ""); len(out) != 0 {
		args = append(args, "--out", out)
	}

	args = append(args, getenv("INPUT_FILTER", "."))

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
