package webapiasset

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
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

	for _, proxySetting := range asset.ProxySettings {
		mTLSType := proxySettingKeyTomTLSType(proxySetting.Key)
		if mTLSType != "" {
			if _, ok := mTLSsSchemaMap[mTLSType]; !ok {
				mTLSsSchemaMap[mTLSType] = models.MTLSSchema{}
			}

			switch proxySetting.Key {
			case mtlsClientEnable, mtlsServerEnable:
				enable := false
				if proxySetting.Value == "true" {
					enable = true
				}

				mTLSsSchemaMap[mTLSType] = models.MTLSSchema{
					FilenameID:      mTLSsSchemaMap[mTLSType].FilenameID,
					Filename:        mTLSsSchemaMap[mTLSType].Filename,
					CertificateType: mTLSsSchemaMap[mTLSType].CertificateType,
					DataID:          mTLSsSchemaMap[mTLSType].DataID,
					Data:            mTLSsSchemaMap[mTLSType].Data,
					Type:            mTLSType,
					EnableID:        proxySetting.ID,
					Enable:          enable,
				}
			case mtlsClientData, mtlsServerData:
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
					fileExtensionsByType = models.MimeTypeToFileExtension(mimeType, true)
				}

				mTLSsSchemaMap[mTLSType] = models.MTLSSchema{
					FilenameID:      mTLSsSchemaMap[mTLSType].FilenameID,
					Filename:        mTLSsSchemaMap[mTLSType].Filename,
					CertificateType: fileExtensionsByType,
					DataID:          proxySetting.ID,
					Data:            decodedData,
					Type:            mTLSType,
					EnableID:        mTLSsSchemaMap[mTLSType].EnableID,
					Enable:          mTLSsSchemaMap[mTLSType].Enable,
				}
			case mtlsClientFileName, mtlsServerFileName:
				mTLSsSchemaMap[mTLSType] = models.MTLSSchema{
					FilenameID:      proxySetting.ID,
					Filename:        proxySetting.Value,
					CertificateType: mTLSsSchemaMap[mTLSType].CertificateType,
					DataID:          mTLSsSchemaMap[mTLSType].DataID,
					Data:            mTLSsSchemaMap[mTLSType].Data,
					Type:            mTLSType,
					EnableID:        mTLSsSchemaMap[mTLSType].EnableID,
					Enable:          mTLSsSchemaMap[mTLSType].Enable,
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
