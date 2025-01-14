package webapipractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

func UsedByWebAPIPractice(ctx context.Context, c *api.Client, id string) (models.DisplayObjects, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			query usedByPractice {
				usedByPractice(id: "`+id+`")
			}
		`, "usedByPractice")

	if err != nil {
		return nil, err
	}

	usedBy, err := utils.UnmarshalAs[models.DisplayObjects](res)
	if err != nil {
		return models.DisplayObjects{}, fmt.Errorf("failed to unmarshal usedByPractice response: %w", err)
	}

	return usedBy, nil
}
