package ratelimitpractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteRateLimitPractice(ctx context.Context, c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			mutation deletePractice {
				deletePractice(id: "`+id+`")
			}
		`, "deletePractice")

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid deletePractice response %#v should be of type bool", res)
	}

	return value, err
}
