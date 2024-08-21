package cmd

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/grafana/k6registry"
	"gopkg.in/yaml.v3"
)

type legacyRegistry struct {
	Extensions []*legacyExtension `json:"extensions"`
}

type legacyExtension struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tiers       []string `json:"tiers"`
	Type        []string `json:"type"`
	Categories  []string `json:"categories"`
}

func legacyConvert(ctx context.Context) error {
	client, err := contextGitHubClient(ctx)
	if err != nil {
		return err
	}

	content, _, _, err := client.Repositories.GetContents(
		ctx,
		"grafana",
		"k6-docs",
		"src/data/doc-extensions/extensions.json",
		nil,
	)
	if err != nil {
		return err
	}

	str, err := content.GetContent()
	if err != nil {
		return err
	}

	var legacyReg legacyRegistry

	if err := json.Unmarshal([]byte(str), &legacyReg); err != nil {
		return err
	}

	reg := make([]*k6registry.Extension, 0, len(legacyReg.Extensions))

	for _, legacyExt := range legacyReg.Extensions {
		ext := new(k6registry.Extension)

		ext.Module = strings.TrimPrefix(legacyExt.URL, "https://")
		ext.Description = legacyExt.Description
		ext.Tier = legacyTierToTier(legacyExt.Tiers)
		ext.Categories = legacyCategoriesToCategories(legacyExt.Categories)

		for _, legacyType := range legacyExt.Type {
			typ := strings.ToLower(legacyType)

			if typ == "javascript" {
				name := strings.TrimPrefix(legacyExt.Name, "xk6-")

				ext.Imports = []string{"k6/x/" + name}

				continue
			}

			name := strings.TrimPrefix(legacyExt.Name, "xk6-output-")
			name = strings.TrimPrefix(name, "xk6-")

			ext.Outputs = []string{name}
		}

		legacyPatch(ext)

		repo, err := loadRepository(ctx, ext.Module)
		if err != nil {
			return err
		}

		tmp := *ext
		tmp.Repo = repo

		if ok, _ := lintExtension(tmp); ok {
			reg = append(reg, ext)
		}
	}

	encoder := yaml.NewEncoder(os.Stdout) //nolint:forbidigo

	if err := encoder.Encode(reg); err != nil {
		return err
	}

	return nil
}

func legacyTierToTier(tiers []string) k6registry.Tier {
	for _, tier := range tiers {
		switch strings.ToLower(tier) {
		case "official":
			return k6registry.TierOfficial
		case "partner":
			return k6registry.TierPartner
		default:
			return ""
		}
	}

	return ""
}

func legacyCategoriesToCategories(legacyCats []string) []k6registry.Category {
	if len(legacyCats) == 0 {
		return nil
	}

	cats := make([]k6registry.Category, 0, len(legacyCats))

	for _, cat := range legacyCats {
		cats = append(cats, k6registry.Category(strings.ToLower(cat)))
	}

	return cats
}

func legacyPatch(ext *k6registry.Extension) {
	override, found := extOverrides[ext.Module]
	if !found {
		panic("new module: " + ext.Module)
	}

	if len(override.imports) != 0 {
		if override.imports == "-" {
			ext.Imports = nil
		} else {
			ext.Imports = []string{override.imports}
		}
	}

	if len(override.outputs) != 0 {
		ext.Outputs = []string{override.outputs}
	}

	if len(override.module) != 0 {
		ext.Module = override.module
	}

	if len(override.categories) != 0 {
		ext.Categories = override.categories
	}

	for from, to := range phrases {
		ext.Description = strings.ReplaceAll(ext.Description, from, to)
	}
}

type extOverride struct {
	imports    string
	outputs    string
	module     string
	categories []k6registry.Category
}

