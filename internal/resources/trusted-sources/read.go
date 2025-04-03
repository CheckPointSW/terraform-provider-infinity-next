package trustedsources

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/trusted-sources"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetTrustedSourceBehavior(ctx context.Context, c *api.Client, id string) (models.TrustedSourceBehavior, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			{
				getTrustedSourceBehavior(id: "`+id+`") {
					id
					name
					visibility
					behaviorType
					numOfSources
					sourcesIdentifiers {
						id
						source
					}
				}
			}
		`, "getTrustedSourceBehavior")

	if err != nil {
		return models.TrustedSourceBehavior{}, fmt.Errorf("failed to get TrustedSourceBehavior: %w", err)
	}

	behavior, err := utils.UnmarshalAs[models.TrustedSourceBehavior](res)
	if err != nil {
		return models.TrustedSourceBehavior{}, fmt.Errorf("failed to convert response to TrustedSourceBehavior struct. Error: %w", err)
	}

	return behavior, nil
}

func ReadTrustedSourceBehaviorToResourceData(behavior models.TrustedSourceBehavior, d *schema.ResourceData) error {
	d.SetId(behavior.ID)
	d.Set("name", behavior.Name)
	d.Set("visibility", behavior.Visibility)
	d.Set("min_num_of_sources", behavior.NumOfSources)

	sourcesIdentifiers := make([]string, len(behavior.SourcesIdentifiers))
	sourcesIdentifiersIDs := make([]string, len(behavior.SourcesIdentifiers))
	for i, sourceIdentifier := range behavior.SourcesIdentifiers {
		sourcesIdentifiers[i] = sourceIdentifier.Source
		sourcesIdentifiersIDs[i] = sourceIdentifier.CreateSchemaID()
	}

	if len(sourcesIdentifiers) > 0 {
		d.Set("sources_identifiers", sourcesIdentifiers)
		d.Set("sources_identifiers_ids", sourcesIdentifiersIDs)
	}

	return nil
}
