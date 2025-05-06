package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func isGitHubAction() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func emitOutput(ctx context.Context, out string, ref string) error {
	ghOutput := os.Getenv("GITHUB_OUTPUT")
	if len(ghOutput) == 0 {
		return nil
	}

	file, err := os.Create(filepath.Clean(ghOutput))
	if err != nil {
		return err
	}

	changed := isChanged(ctx, ref, out)

	slog.Debug("Detect change", "changed", changed, "ref", ref)

	_, err = fmt.Fprintf(file, "changed=%t\n", changed)
	if err != nil {
		return err
	}

	return file.Close()
}

func isChanged(ctx context.Context, refURL string, localFile string) bool {
	client := &http.Client{Timeout: httpTimeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, refURL, nil)
	if err != nil {
		return true
	}

	resp, err := client.Do(req)
	if err != nil {
		return true
	}

	defer resp.Body.Close() //nolint:errcheck

	refData, err := io.ReadAll(resp.Body)
	if err != nil {
		return true
	}

	localData, err := os.ReadFile(filepath.Clean(localFile))
	if err != nil {
		return true
	}

	return !bytes.Equal(refData, localData)
}
