package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	trustedsources "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/trusted-sources"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTrustedSources() *schema.Resource {
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

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the resource, also acts as it's unique ID",
				Required:    true,
			},
			"min_num_of_sources": {
				Type:        schema.TypeInt,
				Description: "Minimum number of users or addresses that must exhibit similar activity for the behavior to be considered benign",
				Required:    true,
			},
			"sources_identifiers": {
				Type:        schema.TypeList,
				Description: "The trusted sources identifier values",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sources_identifiers_ids": {
				Type:     schema.TypeList,
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

	behavior, err := trustedsources.NewTrustedSourceBehavior(c, createInput)
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

		return diag.FromErr(err)
	}

	return diags
}

func resourceTrustedSourcesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	behavior, err := trustedsources.GetTrustedSourceBehavior(c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	if err := trustedsources.ReadTrustedSourceBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
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

	result, err := trustedsources.UpdateTrustedSourceBehavior(c, d.Id(), updateInput)
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

	behavior, err := trustedsources.GetTrustedSourceBehavior(c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	if err := trustedsources.ReadTrustedSourceBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	return diags
}

func resourceTrustedSourcesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	result, err := trustedsources.DeleteTrustedSourceBehavior(c, d.Id())
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
