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
		ID:       practice.SchemaValidation.ID,
		Filename: practice.SchemaValidation.OASSchema.Name,
		Data:     decodedData,
	}

	schemaValidationMap, err := utils.UnmarshalAs[map[string]any](schemaValidation)
	if err != nil {
		return fmt.Errorf("failed to convert FileSchema struct to map. Error: %w", err)
	}

	d.Set("schema_validation", []map[string]any{schemaValidationMap})

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
					}
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
