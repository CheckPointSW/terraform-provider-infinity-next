package webappasset

import (
	"context"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateWebApplicationAssetInputFromResourceData(d *schema.ResourceData, asset models.WebApplicationAsset) (models.UpdateWebApplicationAssetInput, error) {
	var updateInput models.UpdateWebApplicationAssetInput

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
			if oldSetting.Key == mtlsClientEnable || oldSetting.Key == mtlsClientData || oldSetting.Key == mtlsClientFileName || oldSetting.Key == mtlsServerData || oldSetting.Key == mtlsServerFileName || oldSetting.Key == mtlsServerEnable {
				continue
			}
			if _, ok := newProxySettingsIndicators[oldSetting.Key]; !ok {
				updateInput.RemoveProxySetting = append(updateInput.RemoveProxySetting, oldSetting.ID)
			}
		}
	}

	if oldMTLSs, newMTLSs, hasChange := utils.GetChangeWithParse(d, "mtls", parsemTLSs); hasChange {
		oldMTLSsIndicators := oldMTLSs.ToIndicatorMap()
		mTLSsToAdd := models.FileSchemas{}
		for _, newMTLS := range newMTLSs {
			oldMTLS, ok := oldMTLSsIndicators[newMTLS.Type]
			if !ok {
				mTLSsToAdd = append(mTLSsToAdd, newMTLS)
				//proxysettingstoadd := mapMTLSToProxySettingInputs(newMTLS, models.ProxySettingInputs{})
				//
				//updateInput.AddProxySetting = append(updateInput.AddProxySetting, mapMTLSToProxySettingInputs(newMTLS))
				continue
			}
			if oldMTLS.Enable != newMTLS.Enable {
				var enableToString string
				if newMTLS.Enable {
					enableToString = "true"
				} else {
					enableToString = "false"
				}

				key := mtlsClientEnable
				if oldMTLS.Type == mtlsTypeServer {
					key = mtlsServerEnable
				}
				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldMTLS.EnableID,
					Key:   key,
					Value: enableToString,
				})
			}

			if oldMTLS.Data != newMTLS.Data {
				key := mtlsClientData
				if oldMTLS.Type == mtlsTypeServer {
					key = mtlsServerData
				}

				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldMTLS.DataID,
					Key:   key,
					Value: newMTLS.Data,
				})
			}

			if oldMTLS.Filename != newMTLS.Filename {
				key := mtlsClientFileName
				if oldMTLS.Type == mtlsTypeServer {
					key = mtlsServerFileName
				}

				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldMTLS.FilenameID,
					Key:   key,
					Value: newMTLS.Filename,
				})
			}

			//oldMTLS := oldMTLSsIndicators[newMTLS["type"].(string)]
			//if oldMTLS.Data != newMTLS["data"].(string) || oldMTLS.Enable != newMTLS["enable"].(bool) {
			//	updateInput.UpdateMTLS = append(updateInput.UpdateMTLS, models.UpdateMTLS{
			//		ID:     oldMTLS.ID,
			//		Type:   oldMTLS.Type,
			//		Data:   newMTLS["data"].(string),
			//		Enable: newMTLS["enable"].(bool),
			//	})
			//}
		}

		var proxySettingsToAdd models.ProxySettingInputs
		if mTLSsToAdd != nil {
			proxySettingsToAdd = mapMTLSToProxySettingInputs(mTLSsToAdd, models.ProxySettingInputs{})
		}
		for _, proxySettingToAdd := range proxySettingsToAdd {
			updateInput.AddProxySetting = append(updateInput.AddProxySetting, models.AddProxySetting{
				Key:   proxySettingToAdd.Key,
				Value: proxySettingToAdd.Value,
			})
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
	}

	if oldTags, newTags, hasChange := utils.GetChangeWithParse(d, "tags", parseSchemaTags); hasChange {
		oldTagsIndicatorMap := oldTags.ToIndicatorsMap()
		for _, newTag := range newTags {
			oldTag, ok := oldTagsIndicatorMap[newTag.Key]
			// if tag does not exist - add it
			if !ok {
				updateInput.AddTags = append(updateInput.AddTags, models.AddTag{
					Key:   newTag.Key,
					Value: newTag.Value,
				})

				continue
			}

			// tag exist - check if it needs to be updated
			if oldTag.Value != newTag.Value {
				updateInput.RemoveTags = append(updateInput.RemoveTags, oldTag.ID)
				updateInput.AddTags = append(updateInput.AddTags, models.AddTag{
					Key:   newTag.Key,
					Value: newTag.Value,
				})
			}
		}

		newTagsIndicatorMap := newTags.ToIndicatorsMap()
		for _, oldTag := range oldTags {
			if _, ok := newTagsIndicatorMap[oldTag.Key]; !ok {
				updateInput.RemoveTags = append(updateInput.RemoveTags, oldTag.ID)
			}
		}
	}

	return updateInput, nil
}

func UpdateWebApplicationAsset(ctx context.Context, c *api.Client, id any, input models.UpdateWebApplicationAssetInput) (bool, error) {
	vars := map[string]any{"assetInput": input, "id": id}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation updateWebApplicationAsset($assetInput: WebApplicationAssetUpdateInput!, $id: ID!)
					{	
						updateWebApplicationAsset(assetInput: $assetInput, id: $id) 
					}
				`, "updateWebApplicationAsset", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateWebApplicationAsset response %#v should be of type bool", res)
	}

	return isUpdated, err
}

// parseSchemaSourceIdentifiers converts the source identifiers (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.SourceIdentifierInput
func parseSchemaSourceIdentifiers(sourceIdentifiersFromResourceData any) models.SourceIdentifiersInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](sourceIdentifiersFromResourceData), mapToSourceIdentifierInput)
}

// parseSchemaPracticeWrappers converts the practice wrappers (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.PracticeWrapperInput
func parseSchemaPracticeWrappers(practiceWrappersFromResourceData any) []models.PracticeWrapperInput {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](practiceWrappersFromResourceData), mapToPracticeWrapperInput)
}

// parseSchemaProxySettings converts the proxy settings (type schema.TypeSet) to a slice of map[string]any
// and then converts it to a slice of models.PracticeWrapperInput
func parseSchemaProxySettings(proxySettingsInterfaceFromResourceData any) models.ProxySettingInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](proxySettingsInterfaceFromResourceData), mapToProxySettingInput)
}

// validatePracticeWrapperInput validates that there is no empty modes in the input (because this falis the update api call)
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

func parseSchemaTags(tagsFromResourceData any) models.TagsInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](tagsFromResourceData), mapToTagsInputs)
}

func parsemTLSs(mTLSsFromResourceData any) models.FileSchemas {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](mTLSsFromResourceData), mapToMTLSInput)
}
