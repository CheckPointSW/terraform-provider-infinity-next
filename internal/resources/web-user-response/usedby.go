package webuserresponse

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-user-response"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

func UsedByWebUserResponse(ctx context.Context, c *api.Client, id string) (models.DisplayObjects, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			{
				behaviorUsedBy(id: "`+id+`") {
					id
					name
					type
					subType
					objectStatus
				}
			}
		`, "behaviorUsedBy")

	if err != nil {
		return nil, err
	}

	usedBy, err := utils.UnmarshalAs[models.DisplayObjects](res)
	if err != nil {
		return models.DisplayObjects{}, fmt.Errorf("failed to unmarshal behaviorUsedBy response: %w", err)
	}

	return usedBy, nil
}
