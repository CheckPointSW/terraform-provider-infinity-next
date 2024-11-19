package webapipractice

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadWebAPIPracticeToResourceData(practice models.WebAPIPractice, d *schema.ResourceData) error {
	d.SetId(practice.ID)
	d.Set("name", practice.Name)
	d.Set("practice_type", practice.PracticeType)
	d.Set("category", practice.Category)
	d.Set("default", practice.Default)

	ipsSchema := models.SchemaIPS{
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

	advancedSettings := models.SchemaAdvancedSetting{
		ID:             practice.APIAttacks.AdvancedSetting.ID,
		BodySize:       practice.APIAttacks.AdvancedSetting.BodySize,
		URLSize:        practice.APIAttacks.AdvancedSetting.URLSize,
		HeaderSize:     practice.APIAttacks.AdvancedSetting.HeaderSize,
		MaxObjectDepth: practice.APIAttacks.AdvancedSetting.MaxObjectDepth,
	}

	switch practice.APIAttacks.AdvancedSetting.IllegalHttpMethods {
	case "Yes":
		advancedSettings.IllegalHttpMethods = true
	default:
		advancedSettings.IllegalHttpMethods = false
	}

	apiAttacks := models.SchemaAPIAttacks{
		ID:              practice.APIAttacks.ID,
		MinimumSeverity: practice.APIAttacks.MinimumSeverity,
		AdvancedSetting: []models.SchemaAdvancedSetting{advancedSettings},
	}

	apiAttacksMap, err := utils.UnmarshalAs[map[string]any](apiAttacks)
	if err != nil {
		return fmt.Errorf("failed to convert APIAttacks struct to map. Error: %w", err)
	}

	d.Set("api_attacks", []map[string]any{apiAttacksMap})
	var decodedData string
	if strings.Contains(practice.SchemaValidation.OASSchema.Data, "base64,") {
		b64Data := strings.SplitN(practice.SchemaValidation.OASSchema.Data, "base64,", 2)[1]
		bDecodedData, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			return fmt.Errorf("failed decoding base64 string %s: %w", b64Data, err)
		}

		decodedData = string(bDecodedData)
	}
	schemaValidation := models.FileSchema{
		ID:          practice.SchemaValidation.ID,
		Filename:    practice.SchemaValidation.OASSchema.Name,
		Data:        decodedData,
		Size:        practice.SchemaValidation.OASSchema.Size,
		IsFileExist: practice.SchemaValidation.OASSchema.IsFileExist,
	}

	//oasSchema := models.OASSchema{
	//	Data:        decodedData,
	//	Name:        practice.SchemaValidation.OASSchema.Name,
	//	Size:        practice.SchemaValidation.OASSchema.Size,
	//	IsFileExist: practice.SchemaValidation.OASSchema.IsFileExist,
	//}
	//
	//schemaValidation := models.SchemaValidationSchema{
	//	ID:        practice.SchemaValidation.ID,
	//	OASSchema: []models.OASSchema{oasSchema},
	//}

	schemaValidationMap, err := utils.UnmarshalAs[map[string]any](schemaValidation)
	if err != nil {
		return fmt.Errorf("failed to convert SchemaValidation struct to map. Error: %w", err)
	}

	d.Set("schema_validation", []map[string]any{schemaValidationMap})

	fileSecurity := models.WebApplicationFileSecuritySchema{
		ID:                        practice.FileSecurity.ID,
		SeverityLevel:             practice.FileSecurity.SeverityLevel,
		HighConfidence:            practice.FileSecurity.HighConfidence,
		MediumConfidence:          practice.FileSecurity.MediumConfidence,
		LowConfidence:             practice.FileSecurity.LowConfidence,
		AllowFileSizeLimit:        practice.FileSecurity.AllowFileSizeLimit,
		FileSizeLimit:             practice.FileSecurity.FileSizeLimit,
		FileSizeLimitUnit:         practice.FileSecurity.FileSizeLimitUnit,
		FilesWithoutName:          practice.FileSecurity.FilesWithoutName,
		RequiredArchiveExtraction: practice.FileSecurity.RequiredArchiveExtraction,
		ArchiveFileSizeLimit:      practice.FileSecurity.ArchiveFileSizeLimit,
		ArchiveFileSizeLimitUnit:  practice.FileSecurity.ArchiveFileSizeLimitUnit,
		AllowArchiveWithinArchive: practice.FileSecurity.AllowArchiveWithinArchive,
		AllowAnUnopenedArchive:    practice.FileSecurity.AllowAnUnopenedArchive,
		AllowFileType:             practice.FileSecurity.AllowFileType,
		RequiredThreatEmulation:   practice.FileSecurity.RequiredThreatEmulation,
	}

	fileSecurityMap, err := utils.UnmarshalAs[map[string]any](fileSecurity)
	if err != nil {
		return fmt.Errorf("failed to convert FileSecurity struct to map: %w", err)
	}

	d.Set("file_security", []map[string]any{fileSecurityMap})

	return nil
}

func GetWebAPIPractice(ctx context.Context, c *api.Client, id string) (models.WebAPIPractice, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getWebAPIPractice(id: "`+id+`") {
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
				APIAttacks {
					id
					minimumSeverity
					advancedSetting {
						id
						bodySize
						urlSize
						headerSize
						maxObjectDepth
						illegalHttpMethods
					}
				}
				SchemaValidation {
					id
					OasSchema {
						data
						name
						size
						isFileExist
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
	`, "getWebAPIPractice")

	if err != nil {
		return models.WebAPIPractice{}, fmt.Errorf("failed to get WebAPIPractice: %w", err)
	}

	practice, err := utils.UnmarshalAs[models.WebAPIPractice](res)
	if err != nil {
		return models.WebAPIPractice{}, fmt.Errorf("failed to convert response to WebAPIPractice struct. Error: %w", err)
	}

	return practice, nil
}
