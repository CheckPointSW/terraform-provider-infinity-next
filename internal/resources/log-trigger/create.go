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
	var res models.CreateLogTriggerInput
	res.Name = d.Get("name").(string)
	res.Verbosity = d.Get("verbosity").(string)
	res.AccessControlAllowEvents = d.Get("access_control_allow_events").(bool)
	res.AccessControlDropEvents = d.Get("access_control_drop_events").(bool)
	res.ThreaPreventionDetectEvents = d.Get("threat_prevention_detect_events").(bool)
	res.ThreaPreventionPreventEvents = d.Get("threat_prevention_prevent_events").(bool)
	res.WebRequests = d.Get("web_requests").(bool)
	res.WebURLPath = d.Get("web_url_path").(bool)
	res.WebURLQuery = d.Get("web_url_query").(bool)
	res.WebHeaders = d.Get("web_headers").(bool)
	res.WebBody = d.Get("web_body").(bool)
	res.LogToCloud = d.Get("log_to_cloud").(bool)
	res.LogToAgent = d.Get("log_to_agent").(bool)
	res.ExtendLogging = d.Get("extend_logging").(bool)
	res.ExtendLoggingMinSeverity = d.Get("extend_logging_min_severity").(string)
	res.ResponseBody = d.Get("response_body").(bool)
	res.ResponseCode = d.Get("response_code").(bool)
	res.LogToSyslog = d.Get("log_to_syslog").(bool)
	res.SyslogIPAddress = d.Get("syslog_ip_address").(string)
	res.SyslogPort = d.Get("syslog_port").(int)
	res.LogToCEF = d.Get("log_to_cef").(bool)
	res.CEFIPAddress = d.Get("cef_ip_address").(string)
	res.CEFPort = d.Get("cef_port").(int)

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
					syslogPortNum
					logToCef
					cefIpAddress
					cefPortNum
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
