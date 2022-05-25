package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webapiasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-api-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Allowed practice modes
	detectMode              = "Detect"
	preventMode             = "Prevent"
	inactiveMode            = "Inactive"
	accordingToPracticeMode = "AccordingToPractice"
	disabledMode            = "Disabled"
	learnMode               = "Learn"
	activeMode              = "Active"

	// Allowed source identifiers
	sourceIP      = "SourceIP"
	xForwardedFor = "XForwardedFor"
	headerKey     = "HeaderKey"
	cookie        = "Cookie"

	// Allowed states
	suggestedState = "Suggested"
	activeState    = "Active"
	inactiveState  = "InActive"
)

func ResourceWebAPIAsset() *schema.Resource {
	validatePracticeModeFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{detectMode, preventMode, inactiveMode, accordingToPracticeMode, disabledMode, learnMode, activeMode}, false))
	validateSourceIdentifierFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{sourceIP, xForwardedFor, headerKey, cookie}, false))
	validateStateFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{suggestedState, activeState, headerKey, inactiveState}, false))

	return &schema.Resource{
		Description:   "Web API Asset",
		CreateContext: resourceWebApiAssetCreate,
		ReadContext:   resourceWebApiAssetRead,
		UpdateContext: resourceWebApiAssetUpdate,
		DeleteContext: resourceWebApiAssetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			if diff.HasChange("urls") {
				return diff.SetNewComputed("urls_ids")
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
				Description: "The name of the resource, also acts as it's unique ID",
				Required:    true,
			},
			"profiles": {
				Type:        schema.TypeSet,
				Description: "Profiles linked to the asset",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// top level behaviors
			"trusted_sources": {
				Type:        schema.TypeSet,
				Description: "Trusted sources behavior used by the asset",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateStateFunc,
			},
			"upstream_url": {
				Type: schema.TypeString,
				Description: "The URL of the application's backend server to which the reverse proxy redirects " +
					"the relevant traffic sent to the exposed URL",
				Optional: true,
			},
			"urls": {
				Type:        schema.TypeSet,
				Description: "The application URLs",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"urls_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"practice": {
				Type:        schema.TypeSet,
				Description: "The practices used by the asset",
				Optional:    true,
				Elem: &schema.Resource{
					Description: "Practice wrapper",
					Schema: map[string]*schema.Schema{
						"main_mode": {
							Type:             schema.TypeString,
							Description:      "The mode of the practice: Prevent, Inactive, Disabled or Learn",
							Optional:         true,
							ValidateDiagFunc: validatePracticeModeFunc,
						},
						"sub_practices_modes": {
							Type:        schema.TypeMap,
							Description: "The name of the sub practice as the key and its mode as the value. Allowed modes: Detect, Prevent, Inactive, AccordingToPractice, Disabled, Learn or Active",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"practice_wrapper_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"triggers": {
							Type:        schema.TypeSet,
							Description: "The triggers used with the practice",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// practices.behaviors
						"exceptions": {
							Type:        schema.TypeSet,
							Description: "The exceptions used with the practice",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"proxy_setting": {
				Type:        schema.TypeSet,
				Description: "Settings for the proxy",
				Optional:    true,
				// Remove Computed if default for Set/List is supported - manually edit generated docs and move proxy_setting out of read-only section
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"source_identifier": {
				Type:        schema.TypeSet,
				Description: "Defines how the source identifier valuess of the asset are retrieved",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:             schema.TypeString,
							Description:      "The identifier of the source: SourceIP, XForwardedFor, HeaderKey or Cookie",
							Optional:         true,
							ValidateDiagFunc: validateSourceIdentifierFunc,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"values_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"asset_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sources": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intelligence_tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceWebApiAssetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := webapiasset.CreateWebAPIAssetInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebAPIAsset Create", err, diags)
	}

	asset, err := webapiasset.NewWebAPIAsset(c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("Failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform WebAPIAsset Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("Failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Failed to Publish following WebAPIAsset Create", err, diags)
	}

	if err := webapiasset.ReadWebAPIAssetToResourceData(asset, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to read WebAPIAsset to resource data", err, diags)
	}

	return diags
}

func resourceWebApiAssetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	asset, err := webapiasset.GetWebAPIAsset(c, d.Id())
	if err != nil {
		return utils.DiagError("unable to perform get WebAPIAsset", err, diags)
	}

	if err := webapiasset.ReadWebAPIAssetToResourceData(asset, d); err != nil {
		return utils.DiagError("unable to perform WebAPIAsset Read", err, diags)
	}

	return diags
}

func resourceWebApiAssetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	updateInput, err := webapiasset.UpdateWebAPIAssetInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebAPIAsset update", err, diags)
	}

	result, err := webapiasset.UpdateWebAPIAsset(c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAPIAsset Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAPIAsset Update", err, diags)
	}

	// get the newly created asset and read it into the state
	newAsset, err := webapiasset.GetWebAPIAsset(c, d.Id())
	if err != nil {
		return utils.DiagError("unable to perform get WebAPIAsset", err, diags)
	}

	if err := webapiasset.ReadWebAPIAssetToResourceData(newAsset, d); err != nil {
		return utils.DiagError("unable to perform WebAPIAsset Read", err, diags)
	}

	return diags
}

func resourceWebApiAssetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := webapiasset.DeleteWebAPIAsset(c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAPIAsset Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAPIAsset Delete", err, diags)
	}

	d.SetId("")

	return diags
}
