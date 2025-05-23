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

	blockTypeLocation = "location_instructions"
	blockTypeServer   = "server_instructions"

	locationConfigEnable   = "isLocationConfigFile"
	locationConfigData     = "locationConfigFile"
	locationConfigFileName = "locationConfigFileName"

	serverConfigEnable   = "isServerConfigFile"
	serverConfigData     = "serverConfigFile"
	serverConfigFileName = "serverConfigFileName"

	redirectToHTTPSEnable = "redirectToHttps"
	accessLogEnable       = "accessLog"
	customHeaderEnable    = "isSetHeader"
	customHeaderData      = "setHeader"
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

	var mtls models.MTLSSchemas
	mtls = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "mtls"), mapToMTLSInput)

	res.ProxySettings = mapMTLSToProxySettingInputs(mtls, res.ProxySettings)

	var additionalBlocks models.BlockSchemas
	additionalBlocks = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "additional_instructions_blocks"), mapToBlocksInput)

	res.ProxySettings = mapBlocksToProxySettingInputs(additionalBlocks, res.ProxySettings)

	var customHeaders models.CustomHeadersSchemas
	customHeaders = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "custom_headers"), mapToCustomHeaderInput)
	res.ProxySettings = mapAdvancedToProxySettingInputs(d.Get("redirect_to_https").(bool), d.Get("access_log").(bool), customHeaders, res.ProxySettings)

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

	if id, ok := tagsMap["id"]; ok {
		ret.ID = id.(string)
	}

	return ret

}

func convertBooleanToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func mapToMTLSInput(mTLSMap map[string]any) models.MTLSSchema {
	mTLSFile, err := utils.UnmarshalAs[models.MTLSSchema](mTLSMap)
	if err != nil {
		fmt.Printf("Failed to convert input schema validation to MTLSSchema struct. Error: %+v", err)
		return models.MTLSSchema{}
	}

	mTLSFile = models.NewFileSchemaEncode(mTLSFile.Filename, mTLSFile.Data, mTLSFile.Type, mTLSFile.CertificateType, mTLSFile.Enable)

	if mTLSMap["filename_id"] != nil {
		mTLSFile.FilenameID = mTLSMap["filename_id"].(string)
	}

	if mTLSMap["data_id"] != nil {
		mTLSFile.DataID = mTLSMap["data_id"].(string)
	}

	if mTLSMap["enable_id"] != nil {
		mTLSFile.EnableID = mTLSMap["enable_id"].(string)
	}

	return mTLSFile
}

func mapMTLSToProxySettingInputs(mTLS models.MTLSSchemas, proxySettings models.ProxySettingInputs) models.ProxySettingInputs {
	for _, mTLSFile := range mTLS {
		var proxySettingEnable, proxySettingData, proxySettingFileName models.ProxySettingInput
		switch mTLSFile.Type {
		case mtlsTypeClient:
			proxySettingEnable.Key = mtlsClientEnable
			proxySettingData.Key = mtlsClientData
			proxySettingFileName.Key = mtlsClientFileName
		case mtlsTypeServer:
			proxySettingEnable.Key = mtlsServerEnable
			proxySettingData.Key = mtlsServerData
			proxySettingFileName.Key = mtlsServerFileName
		default:
			continue
		}

		proxySettingEnable.Value = convertBooleanToString(mTLSFile.Enable)
		proxySettingData.Value = mTLSFile.Data
		proxySettingFileName.Value = mTLSFile.Filename

		proxySettings = append(proxySettings, proxySettingEnable, proxySettingData, proxySettingFileName)
	}

	return proxySettings
}

func mapToBlocksInput(blocksMap map[string]any) models.BlockSchema {
	blockFile, err := utils.UnmarshalAs[models.BlockSchema](blocksMap)
	if err != nil {
		fmt.Printf("Failed to convert input schema validation to BlockSchema struct. Error: %+v", err)
		return models.BlockSchema{}
	}

	blockFile = models.NewFileSchemaEncodeBlocks(blockFile.Filename, blockFile.Data, blockFile.FilenameType, blockFile.Type, blockFile.Enable)

	if blocksMap["filename_id"] != nil {
		blockFile.FilenameID = blocksMap["filename_id"].(string)
	}

	if blocksMap["data_id"] != nil {
		blockFile.DataID = blocksMap["data_id"].(string)
	}

	if blocksMap["enable_id"] != nil {
		blockFile.EnableID = blocksMap["enable_id"].(string)
	}

	return blockFile
}

func mapBlocksToProxySettingInputs(blocks models.BlockSchemas, proxySettings models.ProxySettingInputs) models.ProxySettingInputs {
	blockTypes := make(map[string]bool)
	for _, block := range blocks {
		blockType := block.Type
		if blockTypes[blockType] {
			continue
		} else {
			blockTypes[blockType] = true
		}

		var proxySettingEnable, proxySettingData, proxySettingFileName models.ProxySettingInput
		switch blockType {
		case blockTypeLocation:
			proxySettingEnable.Key = locationConfigEnable
			proxySettingData.Key = locationConfigData
			proxySettingFileName.Key = locationConfigFileName
		case blockTypeServer:
			proxySettingEnable.Key = serverConfigEnable
			proxySettingData.Key = serverConfigData
			proxySettingFileName.Key = serverConfigFileName
		default:
			continue
		}

		proxySettingEnable.Value = convertBooleanToString(block.Enable)
		proxySettingData.Value = block.Data
		proxySettingFileName.Value = block.Filename

		proxySettings = append(proxySettings, proxySettingEnable, proxySettingData, proxySettingFileName)
	}

	return proxySettings
}

func mapToCustomHeaderInput(customHeadersMap map[string]any) models.CustomHeaderSchema {
	var customHeader models.CustomHeaderSchema

	customHeader.Name = customHeadersMap["name"].(string)
	customHeader.Value = customHeadersMap["value"].(string)

	if id, _ := customHeadersMap["header_id"]; id != nil {
		customHeader.HeaderID = id.(string)
	}

	return customHeader
}

func mapAdvancedToProxySettingInputs(redirectToHTTPS, accessLog bool, customHeaders models.CustomHeadersSchemas, proxySettings models.ProxySettingInputs) models.ProxySettingInputs {
	if redirectToHTTPS {
		proxySettingEnable := models.ProxySettingInput{
			Key:   redirectToHTTPSEnable,
			Value: "true",
		}
		proxySettings = append(proxySettings, proxySettingEnable)
	}

	if accessLog {
		proxySettingEnable := models.ProxySettingInput{
			Key:   accessLogEnable,
			Value: "true",
		}
		proxySettings = append(proxySettings, proxySettingEnable)
	}

	if len(customHeaders) == 0 {
		return proxySettings
	}

	proxySettingEnable := models.ProxySettingInput{
		Key:   customHeaderEnable,
		Value: "true",
	}

	proxySettings = append(proxySettings, proxySettingEnable)
	for _, customHeader := range customHeaders {
		var proxySettingData models.ProxySettingInput

		proxySettingData.Key = customHeaderData
		proxySettingData.Value = fmt.Sprintf("%s:%s", customHeader.Name, customHeader.Value)
		proxySettings = append(proxySettings, proxySettingData)
	}

	return proxySettings
}
