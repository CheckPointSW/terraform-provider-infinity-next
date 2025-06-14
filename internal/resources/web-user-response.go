package resources

import (
	"context"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webAPIAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	webAppAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-user-response"
	webapiasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-api-asset"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"
	webuserresponse "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-user-response"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceWebUserResponse() *schema.Resource {
	validateVisibility := validation.ToDiagFunc(validation.StringInSlice([]string{visibilityShared, visibilityLocal}, false))
	return &schema.Resource{
		Description: "Determine the response returned to the client who initiated a blocked traffic." +
			"The response can be a simple HTTP error code, an HTTP redirect message, or a Block page that a user can view in their browser.",

		CreateContext: resourceWebUserResponseCreate,
		ReadContext:   resourceWebUserResponseRead,
		UpdateContext: resourceWebUserResponseUpdate,
		DeleteContext: resourceWebUserResponseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			if diff.HasChange("sources_identifiers") {
				return diff.SetNewComputed("sources_identifiers_ids")
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the resource, also acts as its unique ID",
				Required:    true,
			},
			"visibility": {
				Type:             schema.TypeString,
				Description:      "The visibility of the web user response object",
				Optional:         true,
				Default:          "Shared",
				ValidateDiagFunc: validateVisibility,
			},
			"mode": {
				Type:             schema.TypeString,
				Description:      "The type of the web user response object",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BlockPage", "Redirect", "ResponseCodeOnly"}, false)),
			},
			"message_title": {
				Type:        schema.TypeString,
				Description: "The title of the web page to be shown to the user sending the malicious traffic",
				Optional:    true,
			},
			"message_body": {
				Type:        schema.TypeString,
				Description: "The body of the message to be shown to the user",
				Optional:    true,
			},
			"http_response_code": {
				Type:        schema.TypeInt,
				Description: "It is recommended to use a 403 (Forbidden) as a response code",
				Optional:    true,
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Description: "The client will be redirected to the provided URL where you can provide any customized web page",
				Optional:    true,
			},
			"x_event_id": {
				Type: schema.TypeBool,
				Description: "When selected the redirect message will include this header with a value that provides an internal reference ID " +
					"that will match a security log generated by the incident, if log triggers are configured",
				Optional: true,
			},
		},
	}
}

func resourceWebUserResponseCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	createInput, err := webuserresponse.CreateWebUserResponseBehaviorInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebUserResponseBehavior Create", err, diags)
	}

	behavior, err := webuserresponse.NewWebUserResponseBehavior(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebUserResponseBehavior Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebUserResponseBehavior Create", err, diags)
	}

	if err := webuserresponse.ReadWebUserResponseBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform WebUserResponseBehavior Read after Create", err, diags)
	}

	return diags
}

func resourceWebUserResponseRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	behavior, err := webuserresponse.GetWebUserResponseBehavior(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform WebUserResponseBehavior Get before read", err, diags)
	}

	if err := webuserresponse.ReadWebUserResponseBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform WebUserResponseBehavior read to state file", err, diags)
	}

	return diags
}

func resourceWebUserResponseUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	updateInput, err := webuserresponse.UpdateWebUserResponseBehaviorInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebUserResponseBehavior Update", err, diags)
	}

	result, err := webuserresponse.UpdateWebUserResponseBehavior(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebUserResponseBehavior Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebUserResponseBehavior Update", err, diags)
	}

	behavior, err := webuserresponse.GetWebUserResponseBehavior(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform WebUserResponseBehavior Get before read after update", err, diags)
	}

	if err := webuserresponse.ReadWebUserResponseBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform WebUserResponseBehavior read to state file after update", err, diags)
	}

	return diags
}

func resourceWebUserResponseDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	result, err := webuserresponse.DeleteWebUserResponseBehavior(ctx, c, d.Id())
	if err != nil || !result {
		// Check if the error is due to the web user response behavior being used by other resources
		if err != nil && strings.Contains(err.Error(), errorMsgPointedObjects) {
			// Get the resources that are using the web user response behavior
			usedBy, err2 := webuserresponse.UsedByWebUserResponse(ctx, c, d.Id())
			if err2 != nil {
				diags = utils.DiagError("unable to perform WebUserResponse UsedBy", err2, diags)
				return utils.DiagError("unable to perform WebUserResponse Delete", err, diags)
			}

			if usedBy != nil || len(usedBy) > 0 {
				// Remove the web user response behavior from the resources that are using it
				if err2 := handleWebUserResponseReferences(ctx, usedBy, c, d.Id()); err2 != nil {
					diags = err2
					return utils.DiagError("unable to perform WebUserResponse Delete", err, diags)
				}

				// Retry to delete the web user response behavior
				result, err := webuserresponse.DeleteWebUserResponseBehavior(ctx, c, d.Id())
				if err != nil || !result {
					if _, discardErr := c.DiscardChanges(); discardErr != nil {
						diags = utils.DiagError("failed to discard changes", discardErr, diags)
					}

					return utils.DiagError("unable to perform WebUserResponse Delete after updating references", err, diags)
				}

			}

		} else {
			if _, discardErr := c.DiscardChanges(); discardErr != nil {
				diags = utils.DiagError("failed to discard changes", discardErr, diags)
			}

			return utils.DiagError("unable to perform WebUserResponseBehavior Delete", err, diags)
		}
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebUserResponseBehavior Delete", err, diags)
	}

	d.SetId("")

	return diags
}

func handleWebUserResponseReferences(ctx context.Context, usedBy models.DisplayObjects, c *api.Client, behaviorID string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, usedByResource := range usedBy {
		if usedByResource.ObjectStatus == "Deleted" {
			continue
		}

		switch usedByResource.SubType {
		case "WebAPI":
			webAPIAsset := webAPIAssetModels.UpdateWebAPIAssetInput{
				RemoveBehaviors: []string{behaviorID},
			}

			updated, err := webapiasset.UpdateWebAPIAsset(ctx, c, usedByResource.ID, webAPIAsset)
			if err != nil || !updated {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("failed to perform UpdateWebAPIAsset to remove behavior", err, diags)
			}

		case "WebApplication":
			webAppAsset := webAppAssetModels.UpdateWebApplicationAssetInput{
				RemoveBehaviors: []string{behaviorID},
			}

			updated, err := webappasset.UpdateWebApplicationAsset(ctx, c, usedByResource.ID, webAppAsset)
			if err != nil || !updated {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("failed to perform UpdateWebApplicationAsset to remove behavior", err, diags)
			}

		default:
			return utils.DiagError("failed to update usedByResource", nil, diags)
		}

	}

	return nil
}
