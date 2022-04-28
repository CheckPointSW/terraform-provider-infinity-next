package webapppractice

import (
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadWebApplicationPracticeToResourceData(practice models.WebApplicationPractice, d *schema.ResourceData) error {
	d.SetId(practice.ID)
	d.Set("name", practice.Name)
	d.Set("practice_type", practice.PracticeType)
	d.Set("category", practice.Category)
	d.Set("default", practice.Default)

	ipsSchema := models.WebApplicationPracticeIPSSchema{
		ID:                  practice.IPS.ID,
		PerformanceImpact:   practice.IPS.PerformanceImpact,
		SeverityLevel:       practice.IPS.SeverityLevel,
		ProtectionsFromYear: practice.IPS.ProtectionsFromYear,
		HighConfidence:      practice.IPS.HighConfidence,
		MediumConfidence:    practice.IPS.MediumConfidence,
		LowConfidence:       practice.IPS.LowConfidence,
	}
	if ipsSchema.ProtectionsFromYear != "" {
		ipsSchema.ProtectionsFromYear = strings.TrimPrefix(ipsSchema.ProtectionsFromYear, "Y")
	}

	ipsMap, err := utils.UnmarshalAs[map[string]any](ipsSchema)
	if err != nil {
		return fmt.Errorf("failed to convert IPS struct to map. Error: %w", err)
	}

	d.Set("ips", []map[string]any{ipsMap})

	advancedSettings := models.WebApplicationPracticeAdvancedSettingSchema{
		ID:              practice.WebAttacks.AdvancedSetting.ID,
		CSRFProtection:  practice.WebAttacks.AdvancedSetting.CSRFProtection,
		OpenRedirect:    practice.WebAttacks.AdvancedSetting.OpenRedirect,
		ErrorDisclosure: practice.WebAttacks.AdvancedSetting.ErrorDisclosure,
		BodySize:        practice.WebAttacks.AdvancedSetting.BodySize,
		URLSize:         practice.WebAttacks.AdvancedSetting.URLSize,
		HeaderSize:      practice.WebAttacks.AdvancedSetting.HeaderSize,
		MaxObjectDepth:  practice.WebAttacks.AdvancedSetting.MaxObjectDepth,
	}

	switch practice.WebAttacks.AdvancedSetting.IllegalHttpMethods {
	case "Yes":
		advancedSettings.IllegalHttpMethods = true
	case "No", "":
		advancedSettings.IllegalHttpMethods = false
	default:
		return fmt.Errorf("invalid illegalHttpMethods %s", practice.WebAttacks.AdvancedSetting.IllegalHttpMethods)
	}

	webAttacks := models.WebApplicationPracticeWebAttacksSchema{
		ID:              practice.WebAttacks.ID,
		MinimumSeverity: practice.WebAttacks.MinimumSeverity,
		AdvancedSetting: []models.WebApplicationPracticeAdvancedSettingSchema{advancedSettings},
	}

	webAttacksMap, err := utils.UnmarshalAs[map[string]any](webAttacks)
	if err != nil {
		return fmt.Errorf("failed to convert WebAttacks struct to map. Error: %w", err)
	}

	d.Set("web_attacks", []map[string]any{webAttacksMap})

	injectURIs := make([]string, len(practice.WebBot.InjectURIs))
	injectURIsIDs := make([]string, len(practice.WebBot.InjectURIs))
	for i, uri := range practice.WebBot.InjectURIs {
		injectURIs[i] = uri.URI
		injectURIsIDs[i] = uri.CreateSchemaID()
	}

	validURIs := make([]string, len(practice.WebBot.ValidURIs))
	validURIsIDs := make([]string, len(practice.WebBot.ValidURIs))
	for i, uri := range practice.WebBot.ValidURIs {
		validURIs[i] = uri.URI
		validURIsIDs[i] = uri.CreateSchemaID()
	}

	webBot := models.WebApplicationPracticeWebBotSchema{
		ID:            practice.WebBot.ID,
		InjectURIs:    injectURIs,
		InjectURIsIDs: injectURIsIDs,
		ValidURIs:     validURIs,
		ValidURIsIDs:  validURIsIDs,
	}

	webBotMap, err := utils.UnmarshalAs[map[string]any](webBot)
	if err != nil {
		return fmt.Errorf("failed to convert WebBot struct to map: %w", err)
	}

	d.Set("web_bot", []map[string]any{webBotMap})

	return nil
}

func GetWebApplicationPractice(c *api.Client, id string) (models.WebApplicationPractice, error) {
	res, err := c.MakeGraphQLRequest(`
		{
			getWebApplicationPractice(id: "`+id+`") {
				id
				name
				practiceType
				category
				default
				IPS {
					id
					performanceImpact
					severityLevel
					protectionsFromYear
					highConfidence
					mediumConfidence
					lowConfidence
				}
				WebAttacks {
					id
					minimumSeverity
					advancedSetting {
						id
						CSRFProtection
						openRedirect
						errorDisclosure
						bodySize
						urlSize
						headerSize
						maxObjectDepth
						illegalHttpMethods
					}
				}
				WebBot {
					id
					injectURIs {
						id
						URI
					}
					validURIs {
						id
						URI
					}
				}
			}
		}
	`, "getWebApplicationPractice")

	if err != nil {
		return models.WebApplicationPractice{}, fmt.Errorf("failed to get WebApplicationPractice: %w", err)
	}

	practice, err := utils.UnmarshalAs[models.WebApplicationPractice](res)
	if err != nil {
		return models.WebApplicationPractice{}, fmt.Errorf("failed to convert response to WebApplicationPractice struct. Error: %w", err)
	}

	return practice, nil
}
