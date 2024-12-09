package webuserresponse

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-user-response"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetWebUserResponseBehavior(ctx context.Context, c *api.Client, id string) (models.WebUserResponseBehavior, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
			{
				getWebUserResponseBehavior(id: "`+id+`") {
					id
					name
					visibility
					mode
					messageTitle
					messageBody
					httpResponseCode
					redirectURL
					xEventId
				}
			}
		`, "getWebUserResponseBehavior")

	if err != nil {
		return models.WebUserResponseBehavior{}, fmt.Errorf("failed to get WebUserResponseBehavior: %w", err)
	}

	behavior, err := utils.UnmarshalAs[models.WebUserResponseBehavior](res)
	if err != nil {
		return models.WebUserResponseBehavior{}, fmt.Errorf("failed to convert response to WebUserResponseBehavior struct. Error: %w", err)
	}

	return behavior, nil
}

func ReadWebUserResponseBehaviorToResourceData(behavior models.WebUserResponseBehavior, d *schema.ResourceData) error {
	d.SetId(behavior.ID)
	d.Set("name", behavior.Name)
	d.Set("visibility", behavior.Visibility)
	d.Set("mode", behavior.Mode)
	d.Set("message_title", behavior.MessageTitle)
	d.Set("message_body", behavior.MessageBody)
	d.Set("http_response_code", behavior.HTTPResponseCode)
	d.Set("redirect_url", behavior.RedirectURL)
	d.Set("x_event_id", behavior.XEventID)

	return nil
}
