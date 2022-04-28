package exceptions

import (
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateExceptionBehaviorInputFromResourceData(d *schema.ResourceData) (CreateExceptionBehaviorInput, error) {
	var res CreateExceptionBehaviorInput

	res.Name = d.Get("name").(string)
	res.Visibility = "Shared"

	if v, ok := d.GetOk("exception"); ok {
		exceptions := v.([]any)
		res.Exceptions = make([]ExceptionObjectInput, 0, len(exceptions))
		for _, exceptionIface := range exceptions {
			var schemaException SchemaExceptionObject
			bValue, err := json.Marshal(exceptionIface)
			if err != nil {
				return CreateExceptionBehaviorInput{}, fmt.Errorf("failed to marshal exception %#v: %w", exceptionIface, err)
			}

			if err := json.Unmarshal(bValue, &schemaException); err != nil {
				return CreateExceptionBehaviorInput{}, fmt.Errorf("failed to unmarshal exception %s to SchemaExceptionObject: %w", string(bValue), err)
			}

			var exception ExceptionObjectInput
			action := NewAction(schemaException.Action)
			bAction, err := json.Marshal(action)
			if err != nil {
				return CreateExceptionBehaviorInput{}, fmt.Errorf("failed to marshal exception action %#v: %w", action, err)
			}

			exception.Actions = []string{string(bAction)}
			exception.Comment = schemaException.Comment
			match := AndMatchFromMap(schemaException.Match)
			bMatch, err := json.Marshal(match)
			if err != nil {
				return CreateExceptionBehaviorInput{}, fmt.Errorf("failed to marshal exception match %#v: %w", match, err)
			}

			exception.Match = string(bMatch)
			res.Exceptions = append(res.Exceptions, exception)
		}
	}

	return res, nil
}

func NewExceptionBehavior(c *api.Client, input CreateExceptionBehaviorInput) (CreateExceptionBehaviorResponse, error) {
	vars := map[string]any{"ownerId": nil, "practiceId": nil, "behaviorInput": input}
	res, err := c.MakeGraphQLRequest(`
				mutation newExceptionBehavior($ownerId: ID, $practiceId: ID, $behaviorInput: ExceptionBehaviorInput)
					{
						newExceptionBehavior(ownerId: $ownerId, practiceId: $practiceId, behaviorInput: $behaviorInput) {
							id
						}
					}
				`, "newExceptionBehavior", vars)

	if err != nil {
		return CreateExceptionBehaviorResponse{}, fmt.Errorf("failed to create new ExceptionBehavior: %w", err)
	}

	createRes, err := utils.UnmarshalAs[CreateExceptionBehaviorResponse](res)
	if err != nil {
		return CreateExceptionBehaviorResponse{}, fmt.Errorf("failed to convert response to CreateExceptionBehaviorResponse struct. Error: %w", err)
	}

	return createRes, err
}
