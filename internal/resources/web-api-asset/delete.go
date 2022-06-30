package webapiasset

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteWebAPIAsset(ctx context.Context, c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			mutation deleteAsset {
				deleteAsset(id: "`+id+`")
			}
		`, "deleteAsset")

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid deleteAsset response %#v should be of type bool", res)
	}

	return value, err
}
