package main

import (
	"errors"
	"os"

	"golang.org/x/term"
)

type formatableError = interface {
	error
	Format(width int, color bool) string
}

func formatError(err error) string {
	width, color := formatOptions(int(os.Stderr.Fd())) //nolint:forbidigo,gosec // CLI tool #nosec G115

	var perr formatableError
	if errors.As(err, &perr) {
		return perr.Format(width, color)
	}

	return err.Error()
}

func formatOptions(fd int) (int, bool) {
	color := false
	width := 0

	if term.IsTerminal(fd) {
		if os.Getenv("NO_COLOR") != "true" { //nolint:forbidigo // CLI tool
			color = true
		}

		if w, _, err := term.GetSize(fd); err == nil {
			width = w
		}
	}

	return width, color
}
