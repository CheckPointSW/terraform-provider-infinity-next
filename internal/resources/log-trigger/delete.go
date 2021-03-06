package logtrigger

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteLogTrigger(ctx context.Context, c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		mutation deleteTrigger {
			deleteTrigger(id: "`+id+`")
		}
	`, "deleteTrigger")

	if err != nil {
		return false, err
	}

	isDeleted, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid deleteTrigger response %#v should be of type bool", res)
	}

	return isDeleted, err
}
