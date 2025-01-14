package webapppractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

func UsedByWebApplicationPractice(ctx context.Context, c *api.Client, id string) (models.DisplayObjects, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			{
				practiceUsedBy(id: "`+id+`") {
					id
					name
					type
					subType
					objectStatus
				}
			}
		`, "practiceUsedBy")

	if err != nil {
		return nil, err
	}

	usedBy, err := utils.UnmarshalAs[models.DisplayObjects](res)
	if err != nil {
		return models.DisplayObjects{}, fmt.Errorf("failed to unmarshal practiceUsedBy response: %w", err)
	}

	return usedBy, nil
}
