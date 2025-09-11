package resources

import (
	"context"
	"strings"

	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/rate-limit-practice"
	webAppAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	ratelimitpractice "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/rate-limit-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	minuteScope = "Minute"
	secondScope = "Second"
)

func ResourceRateLimitPractice() *schema.Resource {
	validationVisibility := validation.ToDiagFunc(
		validation.StringInSlice([]string{visibilityShared, visibilityLocal}, false))
	validateRuleScopeFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{minuteScope, secondScope}, false))
	validateRuleActionFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{detectMode, preventMode, accordingToPracticeMode}, false))

	return &schema.Resource{
		Description: "Rate limit Practice",

		CreateContext: resourceRateLimitPracticeCreate,
		ReadContext:   resourceRateLimitPracticeRead,
		UpdateContext: resourceRateLimitPracticeUpdate,
		DeleteContext: resourceRateLimitPracticeDelete,
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
			"practice_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:             schema.TypeString,
				Description:      "The visibility of the resource, Shared or Local",
				Default:          "Shared",
				Optional:         true,
				ValidateDiagFunc: validationVisibility,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rule": {
				Type:        schema.TypeSet,
				Description: "Rate limit rules",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scope": {
							Type:             schema.TypeString,
							Description:      "The time unit during which the rate limit is enfrorced. Must be one of: Minute, Second",
							Required:         true,
							ValidateDiagFunc: validateRuleScopeFunc,
						},
						"limit": {
							Type:        schema.TypeInt,
							Description: "The actual number of requests to enable",
							Required:    true,
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "A general comment which describes the rate limit rule",
							Optional:    true,
						},
						"action": {
							Type:             schema.TypeString,
							Description:      "The action to perform upon a request which crosses the rate limit. Must be one of: Detect, Prevent, AccordingToPractice (case sensitive!!). Defaults to AccordingToPractice",
							Optional:         true,
							Default:          accordingToPracticeMode,
							ValidateDiagFunc: validateRuleActionFunc,
						},
					},
				},
			},
		},
	}
}

func resourceRateLimitPracticeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	createInput, err := ratelimitpractice.CreateRateLimitPracticeInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("Failed to create rate limit practice input struct", err, diags)
	}

	practice, err := ratelimitpractice.NewRateLimitPracticePractice(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform RateLimitPractice Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following RateLimitPractice Create", err, diags)
	}

	if err := ratelimitpractice.ReadRateLimitPracticeToResourceData(practice, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to perform read after RateLimitPractice create", err, diags)
	}

	return diags
}

func resourceRateLimitPracticeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)
	id := d.Id()

	practice, err := ratelimitpractice.GetRateLimitPractice(ctx, c, id, false)
	if err != nil {
		return utils.DiagError("unable to get rate limit practice to perform RateLimitPractice Read", err, diags)
	}

	if err := ratelimitpractice.ReadRateLimitPracticeToResourceData(practice, d); err != nil {
		return utils.DiagError("unable to perform RateLimitPractice Read", err, diags)
	}

	return diags
}

func resourceRateLimitPracticeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	updateInput, err := ratelimitpractice.UpdateRateLimitPracticeInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("Failed to parse RateLimitPractice Update to struct", err, diags)
	}

	result, err := ratelimitpractice.UpdateRateLimitPractice(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform RateLimitPractice Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following RateLimitPractice Update", err, diags)
	}

	practice, err := ratelimitpractice.GetRateLimitPractice(ctx, c, d.Id(), true)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes after publish", discardErr, diags)
		}

		return utils.DiagError("failed to get following RateLimitPractice after update", err, diags)
	}

	if err := ratelimitpractice.ReadRateLimitPracticeToResourceData(practice, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Read following RateLimitPractice after update", err, diags)
	}

	return diags
}

func resourceRateLimitPracticeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	result, err := ratelimitpractice.DeleteRateLimitPractice(ctx, c, d.Id())
	if err != nil || !result {
		// Check if the error is due to the rate limit practice being used by other resources
		if err != nil && strings.Contains(err.Error(), "can't be deleted since it is pointed from other objects") {
			// Get the resources that are using the rate limit practice
			usedBy, err2 := ratelimitpractice.UsedByRateLimitPractice(ctx, c, d.Id())
			if err2 != nil {
				diags = utils.DiagError("unable to perform RateLimitPractice UsedBy", err2, diags)
				return utils.DiagError("unable to perform RateLimitPractice Delete", err, diags)
			}

			if usedBy != nil || len(usedBy) > 0 {
				// Remove the rate limit practice from the resources that are using it
				if err2 := handleRateLimitPracticeReferences(ctx, usedBy, c, d.Id()); err2 != nil {
					diags = err2
					return utils.DiagError("unable to perform RateLimitPractice Delete", err, diags)
				}

				// Retry to delete the rate limit practice
				result, err := ratelimitpractice.DeleteRateLimitPractice(ctx, c, d.Id())
				if err != nil || !result {
					if _, discardErr := c.DiscardChanges(); discardErr != nil {
						diags = utils.DiagError("failed to discard changes", discardErr, diags)
					}

					return utils.DiagError("unable to perform RateLimitPractice Delete after updating references", err, diags)
				}

			}

		} else {
			if _, discardErr := c.DiscardChanges(); discardErr != nil {
				diags = utils.DiagError("failed to discard changes", discardErr, diags)
			}

			return utils.DiagError("unable to perform RateLimitPractice Delete", err, diags)
		}
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following RateLimitPractice Delete", err, diags)
	}

	d.SetId("")

	return diags
}

func handleRateLimitPracticeReferences(ctx context.Context, usedBy models.DisplayObjects, c *api.Client, practiceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, usedByResource := range usedBy {
		if usedByResource.ObjectStatus == "Deleted" || usedByResource.Type == "Wrapper" {
			continue
		}

		switch usedByResource.SubType {
		case "WebApplication":
			webAppAsset := webAppAssetModels.UpdateWebApplicationAssetInput{
				RemovePracticeWrappers: []string{practiceID},
			}

			updated, err := webappasset.UpdateWebApplicationAsset(ctx, c, usedByResource.ID, webAppAsset)
			if err != nil || !updated {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("failed to perform UpdateWebApplicationAsset to remove practice", err, diags)
			}

		default:
			return utils.DiagError("failed to update usedByResource", nil, diags)
		}

	}

	return nil
}
