package logtrigger

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/log-trigger"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateLogTriggerInputFromResourceData(d *schema.ResourceData) (models.CreateLogTriggerInput, error) {
	accessControlAllowEvents := d.Get("access_control_allow_events").(bool)
	accessControlDropEvents := d.Get("access_control_drop_events").(bool)
	threatPreventionDetectEvents := d.Get("threat_prevention_detect_events").(bool)
	threatPreventionPreventEvents := d.Get("threat_prevention_prevent_events").(bool)
	webRequests := d.Get("web_requests").(bool)
	webURLPath := d.Get("web_url_path").(bool)
	webURLQuery := d.Get("web_url_query").(bool)
	webHeaders := d.Get("web_headers").(bool)
	webBody := d.Get("web_body").(bool)
	logToCloud := d.Get("log_to_cloud").(bool)
	logToAgent := d.Get("log_to_agent").(bool)
	extendLogging := d.Get("extend_logging").(bool)
	responseBody := d.Get("response_body").(bool)
	responseCode := d.Get("response_code").(bool)
	logToSyslog := d.Get("log_to_syslog").(bool)
	syslogPort := d.Get("syslog_port").(int)
	logToCEF := d.Get("log_to_cef").(bool)
	cefPort := d.Get("cef_port").(int)
	complianceWarnings := d.Get("compliance_warnings").(bool)
	complianceViolations := d.Get("compliance_violations").(bool)

	var res models.CreateLogTriggerInput
	res.Name = d.Get("name").(string)
	res.Verbosity = d.Get("verbosity").(string)
	res.AccessControlAllowEvents = &accessControlAllowEvents
	res.AccessControlDropEvents = &accessControlDropEvents
	res.ThreatPreventionDetectEvents = &threatPreventionDetectEvents
	res.ThreatPreventionPreventEvents = &threatPreventionPreventEvents
	res.WebRequests = &webRequests
	res.WebURLPath = &webURLPath
	res.WebURLQuery = &webURLQuery
	res.WebHeaders = &webHeaders
	res.WebBody = &webBody
	res.LogToCloud = &logToCloud
	res.LogToAgent = &logToAgent
	res.ExtendLogging = &extendLogging
	res.ExtendLoggingMinSeverity = d.Get("extend_logging_min_severity").(string)
	res.ResponseBody = &responseBody
	res.ResponseCode = &responseCode
	res.LogToSyslog = &logToSyslog
	res.SyslogIPAddress = d.Get("syslog_ip_address").(string)
	res.SyslogProtocol = d.Get("syslog_protocol").(string)
	res.SyslogPort = &syslogPort
	res.LogToCEF = &logToCEF
	res.CEFIPAddress = d.Get("cef_ip_address").(string)
	res.CEFPort = &cefPort
	res.CEFProtocol = d.Get("cef_protocol").(string)
	res.ComplianceWarnings = &complianceWarnings
	res.ComplianceViolations = &complianceViolations

	return res, nil
}

func NewLogTrigger(ctx context.Context, c *api.Client, triggerInput models.CreateLogTriggerInput) (models.LogTrigger, error) {
	vars := map[string]any{"triggerInput": triggerInput}
	res, err := c.MakeGraphQLRequest(ctx, `
		mutation newLogTrigger($triggerInput: LogTriggerInput)
			{	
				newLogTrigger (triggerInput: $triggerInput) {
					id
					name
					verbosity
					complianceWarnings
					complianceViolations
					acAllow
					acDrop
					tpDetect
					tpPrevent
					webRequests
					webUrlPath
					webUrlQuery
					webHeaders
					webBody
					logToCloud
					logToAgent
					extendLogging
					extendLoggingMinSeverity
					responseBody
					responseCode
					logToSyslog
					syslogIpAddress
					syslogProtocol
					syslogPortNum
					logToCef
					cefIpAddress
					cefPortNum
					cefProtocol
				}
			}
		`, "newLogTrigger", vars)

	if err != nil {
		return models.LogTrigger{}, err
	}

	logTrigger, err := utils.UnmarshalAs[models.LogTrigger](res)
	if err != nil {
		return models.LogTrigger{}, fmt.Errorf("failed to cinvert response to LogTrigger struct. Error: %w", err)
	}

	return logTrigger, err
}
