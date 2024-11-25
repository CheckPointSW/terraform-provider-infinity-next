package webappasset

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateWebApplicationAssetInputFromResourceData(d *schema.ResourceData) (models.CreateWebApplicationAssetInput, error) {
	var res models.CreateWebApplicationAssetInput

	res.Name = d.Get("name").(string)
	res.UpstreamURL = d.Get("upstream_url").(string)
	res.Profiles = utils.MustResourceDataCollectionToSlice[string](d, "profiles")
	res.Behaviors = utils.MustResourceDataCollectionToSlice[string](d, "behaviors")
	res.URLs = utils.MustResourceDataCollectionToSlice[string](d, "urls")
	res.PracticeWrappers = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "practice"), mapToPracticeWrapperInput)
	res.ProxySettings = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "proxy_setting"), mapToProxySettingInput)
	res.SourceIdentifiers = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "source_identifier"), mapToSourceIdentifierInput)
	res.Tags = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "tags"), mapToTagsInputs)
	res.IsSharesURLs = d.Get("is_shares_urls").(bool)

	return res, nil
}

func NewWebApplicationAsset(ctx context.Context, c *api.Client, input models.CreateWebApplicationAssetInput) (models.WebApplicationAsset, error) {
	vars := map[string]any{"assetInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newWebApplicationAsset($assetInput: WebApplicationAssetInput!)
					{
						newWebApplicationAsset(assetInput: $assetInput) {
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
							}
							profiles {
								id
							}
							tags {
								id
								key
								value
							}
							behaviors {
								id
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
							isSharesURLs
						}
					}
				`, "newWebApplicationAsset", vars)

	if err != nil {
		return models.WebApplicationAsset{}, fmt.Errorf("failed to create new WebApplicationAsset: %w", err)
	}

	asset, err := utils.UnmarshalAs[models.WebApplicationAsset](res)
	if err != nil {
		return models.WebApplicationAsset{}, fmt.Errorf("failed to convert graphQL response to WebApplicationAsset struct. Error: %w", err)
	}

	return asset, err
}

// mapToPracticeWrapperInput parses the "practice" input to a practice wrapper
func mapToPracticeWrapperInput(practiceWrapperMap map[string]any) models.PracticeWrapperInput {
	var practiceWrapper models.PracticeWrapperInput

	practiceWrapper.PracticeID = practiceWrapperMap["id"].(string)
	practiceWrapper.PracticeWrapperID = practiceWrapperMap["practice_wrapper_id"].(string)
	practiceWrapper.MainMode = practiceWrapperMap["main_mode"].(string)
	practicesModesMap := make(map[string]string)

	if subPracticesModes, ok := practiceWrapperMap["sub_practices_modes"]; ok {
		for subPractice, mode := range subPracticesModes.(map[string]any) {
			practicesModesMap[subPractice] = mode.(string)
		}
	}

	practiceWrapper.SubPracticeModes = make([]models.PracticeModeInput, 0, len(practicesModesMap))
	for subPractice, mode := range practicesModesMap {
		practiceWrapper.SubPracticeModes = append(practiceWrapper.SubPracticeModes,
			models.PracticeModeInput{Mode: mode, SubPractice: subPractice})
	}

	if triggersInterface, ok := practiceWrapperMap["triggers"]; ok {
		triggersSet := triggersInterface.(*schema.Set)
		practiceWrapper.Triggers = make([]string, 0, triggersSet.Len())
		for _, trigger := range triggersSet.List() {
			practiceWrapper.Triggers = append(practiceWrapper.Triggers, trigger.(string))
		}
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

func mapToTagsInputs(tagsMap map[string]any) models.TagInput {
	var ret models.TagInput
	ret.Key, ret.Value = tagsMap["key"].(string), tagsMap["value"].(string)
	return ret

}
