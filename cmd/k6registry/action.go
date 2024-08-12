package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//nolint:forbidigo
func emitOutput() error {
	out := getenv("INPUT_OUT", "")
	if len(out) == 0 {
		return nil
	}

	ref := getenv("INPUT_REF", "")
	if len(ref) == 0 {
		return nil
	}

	ghOutput := getenv("GITHUB_OUTPUT", "")
	if len(ghOutput) == 0 {
		return nil
	}

	file, err := os.Create(filepath.Clean(ghOutput))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "changed=%t\n", isChanged(ref, out))
	if err != nil {
		return err
	}

	return file.Close()
}

//nolint:forbidigo
func isChanged(refURL string, localFile string) bool {
	client := &http.Client{Timeout: httpTimeout}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, refURL, nil)
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

const httpTimeout = 10 * time.Second
