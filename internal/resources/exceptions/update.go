package exceptions

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/exceptions"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseSchemaExceptions(exceptionsFromResourceData any) []models.ExceptionObjectInput {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](exceptionsFromResourceData), mapToExceptionObjectInput)
}

func UpdateExceptionBehaviorInputFromResourceData(d *schema.ResourceData) (models.UpdateExceptionBehaviorInput, error) {
	var res models.UpdateExceptionBehaviorInput
	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		res.Name = newName
	}

	if oldExceptions, newExceptions, hasChange := utils.GetChangeWithParse(d, "exception", parseSchemaExceptions); hasChange {
		exceptionsToAdd, exceptionsToRemove := utils.SlicesDiff(oldExceptions, newExceptions)
		res.AddExceptions = utils.Map(exceptionsToAdd, utils.MustUnmarshalAs[models.AddExceptionObjectInput, models.ExceptionObjectInput])
		res.RemoveExceptions = utils.Map(exceptionsToRemove, func(toRemove models.ExceptionObjectInput) string { return toRemove.ID })
	}

	return res, nil
}

func UpdateExceptionBehavior(ctx context.Context, c *api.Client, id string, input models.UpdateExceptionBehaviorInput) (bool, error) {
	vars := map[string]any{"behaviorInput": input, "id": id}
	res, err := c.MakeGraphQLRequest(ctx, `
		mutation updateExceptionBehavior($behaviorInput: ExceptionBehaviorUpdateInput, $id: ID!)
		{
			updateExceptionBehavior(behaviorInput: $behaviorInput, id: $id)
		}
	`, "updateExceptionBehavior", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateExceptionBehavior response %#v should be of type bool", res)
	}

	return isUpdated, err
}
