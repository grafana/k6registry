package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
func writeData(filename string, data []byte) error {
	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, permDir); err != nil {
		return err
	}

	return os.WriteFile(filename, data, permFile)
}

func writeJSON(filename string, source interface{}) error {
	var buff bytes.Buffer

	encoder := json.NewEncoder(&buff)

	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(source)
	if err != nil {
		return err
	}

	return writeData(filename, buff.Bytes())
}

func writeAPIGroupGlobal(registry k6registry.Registry, target string) error {
	if err := writeJSON(filepath.Join(target, "registry.json"), registry); err != nil {
		return err
	}

	if err := writeJSON(filepath.Join(target, "catalog.json"), k6registry.RegistryToCatalog(registry)); err != nil {
		return err
	}

	return writeJSON(filepath.Join(target, "metrics.json"), k6registry.CalculateMetrics(registry))
}

const gradesvg = `<svg xmlns="http://www.w3.org/2000/svg" width="17" height="20"><clipPath id="B"><rect width="17" height="20" rx="3" fill="#fff"/></clipPath><g clip-path="url(#B)"><path fill="%s" d="M0 0h17v20H0z"/><path fill="url(#A)" d="M0 0h17v20H0z"/></g><g text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" text-rendering="geometricPrecision" font-size="110"><text x="85" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="70">%s</text><text x="85" y="140" transform="scale(.1)" fill="#fff" textLength="70">%s</text></g></svg>` //nolint:lll

func writeAPIGroupModule(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "module")

	for _, ext := range registry {
		dir := filepath.Join(base, ext.Module)

		if ext.Compliance != nil {
			sgrade := string(ext.Compliance.Grade)

			b, err := badge.RenderBytes("k6 registry", sgrade, badgecolor(ext.Compliance.Grade))
			if err != nil {
				return err
			}

			err = writeData(filepath.Join(dir, "badge.svg"), b)
			if err != nil {
				return err
			}

			grade := fmt.Sprintf(gradesvg, badgecolor(ext.Compliance.Grade), sgrade, sgrade)

			err = writeData(filepath.Join(dir, "grade.svg"), []byte(grade))
			if err != nil {
				return err
			}
		}

		if err := writeJSON(filepath.Join(dir, "extension.json"), ext); err != nil {
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

func writeAPISubsetProduct(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "product")

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
		prefix := string(prod)
		if err := writeJSON(filepath.Join(base, prefix+".json"), reg); err != nil {
			return err
		}

		if err := writeJSON(filepath.Join(base, prefix+"-catalog.json"), k6registry.RegistryToCatalog(reg)); err != nil {
			return err
		}
	}

	return nil
}

func writeAPISubsetTier(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "tier")

	tiers := make(map[k6registry.Tier]k6registry.Registry, len(k6registry.Tiers))

	var k6ext *k6registry.Extension

	for idx, ext := range registry {
		if len(ext.Tier) == 0 {
			continue
		}

		if ext.Module == k6Module {
			k6ext = &registry[idx]
		}

		reg, found := tiers[ext.Tier]
		if !found {
			reg = make(k6registry.Registry, 0)
		}

		reg = append(reg, ext)
		tiers[ext.Tier] = reg
	}

	for _, tier := range k6registry.Tiers {
		if tier == k6registry.TierOfficial {
			continue
		}

		reg, found := tiers[tier]
		if !found {
			reg = make(k6registry.Registry, 0)
		}

		reg = append(reg, *k6ext)
		tiers[tier] = reg
	}

	for tier, reg := range tiers {
		prefix := string(tier)

		if err := writeJSON(filepath.Join(base, prefix+".json"), reg); err != nil {
			return err
		}

		if err := writeJSON(filepath.Join(base, prefix+"-catalog.json"), k6registry.RegistryToCatalog(reg)); err != nil {
			return err
		}
	}

	return writeAPISubsetTierAtLeast(registry, target)
}

func writeAPISubsetTierAtLeast(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "tier", "at-least")

	tiers := make(map[k6registry.Tier]k6registry.Registry, len(k6registry.Tiers))

	for _, tier := range k6registry.Tiers {
		tiers[tier] = make(k6registry.Registry, 0)
	}

	for _, ext := range registry {
		for _, tier := range k6registry.Tiers {
			if ext.Tier.Level() > tier.Level() {
				continue
			}

			reg, found := tiers[tier]
			if !found {
				reg = make(k6registry.Registry, 0)
			}

			reg = append(reg, ext)
			tiers[tier] = reg
		}
	}

	for tier, reg := range tiers {
		prefix := string(tier)

		if err := writeJSON(filepath.Join(base, prefix+".json"), reg); err != nil {
			return err
		}

		if err := writeJSON(filepath.Join(base, prefix+"-catalog.json"), k6registry.RegistryToCatalog(reg)); err != nil {
			return err
		}
	}

	return nil
}

func writeAPISubsetGrade(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "grade")

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
		if err := writeJSON(filepath.Join(base, string(grade)+".json"), reg); err != nil {
			return err
		}
	}

	return writeAPISubsetGradeAtLeast(registry, target)
}

func writeAPISubsetGradeAtLeast(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "grade", "at-least")

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
		if err := writeJSON(filepath.Join(base, string(grade)+".json"), reg); err != nil {
			return err
		}
	}

	return nil
}

func writeAPISubsetCategory(registry k6registry.Registry, target string) error {
	base := filepath.Join(target, "category")

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
		if err := writeJSON(filepath.Join(base, string(cat)+".json"), reg); err != nil {
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
		return "yellow"
	case k6registry.GradeE:
		return "orange"
	case k6registry.GradeF:
		return "red"
	default:
		return "blue"
	}
}

func isCatalog(filename string) bool {
	basename := strings.TrimSuffix(filename, filepath.Ext(filename))

	return strings.HasSuffix(basename, "catalog")
}

func testAPI(paths []string, dir string) error {
	for _, name := range paths {
		slog.Debug("Testing API", "path", name)

		name = filepath.FromSlash(strings.TrimPrefix(name, "/"))
		name = filepath.Join(dir, name)

		if err := testFile(name); err != nil {
			return err
		}
	}

	return nil
}

//nolint:forbidigo
func testFile(filename string) error {
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}

	var catalog k6registry.Catalog

	if isCatalog(filename) {
		if err := json.Unmarshal(data, &catalog); err != nil {
			return err
		}
	} else {
		var registry k6registry.Registry

		if err := json.Unmarshal(data, &registry); err != nil {
			return err
		}

		catalog = k6registry.RegistryToCatalog(registry)
	}

	_, found := catalog[k6ImportPath]
	if !found {
		return fmt.Errorf("%w: %s: missing %s module", errTestFailed, filename, k6Module)
	}

	if len(catalog) == 1 {
		return fmt.Errorf("%w: %s: missing extensions", errTestFailed, filename)
	}

	for _, ext := range catalog {
		if ok, msgs := lintExtension(ext); !ok && len(msgs) != 0 {
			return fmt.Errorf("%w: %s: %s", errTestFailed, filename, msgs[0])
		}
	}

	return nil
}

var errTestFailed = errors.New("test failed")
