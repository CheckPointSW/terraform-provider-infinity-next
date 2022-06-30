package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	embeddedprofile "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/embedded-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceEmbeddedProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Embedded profile",

		CreateContext: resourceEmbeddedProfileCreate,
		ReadContext:   resourceEmbeddedProfileRead,
		UpdateContext: resourceEmbeddedProfileUpdate,
		DeleteContext: resourceEmbeddedProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			if diff.HasChange("additional_settings") {
				return diff.SetNewComputed("additional_settings_ids")
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
			"defined_applications_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"upgrade_mode": {
				Type: schema.TypeString,
				Description: "The upgrade mode of the profile: Automatic, Manual or Scheduled.\n" +
					"The default is Automatic",
				Optional:         true,
				Default:          embeddedprofile.UpgradeModeAutomatic,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{embeddedprofile.UpgradeModeAutomatic, embeddedprofile.UpgradeModeManual, embeddedprofile.UpgradeModeScheduled}, false)),
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

func resourceEmbeddedProfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := embeddedprofile.CreateEmbeddedProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform EmbeddedProfile Create", err, diags)
	}

	profile, err := embeddedprofile.NewEmbeddedProfile(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform EmbeddedProfile Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following EmbeddedProfile Create", err, diags)
	}

	if err = embeddedprofile.ReadEmbeddedProfileToResourceData(profile, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to read EmbeddedProfile to resource data", err, diags)
	}

	return diags
}

func resourceEmbeddedProfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	id := d.Id()

	profile, err := embeddedprofile.GetEmbeddedProfile(ctx, c, id)
	if err != nil {
		return utils.DiagError("unable to perform EmbeddedProfile Read", err, diags)
	}

	if err := embeddedprofile.ReadEmbeddedProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform EmbeddedProfile Read", err, diags)
	}

	return diags
}

func resourceEmbeddedProfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	updateInput, err := embeddedprofile.UpdateEmbeddedProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform EmbeddedProfile Update", err, diags)
	}

	result, err := embeddedprofile.UpdateEmbeddedProfile(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform EmbeddedProfile Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following EmbeddedProfile Update", err, diags)
	}

	profile, err := embeddedprofile.GetEmbeddedProfile(ctx, c, d.Id())
	if err != nil {
		return utils.DiagError("failed get EmbeddedProfile after update", err, diags)
	}

	if err := embeddedprofile.ReadEmbeddedProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform read EmbeddedProfile read after update", err, diags)
	}

	return diags
}

func resourceEmbeddedProfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := embeddedprofile.DeleteEmbeddedProfile(ctx, c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform EmbeddedProfile Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following EmbeddedProfile Delete", err, diags)
	}

	d.SetId("")

	return diags
}
