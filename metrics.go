package k6registry

import "strings"

// CalculateMetrics calculates registry metrics.
func CalculateMetrics(reg Registry) *Metrics {
	const k6Module = "go.k6.io/k6"
	m := new(Metrics)

	for _, ext := range reg {
		if ext.Module == k6Module {
			continue
		}

		m.RegistryExtensionCount++

		m.tier(ext.Tier)

		if ext.Compliance != nil {
			m.grade(ext.Compliance.Grade)
		}

		for _, prod := range ext.Products {
			m.product(prod)
		}

		if len(ext.Products) == 0 {
			m.RegistryProductOSSCount++
		}

		if len(ext.Imports) > 0 {
			m.RegistryTypeJavaScriptCount++
		}

		if len(ext.Outputs) > 0 {
			m.RegistryTypeOutputCount++
		}

		if strings.HasPrefix(ext.Module, "github.com/grafana/") && ext.Tier != TierOfficial {
			m.RegistryTierUnofficialCount++
		}

		for _, cat := range ext.Categories {
			m.category(cat)
		}

		if len(ext.Categories) == 0 {
			m.RegistryCategoryMiscCount++
		}

		for _, issue := range ext.Compliance.Issues {
			m.issue(issue)
		}

		if ext.Cgo {
			m.RegistryCgoCount++
		}
	}

	return m
}

func (m *Metrics) tier(tier Tier) {
	switch tier {
	case TierOfficial:
		m.RegistryTierOfficialCount++
	case TierPartner:
		m.RegistryTierPartnerCount++
	case TierCommunity:
		fallthrough
	default:
		m.RegistryTierCommunityCount++
	}
}

func (m *Metrics) grade(grade Grade) {
	switch grade {
	case GradeA:
		m.RegistryGradeACount++
	case GradeB:
		m.RegistryGradeBCount++
	case GradeC:
		m.RegistryGradeCCount++
	case GradeD:
		m.RegistryGradeDCount++
	case GradeE:
		m.RegistryGradeECount++
	case GradeF:
		m.RegistryGradeFCount++
	default:
	}
}

func (m *Metrics) product(product Product) {
	switch product {
	case ProductCloud:
		m.RegistryProductCloudCount++
	case ProductSynthetic:
		m.RegistryProductSyntheticCount++
	case ProductOSS:
		fallthrough
	default:
		m.RegistryProductOSSCount++
	}
}

func (m *Metrics) category(category Category) {
	switch category {
	case CategoryAuthentication:
		m.RegistryCategoryAuthenticationCount++
	case CategoryBrowser:
		m.RegistryCategoryBrowserCount++
	case CategoryData:
		m.RegistryCategoryDataCount++
	case CategoryKubernetes:
		m.RegistryCategoryKubernetesCount++
	case CategoryMessaging:
		m.RegistryCategoryMessagingCount++
	case CategoryObservability:
		m.RegistryCategoryObservabilityCount++
	case CategoryProtocol:
		m.RegistryCategoryProtocolCount++
	case CategoryReporting:
		m.RegistryCategoryReportingCount++
	case CategoryMisc:
		fallthrough
	default:
		m.RegistryCategoryMiscCount++
	}
}

func (m *Metrics) issue(issue string) {
	switch issue {
	case "module":
		m.RegistryIssueModuleCount++
	case "replace":
		m.RegistryIssueReplaceCount++
	case "readme":
		m.RegistryIssueReadmeCount++
	case "examples":
		m.RegistryIssueExamplesCount++
	case "license":
		m.RegistryIssueLicenseCount++
	case "git":
		m.RegistryIssueGitCount++
	case "versions":
		m.RegistryIssueVersionsCount++
	case "build":
		m.RegistryIssueBuildCount++
	case "smoke":
		m.RegistryIssueSmokeCount++
	case "types":
		m.RegistryIssueTypesCount++
	case "codeowners":
		m.RegistryIssueCodeownersCount++
	default:
	}
}
