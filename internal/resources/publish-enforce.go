package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	publishenforce "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/publish-enforce"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePublishEnforce() *schema.Resource {
	return &schema.Resource{
		Description: "Publish and Enforce resource - triggers publish and / or enforce operations. " +
			"Works the same as running `inext publish` and `inext enforce` CLI commands. " +
			"**Note: Only ONE instance of this resource is allowed per provider/account.**",

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
			"profile_ids": {
				Type:        schema.TypeList,
				Description: "List of profile IDs to enforce. If empty, all profiles will be enforced",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"profile_types": {
				Type:        schema.TypeList,
				Description: "List of profile types to publish (e.g., Kubernetes, Embedded). If empty, all profiles will be published",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"skip_nginx_validation": {
				Type:        schema.TypeBool,
				Description: "When true, skips nginx configuration validation during publish. Useful when publishing policies that include custom nginx configurations that may not pass standard validation",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

// resourcePublishEnforceCreateOrUpdate handles both create and update operations
// since they perform identical logic: trigger publish/enforce based on new values
func resourcePublishEnforceCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)
	if d.Id() == "" {
		d.SetId("publish-enforce")
	}

	shouldPublish := publishenforce.ShouldPublishFromResourceData(d)
	shouldEnforce := publishenforce.ShouldEnforceFromResourceData(d)

	// Execute publish if requested (same as `inext publish`)
	if shouldPublish {
		publishOpts := publishenforce.GetPublishOptionsFromResourceData(d)
		if err := publishenforce.ExecutePublish(ctx, c, publishOpts); err != nil {
			return utils.DiagError("failed to publish changes", err, diags)
		}
	}

	if shouldEnforce {
		profileIDs := publishenforce.GetProfileIDsFromResourceData(d)
		if err := publishenforce.ExecuteEnforce(ctx, c, profileIDs); err != nil {
			return utils.DiagError("failed to enforce policy", err, diags)
		}
	}

	d.Set("publish", false)
	d.Set("enforce", false)

	return diags
}

func resourcePublishEnforceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// This is a trigger-only resource, nothing to read from the API
	// Always report state as false so that config with true triggers an update
	d.Set("publish", false)
	d.Set("enforce", false)
	return nil
}

func resourcePublishEnforceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
