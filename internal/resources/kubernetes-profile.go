package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	kubernetesprofile "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/kubernetes-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Kubernetes profile",

		CreateContext: resourceKubernetesProfileCreate,
		ReadContext:   resourceKubernetesProfileRead,
		UpdateContext: resourceKubernetesProfileUpdate,
		DeleteContext: resourceKubernetesProfileDelete,
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
			"profile_sub_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AccessControl", "AppSec"}, false)),
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

func resourceKubernetesProfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := kubernetesprofile.CreateKubernetesProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform KubernetesProfile Create", err, diags)
	}

	profile, err := kubernetesprofile.NewKubernetesProfile(c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform KubernetesProfile Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following KubernetesProfile Create", err, diags)
	}

	if err = kubernetesprofile.ReadKubernetesProfileToResourceData(profile, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to read KubernetesProfile to resource data", err, diags)
	}

	return diags
}

func resourceKubernetesProfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	id := d.Id()

	profile, err := kubernetesprofile.GetKubernetesProfile(c, id)
	if err != nil {
		return utils.DiagError("unable to perform KubernetesProfile Read", err, diags)
	}

	if err := kubernetesprofile.ReadKubernetesProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform KubernetesProfile Read", err, diags)
	}

	return diags
}

func resourceKubernetesProfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	updateInput, err := kubernetesprofile.UpdateKubernetesProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform KubernetesProfile Update", err, diags)
	}

	result, err := kubernetesprofile.UpdateKubernetesProfile(c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform KubernetesProfile Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following KubernetesProfile Update", err, diags)
	}

	profile, err := kubernetesprofile.GetKubernetesProfile(c, d.Id())
	if err != nil {
		return utils.DiagError("failed get KubernetesProfile after update", err, diags)
	}

	if err := kubernetesprofile.ReadKubernetesProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform read KubernetesProfile read after update", err, diags)
	}

	return diags
}

func resourceKubernetesProfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := kubernetesprofile.DeleteKubernetesProfile(c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform KubernetesProfile Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following KubernetesProfile Delete", err, diags)
	}

	d.SetId("")

	return diags
}
