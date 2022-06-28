package embeddedprofile

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
)

func DeleteEmbeddedProfile(ctx context.Context, c *api.Client, id string) (bool, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			mutation deleteProfile {
				deleteProfile(id: "`+id+`")
			}
		`, "deleteProfile")

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid deleteProfile response %#v should be of type bool", res)
	}

	return value, err
}
