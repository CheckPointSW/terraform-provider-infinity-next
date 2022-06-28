package webappasset

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	d.Set("trusted_sources", asset.Behaviors.ToSchema())
	d.Set("profiles", asset.Profiles.ToSchema())

	proxySettingsSchemaMap, err := utils.UnmarshalAs[[]map[string]any](asset.ProxySettings)
	if err != nil {
		return fmt.Errorf("failed to convert proxy settings to slice of maps. Error: %+v", err)
	}

	d.Set("proxy_setting", proxySettingsSchemaMap)

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
					behaviors {
						id
					}
				}
				profiles {
					id
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
			}
		}
	`, "getWebApplicationAsset")

	if err != nil {
		return models.WebApplicationAsset{}, fmt.Errorf("failed to get WebApplicationAsset: %w", err)
	}

	if res == nil {
		return models.WebApplicationAsset{}, api.ErrorNotFound
	}

	// res, ok := mapRes.(map[string]any)
	// if !ok {
	// 	return models.WebApplicationAsset{}, fmt.Errorf("invalid response field getWebApplicationAsset value, should be of type map[string]any but got %#v", mapRes)
	// }

	asset, err := utils.UnmarshalAs[models.WebApplicationAsset](res)
	if err != nil {
		return models.WebApplicationAsset{}, fmt.Errorf("Failed to convert graphQL response to WebApplicationAsset struct. Error: %w", err)
	}

	return asset, nil
}
