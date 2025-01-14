package resources

import (
	"context"
	webAPIModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	webAppModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	webapiasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-api-asset"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	trustedsources "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/trusted-sources"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTrustedSources() *schema.Resource {
	validateVisibility := validation.ToDiagFunc(
		validation.StringInSlice([]string{visibilityShared, visibilityLocal}, false))
	return &schema.Resource{
		Description: "Trusted sources that serve as a baseline for comparison for benign behavior, " +
			"and how many users or addresses must exhibit similar activity for it to be considered bengin by the learning model",

		CreateContext: resourceTrustedSourcesCreate,
		ReadContext:   resourceTrustedSourcesRead,
		UpdateContext: resourceTrustedSourcesUpdate,
		DeleteContext: resourceTrustedSourcesDelete,
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
				Description:      "The visibility of the resource - Shared or Local",
				Default:          "Shared",
				Optional:         true,
				ValidateDiagFunc: validateVisibility,
			},
			"min_num_of_sources": {
				Type:        schema.TypeInt,
				Description: "Minimum number of users or addresses that must exhibit similar activity for the behavior to be considered benign",
				Required:    true,
			},
			"sources_identifiers": {
				Type:        schema.TypeSet,
				Description: "The trusted sources identifier values",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sources_identifiers_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTrustedSourcesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	createInput, err := trustedsources.CreateTrustedSourceBehaviorInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform TrustedSourceBehavior Create", err, diags)
	}

	behavior, err := trustedsources.NewTrustedSourceBehavior(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform TrustedSourceBehavior Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following TrustedSourceBehavior Create", err, diags)
	}

	if err := trustedsources.ReadTrustedSourceBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform TrustedSourceBehavior Read after Create", err, diags)
	}

	return diags
}

func resourceTrustedSourcesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	behavior, err := trustedsources.GetTrustedSourceBehavior(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform TrustedSourceBehavior Get before read", err, diags)
	}

	if err := trustedsources.ReadTrustedSourceBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform TrustedSourceBehavior read to state file", err, diags)
	}

	return diags
}

func resourceTrustedSourcesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	updateInput, err := trustedsources.UpdateTrustedSourceBehaviorInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform TrustedSourceBehavior Update", err, diags)
	}

	result, err := trustedsources.UpdateTrustedSourceBehavior(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform TrustedSourceBehavior Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following TrustedSourceBehavior Update", err, diags)
	}

	behavior, err := trustedsources.GetTrustedSourceBehavior(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform TrustedSourceBehavior Get before read after update", err, diags)
	}

	if err := trustedsources.ReadTrustedSourceBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform TrustedSourceBehavior read to state file after update", err, diags)
	}

	return diags
}

func resourceTrustedSourcesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	usedBy, err := trustedsources.UsedByTrustedSourceBehavior(ctx, c, d.Id())
	if err != nil {
		return utils.DiagError("unable to perform TrustedSourceBehavior Delete", err, diags)
	}

	if usedBy != nil {
		for _, usedByResource := range usedBy {
			if usedByResource.ObjectStatus == "Deleted" {
				continue
			}

			switch usedByResource.SubType {
			case "WebAPI":
				objectToUpdate, err := webapiasset.GetWebAPIAsset(ctx, c, usedByResource.ID)
				if err != nil {
					return utils.DiagError("unable to perform WebAPIAsset Read", err, diags)
				}

				webAPIAsset := webAPIModels.UpdateWebAPIAssetInput{
					RemovePracticeWrappers: []string{d.Id()},
				}

				updated, err := webapiasset.UpdateWebAPIAsset(ctx, c, objectToUpdate.ID, webAPIAsset)
				if err != nil || !updated {
					if _, discardErr := c.DiscardChanges(); discardErr != nil {
						diags = utils.DiagError("failed to discard changes", discardErr, diags)
					}

					return utils.DiagError("unable to perform TrustedSourceBehavior Delete", err, diags)
				}

			case "WebApplication":
				objectToUpdate, err := webappasset.GetWebApplicationAsset(ctx, c, usedByResource.ID)
				if err != nil {
					return utils.DiagError("unable to perform WebApplicationAsset Read", err, diags)
				}

				webAppAsset := webAppModels.UpdateWebApplicationAssetInput{
					RemovePracticeWrappers: []string{d.Id()},
				}

				updated, err := webappasset.UpdateWebApplicationAsset(ctx, c, objectToUpdate.ID, webAppAsset)
				if err != nil || !updated {
					if _, discardErr := c.DiscardChanges(); discardErr != nil {
						diags = utils.DiagError("failed to discard changes", discardErr, diags)
					}

					return utils.DiagError("unable to perform TrustedSourceBehavior Delete", err, diags)
				}

			default:
				return utils.DiagError("unable to perform TrustedSourceBehavior Delete", err, diags)
			}
		}
	}

	result, err := trustedsources.DeleteTrustedSourceBehavior(ctx, c, d.Id())
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform TrustedSourceBehavior Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following TrustedSourceBehavior Delete", err, diags)
	}

	d.SetId("")

	return diags
}
