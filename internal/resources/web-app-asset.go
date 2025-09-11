package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mTLSServer = "server"
	mTLSClient = "client"

	mTLSFileTypePEM = ".pem"
	mTLSFileTypeCRT = ".crt"
	mTLSFileTypeDER = ".der"
	mTLSFileTypeP12 = ".p12"
	mTLSFileTypePFX = ".pfx"
	mTLSFileTypeP7B = ".p7b"
	mTLSFileTypeP7C = ".p7c"
	mTLSFileTypeCER = ".cer"

	instructionsBlockLocation = "location_instructions"
	instructionsBlockServer   = "server_instructions"

	instructionsBlockFileTypeJSON = ".json"
	instructionsBlockFileTypeYML  = ".yml"
)

func ResourceWebAppAsset() *schema.Resource {
	validateStateFunc := validation.ToDiagFunc(validation.StringInSlice(
		[]string{suggestedState, activeState, inactiveState}, false))
	mTLSTypeValidation := validation.ToDiagFunc(validation.StringInSlice(
		[]string{mTLSServer, mTLSClient}, false))
	mTLSFileTypeValidation := validation.ToDiagFunc(validation.StringInSlice(
		[]string{mTLSFileTypePEM, mTLSFileTypeCRT, mTLSFileTypeDER, mTLSFileTypeP12, mTLSFileTypePFX, mTLSFileTypeP7B, mTLSFileTypeP7C, mTLSFileTypeCER}, false))
	instructionsBlockFileTypeValidation := validation.ToDiagFunc(validation.StringInSlice(
		[]string{instructionsBlockFileTypeJSON, instructionsBlockFileTypeYML}, false))
	instructionsBlockTypeValidation := validation.ToDiagFunc(validation.StringInSlice(
		[]string{instructionsBlockLocation, instructionsBlockServer}, false))
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
			"redirect_to_https": {
				Type:        schema.TypeBool,
				Description: "Advanced Proxy Setting - Redirect incoming HTTP requests to the same URL using HTTPS. (The configured application URLs for this asset must include both the HTTP and the HTTPS version of each URL)",
				Optional:    true,
				Default:     false,
			},
			"redirect_to_https_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_log": {
				Type:        schema.TypeBool,
				Description: "Advanced Proxy Setting - Activate access log on gateway.",
				Optional:    true,
				Default:     false,
			},
			"access_log_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_headers": {
				Type:        schema.TypeSet,
				Description: "Advanced Proxy Settings - The custom headers settings",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_headers_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_instructions_blocks": {
				Type:        schema.TypeSet,
				Description: "The additional instructions blocks settings - location or server blocks",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filename": {
							Description: "The name of the instructions block file",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"filename_type": {
							Description:      "The type of the instructions block file - .json, .yml",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: instructionsBlockFileTypeValidation,
						},
						"data_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Description: "The instructions block data",
							Type:        schema.TypeString,
							Sensitive:   true,
							Optional:    true,
						},
						"type": {
							Description:      "The type of the additional instructions block - location_instructions or server_instructions",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: instructionsBlockTypeValidation,
						},
						"enable_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Description: "Whether the instructions block is enabled",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"mtls": {
				Type:        schema.TypeSet,
				Description: "The mutual TLS settings",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filename": {
							Description: "The name of the certificate file",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"certificate_type": {
							Description:      "The type of the certificate file - .pem, .crt, .der, .p12, .pfx, .p7b, .p7c, .cer",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: mTLSFileTypeValidation,
						},
						"data_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Description: "The certificate data",
							Type:        schema.TypeString,
							Sensitive:   true,
							Optional:    true,
						},
						"type": {
							Description:      "The type of the mTLS - server or client",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: mTLSTypeValidation,
						},
						"enable_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Description: "Whether the mTLS is enabled",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
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

	asset, err := webappasset.NewWebApplicationAsset(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppAsset Create", err, diags)
	}

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

	if err := webappasset.ReadWebApplicationAssetToResourceData(asset, d); err != nil {
		return utils.DiagError("unable to perform WebAppAsset Read", err, diags)
	}

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

	if err := webappasset.ReadWebApplicationAssetToResourceData(asset, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

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
