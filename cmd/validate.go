package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/grafana/k6registry"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

func yaml2json(input []byte) ([]byte, error) {
	var data interface{}

	if err := yaml.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func validateWithSchema(input io.Reader) ([]byte, error) {
	yamlRaw, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	jsonRaw, err := yaml2json(yamlRaw)
	if err != nil {
		return nil, err
	}

	documentLoader := gojsonschema.NewBytesLoader(jsonRaw)
	schemaLoader := gojsonschema.NewBytesLoader(k6registry.Schema)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	if result.Valid() {
		return yamlRaw, nil
	}

	var buff strings.Builder

	for _, desc := range result.Errors() {
		buff.WriteString(fmt.Sprintf(" - %s\n", desc.String()))
	}

	return nil, fmt.Errorf("%w: schema validation failed\n%s", errInvalidRegistry, buff.String())
}

func validateWithLinter(registry k6registry.Registry) error {
	var buff strings.Builder

	for _, ext := range registry {
		if ok, msgs := lintExtension(ext); !ok {
			for _, msg := range msgs {
				buff.WriteString(fmt.Sprintf(" - %s\n", msg))
			}
		}
	}

	if buff.Len() == 0 {
		return nil
	}

	return fmt.Errorf("%w: linter validation failed\n%s", errInvalidRegistry, buff.String())
}

func hasTopic(ext k6registry.Extension) bool {
	for _, topic := range ext.Repo.Topics {
		if topic == "xk6" {
			return true
		}
	}

	for _, product := range ext.Products {
		if product == k6registry.ProductOSS { // for oss, topic is required
			return false
		}
	}

	return true // for non oss, topic isn't required
}

func lintExtension(ext k6registry.Extension) (bool, []string) {
	if ext.Repo == nil {
		return false, []string{"unsupported module: " + ext.Module}
	}

	var msgs []string

	if len(ext.Versions) == 0 {
		msgs = append(msgs, "no released versions: "+ext.Module)
	}

	if ext.Repo.Public {
		if !hasTopic(ext) && ext.Module != k6Module {
			msgs = append(msgs, "missing xk6 topic: "+ext.Module)
		}

		if len(ext.Repo.License) == 0 {
			msgs = append(msgs, "missing license: "+ext.Module)
		} else if _, ok := validLicenses[ext.Repo.License]; !ok {
			msgs = append(msgs, "unsupported license: "+ext.Repo.License+" "+ext.Module)
		}

		if ext.Repo.Archived {
			msgs = append(msgs, "repository is archived: "+ext.Module)
		}
	}

	if len(msgs) == 0 {
		return true, nil
	}

	return false, msgs
}

var errInvalidRegistry = errors.New("invalid registry")

// source: https://spdx.org/licenses/
// both FSF Free and OSI Approved licenses
var validLicenses = map[string]struct{}{ //nolint:gochecknoglobals
	"AFL-1.1":           {},
	"AFL-1.2":           {},
	"AFL-2.0":           {},
	"AFL-2.1":           {},
	"AFL-3.0":           {},
	"AGPL-3.0":          {},
	"AGPL-3.0-only":     {},
	"AGPL-3.0-or-later": {},
	"Apache-1.1":        {},
	"Apache-2.0":        {},
	"APSL-2.0":          {},
	"Artistic-2.0":      {},
	"BSD-2-Clause":      {},
	"BSD-3-Clause":      {},
	"BSL-1.0":           {},
	"CDDL-1.0":          {},
	"CPAL-1.0":          {},
	"CPL-1.0":           {},
	"ECL-2.0":           {},
	"EFL-2.0":           {},
	"EPL-1.0":           {},
	"EPL-2.0":           {},
	"EUDatagrid":        {},
	"EUPL-1.1":          {},
	"EUPL-1.2":          {},
	"GPL-2.0-only":      {},
	"GPL-2.0":           {},
	"GPL-2.0-or-later":  {},
	"GPL-3.0-only":      {},
	"GPL-3.0":           {},
	"GPL-3.0-or-later":  {},
	"HPND":              {},
	"Intel":             {},
	"IPA":               {},
	"IPL-1.0":           {},
	"ISC":               {},
	"LGPL-2.1":          {},
	"LGPL-2.1-only":     {},
	"LGPL-2.1-or-later": {},
	"LGPL-3.0":          {},
	"LGPL-3.0-only":     {},
	"LGPL-3.0-or-later": {},
	"LPL-1.02":          {},
	"MIT":               {},
	"MPL-1.1":           {},
	"MPL-2.0":           {},
	"MS-PL":             {},
	"MS-RL":             {},
	"NCSA":              {},
	"Nokia":             {},
	"OFL-1.1":           {},
	"OSL-1.0":           {},
	"OSL-2.0":           {},
	"OSL-2.1":           {},
	"OSL-3.0":           {},
	"PHP-3.01":          {},
	"Python-2.0":        {},
	"QPL-1.0":           {},
	"RPSL-1.0":          {},
	"SISSL":             {},
	"Sleepycat":         {},
	"SPL-1.0":           {},
	"Unlicense":         {},
	"UPL-1.0":           {},
	"W3C":               {},
	"Zlib":              {},
	"ZPL-2.0":           {},
	"ZPL-2.1":           {},
}
