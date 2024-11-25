package webapiasset

import (
	"context"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateWebAPIAsset(ctx context.Context, c *api.Client, id any, input models.UpdateWebAPIAssetInput) (bool, error) {
	vars := map[string]any{"assetInput": input, "id": id}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation updateWebAPIAsset($assetInput: WebAPIAssetUpdateInput!, $id: ID!)
					{	
						updateWebAPIAsset(assetInput: $assetInput, id: $id) 
					}
				`, "updateWebAPIAsset", vars)

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateWebAPIAsset response %#v should be of type bool", res)
	}

	return value, err
}

func UpdateWebAPIAssetInputFromResourceData(d *schema.ResourceData) (models.UpdateWebAPIAssetInput, error) {
	var updateInput models.UpdateWebAPIAssetInput

	if _, newName, hasChange := utils.GetChangeWithParse(d, "name", utils.MustValueAs[string]); hasChange {
		updateInput.Name = newName
	}

	if _, newUpstreamURL, hasChange := utils.GetChangeWithParse(d, "upstream_url", utils.MustValueAs[string]); hasChange {
		updateInput.UpstreamURL = newUpstreamURL
	}

	if _, newObjectState, hasChange := utils.GetChangeWithParse(d, "state", utils.MustValueAs[string]); hasChange {
		updateInput.State = newObjectState
	}

	if oldProfilesString, newProfilesString, hasChange := utils.GetChangeWithParse(d, "profiles", utils.MustSchemaCollectionToSlice[string]); hasChange {
		updateInput.AddProfiles, updateInput.RemoveProfiles = utils.SlicesDiff(oldProfilesString, newProfilesString)
	}

	if oldBehaviorsStringList, newBehaviorsStringList, hasChange := utils.GetChangeWithParse(d, "behaviors", utils.MustSchemaCollectionToSlice[string]); hasChange {
		updateInput.AddBehaviors, updateInput.RemoveBehaviors = utils.SlicesDiff(oldBehaviorsStringList, newBehaviorsStringList)
	}

	if _, newIsSharesURLs, hasChange := utils.GetChangeWithParse(d, "is_shares_urls", utils.MustValueAs[bool]); hasChange {
		updateInput.IsSharesURLs = newIsSharesURLs
	}

	if oldURLsString, newURLsString, hasChange := utils.GetChangeWithParse(d, "urls", utils.MustSchemaCollectionToSlice[string]); hasChange {
		oldURLsIDs := utils.MustResourceDataCollectionToSlice[string](d, "urls_ids")
		oldURLsToIDsMap := make(map[string]string)
		for _, oldURLID := range oldURLsIDs {
			urlAndID := strings.Split(oldURLID, models.URLIDSeparator)
			oldURLsToIDsMap[urlAndID[0]] = urlAndID[1]
		}

		addedURLs, removedURLs := utils.SlicesDiff(oldURLsString, newURLsString)
		updateInput.AddURLs = addedURLs
		for _, removedURL := range removedURLs {
			updateInput.RemoveURLs = append(updateInput.RemoveURLs, oldURLsToIDsMap[removedURL])
		}
	}

	if oldPracticeWrappers, newPracticeWrappers, hasChange := utils.GetChangeWithParse(d, "practice", parseSchemaPracticeWrappers); hasChange {
		practiceWrappersInputsToAdd, practiceWrappersInputsToRemove := utils.SlicesDiff(oldPracticeWrappers, newPracticeWrappers)
		practiceWrappersInputsToAdd = utils.Filter(practiceWrappersInputsToAdd, validatePracticeWrapperInput)
		practiceWrappersInputsToRemove = utils.Filter(practiceWrappersInputsToRemove, validatePracticeWrapperInput)
		practiceWrappersToAdd := utils.Map(practiceWrappersInputsToAdd, utils.MustUnmarshalAs[models.AddPracticeWrapper, models.PracticeWrapperInput])
		practicesIDsToRemove := utils.Map(practiceWrappersInputsToRemove, func(wrapper models.PracticeWrapperInput) string { return wrapper.PracticeID })
		updateInput.AddPracticeWrappers = practiceWrappersToAdd
		updateInput.RemovePracticeWrappers = practicesIDsToRemove
	}

	if oldProxySettings, newProxySettings, hasChange := utils.GetChangeWithParse(d, "proxy_setting", parseSchemaProxySettings); hasChange {
		oldProxySettingsIndicators := oldProxySettings.ToIndicatorsMap()
		for _, newSetting := range newProxySettings {
			// if key does not exist then this is a new setting to add
			if _, ok := oldProxySettingsIndicators[newSetting.Key]; !ok {
				updateInput.AddProxySetting = append(updateInput.AddProxySetting, models.AddProxySetting{
					Key:   newSetting.Key,
					Value: newSetting.Value,
				})
				continue
			}

			// we know the key exist
			// if the value is different - update the setting
			oldSetting := oldProxySettingsIndicators[newSetting.Key]
			if oldSetting.Value != newSetting.Value {
				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					Key:   newSetting.Key,
					Value: newSetting.Value,
					ID:    oldSetting.ID,
				})
			}
		}

		newProxySettingsIndicators := newProxySettings.ToIndicatorsMap()
		for _, oldSetting := range oldProxySettings {
			if _, ok := newProxySettingsIndicators[oldSetting.Key]; !ok {
				updateInput.RemoveProxySetting = append(updateInput.RemoveProxySetting, oldSetting.ID)
			}
		}
	}

	if oldSourceIdentifiers, newSourceIdentifiers, hasChange := utils.GetChangeWithParse(d, "source_identifier", parseSchemaSourceIdentifiers); hasChange {
		oldSourceIdentifiersIndicatorMap := oldSourceIdentifiers.ToIndicatorsMap()
		for _, newSourceIdentifier := range newSourceIdentifiers {
			// if source identifier does not exist - add it
			if _, ok := oldSourceIdentifiersIndicatorMap[newSourceIdentifier.SourceIdentifier]; !ok {
				updateInput.AddSourceIdentifiers = append(updateInput.AddSourceIdentifiers, models.AddSourceIdentifier{
					SourceIdentifier: newSourceIdentifier.SourceIdentifier,
					Values:           newSourceIdentifier.Values,
				})

				continue
			}

			// source identifier exist - check if it needs to be updated
			oldSourceIdentifier := oldSourceIdentifiersIndicatorMap[newSourceIdentifier.SourceIdentifier]
			valuesToAdd, valuesToRemove := utils.SlicesDiff(oldSourceIdentifier.Values, newSourceIdentifier.Values)
			valuesIDsIndicatorsMap := oldSourceIdentifier.ValuesIDs.ToIndicatorsMap()
			var valuesIDsToRemove []string
			for _, valueToRemove := range valuesToRemove {
				valuesIDsToRemove = append(valuesIDsToRemove, valuesIDsIndicatorsMap[valueToRemove])
			}

			updateInput.UpdateSourceIdentifiers = append(updateInput.UpdateSourceIdentifiers, models.UpdateSourceIdentifier{
				ID:               oldSourceIdentifier.ID,
				SourceIdentifier: oldSourceIdentifier.SourceIdentifier,
				AddValues:        valuesToAdd,
				RemoveValues:     valuesIDsToRemove,
				UpdateValues:     []models.UpdateSourceIdentifierValue{},
			})
		}

		newSourceIdentifiersIndicatorMap := newSourceIdentifiers.ToIndicatorsMap()
		for _, oldSourceIdentifier := range oldSourceIdentifiers {
			if _, ok := newSourceIdentifiersIndicatorMap[oldSourceIdentifier.SourceIdentifier]; !ok {
				updateInput.RemoveSourceIdentifiers = append(updateInput.RemoveSourceIdentifiers, oldSourceIdentifier.ID)
			}
		}

		if oldTags, newTags, hasChange := utils.GetChangeWithParse(d, "tags", parseSchemaTags); hasChange {
			tagsInputsToAdd, tagsInputsToRemove := utils.SlicesDiff(oldTags, newTags)
			tagsInputsToAdd = utils.Filter(tagsInputsToAdd, validateTag)
			tagsInputsToRemove = utils.Filter(tagsInputsToRemove, validateTag)
			tagsToAdd := utils.Map(tagsInputsToAdd, utils.MustUnmarshalAs[models.AddTag, models.TagInput])
			tagsToRemove := utils.Map(tagsInputsToRemove, func(tag models.TagInput) string { return tag.ID })
			updateInput.AddTags = tagsToAdd
			updateInput.RemoveTags = tagsToRemove
		}
	}

	return updateInput, nil
}

// parseSchemaSourceIdentifiers converts the source identifiers (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.SourceIdentifierInput
func parseSchemaSourceIdentifiers(sourceIdentifiersFromResourceData any) models.SourceIdentifiersInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](sourceIdentifiersFromResourceData), mapToSourceIdentifierInput)
}

// parseSchemaPracticeWrappers converts the practice wrappers (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.PracticeWrapperInput
func parseSchemaPracticeWrappers(practiceWrappersFromResourceData any) models.PracticeWrappersInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](practiceWrappersFromResourceData), mapToPracticeWrapperInput)
}

// parseSchemaProxySettings converts the proxy settings (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.PracticeWrapperInput
func parseSchemaProxySettings(proxySettingsInterfaceFromResourceData any) models.ProxySettingInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](proxySettingsInterfaceFromResourceData), mapToProxySettingInput)
}

// parseSchemaTags converts the tags (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.TagInput
func parseSchemaTags(tagsFromResourceData any) models.TagsInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](tagsFromResourceData), mapToTagInput)
}

// validatePracticeWrapperInput validates that there is no empty modes in the input (because this fails the update api call)
// this function is used during update of a practice since the getChange func of the terraform helper package
// sometimes returns an extra empty practice
func validatePracticeWrapperInput(practice models.PracticeWrapperInput) bool {
	if practice.PracticeID == "" || practice.MainMode == "" {
		return false
	}

	for _, mode := range practice.SubPracticeModes {
		if mode.Mode == "" {
			return false
		}
	}

	return true
}

func validateTag(tag models.TagInput) bool {
	return tag.Key != "" && tag.Value != ""
}
