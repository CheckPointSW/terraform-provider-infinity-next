package resources

import (
	"context"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	webAppAssetModels "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-asset"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-practice"
	webappasset "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-asset"
	webapppractice "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/web-app-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	severityLevelLowOrAbove    = "LowOrAbove"
	severityLevelMediumOrAbove = "MediumOrAbove"
	severityLevelHighOrAbove   = "HighOrAbove"
	severityLevelCritical      = "Critical"

	fileSecurityModeDetect              = "Detect"
	fileSecurityModePrevent             = "Prevent"
	fileSecurityModeInactive            = "Inactive"
	fileSecurityModeAccordingToPractice = "AccordingToPractice"

	fileSizeUnitsBytes = "Bytes"
	fileSizeUnitsKB    = "KB"
	fileSizeUnitsMB    = "MB"
	fileSizeUnitsGB    = "GB"

	waapModeDisabled = "Disabled"
	waapModeLearn    = "Learn"
	waapModePrevent  = "Prevent"
	waapModePractice = "AccordingToPractice"

	errorMsgPointedObjects = "can't be deleted since it is pointed from other objects"
)

func ResourceWebAppPractice() *schema.Resource {
	validationSeverityLevel := validation.ToDiagFunc(
		validation.StringInSlice([]string{severityLevelLowOrAbove, severityLevelMediumOrAbove, severityLevelHighOrAbove, severityLevelCritical}, false))
	validationFileSecurityMode := validation.ToDiagFunc(
		validation.StringInSlice([]string{fileSecurityModeDetect, fileSecurityModePrevent, fileSecurityModeInactive, fileSecurityModeAccordingToPractice}, false))
	validationFileSizeUnits := validation.ToDiagFunc(
		validation.StringInSlice([]string{fileSizeUnitsBytes, fileSizeUnitsKB, fileSizeUnitsMB, fileSizeUnitsGB}, false))
	validationVisibility := validation.ToDiagFunc(
		validation.StringInSlice([]string{visibilityShared, visibilityLocal}, false))
	validationPerformanceImpact := validation.ToDiagFunc(
		validation.StringInSlice([]string{performanceImpactVeryLow, performanceImpactLowOrLower, performanceImpactMediumOrLower, performanceImpactHighOrLower}, false))
	validationMinimumSeverity := validation.ToDiagFunc(
		validation.StringInSlice([]string{severityLevelCritical, severityLevelHigh, severityLevelMedium}, false))
	validationWAAPMode := validation.ToDiagFunc(
		validation.StringInSlice([]string{waapModeDisabled, waapModeLearn, waapModePrevent, waapModePractice}, false))
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
			"visibility": {
				Type:             schema.TypeString,
				Description:      "The visibility of the resource, Shared or Local",
				Default:          "Shared",
				Optional:         true,
				ValidateDiagFunc: validationVisibility,
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
							Description:      "The performance impact: VeryLow, LowOrLower, MediumOrLower or HighOrLower",
							Default:          "MediumOrLower",
							Optional:         true,
							ValidateDiagFunc: validationPerformanceImpact,
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
							ValidateDiagFunc: validationMinimumSeverity,
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
										ValidateDiagFunc: validationWAAPMode,
									},
									"open_redirect": {
										Type:             schema.TypeString,
										Description:      "Prevent, AccordingToPractice, Disabled or Learn",
										Default:          "Disabled",
										Optional:         true,
										ValidateDiagFunc: validationWAAPMode,
									},
									"error_disclosure": {
										Type:             schema.TypeString,
										Description:      "Prevent, AccordingToPractice, Disabled or Learn",
										Default:          "Disabled",
										Optional:         true,
										ValidateDiagFunc: validationWAAPMode,
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
			"file_security": {
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
						"severity_level": {
							Type:             schema.TypeString,
							Description:      "LowOrAbove, MediumOrAbove, HighOrAbove or Critical",
							Default:          "MediumOrAbove",
							Optional:         true,
							ValidateDiagFunc: validationSeverityLevel,
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
						"allow_file_size_limit": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "AccordingToPractice",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
						"file_size_limit": {
							Type:     schema.TypeInt,
							Default:  10,
							Optional: true,
						},
						"file_size_limit_unit": {
							Type:             schema.TypeString,
							Description:      "Bytes, KB, MB or GB",
							Default:          "MB",
							Optional:         true,
							ValidateDiagFunc: validationFileSizeUnits,
						},
						"files_without_name": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "AccordingToPractice",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
						"required_archive_extraction": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"archive_file_size_limit": {
							Type:     schema.TypeInt,
							Default:  10,
							Optional: true,
						},
						"archive_file_size_limit_unit": {
							Type:             schema.TypeString,
							Description:      "Bytes, KB, MB or GB",
							Default:          "MB",
							Optional:         true,
							ValidateDiagFunc: validationFileSizeUnits,
						},
						"allow_archive_within_archive": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "AccordingToPractice",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
						"allow_an_unopened_archive": {
							Type:             schema.TypeString,
							Description:      "Detect, Prevent, Inactive or AccordingToPractice",
							Default:          "AccordingToPractice",
							Optional:         true,
							ValidateDiagFunc: validationFileSecurityMode,
						},
						"allow_file_type": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"required_threat_emulation": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
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
		// Check if the error is due to the web app practice being used by other resources
		if err != nil && strings.Contains(err.Error(), errorMsgPointedObjects) {
			// Get the resources that are using the web app practice
			usedBy, err := webapppractice.UsedByWebApplicationPractice(ctx, c, d.Id())
			if err != nil {
				return utils.DiagError("unable to perform WebAppPractice Delete", err, diags)
			}

			if usedBy != nil || len(usedBy) > 0 {
				// Remove the web app practice from the resources that are using it
				if err2 := handleWebAppPracticeReferences(ctx, usedBy, c, d.Id()); err2 != nil {
					return err2
				}

				// Retry to delete the web app practice
				result, err := webapppractice.DeleteWebApplicationPractice(ctx, c, d.Id())
				if err != nil || !result {
					if _, discardErr := c.DiscardChanges(); discardErr != nil {
						diags = utils.DiagError("failed to discard changes", discardErr, diags)
					}

					return utils.DiagError("unable to perform WebAppPractice Delete", err, diags)
				}

			}

		} else {
			if _, discardErr := c.DiscardChanges(); discardErr != nil {
				diags = utils.DiagError("failed to discard changes", discardErr, diags)
			}

			return utils.DiagError("failed to Publish following WebAppPractice Delete", err, diags)
		}
	}

	d.SetId("")

	return diags
}

func handleWebAppPracticeReferences(ctx context.Context, usedBy models.DisplayObjects, c *api.Client, practiceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, usedByResource := range usedBy {
		if usedByResource.ObjectStatus == "Deleted" || usedByResource.Type == "Wrapper" {
			continue
		}

		switch usedByResource.SubType {
		case "WebApp":
			webAppAsset := webAppAssetModels.UpdateWebApplicationAssetInput{
				RemovePracticeWrappers: []string{practiceID},
			}

			updated, err := webappasset.UpdateWebApplicationAsset(ctx, c, usedByResource.ID, webAppAsset)
			if err != nil || !updated {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("failed to update usedByResource", err, diags)
			}

		default:
			return utils.DiagError("failed to update usedByResource", nil, diags)
		}

	}

	return nil
}
