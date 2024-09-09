// Package main contains CLI documentation generator tool.
package main

import (
	"github.com/grafana/clireadme"
	"github.com/grafana/k6registry/cmd"
)

func main() {
	root, _ := cmd.New(nil)
	clireadme.Main(root, 1)
}
