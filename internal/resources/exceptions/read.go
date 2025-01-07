package exceptions

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/exceptions"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetExceptionBehavior(ctx context.Context, c *api.Client, id string) (models.ExceptionBehavior, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			{
				getExceptionBehavior(id: "`+id+`") {
					id
					name
					visibility
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
		return models.ExceptionBehavior{}, fmt.Errorf("failed to get ExceptionBehavior: %w", err)
	}

	behavior, err := utils.UnmarshalAs[models.ExceptionBehavior](res)
	if err != nil {
		return models.ExceptionBehavior{}, fmt.Errorf("failed to convert response to ExceptionBehavior struct. Error: %w", err)
	}

	return behavior, nil
}

func ReadExceptionBehaviorToResourceData(behavior models.ExceptionBehavior, d *schema.ResourceData) error {
	d.SetId(behavior.ID)
	d.Set("name", behavior.Name)
	d.Set("visibility", behavior.Visibility)
	schemaExceptions := behavior.Exceptions.ToSchema()
	schemaExceptionsMap, err := utils.UnmarshalAs[[]map[string]any](schemaExceptions)
	if err != nil {
		return fmt.Errorf("failed to convert exceptions to slice of map: %w", err)
	}

	d.Set("exception", schemaExceptionsMap)

	return nil
}
