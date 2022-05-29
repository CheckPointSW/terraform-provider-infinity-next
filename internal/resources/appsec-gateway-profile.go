package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	appsecgatewayprofile "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/appsec-gateway-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAppSecGatewayProfile() *schema.Resource {
	return &schema.Resource{
		Description: "CloudGuard Application Security Gateway profile is deployed as a VM that runs on a Check Point Gaia OS " +
			"with a reverse proxy and Check Point Nano-Agent",

		CreateContext: resourceAppSecGatewayProfileCreate,
		ReadContext:   resourceAppSecGatewayProfileRead,
		UpdateContext: resourceAppSecGatewayProfileUpdate,
		DeleteContext: resourceAppSecGatewayProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			if diff.HasChange("additional_settings") {
				if err := diff.SetNewComputed("additional_settings_ids"); err != nil {
					return err
				}
			}

			if diff.HasChange("reverseproxy_additional_settings") {
				if err := diff.SetNewComputed("reverseproxy_additional_settings_ids"); err != nil {
					return err
				}
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the resource, also acts as its unique ID",
				Required:    true,
			},
			"profile_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_sub_type": {
				Type:         schema.TypeString,
				Description:  "The environment of deployment for the AppSec VM: Aws, Azure, VMware or HyperV",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{appsecgatewayprofile.ProfileSubTypeAws, appsecgatewayprofile.ProfileSubTypeAzure, appsecgatewayprofile.ProfileSubTypeVMware, appsecgatewayprofile.ProfileSubTypeHyperV}, false),
			},
			"additional_settings": {
				Type:        schema.TypeMap,
				Description: "Controls the settings of the connected agents",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"additional_settings_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"upgrade_mode": {
				Type: schema.TypeString,
				Description: "The upgrade mode of the profile: Automatic, Manual or Scheduled.\n" +
					"The default is Automatic",
				Optional:         true,
				Default:          appsecgatewayprofile.UpgradeModeAutomatic,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{appsecgatewayprofile.UpgradeModeAutomatic, appsecgatewayprofile.UpgradeModeManual, appsecgatewayprofile.UpgradeModeScheduled}, false)),
			},
			"upgrade_time_schedule_type": {
				Type:             schema.TypeString,
				Description:      "The schedule type in case upgrade mode is scheduled: DaysInWeek",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DaysInWeek"}, false)),
			},
			"upgrade_time_hour": {
				Type:        schema.TypeString,
				Description: "The hour of the upgrade time start, for example: 10:00 or 20:00",
				Optional:    true,
			},
			"upgrade_time_duration": {
				Type:        schema.TypeInt,
				Description: "The duration of the upgrade in hours",
				Optional:    true,
			},
			"upgrade_time_week_days": {
				Type:        schema.TypeSet,
				Description: "The week days of the upgrade time schedule: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reverseproxy_upstream_timeout": {
				Type:        schema.TypeInt,
				Description: "Sets the reverse proxy upstream timeout in seconds",
				Optional:    true,
			},
			"reverseproxy_additional_settings": {
				Type:        schema.TypeMap,
				Description: "Sets the reverse proxy settings of linked assets",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reverseproxy_additional_settings_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"max_number_of_agents": {
				Type:             schema.TypeInt,
				Description:      "Sets the maximum number of agents that can be connected to this profile",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtMost(1000)),
			},
			"authentication_token": {
				Type:        schema.TypeString,
				Description: "The token used to register an agent to the profile",
				Computed:    true,
			},
		},
	}
}

func resourceAppSecGatewayProfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := appsecgatewayprofile.CreateCloudGuardAppSecGatewayProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform AppSecGatewayProfile Create", err, diags)
	}

	profile, err := appsecgatewayprofile.NewAppSecGatewayProfile(c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform AppSecGatewayProfile Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following AppSecGatewayProfile Create", err, diags)
	}

	if err = appsecgatewayprofile.ReadCloudGuardAppSecGatewayProfileToResourceData(profile, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to read AppSecGatewayProfile to resource data", err, diags)
	}

	return diags
}

func resourceAppSecGatewayProfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	id := d.Id()

	profile, err := appsecgatewayprofile.GetCloudGuardAppSecGatewayProfile(c, id)
	if err != nil {
		return utils.DiagError("unable to perform AppSecGatewayProfile Read", err, diags)
	}

	if err := appsecgatewayprofile.ReadCloudGuardAppSecGatewayProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform AppSecGatewayProfile Read", err, diags)
	}

	return diags
}

func resourceAppSecGatewayProfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	updateInput, err := appsecgatewayprofile.UpdateCloudGuardAppSecGatewayProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform AppSecGatewayProfile Update", err, diags)
	}

	result, err := appsecgatewayprofile.UpdateAppSecGatewayProfile(c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform AppSecGatewayProfile Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following AppSecGatewayProfile Update", err, diags)
	}

	profile, err := appsecgatewayprofile.GetCloudGuardAppSecGatewayProfile(c, d.Id())
	if err != nil {
		return utils.DiagError("failed get AppSecGatewayProfile after update", err, diags)
	}

	if err := appsecgatewayprofile.ReadCloudGuardAppSecGatewayProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform read AppSecGatewayProfile read after update", err, diags)
	}

	return diags
}

func resourceAppSecGatewayProfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := appsecgatewayprofile.DeleteAppSecGatewayProfile(c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform AppSecGatewayProfile Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following AppSecGatewayProfile Delete", err, diags)
	}

	d.SetId("")

	return diags
}
