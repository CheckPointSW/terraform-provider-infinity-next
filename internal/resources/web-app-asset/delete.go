package webappasset

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteWebApplicationAsset(c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(`
			mutation deleteAsset {
				deleteAsset(id: "`+id+`")
			}
		`, "deleteAsset")

	if err != nil {
		return false, err
	}

	isDeleted, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid deleteAsset response %#v should be of type bool", res)
	}

	return isDeleted, err
}
