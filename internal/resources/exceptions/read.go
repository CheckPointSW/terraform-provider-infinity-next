package exceptions

import (
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetExceptionBehavior(c *api.Client, id string) (ExceptionBehavior, error) {
	res, err := c.MakeGraphQLRequest(`
			{
				getExceptionBehavior(id: "`+id+`") {
					id
					name
					exceptions {
						id
						match
						actions {
							id
							action
						}
						comment
					}
				}
			}
		`, "getExceptionBehavior")

	if err != nil {
		return ExceptionBehavior{}, fmt.Errorf("failed to get ExceptionBehavior: %w", err)
	}

	behavior, err := utils.UnmarshalAs[ExceptionBehavior](res)
	if err != nil {
		return ExceptionBehavior{}, fmt.Errorf("failed to convert response to ExceptionBehavior struct. Error: %w", err)
	}

	return behavior, nil
}

func ReadExceptionBehaviorToResourceData(behavior ExceptionBehavior, d *schema.ResourceData) error {
	d.SetId(behavior.ID)
	d.Set("name", behavior.Name)

	exceptions := make([]SchemaExceptionBehavior, len(behavior.Exceptions))
	for i, exception := range behavior.Exceptions {
		var schemaException SchemaExceptionBehavior
		schemaException.ID = exception.ID
		if len(exception.Actions) > 0 {
			var action Action
			if err := json.Unmarshal([]byte(exception.Actions[0].Action), &action); err != nil {
				return fmt.Errorf("failed to unmarshal exception action %s to Action: %w", exception.Actions[0].Action, err)
			}

			schemaException.Action = action.Value
			schemaException.ActionID = exception.Actions[0].ID
		}

		var match Match
		if err := json.Unmarshal([]byte(exception.Match), &match); err != nil {
			return fmt.Errorf("failed to unmarshal exception match %s to Match: %w", exception.Match, err)
		}

		schemaException.Match = match.ToSchemaMap()
		schemaException.Comment = exception.Comment

		exceptions[i] = schemaException
	}

	if v, ok := d.GetOk("exception"); ok {
		var resourceDataExceptions []SchemaExceptionBehavior
		bExceptions, err := json.Marshal(v.([]any))
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bExceptions, &resourceDataExceptions); err != nil {
			return err
		}

		added, removed := ExceptionsDiff(resourceDataExceptions, exceptions)
		if len(added) == 0 && len(removed) == 0 {
			exceptions = resourceDataExceptions
		}
	}

	var exceptionsMaps []map[string]any
	bExceptions, err := json.Marshal(exceptions)
	if err != nil {
		return fmt.Errorf("failed to marshal exceptions %#v: %w", exceptions, err)
	}

	if err := json.Unmarshal(bExceptions, &exceptionsMaps); err != nil {
		return fmt.Errorf("failed to unmarshal exceptions %s to []map[string]any: %w", string(bExceptions), err)
	}

	d.Set("exception", exceptionsMaps)

	return nil
}
