package webapipractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateWebAPIPracticeInputFromResourceData(d *schema.ResourceData) (models.CreateWebAPIPracticeInput, error) {
	var res models.CreateWebAPIPracticeInput

	res.Name = d.Get("name").(string)
	res.Visibility = "Shared"
	ipsSlice := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "ips"), mapToIPSInput)
	if len(ipsSlice) > 0 {
		res.IPS = ipsSlice[0]
	}

	apiAttacksSlice := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "api_attacks"), mapToAPIAttacksInput)
	if len(apiAttacksSlice) > 0 {
		res.APIAttacks = apiAttacksSlice[0]
	}

	schemaValidationSlice := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "schema_validation"), mapToSchemaValidationInput)
	if len(schemaValidationSlice) > 0 {
		res.SchemaValidation = schemaValidationSlice[0]
	}

	return res, nil
}

func NewWebAPIPractice(ctx context.Context, c *api.Client, input models.CreateWebAPIPracticeInput) (models.WebAPIPractice, error) {
	vars := map[string]any{"ownerId": nil, "mainMode": nil, "subPracticeModes": nil, "practiceInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newWebAPIPractice($ownerId: ID, $mainMode: PracticeMode, $subPracticeModes: [PracticeModeInput], $practiceInput: WebAPIPracticeInput)
					{
						newWebAPIPractice(ownerId: $ownerId, subPracticeModes: $subPracticeModes, mainMode: $mainMode, practiceInput: $practiceInput) {
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
				`, "newWebAPIPractice", vars)

	if err != nil {
		return models.WebAPIPractice{}, fmt.Errorf("failed to create new WebAPIPractice: %w", err)
	}

	practice, err := utils.UnmarshalAs[models.WebAPIPractice](res)
	if err != nil {
		return models.WebAPIPractice{}, fmt.Errorf("failed to convert response to WebAPIPractice struct. Error: %w", err)
	}

	return practice, err
}

func mapToIPSInput(ipsMap map[string]any) models.IPSInput {
	//var ret models.IPSInput
	//ret.PerformanceImpact = ipsMap["performance_impact"].(string)
	//ret.SeverityLevel = ipsMap["severity_level"].(string)
	//ret.ProtectionsFromYear = "Y" + ipsMap["protections_from_year"].(string)
	//ret.HighConfidence = ipsMap["high_confidence"].(string)
	//ret.MediumConfidence = ipsMap["medium_confidence"].(string)
	//ret.LowConfidence = ipsMap["low_confidence"].(string)
	//if id, ok := ipsMap["id"]; ok {
	//	ret.ID = id.(string)
	//}
	//
	//return ret
	return models.IPSInput{
		PerformanceImpact:   ipsMap["performance_impact"].(string),
		SeverityLevel:       ipsMap["severity_level"].(string),
		ProtectionsFromYear: "Y" + ipsMap["protections_from_year"].(string),
		HighConfidence:      ipsMap["high_confidence"].(string),
		MediumConfidence:    ipsMap["medium_confidence"].(string),
		LowConfidence:       ipsMap["low_confidence"].(string),
	}
}

func mapToAdvancedSettingInput(advancedSettingMap map[string]any) models.AdvancedSettingInput {
	//var ret models.AdvancedSettingInput
	//
	illegalHttpMethods := "No"
	if advancedSettingMap["illegal_http_methods"].(bool) {
		illegalHttpMethods = "Yes"
	}
	//
	//ret.BodySize = advancedSettingMap["body_size"].(int)
	//ret.URLSize = advancedSettingMap["url_size"].(int)
	//ret.HeaderSize = advancedSettingMap["header_size"].(int)
	//ret.MaxObjectDepth = advancedSettingMap["max_object_depth"].(int)
	//ret.IllegalHttpMethods = illegalHttpMethods
	//if id, ok := advancedSettingMap["id"]; ok {
	//	ret.ID = id.(string)
	//}
	//
	//return ret

	return models.AdvancedSettingInput{
		BodySize:           advancedSettingMap["body_size"].(int),
		URLSize:            advancedSettingMap["url_size"].(int),
		HeaderSize:         advancedSettingMap["header_size"].(int),
		MaxObjectDepth:     advancedSettingMap["max_object_depth"].(int),
		IllegalHttpMethods: illegalHttpMethods,
	}
}

func mapToAPIAttacksInput(apiAttacksMap map[string]any) models.APIAttacksInput {
	var res models.APIAttacksInput
	res.MinimumSeverity = apiAttacksMap["minimum_severity"].(string)
	advancedSetting := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](apiAttacksMap["advanced_setting"]), mapToAdvancedSettingInput)
	if len(advancedSetting) > 0 {
		res.AdvancedSetting = advancedSetting[0]
	}

	return res
}

func mapToSchemaValidationInput(schemaValidationMap map[string]any) models.SchemaValidationInput {
	var ret models.SchemaValidationInput

	if id, ok := schemaValidationMap["id"]; ok {
		ret.ID = id.(string)
	}

	ret.OASSchema = schemaValidationMap["OasSchema"].(string)

	return ret
}

//func createSchemaValidationInput(schemaValidtionFromResourceData any) models.SchemaValidationInput {
//	schemaValidation, err := utils.UnmarshalAs[models.FileSchema](schemaValidtionFromResourceData)
//	if err != nil {
//		fmt.Printf("Failed to convert input schema validation to FileSchema struct. Error: %+v", err)
//		return models.SchemaValidationInput{}
//	}
//
//	schemaValidation = models.NewFileSchemaEncode(schemaValidation.Filename, schemaValidation.Data)
//	return models.SchemaValidationInput{
//		OASSchema: schemaValidation.Data,
//	}
//}
