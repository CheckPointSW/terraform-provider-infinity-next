package webuserresponse

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteWebUserResponseBehavior(ctx context.Context, c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			mutation deleteBehavior {
				deleteBehavior(id: "`+id+`")
			}
		`, "deleteBehavior")

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid deleteBehavior response %#v should be of type bool", res)
	}

	return value, err
}
