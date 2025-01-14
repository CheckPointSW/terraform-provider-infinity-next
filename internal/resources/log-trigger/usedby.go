package logtrigger

import (
	"context"
	"fmt"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/log-trigger"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

func UsedByLogTrigger(ctx context.Context, c *api.Client, id string) (models.TriggersUsedBy, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			query usedByTrigger($id: ID!){
				usedByTrigger(id: "`+id+`")
			}
		`, "usedByTrigger")

	if err != nil {
		return nil, err
	}

	usedBy, err := utils.UnmarshalAs[models.TriggersUsedBy](res)
	if err != nil {
		return models.TriggersUsedBy{}, fmt.Errorf("failed to unmarshal usedByTrigger response: %w", err)
	}

	return usedBy, nil
}
