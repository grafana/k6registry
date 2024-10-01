package cmd

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/grafana/k6registry"
)

func loadOrigin(ctx context.Context, from string) (map[string]k6registry.Extension, error) {
	dict := make(map[string]k6registry.Extension, 0)

	if len(from) == 0 {
		return dict, nil
	}

	client := &http.Client{Timeout: httpTimeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, from, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reg k6registry.Registry

	err = json.Unmarshal(data, &reg)
	if err != nil {
		return nil, err
	}

	for _, ext := range reg {
		dict[ext.Module] = ext
	}

	return dict, nil
}

func fromOrigin(ext *k6registry.Extension, origin map[string]k6registry.Extension, loc string) bool {
	oext, found := origin[ext.Module]
	if !found {
		return false
	}

	slog.Debug("Import extension", "module", ext.Module, "origin", loc)

	if ext.Compliance == nil {
		ext.Compliance = oext.Compliance
	}

	if ext.Repo == nil {
		ext.Repo = oext.Repo
	}

	if !ext.Cgo {
		ext.Cgo = oext.Cgo
	}

	if len(ext.Versions) == 0 {
		ext.Versions = oext.Versions
	}

	if len(ext.Constraints) == 0 {
		ext.Constraints = oext.Constraints
	}

	if len(ext.Description) == 0 {
		ext.Description = oext.Description
	}

	if len(ext.Tier) == 0 {
		ext.Tier = oext.Tier
	}

	if len(ext.Categories) == 0 {
		ext.Categories = oext.Categories
	}

	if len(ext.Products) == 0 {
		ext.Products = oext.Products
	}

	if len(ext.Imports) == 0 {
		ext.Imports = oext.Imports
	}

	if len(ext.Outputs) == 0 {
		ext.Outputs = oext.Outputs
	}

	return true
}

const httpTimeout = 10 * time.Second
