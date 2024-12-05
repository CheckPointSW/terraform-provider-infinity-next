package webappasset

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mtlsTypeClient = "client"
	mtlsTypeServer = "server"

	mtlsClientEnable   = "isUpstreamTrustedCAFile"
	mtlsClientData     = "upstreamTrustedCAFile"
	mtlsClientFileName = "upstreamTrustedCAFileName"

	mtlsServerEnable   = "isTrustedCAListFile"
	mtlsServerData     = "trustedCAListFile"
	mtlsServerFileName = "trustedCAListFileName"
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

	var mtls []models.FileSchema
	mtls = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "mtls"), mapToMTLSInput)

	res.ProxySettings = mapMTLSToProxySettingInputs(mtls, res.ProxySettings)

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

func mapToMTLSInput(mTLSMap map[string]any) models.FileSchema {
	mTLSFile, err := utils.UnmarshalAs[models.FileSchema](mTLSMap["file"])
	if err != nil {
		fmt.Printf("Failed to convert input schema validation to FileSchema struct. Error: %+v", err)
		return models.FileSchema{}
	}

	mTLSFile = models.NewFileSchemaEncode(mTLSFile.Filename, mTLSFile.Data, mTLSFile.Type, mTLSFile.Enable)

	return mTLSFile
}

func mapMTLSToProxySettingInputs(mTLS []models.FileSchema, proxySettings models.ProxySettingInputs) models.ProxySettingInputs {
	for _, mTLSFile := range mTLS {
		var proxySettingEnable, proxySettingData, proxySettingFileName models.ProxySettingInput
		switch mTLSFile.Type {
		case "client":
			proxySettingEnable.Key = "isUpstreamTrustedCAFile"
			proxySettingData.Key = "upstreamTrustedCAFile"
			proxySettingFileName.Key = "upstreamTrustedCAFileName"
		case "server":
			proxySettingEnable.Key = "isTrustedCAListFile"
			proxySettingData.Key = "trustedCAListFile"
			proxySettingFileName.Key = "trustedCAListFileName"
		default:
			continue
		}

		if mTLSFile.Enable {
			proxySettingEnable.Value = "true"
		} else {
			proxySettingEnable.Value = "false"
		}

		proxySettingData.Value = mTLSFile.Data
		proxySettingFileName.Value = mTLSFile.Filename

		proxySettings = append(proxySettings, proxySettingEnable, proxySettingData, proxySettingFileName)
	}

	return proxySettings
}
