package resources

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
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
				Description: "The name of the resource, also acts as it's unique ID",
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
			"access_control_allow_events": {
				Description: "Log Access Control accepts",
				Type:        schema.TypeBool,
				Optional:    true,
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
			},
			"web_body": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_to_cloud": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"log_to_agent": {
				Type:     schema.TypeBool,
				Optional: true,
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
			},
			"syslog_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"syslog_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"log_to_cef": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cef_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cef_port": {
				Type:     schema.TypeInt,
				Optional: true,
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

	logTrigger, err := logtrigger.NewLogTrigger(c, input)
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
	logTrigger, err := logtrigger.GetLogTrigger(c, d.Id())
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

	result, err := logtrigger.UpdateLogTrigger(c, d.Id(), updateInput)
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

	logTrigger, err := logtrigger.GetLogTrigger(c, d.Id())
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
	result, err := logtrigger.DeleteLogTrigger(c, ID)
	if err != nil || !result {
		if _, discardErr := c.DiscardChanges(); discardErr != nil {
			diags = utils.DiagError("failed to discard changes", discardErr, diags)
		}

		return utils.DiagError("Unable to perform LogTrigger Delete", err, diags)
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
