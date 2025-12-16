package webuserresponse

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-user-response"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateWebUserResponseBehaviorInputFromResourceData(d *schema.ResourceData) (models.UpdateWebUserResponseBehaviorInput, error) {
	var res models.UpdateWebUserResponseBehaviorInput

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		res.Name = newName
	}

	if _, newVisibility, hasChange := utils.MustGetChange[string](d, "visibility"); hasChange {
		res.Visibility = newVisibility
	}

	if _, newMode, hasChange := utils.MustGetChange[string](d, "mode"); hasChange {
		res.Mode = newMode
	}

	if _, newVal, hasChange := utils.MustGetChange[string](d, "message_title"); hasChange {
		res.MessageTitle = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[string](d, "message_body"); hasChange {
		res.MessageBody = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[int](d, "http_response_code"); hasChange {
		res.HTTPResponseCode = &newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[string](d, "redirect_url"); hasChange {
		res.RedirectURL = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "x_event_id"); hasChange {
		res.XEventID = &newVal
	}

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
