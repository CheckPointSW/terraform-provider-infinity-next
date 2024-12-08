package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceWebAppAsset() *schema.Resource {
	validateStateFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{suggestedState, activeState, headerKey, inactiveState}, false))
	return &schema.Resource{
		Description: "Web Application Asset",

		CreateContext: resourceWebAppAssetCreate,
		ReadContext:   resourceWebAppAssetRead,
		UpdateContext: resourceWebAppAssetUpdate,
		DeleteContext: resourceWebAppAssetDelete,
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
				Description: "The name of the resource, also acts as its unique ID",
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
			"behaviors": {
				Type:        schema.TypeSet,
				Description: "behaviors used by the asset",
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
			"tags": {
				Type:        schema.TypeSet,
				Description: "The tags used by the asset",
				Optional:    true,
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
			"practice": {
				Type:        schema.TypeSet,
				Description: "The practices used by the asset",
				Optional:    true,
				Elem: &schema.Resource{
					Description: "Practice wrapper",
					CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
						if diff.HasChange("practice_wrapper_id") {
							return diff.SetNewComputed("practice_wrapper_id")
						}

						return nil
					},
					Schema: map[string]*schema.Schema{
						"main_mode": {
							Type:        schema.TypeString,
							Description: "The mode of the practice: Prevent, Inactive, Disabled or Learn",
							Required:    true,
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
				Description: "Defines how the source identifier values of the asset are retrieved",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Description: "The identifier of the source: SourceIP, XForwardedFor, HeaderKey Cookie or JWTKey",
							Optional:    true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"values_ids": {
							Type:     schema.TypeSet,
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
			"is_shares_urls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mtls": {
				Type:        schema.TypeSet,
				Description: "The MTLS settings",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filename": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"data_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enable_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceWebAppAssetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := webappasset.CreateWebApplicationAssetInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebAppAsset Create", err, diags)
	}

	fmt.Printf("created input: %v\n", createInput)

	asset, err := webappasset.NewWebApplicationAsset(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppAsset Create", err, diags)
	}

	fmt.Printf("created asset: %v\n", asset)

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAppAsset Create", err, diags)
	}

	if err := webappasset.ReadWebApplicationAssetToResourceData(asset, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	return diags
}

func resourceWebAppAssetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	asset, err := webappasset.GetWebApplicationAsset(ctx, c, d.Id())
	if err != nil {
		return utils.DiagError("unable to perform WebAppAsset Read", err, diags)
	}

	fmt.Printf("read asset: %v\n", asset)

	if err := webappasset.ReadWebApplicationAssetToResourceData(asset, d); err != nil {
		return utils.DiagError("unable to perform WebAppAsset Read", err, diags)
	}

	fmt.Printf("read resource data: %v\n", d)

	return diags
}

func resourceWebAppAssetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	oldAsset, err := webappasset.GetWebApplicationAsset(ctx, c, d.Id())
	if err != nil {
		return utils.DiagError("unable to perform get WebApplicationAsset for updating", err, diags)
	}

	updateInput, err := webappasset.UpdateWebApplicationAssetInputFromResourceData(d, oldAsset)
	if err != nil {
		return utils.DiagError("unable to perform WebAppAsset Update", err, diags)
	}

	fmt.Printf("update input: %v\n", updateInput)

	result, err := webappasset.UpdateWebApplicationAsset(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppAsset Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAppAsset Update", err, diags)
	}

	asset, err := webappasset.GetWebApplicationAsset(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	fmt.Printf("updated asset: %v\n", asset)

	if err := webappasset.ReadWebApplicationAssetToResourceData(asset, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	fmt.Printf("updated resource data: %v\n", d)

	return diags
}

func resourceWebAppAssetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := webappasset.DeleteWebApplicationAsset(ctx, c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppAsset Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAppAsset Delete", err, diags)
	}

	d.SetId("")

	return diags
}
