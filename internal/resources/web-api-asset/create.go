package webapiasset

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateWebAPIAssetInputFromResourceData(d *schema.ResourceData) (models.CreateWebAPIAssetInput, error) {
	var res models.CreateWebAPIAssetInput

	res.Name = d.Get("name").(string)
	res.UpstreamURL = d.Get("upstream_url").(string)
	res.Profiles = utils.MustResourceDataCollectionToSlice[string](d, "profiles")
	res.Behaviors = utils.MustResourceDataCollectionToSlice[string](d, "trusted_sources")
	res.URLs = utils.MustResourceDataCollectionToSlice[string](d, "urls")
	res.PracticeWrappers = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "practice"), mapToPracticeWrapperInput)
	res.ProxySettings = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "proxy_setting"), mapToProxySettingInput)
	res.SourceIdentifiers = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "source_identifier"), mapToSourceIdentifierInput)

	return res, nil
}

// NewWebAPIAsset sends a request to create the WebAPIAsset and retruns the newly created asset
func NewWebAPIAsset(c *api.Client, input models.CreateWebAPIAssetInput) (models.WebAPIAsset, error) {
	vars := map[string]any{"assetInput": input}
	res, err := c.MakeGraphQLRequest(`
				mutation newWebAPIAsset($assetInput: WebAPIAssetInput!)
					{
						newWebAPIAsset(assetInput: $assetInput) {
							id
							name
							state
							upstreamURL
							practices {
								id
								mainMode
								subPracticeModes {
									mode
									subPractice
								}
								practice {
									id
								}
								type
								status
								triggers {
									id
								}
								behaviors {
									id
								}
							}
							profiles {
								id
							}
							behaviors {
								id
								name
							}
							sourceIdentifiers {
								id
								sourceIdentifier
								values {
									id
									IdentifierValue
								}
							}
							proxySetting {
								id
								key
								value
							}
							URLs {
								id
								URL
							}
							assetType
							sources
							class
							category
							family
							group
							order
							kind
							mainAttributes
							intelligenceTags
							readOnly
						}
					}
				`, "newWebAPIAsset", vars)

	if err != nil {
		return models.WebAPIAsset{}, fmt.Errorf("failed to create new WebAPIAsset: %w", err)
	}

	asset, err := utils.UnmarshalAs[models.WebAPIAsset](res)
	if err != nil {
		return models.WebAPIAsset{}, fmt.Errorf("failed to convert graphQL response to WebAPIAsset struct. Error: %+v", err)
	}

	return asset, nil
}

// mapToPracticeWrapperInput parses the "practice" input to a practice wrapper
func mapToPracticeWrapperInput(practiceWrapperMap map[string]any) models.PracticeWrapperInput {
	var practiceWrapper models.PracticeWrapperInput

	practiceWrapper.PracticeID = practiceWrapperMap["id"].(string)
	practiceWrapper.PracticeWrapperID = practiceWrapperMap["practice_wrapper_id"].(string)
	practiceWrapper.MainMode = practiceWrapperMap["main_mode"].(string)

	if subPracticesModes, ok := practiceWrapperMap["sub_practices_modes"]; ok {
		subPracticesModesMap := subPracticesModes.(map[string]any)
		practiceWrapper.SubPracticeModes = make([]models.PracticeModeInput, 0, len(subPracticesModesMap))
		for subPratice, mode := range subPracticesModesMap {
			practiceWrapper.SubPracticeModes = append(practiceWrapper.SubPracticeModes,
				models.PracticeModeInput{Mode: mode.(string), SubPractice: subPratice})
		}
	}

	if triggersInterface, ok := practiceWrapperMap["triggers"]; ok {
		practiceWrapper.Triggers = utils.MustSchemaCollectionToSlice[string](triggersInterface)
	}

	if behaviorsInterface, ok := practiceWrapperMap["exceptions"]; ok {
		practiceWrapper.Behaviors = utils.MustSchemaCollectionToSlice[string](behaviorsInterface)
	}

	return practiceWrapper
}

func mapToProxySettingInput(proxySettingMap map[string]any) models.ProxySettingInput {
	var ret models.ProxySettingInput
	ret.Key, ret.Value = proxySettingMap["key"].(string), proxySettingMap["value"].(string)
	if id, ok := proxySettingMap["id"]; ok {
		ret.ID = id.(string)
	}

	return ret
}

func mapToSourceIdentifierInput(sourceIdentifierMap map[string]any) models.SourceIdentifierInput {
	var ret models.SourceIdentifierInput
	ret.SourceIdentifier = sourceIdentifierMap["identifier"].(string)
	ret.Values = utils.MustSchemaCollectionToSlice[string](sourceIdentifierMap["values"])
	if valuesIDs, ok := sourceIdentifierMap["values_ids"]; ok {
		ret.ValuesIDs = utils.MustSchemaCollectionToSlice[string](valuesIDs)
	}

	if sourceIdentifierID, ok := sourceIdentifierMap["id"]; ok {
		ret.ID = sourceIdentifierID.(string)
	}

	return ret
}
