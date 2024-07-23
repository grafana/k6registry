// Package main contains CLI documentation generator tool.
package main

import (
	"github.com/grafana/clireadme"
	"github.com/grafana/k6registry/cmd"
)

func main() {
	root := cmd.New()
	clireadme.Main(root, 1)
}
