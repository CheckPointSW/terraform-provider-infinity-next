package exceptions

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteExceptionBehavior(c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(`
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
