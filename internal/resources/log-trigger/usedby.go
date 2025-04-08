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
			{
				triggerUsedBy(id: "`+id+`") {
					container
					practices
				}
			}
		`, "triggerUsedBy")

	if err != nil {
		return nil, err
	}

	usedBy, err := utils.UnmarshalAs[models.TriggersUsedBy](res)
	if err != nil {
		return models.TriggersUsedBy{}, fmt.Errorf("failed to unmarshal triggerUsedBy response: %w", err)
	}

	return usedBy, nil
}

func UpdatePracticeTriggers(ctx context.Context, c *api.Client, triggerID string, practiceID string, containerID string) (bool, error) {
	vars := map[string]any{"addTriggers": nil, "removeTriggers": []string{triggerID}, "practiceId": practiceID, "containerId": containerID}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation updatePracticeTriggers($addTriggers: [ID], $removeTriggers: [ID], $practiceId: ID!, $containerId: ID!)
					{
						updatePracticeTriggers(addTriggers: $addTriggers, removeTriggers: $removeTriggers, practiceId: $practiceId, containerId: $containerId)
					}
				`, "updatePracticeTriggers", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid UpdatePracticeTriggers %#v should be of type bool", res)
	}

	return isUpdated, err
}
