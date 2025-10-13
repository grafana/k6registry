// Package k6registry contains the data model of the k6 extensions registry.
package k6registry

import (
	_ "embed"
)

// Schema contains JSON schema for Grafana k6 Extension Registry JSON.
//
//go:embed registry.schema.json
var Schema []byte

//nolint:gochecknoglobals
var (
	// Tiers contains possible values for Tier.
	Tiers = []Tier{TierOfficial, TierCommunity}
)

// Level returns level of support (less is better).
//
//nolint:mnd
func (t Tier) Level() int {
	switch t {
	case TierOfficial:
		return 1
	case TierCommunity:
		return 2
	}

	return 0
}
