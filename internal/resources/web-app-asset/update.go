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
			// if the key is one of the advanced proxy settings - skip it
			if proxySettingKeyToBlockType(oldSetting.Key) != "" {
				continue
			}

			if _, ok := newProxySettingsIndicators[oldSetting.Key]; !ok {
				updateInput.RemoveProxySetting = append(updateInput.RemoveProxySetting, oldSetting.ID)
			}
		}
	}

	if oldMTLSs, newMTLSs, hasChange := utils.GetChangeWithParse(d, "mtls", parsemTLSs); hasChange {
		oldMTLSsIndicators := oldMTLSs.ToIndicatorMap()
		mTLSsToAdd := models.MTLSSchemas{}
		for _, newMTLS := range newMTLSs {
			oldMTLS, ok := oldMTLSsIndicators[newMTLS.Type]
			if !ok {
				mTLSsToAdd = append(mTLSsToAdd, newMTLS)
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

	if _, newRedirectToHTTPS, hasChange := utils.GetChangeWithParse(d, "redirect_to_https", utils.MustValueAs[bool]); hasChange {
		value := "false"
		if newRedirectToHTTPS {
			value = "true"
		}

		if id := d.Get("redirect_to_https_id").(string); id != "" {
			updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
				ID:    id,
				Key:   redirectToHTTPSEnable,
				Value: value,
			})
		} else {
			updateInput.AddProxySetting = append(updateInput.AddProxySetting, models.AddProxySetting{
				Key:   redirectToHTTPSEnable,
				Value: value,
			})
		}

	}

	if _, newAccessLog, hasChange := utils.GetChangeWithParse(d, "access_log", utils.MustValueAs[bool]); hasChange {
		value := "false"
		if newAccessLog {
			value = "true"
		}

		if id := d.Get("access_log_id").(string); id != "" {
			updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
				ID:    id,
				Key:   accessLogEnable,
				Value: value,
			})
		} else {
			updateInput.AddProxySetting = append(updateInput.AddProxySetting, models.AddProxySetting{
				Key:   accessLogEnable,
				Value: value,
			})
		}

	}

	if oldCustomHeaders, newCustomHeaders, hasChange := utils.GetChangeWithParse(d, "custom_headers", parseCustomHeaders); hasChange {
		if len(newCustomHeaders) == 0 {
			updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
				ID:    d.Get("custom_headers_id").(string),
				Key:   customHeaderEnable,
				Value: "false",
			})
		}

		if len(oldCustomHeaders) == 0 {
			if id := d.Get("custom_headers_id").(string); id != "" {
				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    id,
					Key:   customHeaderEnable,
					Value: "true",
				})
			} else {
				updateInput.AddProxySetting = append(updateInput.AddProxySetting, models.AddProxySetting{
					Key:   customHeaderEnable,
					Value: "true",
				})
			}
		}

		oldCustomHeadersIndicatorMap := oldCustomHeaders.ToIndicatorMap()
		customHeadersToAdd := models.CustomHeadersSchemas{}
		for _, newCustomHeader := range newCustomHeaders {
			nameAndValue := fmt.Sprintf("%s:%s", newCustomHeader.Name, newCustomHeader.Value)
			_, ok := oldCustomHeadersIndicatorMap[nameAndValue]
			if !ok {
				customHeadersToAdd = append(customHeadersToAdd, newCustomHeader)
				continue
			}

		}

		newCustomHeadersIndicatorMap := newCustomHeaders.ToIndicatorMap()
		for _, oldCustomHeader := range oldCustomHeaders {
			nameAndValue := fmt.Sprintf("%s:%s", oldCustomHeader.Name, oldCustomHeader.Value)
			if _, ok := newCustomHeadersIndicatorMap[nameAndValue]; !ok {
				updateInput.RemoveProxySetting = append(updateInput.RemoveProxySetting, oldCustomHeader.HeaderID)
			}

		}

		for _, customHeaderToAdd := range customHeadersToAdd {
			updateInput.AddProxySetting = append(updateInput.AddProxySetting, models.AddProxySetting{
				Key:   customHeaderData,
				Value: fmt.Sprintf("%s:%s", customHeaderToAdd.Name, customHeaderToAdd.Value),
			})
		}
	}

	if oldBlocks, newBlocks, hasChange := utils.GetChangeWithParse(d, "additional_instructions_blocks", parseBlocks); hasChange {
		oldBlocksIndicatorMap := oldBlocks.ToIndicatorMap()
		additionalBlocksToAdd := models.BlockSchemas{}
		fmt.Printf("oldBlocks: %+v\n", oldBlocks)
		fmt.Printf("newBlocks: %+v\n", newBlocks)
		for _, newBlock := range newBlocks {
			oldBlock, ok := oldBlocksIndicatorMap[newBlock.Type]
			fmt.Printf("oldBlock=%+v newBlock=%+v. oldBlocksIndicatorMap=%+v\n", oldBlock, newBlock, oldBlocksIndicatorMap)
			if !ok {
				if newBlock.Enable {
					additionalBlocksToAdd = append(additionalBlocksToAdd, newBlock)
				}
				fmt.Printf("Continuing\n")
				continue
			}

			if !newBlock.Enable {
				key := serverConfigEnable
				if newBlock.Type == blockTypeLocation {
					key = locationConfigEnable
				}

				fmt.Printf("adding to updateProxySetting 1: %s", oldBlock.EnableID)
				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldBlock.EnableID,
					Key:   key,
					Value: "false",
				})

				if oldBlock.Enable {
					if oldBlock.FilenameID != "" {
						updateInput.RemoveProxySetting = append(updateInput.RemoveProxySetting, oldBlock.FilenameID)
					}

					if oldBlock.DataID != "" {
						updateInput.RemoveProxySetting = append(updateInput.RemoveProxySetting, oldBlock.DataID)
					}

				}

				continue
			}

			if oldBlock.Enable != newBlock.Enable {
				enableKey := serverConfigEnable
				dataKey := serverConfigData
				filenameKey := serverConfigFileName
				if newBlock.Type == blockTypeLocation {
					enableKey = locationConfigEnable
					dataKey = locationConfigData
					filenameKey = locationConfigFileName
				}
				fmt.Printf("adding to updateProxySetting 2: %s", oldBlock.EnableID)
				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldBlock.EnableID,
					Key:   enableKey,
					Value: "true",
				})

				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldBlock.DataID,
					Key:   dataKey,
					Value: newBlock.Data,
				})

				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldBlock.FilenameID,
					Key:   filenameKey,
					Value: newBlock.Filename,
				})

			}

			if oldBlock.Data != newBlock.Data || oldBlock.Filename != newBlock.Filename {
				dataKey := serverConfigData
				filenameKey := serverConfigFileName
				if newBlock.Type == blockTypeLocation {
					dataKey = locationConfigData
					filenameKey = locationConfigFileName
				}

				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldBlock.DataID,
					Key:   dataKey,
					Value: newBlock.Data,
				})

				updateInput.UpdateProxySetting = append(updateInput.UpdateProxySetting, models.UpdateProxySetting{
					ID:    oldBlock.FilenameID,
					Key:   filenameKey,
					Value: newBlock.Filename,
				})

			}

		}

		var proxySettingsToAdd models.ProxySettingInputs
		if additionalBlocksToAdd != nil {
			proxySettingsToAdd = mapBlocksToProxySettingInputs(additionalBlocksToAdd, models.ProxySettingInputs{})
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

func parsemTLSs(mTLSsFromResourceData any) models.MTLSSchemas {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](mTLSsFromResourceData), mapToMTLSInput)
}

func parseBlocks(blocksFromResourceData any) models.BlockSchemas {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](blocksFromResourceData), mapToBlocksInput)
}

func parseCustomHeaders(customHeadersFromResourceData any) models.CustomHeadersSchemas {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](customHeadersFromResourceData), mapToCustomHeaderInput)
}
