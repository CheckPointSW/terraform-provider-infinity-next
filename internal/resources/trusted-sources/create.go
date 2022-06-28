package trustedsources

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/trusted-sources"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateTrustedSourceBehaviorInputFromResourceData(d *schema.ResourceData) (models.CreateTrustedSourceBehaviorInput, error) {
	var input models.CreateTrustedSourceBehaviorInput

	input.Name = d.Get("name").(string)
	input.Visibility = "Shared"
	input.NumOfSources = d.Get("min_num_of_sources").(int)
	input.SourcesIdentifiers = utils.MustResourceDataCollectionToSlice[string](d, "sources_identifiers")

	return input, nil
}

func NewTrustedSourceBehavior(ctx context.Context, c *api.Client, input models.CreateTrustedSourceBehaviorInput) (models.TrustedSourceBehavior, error) {
	vars := map[string]any{"ownerId": nil, "practiceId": nil, "behaviorInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newTrustedSourceBehavior($ownerId: ID, $practiceId: ID, $behaviorInput: TrustedSourceBehaviorInput)
					{
						newTrustedSourceBehavior(ownerId: $ownerId, practiceId: $practiceId, behaviorInput: $behaviorInput) {
							id
							name
							behaviorType
							numOfSources
							sourcesIdentifiers {
								id
								source
							}
						}
					}
				`, "newTrustedSourceBehavior", vars)

	if err != nil {
		return models.TrustedSourceBehavior{}, fmt.Errorf("failed to create new TrustedSourceBehavior: %w", err)
	}

	behavior, err := utils.UnmarshalAs[models.TrustedSourceBehavior](res)
	if err != nil {
		return models.TrustedSourceBehavior{}, fmt.Errorf("failed to convert response to CreateTrustedSourceBehaviorResponse struct. Error: %w", err)
	}

	return behavior, err
}
