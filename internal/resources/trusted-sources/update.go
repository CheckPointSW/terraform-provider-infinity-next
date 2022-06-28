package trustedsources

import (
	"context"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/trusted-sources"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateTrustedSourceBehaviorInputFromResourceData(d *schema.ResourceData) (models.UpdateTrustedSourceBehaviorInput, error) {
	var res models.UpdateTrustedSourceBehaviorInput
	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		res.Name = newName
	}

	if _, newMinNumberOfSources, hasChange := utils.MustGetChange[int](d, "min_num_of_sources"); hasChange {
		res.NumOfSources = newMinNumberOfSources
	}

	if oldSources, newSources, hasChange := utils.GetChangeWithParse(d, "sources_identifiers", utils.MustSchemaCollectionToSlice[string]); hasChange {
		added, removed := utils.SlicesDiff(oldSources, newSources)
		oldSourcesIDs := utils.MustResourceDataCollectionToSlice[string](d, "sources_identifiers_ids")
		oldSourcesToIDsMap := make(map[string]string)
		for _, sourceID := range oldSourcesIDs {
			sourceAndID := strings.Split(sourceID, models.TrustedSourceIDSeparator)
			oldSourcesToIDsMap[sourceAndID[0]] = sourceAndID[1]
		}

		var idsToRemove []string
		for _, sourceToRemove := range removed {
			if idToRemove, ok := oldSourcesToIDsMap[sourceToRemove]; ok {
				idsToRemove = append(idsToRemove, idToRemove)
			}
		}

		res.AddSourcesIdentifiers = added
		res.RemoveSourcesIdentifiersIDs = idsToRemove
	}

	return res, nil
}

func UpdateTrustedSourceBehavior(ctx context.Context, c *api.Client, id string, input models.UpdateTrustedSourceBehaviorInput) (bool, error) {
	vars := map[string]any{"behaviorInput": input, "id": id}

	res, err := c.MakeGraphQLRequest(ctx, `
		mutation updateTrustedSourceBehavior($behaviorInput: TrustedSourceBehaviorUpdateInput, $id: ID!)
		{
			updateTrustedSourceBehavior(behaviorInput: $behaviorInput, id: $id)
		}
	`, "updateTrustedSourceBehavior", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateTrustedSourceBehavior response %#v should be of type bool", res)
	}

	return isUpdated, err
}
