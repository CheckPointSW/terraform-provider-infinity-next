package webapppractice

import (
	"context"
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
	d.Set("visibility", practice.Visibility)

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
	default:
		advancedSettings.IllegalHttpMethods = false
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

	//fileSecurity := models.FileSecurity{
	//	ID:                        practice.FileSecurity.ID,
	//	SeverityLevel:             practice.FileSecurity.SeverityLevel,
	//	HighConfidence:            practice.FileSecurity.HighConfidence,
	//	MediumConfidence:          practice.FileSecurity.MediumConfidence,
	//	LowConfidence:             practice.FileSecurity.LowConfidence,
	//	AllowFileSizeLimit:        practice.FileSecurity.AllowFileSizeLimit,
	//	FileSizeLimit:             practice.FileSecurity.FileSizeLimit,
	//	FileSizeLimitUnit:         practice.FileSecurity.FileSizeLimitUnit,
	//	FilesWithoutName:          practice.FileSecurity.FilesWithoutName,
	//	RequiredArchiveExtraction: practice.FileSecurity.RequiredArchiveExtraction,
	//	ArchiveFileSizeLimit:      practice.FileSecurity.ArchiveFileSizeLimit,
	//	ArchiveFileSizeLimitUnit:  practice.FileSecurity.ArchiveFileSizeLimitUnit,
	//	//AllowArchiveWithinArchive: practice.FileSecurity.AllowArchiveWithinArchive,
	//	AllowAnUnopenedArchive:  practice.FileSecurity.AllowAnUnopenedArchive,
	//	AllowFileType:           practice.FileSecurity.AllowFileType,
	//	RequiredThreatEmulation: practice.FileSecurity.RequiredThreatEmulation,
	//}
	//
	//fileSecurityMap, err := utils.UnmarshalAs[map[string]any](fileSecurity)
	//if err != nil {
	//	return fmt.Errorf("failed to convert FileSecurity struct to map. Error: %w", err)
	//}
	//
	//d.Set("file_security", []map[string]any{fileSecurityMap})

	return nil
}

func GetWebApplicationPractice(ctx context.Context, c *api.Client, id string) (models.WebApplicationPractice, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getWebApplicationPractice(id: "`+id+`") {
				id
				name
				visibility
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
				FileSecurity {
					id
					severityLevel
					highConfidence
					mediumConfidence
					lowConfidence
					allowFileSizeLimit
					fileSizeLimit
					fileSizeLimitUnit
					filesWithoutName
					requiredArchiveExtraction
					archiveFileSizeLimit
					archiveFileSizeLimitUnit
					allowArchiveWithinArchive
					allowAnUnopenedArchive
					allowFileType
					requiredThreatEmulation
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
