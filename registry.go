// Package k6registry contains the data model of the k6 extensions registry.
package k6registry

import (
	_ "embed"
)

// Schema contains JSON schema for registry JSON.
//
//go:embed docs/registry.schema.json
var Schema []byte

//nolint:gochecknoglobals
var (
	// Categories contains possible values for Category
	Categories = []Category{
		CategoryAuthentication,
		CategoryBrowser,
		CategoryData,
		CategoryKubernetes,
		CategoryMessaging,
		CategoryMisc,
		CategoryObservability,
		CategoryProtocol,
		CategoryReporting,
	}

	// Products contains possible values for Product
	Products = []Product{ProductCloud, ProductOSS}

	// Grades contains possible values for Grade
	Grades = []Grade{GradeA, GradeB, GradeC, GradeD, GradeE, GradeF}

	// Tiers contains possible values for Tier
	Tiers = []Tier{TierOfficial, TierPartner, TierCommunity}
)
