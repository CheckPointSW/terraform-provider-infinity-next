package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webapppractice "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceWebAppPractice() *schema.Resource {
	return &schema.Resource{
		Description: "Web Application Practice",

		CreateContext: resourceWebAppPracticeCreate,
		ReadContext:   resourceWebAppPracticeRead,
		UpdateContext: resourceWebAppPracticeUpdate,
		DeleteContext: resourceWebAppPracticeDelete,
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
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ips": {
				Type:        schema.TypeSet,
				Description: "IPS protection",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"performance_impact": {
							Type:             schema.TypeString,
							Description:      "The performance impact: LowOrLower, MediumOrLower or HighOrLower",
							Default:          "MediumOrLower",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LowOrLower", "MediumOrLower", "HighOrLower"}, false)),
						},
						"severity_level": {
							Type:             schema.TypeString,
							Description:      "The severity level: LowOrAbove, MediumOrAbove, HighOrAbove or Critical",
							Default:          "MediumOrAbove",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"}, false)),
						},
						"protections_from_year": {
							Type:        schema.TypeString,
							Description: "The year to apply protections from: 1999, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020",
							Default:     "2020",
							Optional:    true,
						},
						"high_confidence": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent or Inactive",
							Default:          "Prevent",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Detect", "Prevent", "Inactive"}, false)),
						},
						"medium_confidence": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent or Inactive",
							Default:          "Prevent",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Detect", "Prevent", "Inactive"}, false)),
						},
						"low_confidence": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent or Inactive",
							Default:          "Detect",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Detect", "Prevent", "Inactive"}, false)),
						},
						// "advanced_policy": {
						// 	Type:     schema.TypeSet,
						// 	Optional: true,
						// 	MaxItems: 1,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 			"filename": {
						// 				Type:     schema.TypeString,
						// 				Required: true,
						// 			},
						// 			"data": {
						// 				Type:      schema.TypeString,
						// 				Sensitive: true,
						// 				Required:  true,
						// 			},
						// 			"size": {
						// 				Type:     schema.TypeInt,
						// 				Optional: true,
						// 				Computed: true,
						// 			},
						// 			"override_setting": {
						// 				Type:     schema.TypeBool,
						// 				Default:  false,
						// 				Optional: true,
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
			"web_attacks": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"minimum_severity": {
							Type:             schema.TypeString,
							Description:      "Medium, High or Critical",
							Default:          "High",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Critical", "High", "Medium"}, false)),
						},
						"advanced_setting": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"csrf_protection": {
										Type:             schema.TypeString,
										Description:      "Prevent, AccordingToPractice, Disabled or Learn",
										Default:          "Disabled",
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Disabled", "Learn", "Prevent", "AccordingToPractice"}, false)),
									},
									"open_redirect": {
										Type:             schema.TypeString,
										Description:      "Prevent, AccordingToPractice, Disabled or Learn",
										Default:          "Disabled",
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Disabled", "Learn", "Prevent", "AccordingToPractice"}, false)),
									},
									"error_disclosure": {
										Type:             schema.TypeString,
										Description:      "Prevent, AccordingToPractice, Disabled or Learn",
										Default:          "Disabled",
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Disabled", "Learn", "Prevent", "AccordingToPractice"}, false)),
									},
									"body_size": {
										Type:     schema.TypeInt,
										Default:  1000000,
										Optional: true,
									},
									"url_size": {
										Type:     schema.TypeInt,
										Default:  32768,
										Optional: true,
									},
									"header_size": {
										Type:     schema.TypeInt,
										Default:  102400,
										Optional: true,
									},
									"max_object_depth": {
										Type:     schema.TypeInt,
										Default:  40,
										Optional: true,
									},
									"illegal_http_methods": {
										Type:     schema.TypeBool,
										Default:  false,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"web_bot": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inject_uris": {
							Type:        schema.TypeSet,
							Description: "Defines where to inject the Anti-Bot script. The input is the path of the URI",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"inject_uris_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"valid_uris": {
							Type: schema.TypeSet,
							Description: "Defines which requests must be validated after the script is injected into a specific URI.\n" +
								"The input is the path of the URI",
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"valid_uris_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceWebAppPracticeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := webapppractice.CreateWebApplicationPracticeInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebAppPractice Create", err, diags)
	}

	practice, err := webapppractice.NewWebApplicationPractice(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppPractice Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAppPractice Create", err, diags)
	}

	if err := webapppractice.ReadWebApplicationPracticeToResourceData(practice, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	return diags
}

func resourceWebAppPracticeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	id := d.Id()

	practice, err := webapppractice.GetWebApplicationPractice(ctx, c, id)
	if err != nil {
		return utils.DiagError("unable to perform WebAppPractice Read", err, diags)
	}

	if err := webapppractice.ReadWebApplicationPracticeToResourceData(practice, d); err != nil {
		return utils.DiagError("unable to perform WebAppPractice Read", err, diags)
	}

	return diags
}

func resourceWebAppPracticeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	updateInput, err := webapppractice.UpdateWebApplicationPracticeInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("Failed to parse WebAppPractice Update to struct", err, diags)
	}

	result, err := webapppractice.UpdateWebApplicationPractice(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppPractice Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAppPractice Update", err, diags)
	}

	practice, err := webapppractice.GetWebApplicationPractice(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	if err := webapppractice.ReadWebApplicationPracticeToResourceData(practice, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	return diags
}

func resourceWebAppPracticeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	result, err := webapppractice.DeleteWebApplicationPractice(ctx, c, d.Id())
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAppPractice Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAppPractice Delete", err, diags)
	}

	d.SetId("")

	return diags
}