var extOverrides = map[string]extOverride{ //nolint:gochecknoglobals
	"github.com/AckeeCZ/xk6-google-iap":               {imports: "k6/x/googleIap"},
	"github.com/BarthV/xk6-es":                        {outputs: "xk6-es"},
	"github.com/GhMartingit/xk6-mongo":                {},
	"github.com/JorTurFer/xk6-input-prometheus":       {imports: "k6/x/prometheusread"},
	"github.com/Juandavi1/xk6-prompt":                 {categories: []k6registry.Category{k6registry.CategoryMisc}},
	"github.com/LeonAdato/xk6-output-statsd":          {outputs: "output-statsd"},
	"github.com/Maksimall89/xk6-output-clickhouse":    {},
	"github.com/NAlexandrov/xk6-tcp":                  {},
	"github.com/SYM01/xk6-proxy":                      {categories: []k6registry.Category{k6registry.CategoryProtocol}},
	"github.com/acuenca-facephi/xk6-read":             {},
	"github.com/akiomik/xk6-nostr":                    {},
	"github.com/anycable/xk6-cable":                   {},
	"github.com/avitalique/xk6-file":                  {},
	"github.com/deejiw/xk6-gcp":                       {},
	"github.com/deejiw/xk6-interpret":                 {},
	"github.com/distribworks/xk6-ethereum":            {categories: []k6registry.Category{k6registry.CategoryProtocol}},
	"github.com/domsolutions/xk6-fasthttp":            {categories: []k6registry.Category{k6registry.CategoryProtocol}},
	"github.com/dynatrace/xk6-output-dynatrace":       {outputs: "output-dynatrace"},
	"github.com/elastic/xk6-output-elasticsearch":     {outputs: "output-elasticsearch"},
	"github.com/fornfrey/xk6-celery":                  {},
	"github.com/frankhefeng/xk6-oauth-pkce":           {},
	"github.com/gjergjsheldija/xk6-mllp":              {},
	"github.com/golioth/xk6-coap":                     {},
	"github.com/gpiechnik2/xk6-httpagg":               {},
	"github.com/gpiechnik2/xk6-smtp":                  {},
	"github.com/grafana/xk6-client-prometheus-remote": {imports: "k6/x/remotewrite"},
	"github.com/grafana/xk6-client-tracing":           {imports: "k6/x/tracing"},
	"github.com/grafana/xk6-dashboard":                {},
	"github.com/grafana/xk6-disruptor":                {categories: []k6registry.Category{k6registry.CategoryKubernetes}},
	"github.com/grafana/xk6-exec":                     {},
	"github.com/grafana/xk6-kubernetes":               {categories: []k6registry.Category{k6registry.CategoryKubernetes}},
	"github.com/grafana/xk6-loki":                     {},
	"github.com/grafana/xk6-notification":             {},
	"github.com/grafana/xk6-output-influxdb":          {outputs: "xk6-influxdb"},
	"github.com/grafana/xk6-output-kafka":             {outputs: "xk6-kafka"},
	"github.com/grafana/xk6-output-timescaledb":       {},
	"github.com/grafana/xk6-sql":                      {},
	"github.com/grafana/xk6-ssh":                      {},
	"github.com/goharbor/xk6-harbor":                  {},
	"github.com/heww/xk6-harbor":                      {module: "github.com/goharbor/xk6-harbor"},
	"github.com/kelseyaubrecht/xk6-webtransport": {
		categories: []k6registry.Category{k6registry.CategoryMessaging, k6registry.CategoryProtocol},
	},
	"github.com/kubeshop/xk6-tracetest":                        {},
	"github.com/leonyork/xk6-output-timestream":                {},
	"github.com/maksimall89/xk6-telegram":                      {},
	"github.com/martymarron/xk6-output-prometheus-pushgateway": {outputs: "output-prometheus-pushgateway"},
	"github.com/mcosta74/xk6-plist":                            {},
	"github.com/mostafa/xk6-kafka":                             {},
	"github.com/mridehalgh/xk6-sqs":                            {},
	"github.com/oleiade/xk6-kv":                                {},
	"github.com/patrick-janeiro/xk6-neo4j":                     {},
	"github.com/phymbert/xk6-sse":                              {},
	"github.com/pmalhaire/xk6-mqtt":                            {},
	"github.com/skibum55/xk6-git":                              {},
	"github.com/szkiba/xk6-ansible-vault":                      {},
	"github.com/szkiba/xk6-cache":                              {outputs: "cache", imports: "-"},
	"github.com/szkiba/xk6-chai":                               {},
	"github.com/szkiba/xk6-csv":                                {},
	"github.com/szkiba/xk6-dotenv":                             {},
	"github.com/szkiba/xk6-faker":                              {},
	"github.com/szkiba/xk6-g0":                                 {},
	"github.com/szkiba/xk6-mock":                               {},
	"github.com/szkiba/xk6-output-plugin":                      {},
	"github.com/szkiba/xk6-prometheus":                         {},
	"github.com/szkiba/xk6-toml":                               {},
	"github.com/szkiba/xk6-top":                                {},
	"github.com/grafana/xk6-ts":                                {},
	"github.com/szkiba/xk6-ts":                                 {module: "github.com/grafana/xk6-ts"},
	"github.com/szkiba/xk6-yaml":                               {},
	"github.com/thmshmm/xk6-opentelemetry":                     {},
	"github.com/thotasrinath/xk6-couchbase":                    {},
	"github.com/tmieulet/xk6-cognito":                          {},
	"github.com/walterwanderley/xk6-stomp":                     {},
	"github.com/nicholasvuono/xk6-playwright":                  {},
	"github.com/ydarias/xk6-nats":                              {},
	"go.k6.io/k6":                                              {},
	"github.com/wosp-io/xk6-playwright":                        {module: "github.com/nicholasvuono/xk6-playwright"},
}

var phrases = map[string]string{ //nolint:gochecknoglobals
	"mqtt": "MQTT",
}
