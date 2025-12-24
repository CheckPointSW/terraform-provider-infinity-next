package resources

import (
	"context"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/log-trigger"
	logtrigger "github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources/log-trigger"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceLogTrigger() *schema.Resource {
	return &schema.Resource{
		Description: "Granular log setting and destination of logging",

		CreateContext: resourceLogTriggerCreate,
		ReadContext:   resourceLogTriggerRead,
		UpdateContext: resourceLogTriggerUpdate,
		DeleteContext: resourceLogTriggerDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "The name of the resource, also acts as its unique ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"verbosity": {
				Description:  "The verbosity of the log: Standard, Minimal or Extended",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Standard",
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Minimal", "Extended"}, false),
			},
			"compliance_warnings": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"compliance_violations": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"access_control_allow_events": {
				Description: "Log Access Control accepts",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"access_control_drop_events": {
				Description: "Log Access Control drops",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"threat_prevention_detect_events": {
				Description: "Log Threat Prevention Prevents",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"threat_prevention_prevent_events": {
				Description: "Log Threat Prevention Detects",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"web_requests": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"web_url_path": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"web_url_query": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"web_headers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"web_body": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"log_to_cloud": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"log_to_agent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"extend_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"extend_logging_min_severity": {
				Description:  "Minimum severity of events that will trigger extended logging: High or Critical",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "High",
				ValidateFunc: validation.StringInSlice([]string{"High", "Critical"}, false),
			},
			"response_body": {
				Type:        schema.TypeBool,
				Description: "Add response body to log if true",
				Optional:    true,
				Default:     false,
			},
			"response_code": {
				Type:        schema.TypeBool,
				Description: "Add response code to log if true",
				Optional:    true,
				Default:     true,
			},
			"log_to_syslog": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"syslog_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"syslog_protocol": {
				Description:  "Syslog protocol: UDP or TCP",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "UDP",
				ValidateFunc: validation.StringInSlice([]string{"UDP", "TCP"}, false),
			},
			"syslog_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"log_to_cef": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cef_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cef_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cef_protocol": {
				Description:  "CEF protocol: UDP or TCP",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "UDP",
				ValidateFunc: validation.StringInSlice([]string{"UDP", "TCP"}, false),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLogTriggerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)

	input, err := logtrigger.CreateLogTriggerInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("Failed to create log trigger input struct from resource data", err, diags)
	}

	logTrigger, err := logtrigger.NewLogTrigger(ctx, c, input)
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger Create", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Failed to Publish following LogTrigger Create", err, diags)
	}

	if err := logtrigger.ReadLogTriggerToResourceData(logTrigger, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger Read after Create", err, diags)
	}

	return diags
}

func resourceLogTriggerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*api.Client)
	logTrigger, err := logtrigger.GetLogTrigger(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger Get before read", err, diags)
	}

	if err := logtrigger.ReadLogTriggerToResourceData(logTrigger, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger read to state file", err, diags)
	}

	return diags
}

func resourceLogTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	updateInput, err := logtrigger.UpdateLogTriggerInputFromResourceData(d)
	if err != nil {
		return utils.DiagError("Unable to create log trigger update input struct from resource data", err, diags)
	}

	result, err := logtrigger.UpdateLogTrigger(ctx, c, d.Id(), updateInput)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger Update", err, diags)
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Failed to Publish following LogTrigger Update", err, diags)
	}

	logTrigger, err := logtrigger.GetLogTrigger(ctx, c, d.Id())
	if err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger Get before read after update", err, diags)
	}

	if err := logtrigger.ReadLogTriggerToResourceData(logTrigger, d); err != nil {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger read to state file after update", err, diags)
	}

	return diags
}

func resourceLogTriggerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)

	ID := d.Id()
	result, err := logtrigger.DeleteLogTrigger(ctx, c, ID)
	if err != nil || !result {
		// If the error is due to the log trigger being used by other objects, discard changes and return
		if err != nil && strings.Contains(err.Error(), errorMsgPointedObjects) {
			usedBy, err2 := logtrigger.UsedByLogTrigger(ctx, c, ID)
			if err2 != nil {
				diags = utils.DiagError("Unable to perform LogTrigger UsedBy", err2, diags)
				return utils.DiagError("Unable to perform LogTrigger Delete", err, diags)
			}

			// Update the practices that use the log trigger
			if err2 := handleLogTriggerReferences(ctx, usedBy, c, ID); err2 != nil {
				diags = err2
				return utils.DiagError("Unable to perform LogTrigger Delete", err, diags)
			}

			// Retry the delete operation
			result, err = logtrigger.DeleteLogTrigger(ctx, c, ID)
			if err != nil || !result {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("Unable to perform LogTrigger Delete after updating references", err, diags)
			}
		} else {
			if _, discardErr := c.DiscardChanges(); discardErr != nil {
				diags = utils.DiagError("failed to discard changes", discardErr, diags)
			}

			return utils.DiagError("Unable to perform LogTrigger Delete", err, diags)
		}
	}

	isValid, err := c.PublishChanges()
	if err != nil || !isValid {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Failed to Publish following LogTrigger Delete", err, diags)
	}

	d.SetId("")
	return diags
}

func handleLogTriggerReferences(ctx context.Context, triggersUsedBy models.TriggersUsedBy, c *api.Client, triggerID string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, triggerUsedBy := range triggersUsedBy {
		for _, practice := range triggerUsedBy.Practices {
			result, err := logtrigger.UpdatePracticeTriggers(ctx, c, triggerID, practice, triggerUsedBy.Container)
			if err != nil || !result {
				if _, discardErr := c.DiscardChanges(); discardErr != nil {
					diags = utils.DiagError("failed to discard changes", discardErr, diags)
				}

				return utils.DiagError("Unable to perform UpdatePracticeTriggers", err, diags)
			}
		}

	}

	return nil
}
