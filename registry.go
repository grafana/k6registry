// Package k6registry contains the data model of the k6 extensions registry.
package k6registry

import (
	_ "embed"
)

// Schema contains JSON schema for Grafana k6 Extension Registry JSON.
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

// Level returns level of support (less is better).
func (t Tier) Level() int {
	switch t {
	case TierOfficial:
		return 1
	case TierPartner:
		return 2
	case TierCommunity:
		return 3
	}

	return 0
}

// RegistryToCatalog converts Registry to Catalog.
func RegistryToCatalog(reg Registry) Catalog {
	catalog := make(Catalog, len(reg))

	for _, ext := range reg {
		for _, importPath := range ext.Imports {
			catalog[importPath] = ext
		}

		for _, output := range ext.Outputs {
			catalog[output] = ext
		}
	}

	return catalog
}

// ModuleReference returns true if only the "Module" property has value.
func (ext *Extension) ModuleReference() bool {
	return len(ext.Module) > 0 &&
		len(ext.Description) == 0 &&
		len(ext.Imports) == 0 &&
		len(ext.Outputs) == 0 &&
		len(ext.Versions) == 0 &&
		len(ext.Categories) == 0 &&
		len(ext.Products) == 0 &&
		len(ext.Tier) == 0 &&
		ext.Repo == nil &&
		ext.Compliance == nil
}
