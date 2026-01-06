package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	publishenforce "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/publish-enforce"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// To ensure a single resource like this we use the same id for all resources
const publishEnforceSingletonID = "publish-enforce-singleton"

func ResourcePublishEnforce() *schema.Resource {
	return &schema.Resource{
		Description: "Publish and Enforce resource - triggers publish and/or enforce operations. " +
			"Works the same as running `inext publish` and `inext enforce` CLI commands.",

		CreateContext: resourcePublishEnforceCreateOrUpdate,
		ReadContext:   resourcePublishEnforceRead,
		UpdateContext: resourcePublishEnforceCreateOrUpdate,
		DeleteContext: resourcePublishEnforceDelete,

		Schema: map[string]*schema.Schema{
			"publish": {
				Type:        schema.TypeBool,
				Description: "When true, triggers a publish operation (same as `inext publish`)",
				Optional:    true,
				Default:     false,
			},
			"enforce": {
				Type:        schema.TypeBool,
				Description: "When true, triggers an enforce operation (same as `inext enforce`)",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

// resourcePublishEnforceCreateOrUpdate handles both create and update operations
// since they perform identical logic: trigger publish/enforce based on current values
func resourcePublishEnforceCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	publish := d.Get("publish").(bool)
	enforce := d.Get("enforce").(bool)

	// Use singleton ID - only one instance of this resource allowed
	d.SetId(publishEnforceSingletonID)

	// Execute publish if requested (same as `inext publish`)
	if publish {
		if err := publishenforce.ExecutePublish(ctx, c); err != nil {
			return utils.DiagError("failed to publish changes", err, diags)
		}
	}

	// Execute enforce if requested (same as `inext enforce`)
	if enforce {
		if err := publishenforce.ExecuteEnforce(ctx, c); err != nil {
			return utils.DiagError("failed to enforce policy", err, diags)
		}
	}

	return diags
}

func resourcePublishEnforceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// This is a trigger-only resource, nothing to read from the API
	// Just return current state
	return nil
}

func resourcePublishEnforceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
