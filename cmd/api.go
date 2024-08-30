package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/grafana/k6registry"
	"github.com/narqo/go-badge"
)

func writeAPI(registry k6registry.Registry, target string) error {
	if err := writeAPIGroupGlobal(registry, target); err != nil {
		return err
	}

	if err := writeAPIGroupModule(registry, target); err != nil {
		return err
	}

	if err := writeAPIGroupSubset(registry, target); err != nil {
		return err
	}

	return nil
}

//nolint:forbidigo
func writeAPIGroupGlobal(registry k6registry.Registry, target string) error {
	filename := filepath.Join(target, "registry.schema.json")

	if err := os.WriteFile(filename, k6registry.Schema, 0o600); err != nil {
		return err
	}

	filename = filepath.Join(target, "registry.json")

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o600)
}

//nolint:forbidigo
func writeAPIGroupModule(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "module")

	if err := os.MkdirAll(base, 0o750); err != nil {
		return err
	}

	for _, ext := range registry {
		dir := filepath.Join(base, ext.Module)

		if err := os.MkdirAll(dir, 0o750); err != nil {
			return err
		}

		filename := filepath.Join(dir, "badge.svg")

		if ext.Compliance != nil {
			b, err := badge.RenderBytes("k6 registry", string(ext.Compliance.Grade), badgecolor(ext.Compliance.Grade))
			if err != nil {
				return err
			}

			err = os.WriteFile(filename, b, 0o600)
			if err != nil {
				return err
			}
		}

		filename = filepath.Join(dir, "extension.json")

		data, err := json.MarshalIndent(ext, "", "  ")
		if err != nil {
			return err
		}

		if err = os.WriteFile(filename, data, 0o600); err != nil {
			return err
		}
	}

	return nil
}

func writeAPIGroupSubset(registry k6registry.Registry, target string) error {
	if err := writeAPISubsetProduct(registry, target); err != nil {
		return err
	}

	if err := writeAPISubsetTier(registry, target); err != nil {
		return err
	}

	if err := writeAPISubsetGrade(registry, target); err != nil {
		return err
	}

	return writeAPISubsetCategory(registry, target)
}

//nolint:forbidigo
func writeAPISubsetProduct(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "product")

	if err := os.MkdirAll(base, 0o750); err != nil {
		return err
	}

	products := make(map[k6registry.Product]k6registry.Registry, len(k6registry.Products))

	for _, ext := range registry {
		for _, prod := range ext.Products {
			reg, found := products[prod]
			if !found {
				reg = make(k6registry.Registry, 0)
			}

			reg = append(reg, ext)
			products[prod] = reg
		}
	}

	for prod, reg := range products {
		data, err := json.MarshalIndent(reg, "", " ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(base, string(prod)+".json"), data, 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:forbidigo
func writeAPISubsetTier(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "tier")

	if err := os.MkdirAll(base, 0o750); err != nil {
		return err
	}

	tiers := make(map[k6registry.Tier]k6registry.Registry, len(k6registry.Tiers))

	for _, ext := range registry {
		if len(ext.Tier) == 0 {
			continue
		}

		reg, found := tiers[ext.Tier]
		if !found {
			reg = make(k6registry.Registry, 0)
		}

		reg = append(reg, ext)
		tiers[ext.Tier] = reg
	}

	for tier, reg := range tiers {
		data, err := json.MarshalIndent(reg, "", " ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(base, string(tier)+".json"), data, 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:forbidigo
func writeAPISubsetGrade(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "grade")

	if err := os.MkdirAll(base, 0o750); err != nil {
		return err
	}

	grades := make(map[k6registry.Grade]k6registry.Registry, len(k6registry.Grades))

	for _, grade := range k6registry.Grades {
		grades[grade] = make(k6registry.Registry, 0)
	}

	for _, ext := range registry {
		if ext.Compliance == nil || len(ext.Compliance.Grade) == 0 {
			continue
		}

		reg, found := grades[ext.Compliance.Grade]
		if !found {
			reg = make(k6registry.Registry, 0)
		}

		reg = append(reg, ext)
		grades[ext.Compliance.Grade] = reg
	}

	for grade, reg := range grades {
		data, err := json.MarshalIndent(reg, "", " ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(base, string(grade)+".json"), data, 0o600)
		if err != nil {
			return err
		}
	}

	return writeAPISubsetGradePassing(registry, target)
}

//nolint:forbidigo
func writeAPISubsetGradePassing(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "grade", "passing")

	if err := os.MkdirAll(base, 0o750); err != nil {
		return err
	}

	grades := make(map[k6registry.Grade]k6registry.Registry, len(k6registry.Grades))

	for _, grade := range k6registry.Grades {
		grades[grade] = make(k6registry.Registry, 0)
	}

	for _, ext := range registry {
		if ext.Compliance == nil || len(ext.Compliance.Grade) == 0 {
			continue
		}

		for _, grade := range k6registry.Grades {
			if ext.Compliance.Grade > grade {
				continue
			}

			reg, found := grades[grade]
			if !found {
				reg = make(k6registry.Registry, 0)
			}

			reg = append(reg, ext)
			grades[grade] = reg
		}
	}

	for grade, reg := range grades {
		data, err := json.MarshalIndent(reg, "", " ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(base, string(grade)+".json"), data, 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:forbidigo
func writeAPISubsetCategory(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "category")

	if err := os.MkdirAll(base, 0o750); err != nil {
		return err
	}

	categories := make(map[k6registry.Category]k6registry.Registry, len(k6registry.Categories))

	for _, category := range k6registry.Categories {
		categories[category] = make(k6registry.Registry, 0)
	}

	for _, ext := range registry {
		for _, cat := range ext.Categories {
			reg, found := categories[cat]
			if !found {
				reg = make(k6registry.Registry, 0)
			}

			reg = append(reg, ext)
			categories[cat] = reg
		}
	}

	for cat, reg := range categories {
		data, err := json.MarshalIndent(reg, "", " ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(base, string(cat)+".json"), data, 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

func badgecolor(grade k6registry.Grade) badge.Color {
	switch grade {
	case k6registry.GradeA:
		return "brightgreen"
	case k6registry.GradeB:
		return "green"
	case k6registry.GradeC:
		return "yellowgreen"
	case k6registry.GradeD:
		return "orange"
	case k6registry.GradeE:
		return "yellow"
	case k6registry.GradeF:
		return "red"
	default:
		return "blue"
	}
}
