package webappasset

import (
	"context"
	"encoding/base64"
	"fmt"
	webAPIAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func proxySettingKeyToBlockType(proxySettingKey string) string {
	if proxySettingKey == mtlsClientEnable || proxySettingKey == mtlsClientData || proxySettingKey == mtlsClientFileName {
		return mtlsTypeClient
	}
	if proxySettingKey == mtlsServerEnable || proxySettingKey == mtlsServerData || proxySettingKey == mtlsServerFileName {
		return mtlsTypeServer
	}

	if proxySettingKey == locationConfigEnable || proxySettingKey == locationConfigData || proxySettingKey == locationConfigFileName {
		return blockTypeLocation
	}

	if proxySettingKey == serverConfigEnable || proxySettingKey == serverConfigData || proxySettingKey == serverConfigFileName {
		return blockTypeServer
	}
	if proxySettingKey == redirectToHTTPSEnable || proxySettingKey == accessLogEnable || proxySettingKey == customHeaderEnable || proxySettingKey == customHeaderData {
		return proxySettingKey
	}

	return ""
}

func ReadWebApplicationAssetToResourceData(asset models.WebApplicationAsset, d *schema.ResourceData) error {
	d.SetId(asset.ID)
	d.Set("name", asset.Name)
	d.Set("asset_type", asset.AssetType)
	d.Set("main_attributes", asset.MainAttributes)
	d.Set("sources", asset.Sources)
	d.Set("family", asset.Family)
	d.Set("category", asset.Category)
	d.Set("class", asset.Class)
	d.Set("order", asset.Order)
	d.Set("group", asset.Group)
	d.Set("kind", asset.Kind)
	d.Set("intelligence_tags", asset.IntelligenceTags)
	d.Set("read_only", asset.ReadOnly)
	d.Set("upstream_url", asset.UpstreamURL)
	d.Set("behaviors", asset.Behaviors.ToSchema())
	d.Set("profiles", asset.Profiles.ToSchema())
	d.Set("is_shares_urls", asset.IsSharesURLs)

	var proxySettingsSchemaMap []map[string]any
	mTLSsSchemaMap := make(map[string]models.MTLSSchema)
	var mTLSsMap []map[string]any
	blocksSchemaMap := make(map[string]models.BlockSchema)
	var additionalBlocksMap []map[string]any
	customHeadersSchemaMap := make(map[string]models.CustomHeaderSchema)
	var customHeadersMap []map[string]any

	for _, proxySetting := range asset.ProxySettings {
		blockType := proxySettingKeyToBlockType(proxySetting.Key)
		if blockType != "" {
			if blockType == redirectToHTTPSEnable {
				d.Set("redirect_to_https", proxySetting.Value == "true")
				d.Set("redirect_to_https_id", proxySetting.ID)
				continue
			}

			if blockType == accessLogEnable {
				d.Set("access_log", proxySetting.Value == "true")
				d.Set("access_log_id", proxySetting.ID)
				continue
			}

			if blockType == customHeaderEnable {
				d.Set("custom_header_id", proxySetting.ID)
				continue
			}

			if blockType == customHeaderData {
				nameAndValue := strings.SplitN(proxySetting.Value, ":", 2)
				customHeaderSchema := models.CustomHeaderSchema{
					HeaderID: proxySetting.ID,
					Name:     nameAndValue[0],
					Value:    nameAndValue[1],
				}

				customHeadersSchemaMap[proxySetting.ID] = customHeaderSchema
				continue
			}

			if blockType == blockTypeLocation || blockType == blockTypeServer {
				if _, ok := blocksSchemaMap[blockType]; !ok {
					blocksSchemaMap[blockType] = models.BlockSchema{}
				}

			}

			if _, ok := mTLSsSchemaMap[blockType]; !ok {
				mTLSsSchemaMap[blockType] = models.MTLSSchema{}
			}

			switch proxySetting.Key {
			case mtlsClientEnable, mtlsServerEnable, locationConfigEnable, serverConfigEnable:
				enable := false
				if proxySetting.Value == "true" {
					enable = true
				}

				if blockType == blockTypeServer || blockType == blockTypeLocation {
					blocksSchemaMap[blockType] = models.BlockSchema{
						FilenameID:   proxySetting.ID,
						Filename:     proxySetting.Value,
						FilenameType: blocksSchemaMap[blockType].FilenameType,
						DataID:       blocksSchemaMap[blockType].DataID,
						Data:         blocksSchemaMap[blockType].Data,
						Type:         blockType,
						EnableID:     proxySetting.ID,
						Enable:       enable,
					}
					continue
				}

				mTLSsSchemaMap[blockType] = models.MTLSSchema{
					FilenameID:      mTLSsSchemaMap[blockType].FilenameID,
					Filename:        mTLSsSchemaMap[blockType].Filename,
					CertificateType: mTLSsSchemaMap[blockType].CertificateType,
					DataID:          mTLSsSchemaMap[blockType].DataID,
					Data:            mTLSsSchemaMap[blockType].Data,
					Type:            blockType,
					EnableID:        proxySetting.ID,
					Enable:          enable,
				}
			case mtlsClientData, mtlsServerData, locationConfigData, serverConfigData:
				var decodedData string
				var fileExtensionsByType string
				// proxySetting.Value format is "data:<mimeType>;base64,<base64Data>"
				if strings.Contains(proxySetting.Value, "base64,") {
					b64Data := strings.SplitN(proxySetting.Value, "base64,", 2)[1]
					bDecodedData, err := base64.StdEncoding.DecodeString(b64Data)
					if err != nil {
						return fmt.Errorf("failed decoding base64 string %s: %w", b64Data, err)
					}

					decodedData = string(bDecodedData)

					mimeType := strings.SplitN(proxySetting.Value, ":", 2)[1]
					mimeType = strings.SplitN(mimeType, ";", 2)[0]
					fileExtensionsByType = webAPIAssetModels.MimeTypeToFileExtension(mimeType)
				}

				if blockType == blockTypeServer || blockType == blockTypeLocation {
					blocksSchemaMap[blockType] = models.BlockSchema{
						FilenameID:   proxySetting.ID,
						Filename:     proxySetting.Value,
						FilenameType: blocksSchemaMap[blockType].FilenameType,
						DataID:       proxySetting.ID,
						Data:         decodedData,
						Type:         blockType,
						EnableID:     blocksSchemaMap[blockType].EnableID,
						Enable:       blocksSchemaMap[blockType].Enable,
					}
					continue
				}

				mTLSsSchemaMap[blockType] = models.MTLSSchema{
					FilenameID:      mTLSsSchemaMap[blockType].FilenameID,
					Filename:        mTLSsSchemaMap[blockType].Filename,
					CertificateType: fileExtensionsByType,
					DataID:          proxySetting.ID,
					Data:            decodedData,
					Type:            blockType,
					EnableID:        mTLSsSchemaMap[blockType].EnableID,
					Enable:          mTLSsSchemaMap[blockType].Enable,
				}
			case mtlsClientFileName, mtlsServerFileName, locationConfigFileName, serverConfigFileName:
				if blockType == blockTypeServer || blockType == blockTypeLocation {
					blocksSchemaMap[blockType] = models.BlockSchema{
						FilenameID:   proxySetting.ID,
						Filename:     proxySetting.Value,
						FilenameType: blocksSchemaMap[blockType].FilenameType,
						DataID:       blocksSchemaMap[blockType].DataID,
						Data:         blocksSchemaMap[blockType].Data,
						Type:         blockType,
						EnableID:     blocksSchemaMap[blockType].EnableID,
						Enable:       blocksSchemaMap[blockType].Enable,
					}
					continue
				}

				mTLSsSchemaMap[blockType] = models.MTLSSchema{
					FilenameID:      proxySetting.ID,
					Filename:        proxySetting.Value,
					CertificateType: mTLSsSchemaMap[blockType].CertificateType,
					DataID:          mTLSsSchemaMap[blockType].DataID,
					Data:            mTLSsSchemaMap[blockType].Data,
					Type:            blockType,
					EnableID:        mTLSsSchemaMap[blockType].EnableID,
					Enable:          mTLSsSchemaMap[blockType].Enable,
				}
			default:
				continue
			}
		} else {
			proxySettingSchemaMap, err := utils.UnmarshalAs[map[string]any](proxySetting)
			if err != nil {
				return fmt.Errorf("failed to convert proxy setting to map. Error: %+v", err)
			}

			proxySettingsSchemaMap = append(proxySettingsSchemaMap, proxySettingSchemaMap)
		}
	}

	for _, mTLSSchema := range mTLSsSchemaMap {
		mTLS, err := utils.UnmarshalAs[map[string]any](mTLSSchema)
		if err != nil {
			return fmt.Errorf("failed to convert mTLS to map. Error: %+v", err)
		}

		mTLSsMap = append(mTLSsMap, mTLS)
	}

	for _, blockSchema := range blocksSchemaMap {
		block, err := utils.UnmarshalAs[map[string]any](blockSchema)
		if err != nil {
			return fmt.Errorf("failed to convert %s block to map. Error: %+v", blockSchema.Type, err)
		}

		additionalBlocksMap = append(additionalBlocksMap, block)
	}

	for _, customHeaderSchema := range customHeadersSchemaMap {
		customHeader, err := utils.UnmarshalAs[map[string]any](customHeaderSchema)
		if err != nil {
			return fmt.Errorf("failed to convert custom header to map. Error: %+v", err)
		}

		customHeadersMap = append(customHeadersMap, customHeader)
	}

	d.Set("proxy_setting", proxySettingsSchemaMap)
	d.Set("mtls", mTLSsMap)
	d.Set("additional_instructions_blocks", additionalBlocksMap)
	d.Set("custom_headers", customHeadersMap)

	sourceIdentifiersSchema := asset.SourceIdentifiers.ToSchema()
	sourceIdentifiersSchemaMap, err := utils.UnmarshalAs[[]map[string]any](sourceIdentifiersSchema)
	if err != nil {
		return fmt.Errorf("failed to convert source identifiers to slice of maps. Error: %+v", err)
	}

	d.Set("source_identifier", sourceIdentifiersSchemaMap)

	schemaURLs, schemaURLsIDs := asset.URLs.ToSchema()
	d.Set("urls", schemaURLs)
	d.Set("urls_ids", schemaURLsIDs)

	schemaPracticeWrappers := asset.Practices.ToSchema()
	schemaPracticeWrappersMap, err := utils.UnmarshalAs[[]map[string]any](schemaPracticeWrappers)
	if err != nil {
		return fmt.Errorf("failed to convert practices to slice of maps. Error: %+v", err)
	}

	d.Set("practice", schemaPracticeWrappersMap)

	tagsSchemaMap, err := utils.UnmarshalAs[[]map[string]any](asset.Tags)
	if err != nil {
		return fmt.Errorf("failed to convert tags to slice of maps. Error: %+v", err)
	}

	d.Set("tags", tagsSchemaMap)

	return nil
}

func GetWebApplicationAsset(ctx context.Context, c *api.Client, id string) (models.WebApplicationAsset, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getWebApplicationAsset(id: "`+id+`") {
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
				behaviors {
					id
				}
				tags {
					id
					key
					value
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
	`, "getWebApplicationAsset")

	if err != nil {
		return models.WebApplicationAsset{}, fmt.Errorf("failed to get WebApplicationAsset: %w", err)
	}

	if res == nil {
		return models.WebApplicationAsset{}, api.ErrorNotFound
	}

	asset, err := utils.UnmarshalAs[models.WebApplicationAsset](res)
	if err != nil {
		return models.WebApplicationAsset{}, fmt.Errorf("failed to convert graphQL response to WebApplicationAsset struct. Error: %w", err)
	}

	return asset, nil
}
