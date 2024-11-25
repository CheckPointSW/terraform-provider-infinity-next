package webapiasset

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
