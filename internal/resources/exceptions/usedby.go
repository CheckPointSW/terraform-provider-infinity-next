package exceptions

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/exceptions"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

func UsedByExceptionBehavior(ctx context.Context, c *api.Client, id string) (models.DisplayObjects, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			query usedByBehavior($id: ID!){
				usedByBehavior(id: "`+id+`")
			}
		`, "usedByBehavior")

	if err != nil {
		return nil, err
	}

	usedBy, err := utils.UnmarshalAs[models.DisplayObjects](res)
	if err != nil {
		return models.DisplayObjects{}, fmt.Errorf("failed to unmarshal usedByBehavior response: %w", err)
	}

	return usedBy, nil
}
