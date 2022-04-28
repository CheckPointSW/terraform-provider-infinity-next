package exceptions

import (
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateExceptionBehaviorInputFromResourceData(d *schema.ResourceData) (UpdateExceptionBehaviorInput, error) {
	var res UpdateExceptionBehaviorInput
	res.ID = d.Id()
	res.BehaviorInput.Name = d.Get("name").(string)

	old, new := d.GetChange("exception")
	oldExceptionsIface := old.([]any)
	var oldExceptions []SchemaExceptionObject
	bExceptions, err := json.Marshal(oldExceptionsIface)
	if err != nil {
		return UpdateExceptionBehaviorInput{}, err
	}

	if err := json.Unmarshal(bExceptions, &oldExceptions); err != nil {
		return UpdateExceptionBehaviorInput{}, err
	}

	newExceptionsIface := new.([]any)
	var newExceptions []SchemaExceptionObject
	bExceptions, err = json.Marshal(newExceptionsIface)
	if err != nil {
		return UpdateExceptionBehaviorInput{}, err
	}

	if err := json.Unmarshal(bExceptions, &newExceptions); err != nil {
		return UpdateExceptionBehaviorInput{}, err
	}

	res.BehaviorInput.RemoveExceptions = make([]string, len(oldExceptions))
	for i, exception := range oldExceptions {
		res.BehaviorInput.RemoveExceptions[i] = exception.ID
	}

	res.BehaviorInput.AddExceptions = make([]ExceptionObjectInput, 0, len(newExceptions))
	for _, exception := range newExceptions {
		var exceptionInput ExceptionObjectInput
		action := NewAction(exception.Action)
		bAction, err := json.Marshal(action)
		if err != nil {
			return UpdateExceptionBehaviorInput{}, fmt.Errorf("failed to marshal exception action %#v: %w", action, err)
		}

		exceptionInput.Actions = []string{string(bAction)}
		exceptionInput.Comment = exception.Comment
		match := AndMatchFromMap(exception.Match)
		bMatch, err := json.Marshal(match)
		if err != nil {
			return UpdateExceptionBehaviorInput{}, fmt.Errorf("failed to marshal exception match %#v: %w", match, err)
		}

		exceptionInput.Match = string(bMatch)
		res.BehaviorInput.AddExceptions = append(res.BehaviorInput.AddExceptions, exceptionInput)
	}

	return res, nil
}

func UpdateExceptionBehavior(c *api.Client, input UpdateExceptionBehaviorInput) (bool, error) {
	vars := map[string]any{"behaviorInput": input.BehaviorInput, "id": input.ID}
	res, err := c.MakeGraphQLRequest(`
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
