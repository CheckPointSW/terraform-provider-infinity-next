package webappasset

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func proxySettingKeyTomTLSType(proxySettingKey string) string {
	if proxySettingKey == mtlsClientEnable || proxySettingKey == mtlsClientData || proxySettingKey == mtlsClientFileName {
		return mtlsTypeClient
	}
	if proxySettingKey == mtlsServerEnable || proxySettingKey == mtlsServerData || proxySettingKey == mtlsServerFileName {
		return mtlsTypeServer
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

	//proxySettingsSchemaMap, err := utils.UnmarshalAs[[]map[string]any](asset.ProxySettings)
	//if err != nil {
	//	return fmt.Errorf("failed to convert proxy settings to slice of maps. Error: %+v", err)
	//}
	//
	//d.Set("proxy_setting", proxySettingsSchemaMap)

	var proxySettingsSchemaMap []map[string]any
	mTLSsSchemaMap := make(map[string]models.FileSchema)
	var mTLSsMap []map[string]any

	for _, proxySetting := range asset.ProxySettings {
		mTLSType := proxySettingKeyTomTLSType(proxySetting.Key)
		if mTLSType != "" {
			if _, ok := mTLSsSchemaMap[mTLSType]; !ok {
				mTLSsSchemaMap[mTLSType] = models.FileSchema{}
			}
			switch proxySetting.Key {
			case mtlsClientEnable, mtlsServerEnable:
				if proxySetting.Value == "true" {
					mTLSsSchemaMap[mTLSType] = models.FileSchema{
						FilenameID: mTLSsSchemaMap[mTLSType].FilenameID,
						Filename:   mTLSsSchemaMap[mTLSType].Filename,
						DataID:     mTLSsSchemaMap[mTLSType].DataID,
						Data:       mTLSsSchemaMap[mTLSType].Data,
						Type:       mTLSType,
						EnableID:   proxySetting.ID,
						Enable:     true,
					}
				}
				if proxySetting.Value == "false" {
					mTLSsSchemaMap[mTLSType] = models.FileSchema{
						FilenameID: mTLSsSchemaMap[mTLSType].FilenameID,
						Filename:   mTLSsSchemaMap[mTLSType].Filename,
						DataID:     mTLSsSchemaMap[mTLSType].DataID,
						Data:       mTLSsSchemaMap[mTLSType].Data,
						Type:       mTLSType,
						EnableID:   proxySetting.ID,
						Enable:     false,
					}
				}
			case mtlsClientData, mtlsServerData:
				mTLSsSchemaMap[mTLSType] = models.FileSchema{
					FilenameID: mTLSsSchemaMap[mTLSType].FilenameID,
					Filename:   mTLSsSchemaMap[mTLSType].Filename,
					DataID:     proxySetting.ID,
					Data:       proxySetting.Value,
					Type:       mTLSType,
					EnableID:   mTLSsSchemaMap[mTLSType].EnableID,
					Enable:     mTLSsSchemaMap[mTLSType].Enable,
				}
			case mtlsClientFileName, mtlsServerFileName:
				mTLSsSchemaMap[mTLSType] = models.FileSchema{
					FilenameID: proxySetting.ID,
					Filename:   proxySetting.Value,
					DataID:     mTLSsSchemaMap[mTLSType].DataID,
					Data:       mTLSsSchemaMap[mTLSType].Data,
					Type:       mTLSType,
					EnableID:   mTLSsSchemaMap[mTLSType].EnableID,
					Enable:     mTLSsSchemaMap[mTLSType].Enable,
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
	//case mtlsServerEnable:
	//	if proxySetting.Value == "true" {
	//		mTLSsSchemaMap[mTLSType] = models.FileSchema{
	//			FilenameID: mTLSsSchemaMap[mTLSType].FilenameID,
	//			Filename:   mTLSsSchemaMap[mTLSType].Filename,
	//			DataID:     mTLSsSchemaMap[mTLSType].DataID,
	//			Data:       mTLSsSchemaMap[mTLSType].Data,
	//			Type:       mTLSType,
	//			EnableID:   proxySetting.ID,
	//			Enable:     true,
	//		}
	//	}
	//	if proxySetting.Value == "false" {
	//		mTLSsSchemaMap[mTLSType] = models.FileSchema{
	//			FilenameID: mTLSsSchemaMap[mTLSType].FilenameID,
	//			Filename:   mTLSsSchemaMap[mTLSType].Filename,
	//			DataID:     mTLSsSchemaMap[mTLSType].DataID,
	//			Data:       mTLSsSchemaMap[mTLSType].Data,
	//			Type:       mTLSType,
	//			EnableID:   proxySetting.ID,
	//			Enable:     false,
	//		}
	//	}
	//case mtlsServerData:
	//	mTLSsSchemaMap[mTLSType] = models.FileSchema{
	//		FilenameID: mTLSsSchemaMap[mTLSType].FilenameID,
	//		Filename:   mTLSsSchemaMap[mTLSType].Filename,
	//		DataID:     proxySetting.ID,
	//		Data:       proxySetting.Value,
	//		Type:       mTLSType,
	//		EnableID:   mTLSsSchemaMap[mTLSType].EnableID,
	//		Enable:     mTLSsSchemaMap[mTLSType].Enable,
	//	}
	//case mtlsServerFileName:
	//	mTLSsSchemaMap[mTLSType] = models.FileSchema{
	//		FilenameID: proxySetting.ID,
	//		Filename:   proxySetting.Value,
	//		DataID:     mTLSsSchemaMap[mTLSType].DataID,
	//		Data:       mTLSsSchemaMap[mTLSType].Data,
	//		Type:       mTLSType,
	//		EnableID:   mTLSsSchemaMap[mTLSType].EnableID,
	//		Enable:     mTLSsSchemaMap[mTLSType].Enable,
	//	}

	//proxySettingsSchemaMap, err := utils.UnmarshalAs[[]map[string]any](asset.ProxySettings)
	//if err != nil {
	//	return fmt.Errorf("failed to convert proxy settings to slice of maps. Error: %+v", err)
	//}

	for _, mTLSscehma := range mTLSsSchemaMap {
		mTLS, err := utils.UnmarshalAs[map[string]any](mTLSscehma)
		if err != nil {
			return fmt.Errorf("failed to convert mTLS to map. Error: %+v", err)
		}

		mTLSsMap = append(mTLSsMap, mTLS)
	}

	d.Set("proxy_setting", proxySettingsSchemaMap)
	d.Set("mtls", mTLSsMap)

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
