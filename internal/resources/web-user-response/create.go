package webuserresponse

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-user-response"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateWebUserResponseBehaviorInputFromResourceData(d *schema.ResourceData) (models.CreateWebUserResponseBehaviorInput, error) {
	var input models.CreateWebUserResponseBehaviorInput

	input.Name = d.Get("name").(string)
	input.Visibility = "Shared"
	input.Mode = d.Get("mode").(string)
	input.MessageTitle = d.Get("message_title").(string)
	input.MessageBody = d.Get("message_body").(string)
	input.HTTPResponseCode = d.Get("http_response_code").(int)
	input.RedirectURL = d.Get("redirect_url").(string)
	input.XEventID = d.Get("x_event_id").(bool)

	return input, nil
}

func NewWebUserResponseBehavior(ctx context.Context, c *api.Client, input models.CreateWebUserResponseBehaviorInput) (models.WebUserResponseBehavior, error) {
	vars := map[string]any{"ownerId": nil, "practiceId": nil, "behaviorInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newWebUserResponseBehavior($ownerId: ID, $practiceId: ID, $behaviorInput: WebUserResponseBehaviorInput)
					{
						newWebUserResponseBehavior(ownerId: $ownerId, practiceId: $practiceId, behaviorInput: $behaviorInput) {
							id
							name
							mode
							messageTitle
							messageBody
							httpResponseCode
							redirectURL
							xEventId
						}
					}
				`, "newWebUserResponseBehavior", vars)

	if err != nil {
		return models.WebUserResponseBehavior{}, fmt.Errorf("failed to create new WebUserResponseBehavior: %w", err)
	}

	behavior, err := utils.UnmarshalAs[models.WebUserResponseBehavior](res)
	if err != nil {
		return models.WebUserResponseBehavior{}, fmt.Errorf("failed to convert response to CreateWebUserResponseBehaviorResponse struct. Error: %w", err)
	}

	return behavior, err
}
