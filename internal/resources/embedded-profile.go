package resources

import (
	"context"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/embedded-profile"
	webAPIAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-asset"
	webAppAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	embeddedprofile "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/embedded-profile"
	webapiasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-api-asset"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceEmbeddedProfile() *schema.Resource {
	validateUpgradeMode := validation.ToDiagFunc(validation.StringInSlice([]string{embeddedprofile.UpgradeModeAutomatic, embeddedprofile.UpgradeModeManual, embeddedprofile.UpgradeModeScheduled}, false))
	validateUpgradeTimeType := validation.ToDiagFunc(validation.StringInSlice([]string{embeddedprofile.ScheduleTypeDaily, embeddedprofile.ScheduleTypeDaysInWeek, embeddedprofile.ScheduleTypeDaysInMonth}, false))
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
				Type:        schema.TypeBool,
				Description: "Sets whether reverse proxy will block undefined applications or not",
				Optional:    true,
				Default:     false,
			},
			"upgrade_mode": {
				Type: schema.TypeString,
				Description: "The upgrade mode of the profile: Automatic, Manual or Scheduled.\n" +
					"The default is Automatic",
				Optional:         true,
				Default:          embeddedprofile.UpgradeModeAutomatic,
				ValidateDiagFunc: validateUpgradeMode,
			},
			"upgrade_time_schedule_type": {
				Type:             schema.TypeString,
				Description:      "The schedule type in case upgrade mode is scheduled: DaysInWeek, DaysInMonth or Daily",
				Optional:         true,
				Default:          embeddedprofile.ScheduleTypeDaysInWeek,
				ValidateDiagFunc: validateUpgradeTimeType,
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	// Ignore differences when upgrade_mode is not Scheduled
				//	return d.Get("upgrade_mode").(string) != embeddedprofile.UpgradeModeScheduled
				//},
			},
			"upgrade_time_hour": {
				Type:        schema.TypeString,
				Description: "The hour of the upgrade time start, for example: 10:00 or 20:00",
				Optional:    true,
				Default:     "0:00",
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	// Ignore differences when upgrade_mode is not Scheduled
				//	return d.Get("upgrade_mode").(string) != embeddedprofile.UpgradeModeScheduled
				//},
			},
			"upgrade_time_duration": {
				Type:        schema.TypeInt,
				Description: "The duration of the upgrade in hours",
				Optional:    true,
				Default:     4,
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	// Ignore differences when upgrade_mode is not Scheduled
				//	return d.Get("upgrade_mode").(string) != embeddedprofile.UpgradeModeScheduled
				//},
			},
			"upgrade_time_week_days": {
				Type:        schema.TypeSet,
				Description: "The week days of the upgrade time schedule: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	// Ignore differences when upgrade_mode is not Scheduled
				//	return d.Get("upgrade_mode").(string) != embeddedprofile.UpgradeModeScheduled
				//},
			},
			"upgrade_time_days": {
				Type:        schema.TypeSet,
				Description: "The days of the month of the upgrade time schedule",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	// Ignore differences when upgrade_mode is not Scheduled
				//	return d.Get("upgrade_mode").(string) != embeddedprofile.UpgradeModeScheduled
				//},
			},
			"max_number_of_agents": {
				Type:             schema.TypeInt,
				Description:      "Sets the maximum number of agents that can be connected to this profile",
				Optional:         true,
				Default:          10,
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
		// Check if the error is due to the profile being used by other resources
		if err != nil && strings.Contains(err.Error(), errorMsgPointedObjects) {
			// Get EmbeddedProfile to check if it is used by other resources
			profile, err2 := embeddedprofile.GetEmbeddedProfile(ctx, c, ID)
			if err2 != nil {
				diags = utils.DiagError("unable to Get EmbeddedProfile references", err2, diags)
				return utils.DiagError("unable to perform EmbeddedProfile Delete", err, diags)
			}

			// Remove references
			if err2 := handleEmbeddedProfileReferences(ctx, profile.UsedBy, c, ID); err2 != nil {
				diags = err2
				return utils.DiagError("unable to perform EmbeddedProfile Delete", err, diags)
			}

			// Retry delete after removing references
			result, err := embeddedprofile.DeleteEmbeddedProfile(ctx, c, ID)
			if err != nil || !result {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("unable to perform EmbeddedProfile Delete after updating references", err, diags)
			}
		} else {
			if _, discardErr := c.DiscardChanges(); discardErr != nil {
				diags = utils.DiagError("failed to discard changes", discardErr, diags)
			}

			return utils.DiagError("unable to perform EmbeddedProfile Delete", err, diags)
		}

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

func handleEmbeddedProfileReferences(ctx context.Context, usedBy models.DisplayObjects, c *api.Client, profileID string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, usedByResource := range usedBy {
		switch usedByResource.SubType {
		case "WebAPI":
			webAPIAsset := webAPIAssetModels.UpdateWebAPIAssetInput{
				RemoveProfiles: []string{profileID},
			}

			updated, err := webapiasset.UpdateWebAPIAsset(ctx, c, usedByResource.ID, webAPIAsset)
			if err != nil || !updated {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("failed to perform UpdateWebAPIAsset to remove profile", err, diags)
			}

		case "WebApplication":
			webAppAsset := webAppAssetModels.UpdateWebApplicationAssetInput{
				RemoveProfiles: []string{profileID},
			}

			updated, err := webappasset.UpdateWebApplicationAsset(ctx, c, usedByResource.ID, webAppAsset)
			if err != nil || !updated {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("failed to perform UpdateWebApplicationAsset to remove profile", err, diags)
			}

		default:
			return utils.DiagError("failed to update usedByResource", nil, diags)
		}

	}

	return nil
}
