package webuserresponse

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-user-response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateWebUserResponseBehaviorInputFromResourceData(d *schema.ResourceData) (models.UpdateWebUserResponseBehaviorInput, error) {
	var res models.UpdateWebUserResponseBehaviorInput
	res.Name = d.Get("name").(string)
	res.Visibility = d.Get("visibility").(string)
	res.Mode = d.Get("mode").(string)
	res.MessageTitle = d.Get("message_title").(string)
	res.MessageBody = d.Get("message_body").(string)
	res.HTTPResponseCode = d.Get("http_response_code").(int)
	res.RedirectURL = d.Get("redirect_url").(string)
	res.XEventID = d.Get("x_event_id").(bool)

	return res, nil
}

func UpdateWebUserResponseBehavior(ctx context.Context, c *api.Client, id string, input models.UpdateWebUserResponseBehaviorInput) (bool, error) {
	vars := map[string]any{"behaviorInput": input, "id": id}

	res, err := c.MakeGraphQLRequest(ctx, `
		mutation updateWebUserResponseBehavior($behaviorInput: WebUserResponseBehaviorUpdateInput, $id: ID!)
		{
			updateWebUserResponseBehavior(behaviorInput: $behaviorInput, id: $id)
		}
	`, "updateWebUserResponseBehavior", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateWebUserResponseBehavior response %#v should be of type bool", res)
	}

	return isUpdated, err
}
