package webapiasset

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func proxySettingKeyToBlockType(proxySettingKey string) string {
	switch proxySettingKey {
	case mtlsClientEnable, mtlsClientData, mtlsClientFileName:
		return mtlsTypeClient
	case mtlsServerEnable, mtlsServerData, mtlsServerFileName:
		return mtlsTypeServer
	case locationConfigEnable, locationConfigData, locationConfigFileName:
		return blockTypeLocation
	case serverConfigEnable, serverConfigData, serverConfigFileName:
		return blockTypeServer
	case redirectToHTTPSEnable, accessLogEnable, customHeaderEnable, customHeaderData:
		return proxySettingKey
	default:
		return ""
	}

}

func ReadWebAPIAssetToResourceData(asset models.WebAPIAsset, d *schema.ResourceData) error {
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
	d.Set("state", asset.State)

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
			switch blockType {
			case redirectToHTTPSEnable:
				d.Set("redirect_to_https", proxySetting.Value == "true")
				d.Set("redirect_to_https_id", proxySetting.ID)
			case accessLogEnable:
				d.Set("access_log", proxySetting.Value == "true")
				d.Set("access_log_id", proxySetting.ID)
			case customHeaderEnable:
				d.Set("custom_headers_id", proxySetting.ID)
			case customHeaderData:
				nameAndValue := strings.SplitN(proxySetting.Value, ":", 2)
				customHeaderSchema := models.CustomHeaderSchema{
					HeaderID: proxySetting.ID,
					Name:     nameAndValue[0],
					Value:    nameAndValue[1],
				}

				customHeadersSchemaMap[proxySetting.ID] = customHeaderSchema
			case blockTypeLocation, blockTypeServer:
				if _, ok := blocksSchemaMap[blockType]; !ok {
					blocksSchemaMap[blockType] = models.BlockSchema{}
				}

			default:
				if _, ok := mTLSsSchemaMap[blockType]; !ok {
					mTLSsSchemaMap[blockType] = models.MTLSSchema{}
				}

			}

			switch proxySetting.Key {
			case mtlsClientEnable, mtlsServerEnable, locationConfigEnable, serverConfigEnable:
				enable := proxySetting.Value == "true"
				if blockType == blockTypeLocation || blockType == blockTypeServer {
					blocksSchemaMap[blockType] = models.BlockSchema{
						FilenameID:   blocksSchemaMap[blockType].FilenameID,
						Filename:     blocksSchemaMap[blockType].Filename,
						FilenameType: blocksSchemaMap[blockType].FilenameType,
						DataID:       blocksSchemaMap[blockType].DataID,
						Data:         blocksSchemaMap[blockType].Data,
						Type:         blockType,
						EnableID:     proxySetting.ID,
						Enable:       enable,
					}

				} else {
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
					if blockType == blockTypeLocation || blockType == blockTypeServer {
						fileExtensionsByType = models.MimeTypeToFileExtension(mimeType, false)
					} else {
						fileExtensionsByType = models.MimeTypeToFileExtension(mimeType, true)
					}

				}

				if blockType == blockTypeLocation || blockType == blockTypeServer {
					blocksSchemaMap[blockType] = models.BlockSchema{
						FilenameID:   blocksSchemaMap[blockType].FilenameID,
						Filename:     blocksSchemaMap[blockType].Filename,
						FilenameType: fileExtensionsByType,
						DataID:       proxySetting.ID,
						Data:         decodedData,
						Type:         blockType,
						EnableID:     blocksSchemaMap[blockType].EnableID,
						Enable:       blocksSchemaMap[blockType].Enable,
					}

				} else {
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

				}

			case mtlsClientFileName, mtlsServerFileName, locationConfigFileName, serverConfigFileName:
				if blockType == blockTypeLocation || blockType == blockTypeServer {
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

				} else {
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

	var blockSchemas models.BlockSchemas
	for _, blockSchema := range blocksSchemaMap {
		blockSchemas = append(blockSchemas, blockSchema)
	}

	// Sort the blocks by their type
	sort.Slice(blockSchemas, func(i, j int) bool {
		return blockSchemas[i].Type < blockSchemas[j].Type
	})

	for _, blockSchema := range blockSchemas {
		block, err := utils.UnmarshalAs[map[string]any](blockSchema)
		if err != nil {
			return fmt.Errorf("failed to convert %s block to map. Error: %+v", blockSchema.Type, err)
		}

		additionalBlocksMap = append(additionalBlocksMap, block)
	}

	var customHeaderSchemas models.CustomHeadersSchemas
	for _, customHeaderSchema := range customHeadersSchemaMap {
		customHeaderSchemas = append(customHeaderSchemas, customHeaderSchema)
	}

	// Sort the custom headers by their name
	sort.Slice(customHeaderSchemas, func(i, j int) bool {
		return customHeaderSchemas[i].Name < customHeaderSchemas[j].Name
	})

	for _, customHeaderSchema := range customHeaderSchemas {
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

	if len(tagsSchemaMap) > 0 {
		d.Set("tags", tagsSchemaMap)
	}

	return nil
}

func GetWebAPIAsset(ctx context.Context, c *api.Client, id string) (models.WebAPIAsset, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getWebAPIAsset(id: "`+id+`") {
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
					name
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
			}
		}
	`, "getWebAPIAsset")

	if err != nil {
		return models.WebAPIAsset{}, fmt.Errorf("failed to get WebAPIAsset: %w", err)
	}

	if res == nil {
		return models.WebAPIAsset{}, api.ErrorNotFound
	}

	asset, err := utils.UnmarshalAs[models.WebAPIAsset](res)
	if err != nil {
		return models.WebAPIAsset{}, fmt.Errorf("failed to convert graphQL response to WebAPIAsset struct. Error: %w", err)
	}

	return asset, nil
}
