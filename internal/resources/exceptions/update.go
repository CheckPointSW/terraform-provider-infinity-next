package exceptions

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/exceptions"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseSchemaExceptions(exceptionsFromResourceData any) models.ExceptionObjectInputs {
	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](exceptionsFromResourceData), mapToExceptionObjectInput)
}

func UpdateExceptionBehaviorInputFromResourceData(d *schema.ResourceData) (models.UpdateExceptionBehaviorInput, error) {
	var res models.UpdateExceptionBehaviorInput
	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		res.Name = newName
	}

	if _, newVisibility, hasChange := utils.MustGetChange[string](d, "visibility"); hasChange {
		res.Visibility = newVisibility
	}

	if oldExceptions, newExceptions, hasChange := utils.GetChangeWithParse(d, "exception", parseSchemaExceptions); hasChange {
		oldExceptionsIndicators := oldExceptions.ToIndicatorsMap()
		for _, newException := range newExceptions {
			// if key does not exist then this is a new Exception to add
			if _, ok := oldExceptionsIndicators[newException.ID]; !ok {
				res.AddExceptions = append(res.AddExceptions, models.AddExceptionObjectInput{
					Match:   newException.Match,
					Actions: newException.Actions,
					Comment: newException.Comment,
				})

			}

			// we know the key exist
			// if the value is different - update the Exception
			oldException := oldExceptionsIndicators[newException.Match]
			actionsToAdd, actionsToRemove := utils.SlicesDiff(oldException.Actions, newException.Actions)
			res.UpdateExceptions = append(res.UpdateExceptions, models.ExceptionObjectActionUpdate{
				ID:            newException.ID,
				Match:         newException.Match,
				AddActions:    actionsToAdd,
				RemoveActions: actionsToRemove,
				UpdateActions: models.UpdateExceptionsObjectInputs{},
				Comment:       newException.Comment,
			})
		}

		newExceptionsIndicators := newExceptions.ToIndicatorsMap()
		for _, oldException := range oldExceptions {
			if _, ok := newExceptionsIndicators[oldException.ID]; !ok {
				res.RemoveExceptions = append(res.RemoveExceptions, oldException.ID)
			}
		}
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
