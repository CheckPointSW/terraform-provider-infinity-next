package webapppractice

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func mapToIPSInput(ipsMap map[string]any) models.WebApplicationPractcieIPSInput {
	return models.WebApplicationPractcieIPSInput{
		PerformanceImpact:   ipsMap["performance_impact"].(string),
		SeverityLevel:       ipsMap["severity_level"].(string),
		ProtectionsFromYear: "Y" + ipsMap["protections_from_year"].(string),
		HighConfidence:      ipsMap["high_confidence"].(string),
		MediumConfidence:    ipsMap["medium_confidence"].(string),
		LowConfidence:       ipsMap["low_confidence"].(string),
	}
}

func mapToWebAttacksInput(webAttacksMap map[string]any) models.WebApplicationPracticeWebAttacksInput {
	var res models.WebApplicationPracticeWebAttacksInput
	res.MinimumSeverity = webAttacksMap["minimum_severity"].(string)
	advancedSetting := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](webAttacksMap["advanced_setting"]), mapToAdvancedSettingInput)
	if len(advancedSetting) > 0 {
		res.AdvancedSetting = advancedSetting[0]
	}

	return res
}

func mapToAdvancedSettingInput(advancedSettingMap map[string]any) models.WebApplicationPracticeAdvancedSettingInput {
	illegalHttpMethods := "No"
	if advancedSettingMap["illegal_http_methods"].(bool) {
		illegalHttpMethods = "Yes"
	}

	return models.WebApplicationPracticeAdvancedSettingInput{
		CSRFProtection:     advancedSettingMap["csrf_protection"].(string),
		OpenRedirect:       advancedSettingMap["open_redirect"].(string),
		ErrorDisclosure:    advancedSettingMap["error_disclosure"].(string),
		BodySize:           advancedSettingMap["body_size"].(int),
		URLSize:            advancedSettingMap["url_size"].(int),
		HeaderSize:         advancedSettingMap["header_size"].(int),
		MaxObjectDepth:     advancedSettingMap["max_object_depth"].(int),
		IllegalHttpMethods: illegalHttpMethods,
	}
}

func mapToWebBotInput(webBotMap map[string]any) models.WebApplicationPracticeWebBotInput {
	var webBotInput models.WebApplicationPracticeWebBotInput
	webBotInput.InjectURIs = utils.MustSchemaCollectionToSlice[string](webBotMap["inject_uris"])
	webBotInput.ValidURIs = utils.MustSchemaCollectionToSlice[string](webBotMap["valid_uris"])

	return webBotInput
}

func CreateWebApplicationPracticeInputFromResourceData(d *schema.ResourceData) (models.CreateWebApplicationPracticeInput, error) {
	var res models.CreateWebApplicationPracticeInput

	res.Name = d.Get("name").(string)
	res.Visibility = "Shared"
	ipsSlice := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "ips"), mapToIPSInput)
	if len(ipsSlice) > 0 {
		res.IPS = ipsSlice[0]
	}

	webAttacksSlice := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "web_attacks"), mapToWebAttacksInput)
	if len(webAttacksSlice) > 0 {
		res.WebAttacks = webAttacksSlice[0]
	}

	webBotSlice := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "web_bot"), mapToWebBotInput)
	if len(webBotSlice) > 0 {
		res.WebBot = webBotSlice[0]
	}

	return res, nil
}

func NewWebApplicationPractice(c *api.Client, input models.CreateWebApplicationPracticeInput) (models.WebApplicationPractice, error) {
	vars := map[string]any{"ownerId": nil, "mainMode": nil, "subPracticeModes": []any{}, "practiceInput": input}
	res, err := c.MakeGraphQLRequest(`
				mutation newWebApplicationPractice($ownerId: ID, $mainMode: PracticeMode, $subPracticeModes: [PracticeModeInput], $practiceInput: WebApplicationPracticeInput)
					{
						newWebApplicationPractice(ownerId: $ownerId, subPracticeModes: $subPracticeModes, mainMode: $mainMode, practiceInput: $practiceInput) {
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
				`, "newWebApplicationPractice", vars)

	if err != nil {
		return models.WebApplicationPractice{}, fmt.Errorf("failed to create new WebApplicationPractice: %w", err)
	}

	practice, err := utils.UnmarshalAs[models.WebApplicationPractice](res)
	if err != nil {
		return models.WebApplicationPractice{}, fmt.Errorf("failed to convert response to WebApplicationPractice structL %w", err)
	}

	return practice, err
}
