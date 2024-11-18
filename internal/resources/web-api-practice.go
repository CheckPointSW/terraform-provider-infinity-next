package resources

import (
	"context"
	"fmt"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webapipractice "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-api-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceWebAPIPractice() *schema.Resource {
	validationSeverityLevel := validation.ToDiagFunc(
		validation.StringInSlice([]string{severityLevelLowOrAbove, severityLevelMediumOrAbove, severityLevelHighOrAbove, severityLevelCritical}, false))
	validationFileSecurityMode := validation.ToDiagFunc(
		validation.StringInSlice([]string{fileSecurityModeDetect, fileSecurityModePrevent, fileSecurityModeInactive, fileSecurityModeAccordingToPractice}, false))
	//validationFileSizeUnits := validation.ToDiagFunc(
	//	validation.StringInSlice([]string{fileSizeUnitsBytes, fileSizeUnitsKB, fileSizeUnitsMB, fileSizeUnitsGB}, false))
	//validationWAAPMode := validation.ToDiagFunc(
	//	validation.StringInSlice([]string{waapModeDisabled, waapModeLearn, waapModePrevent, waapModePractice}, false))
	return &schema.Resource{
		Description: "Practice for securing a web API",

		CreateContext: resourceWebAPIPracticeCreate,
		ReadContext:   resourceWebAPIPracticeRead,
		UpdateContext: resourceWebAPIPracticeUpdate,
		DeleteContext: resourceWebAPIPracticeDelete,
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
						},
						"performance_impact": {
							Type:             schema.TypeString,
							Description:      "The performance impact: VeryLow, LowOrLower, MediumOrLower or HighOrLower",
							Default:          "MediumOrLower",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VeryLow", "LowOrLower", "MediumOrLower", "HighOrLower"}, false)),
						},
						"severity_level": {
							Type:             schema.TypeString,
							Description:      "The severity level: LowOrAbove, MediumOrAbove, HighOrAbove or Critical",
							Default:          "MediumOrAbove",
							Optional:         true,
							ValidateDiagFunc: validationSeverityLevel,
						},
						"protections_from_year": {
							Type:        schema.TypeString,
							Description: "The year to apply protections from: 1999, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020",
							Default:     "2020",
							Optional:    true,
						},
						"high_confidence": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "AccordingToPractice",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
						"medium_confidence": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "AccordingToPractice",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
						"low_confidence": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "Detect",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
					},
				},
			},
			"api_attacks": {
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
			"schema_validation": {
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
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_file_exist": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						//"oas_schema": {
						//	Type:     schema.TypeSet,
						//	Computed: true,
						//	Optional: true,
						//	MaxItems: 1,
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"data": {
						//				Type:      schema.TypeString,
						//				Sensitive: true,
						//				Required:  true,
						//			},
						//			"name": {
						//				Type:     schema.TypeString,
						//				Required: true,
						//			},
						//			"size": {
						//				Type:     schema.TypeInt,
						//				Required: true,
						//			},
						//			"is_file_exist": {
						//				Type:     schema.TypeBool,
						//				Optional: true,
						//			},
						//		},
						//	},
						//},
					},
				},
			},
		},
	}
}

func resourceWebAPIPracticeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	createInput, err := webapipractice.CreateWebAPIPracticeInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebAPIPractice Create", err, diags)
	}

	practice, err := webapipractice.NewWebAPIPractice(ctx, c, createInput)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAPIPractice Create", err, diags)
	}

	fmt.Printf("Created new WebAPIPractice: %+v \n", practice)

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAPIPractice Create", err, diags)
	}

	if err := webapipractice.ReadWebAPIPracticeToResourceData(practice, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	return diags
}

func resourceWebAPIPracticeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)
	id := d.Id()

	fmt.Printf("Reading WebAPIPractice: %s\n", id)

	practice, err := webapipractice.GetWebAPIPractice(ctx, c, id)
	if err != nil {
		return utils.DiagError("unable to perform WebAPIPractice Read", err, diags)
	}

	if err := webapipractice.ReadWebAPIPracticeToResourceData(practice, d); err != nil {
		return utils.DiagError("unable to perform WebAPIPractice Read", err, diags)
	}

	return diags
}

func resourceWebAPIPracticeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	fmt.Printf("Updating WebAPIPractice: %s\n", d.Id())

	updateInput, err := webapipractice.UpdateWebAPIPracticeInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("unable to perform WebAPIPractice Update", err, diags)
	}

	result, err := webapipractice.UpdateWebAPIPractice(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAPIPractice Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAPIPractice Update", err, diags)
	}

	practice, err := webapipractice.GetWebAPIPractice(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	if err := webapipractice.ReadWebAPIPracticeToResourceData(practice, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return diag.FromErr(err)
	}

	return diags
}

func resourceWebAPIPracticeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	fmt.Printf("Deleting WebAPIPractice: %s\n", d.Id())

	result, err := webapipractice.DeleteWebAPIPractice(ctx, c, d.Id())
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("unable to perform WebAPIPractice Delete", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("failed to Publish following WebAPIPractice Delete", err, diags)
	}

	d.SetId("")

	return diags
}
