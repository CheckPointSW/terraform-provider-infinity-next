package trustedsources

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/trusted-sources"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

func UsedByTrustedSourceBehavior(ctx context.Context, c *api.Client, id string) (models.DisplayObjects, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			query behaviorUsedBy($id: ID!){
				behaviorUsedBy(id: "`+id+`")
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
