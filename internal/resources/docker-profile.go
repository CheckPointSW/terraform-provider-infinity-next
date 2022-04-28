package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	dockerprofile "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/docker-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDockerProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Docker profile",

		CreateContext: resourceDockerProfileCreate,
		ReadContext:   resourceDockerProfileRead,
		UpdateContext: resourceDockerProfileUpdate,
		DeleteContext: resourceDockerProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:     schema.TypeList,
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

func resourceDockerProfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := dockerprofile.CreateDockerProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform DockerProfile Create", err, diags)
	}

	profile, err := dockerprofile.NewDockerProfile(c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform DockerProfile Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following DockerProfile Create", err, diags)
	}

	if err = dockerprofile.ReadDockerProfileToResourceData(profile, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to read DockerProfile to resource data", err, diags)
	}

	return diags
}

func resourceDockerProfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	id := d.Id()

	profile, err := dockerprofile.GetDockerProfile(c, id)
	if err != nil {
		return utils.DiagError("unable to perform DockerProfile Read", err, diags)
	}

	if err := dockerprofile.ReadDockerProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform DockerProfile Read", err, diags)
	}

	return diags
}

func resourceDockerProfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	updateInput, err := dockerprofile.UpdateDockerProfileInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform DockerProfile Update", err, diags)
	}

	result, err := dockerprofile.UpdateDockerProfile(c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform DockerProfile Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following DockerProfile Update", err, diags)
	}

	profile, err := dockerprofile.GetDockerProfile(c, d.Id())
	if err != nil {
		return utils.DiagError("failed get DockerProfile after update", err, diags)
	}

	if err := dockerprofile.ReadDockerProfileToResourceData(profile, d); err != nil {
		return utils.DiagError("unable to perform read DockerProfile read after update", err, diags)
	}

	return diags
}

func resourceDockerProfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := dockerprofile.DeleteDockerProfile(c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform DockerProfile Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following DockerProfile Delete", err, diags)
	}

	d.SetId("")

	return diags
}
