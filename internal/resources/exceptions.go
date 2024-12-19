package resources

import (
	"context"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/exceptions"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// maxNestLevel is the max nesting level of the matchSchema
	// this is used to avoid infinite recursion when creating the schema
	maxNestLevel = 20
)

func matchSchema(nestLevel int) *schema.Resource {
	if nestLevel == 0 {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"operator": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"and", "or", "not-equals", "equals", "in", "not-in", "exist"}, false)),
				},
				"operand": {
					Optional: true,
					Type:     schema.TypeSet,
					Elem:     &schema.Resource{},
				},
				"key": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"hostName", "sourceIdentifier", "url", "countryCode", "countryName", "manufacturer", "paramName", "paramValue", "protectionName", "sourceIP"}, false)),
				},
				"value": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		}
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"and", "or", "not-equals", "equals", "in", "not-in", "exist"}, false)),
				Default:          "equals",
			},
			"operand": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     matchSchema(nestLevel - 1),
			},
			"key": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"hostName", "sourceIdentifier", "url", "countryCode", "countryName", "manufacturer", "paramName", "paramValue", "protectionName", "sourceIP"}, false)),
			},
			"value": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ResourceExceptions() *schema.Resource {
	validateVisibility := validation.ToDiagFunc(
		validation.StringInSlice([]string{visibilityShared, visibilityLocal}, false))
	return &schema.Resource{
		Description: "Exceptions allows overriding the AppSec ML engine decision based on specific parameters",

		CreateContext: resourceExceptionsCreate,
		ReadContext:   resourceExceptionsRead,
		UpdateContext: resourceExceptionsUpdate,
		DeleteContext: resourceExceptionsDelete,
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
				Description: "The name of the resource, also acts as its unique ID",
				Required:    true,
			},
			"visibility": {
				Type:             schema.TypeString,
				Description:      "The visibility of the exception: Shared or Local",
				Default:          "Shared",
				Optional:         true,
				ValidateDiagFunc: validateVisibility,
			},
			"exception": {
				Type:        schema.TypeSet,
				Description: "Overrides AppSec ML engine decision based on match and action",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     matchSchema(maxNestLevel),
						},
						"action": {
							Type:             schema.TypeString,
							Description:      "The action of the exception: accept, drop, skip or suppressLog",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"accept", "drop", "skip", "suppressLog"}, false)),
						},
						"action_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "Comment for the exception",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceExceptionsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	createInput, err := exceptions.CreateExceptionBehaviorInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform ExceptionBehavior Create", err, diags)
	}

	behavior, err := exceptions.NewExceptionBehavior(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform ExceptionBehavior Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following ExceptionBehavior Create", err, diags)
	}

	if err := exceptions.ReadExceptionBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to read ExceptionBehavior into state file after create and publish", err, diags)
	}

	return diags
}

func resourceExceptionsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	behavior, err := exceptions.GetExceptionBehavior(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to get ExceptionBehavior for read into state file", err, diags)
	}

	if err := exceptions.ReadExceptionBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to read ExceptionBehavior into state file", err, diags)
	}

	return diags
}

func resourceExceptionsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	updateInput, err := exceptions.UpdateExceptionBehaviorInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform ExceptionBehavior Update", err, diags)
	}

	result, err := exceptions.UpdateExceptionBehavior(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform ExceptionBehavior Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following ExceptionBehavior Update", err, diags)
	}

	behavior, err := exceptions.GetExceptionBehavior(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Get ExceptionBehavior following Update", err, diags)
	}

	if err := exceptions.ReadExceptionBehaviorToResourceData(behavior, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to read ExceptionBehavior into state file after update and publish", err, diags)
	}

	return diags
}

func resourceExceptionsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	result, err := exceptions.DeleteExceptionBehavior(ctx, c, d.Id())
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform ExceptionBehavior Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following ExceptionBehavior Delete", err, diags)
	}

	d.SetId("")

	return diags
}
