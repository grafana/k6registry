package k6registry

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// CalculateMetrics calculates registry metrics for all extensions.
func CalculateMetrics(reg Registry) *Metrics {
	return CalculateMetricsCond(reg, func(_ Extension) bool { return true })
}

func calculateMetricsOne(ext Extension, m *Metrics) {
	m.ExtensionCount++

	m.tier(ext.Tier)

	if ext.Compliance != nil {
		m.grade(ext.Compliance.Grade)
	}

	for _, prod := range ext.Products {
		m.product(prod)
	}

	if len(ext.Products) == 0 {
		m.ProductOSSCount++
	}

	if len(ext.Imports) > 0 {
		m.TypeJavaScriptCount++
	}

	if len(ext.Outputs) > 0 {
		m.TypeOutputCount++
	}

	if strings.HasPrefix(ext.Module, "github.com/grafana/") && ext.Tier != TierOfficial {
		m.TierUnofficialCount++
	}

	for _, cat := range ext.Categories {
		m.category(cat)
	}

	if len(ext.Categories) == 0 {
		m.CategoryMiscCount++
	}

	for _, issue := range ext.Compliance.Issues {
		m.issue(issue)
	}

	if ext.Cgo {
		m.CgoCount++
	}
}

// CalculateMetricsCond calculates registry metrics for subset of extensions.
func CalculateMetricsCond(reg Registry, predicate func(Extension) bool) *Metrics {
	const k6Module = "go.k6.io/k6"

	m := new(Metrics)

	for _, ext := range reg {
		if ext.Module == k6Module || !predicate(ext) {
			continue
		}

		calculateMetricsOne(ext, m)
	}

	return m
}

func (m *Metrics) tier(tier Tier) {
	switch tier {
	case TierOfficial:
		m.TierOfficialCount++
	case TierPartner:
		m.TierPartnerCount++
	case TierCommunity:
		fallthrough
	default:
		m.TierCommunityCount++
	}
}

func (m *Metrics) grade(grade Grade) {
	switch grade {
	case GradeA:
		m.GradeACount++
	case GradeB:
		m.GradeBCount++
	case GradeC:
		m.GradeCCount++
	case GradeD:
		m.GradeDCount++
	case GradeE:
		m.GradeECount++
	case GradeF:
		m.GradeFCount++
	case GradeG:
	default:
	}
}

func (m *Metrics) product(product Product) {
	switch product {
	case ProductCloud:
		m.ProductCloudCount++
	case ProductSynthetic:
		m.ProductSyntheticCount++
	case ProductOSS:
		fallthrough
	default:
		m.ProductOSSCount++
	}
}

func (m *Metrics) category(category Category) {
	switch category {
	case CategoryAuthentication:
		m.CategoryAuthenticationCount++
	case CategoryBrowser:
		m.CategoryBrowserCount++
	case CategoryData:
		m.CategoryDataCount++
	case CategoryKubernetes:
		m.CategoryKubernetesCount++
	case CategoryMessaging:
		m.CategoryMessagingCount++
	case CategoryObservability:
		m.CategoryObservabilityCount++
	case CategoryProtocol:
		m.CategoryProtocolCount++
	case CategoryReporting:
		m.CategoryReportingCount++
	case CategoryMisc:
		fallthrough
	default:
		m.CategoryMiscCount++
	}
}

func (m *Metrics) issue(issue string) {
	switch issue {
	case "module":
		m.IssueModuleCount++
	case "replace":
		m.IssueReplaceCount++
	case "readme":
		m.IssueReadmeCount++
	case "examples":
		m.IssueExamplesCount++
	case "license":
		m.IssueLicenseCount++
	case "git":
		m.IssueGitCount++
	case "versions":
		m.IssueVersionsCount++
	case "build":
		m.IssueBuildCount++
	case "smoke":
		m.IssueSmokeCount++
	case "types":
		m.IssueTypesCount++
	case "codeowners":
		m.IssueCodeownersCount++
	default:
	}
}

// WritePrometheus marshals metrics in Prometheus text format.
func (m *Metrics) WritePrometheus(out io.Writer, namePrefix string, helpPrefix string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var dict map[string]int

	if err := json.Unmarshal(data, &dict); err != nil {
		return err
	}

	now := time.Now().UnixMilli()

	for name, value := range dict {
		fullname := namePrefix + name

		if help, hasHelp := metricsHelps[name]; hasHelp {
			fmt.Fprintf(out, "# HELP %s %s\n", fullname, helpPrefix+help) //nolint:errcheck
		}

		fmt.Fprintf(out, "# TYPE %s counter\n", fullname)      //nolint:errcheck
		fmt.Fprintf(out, "%s %d %d\n\n", fullname, value, now) //nolint:errcheck
	}

	return nil
}
