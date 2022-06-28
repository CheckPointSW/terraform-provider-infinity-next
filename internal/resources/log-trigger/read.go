package logtrigger

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/log-trigger"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadLogTriggerToResourceData(logTrigger models.LogTrigger, d *schema.ResourceData) error {
	d.SetId(logTrigger.ID)
	d.Set("name", logTrigger.Name)
	d.Set("verbosity", logTrigger.Verbosity)
	d.Set("access_control_allow_events", logTrigger.AccessControlAllowEvents)
	d.Set("access_control_drop_events", logTrigger.AccessControlDropEvents)
	d.Set("threat_prevention_detect_events", logTrigger.ThreaPreventionDetectEvents)
	d.Set("threat_prevention_prevent_events", logTrigger.ThreaPreventionPreventEvents)
	d.Set("web_requests", logTrigger.WebRequests)
	d.Set("web_url_path", logTrigger.WebURLPath)
	d.Set("web_url_query", logTrigger.WebURLQuery)
	d.Set("web_headers", logTrigger.WebHeaders)
	d.Set("web_body", logTrigger.WebBody)
	d.Set("log_to_cloud", logTrigger.LogToCloud)
	d.Set("log_to_agent", logTrigger.LogToAgent)
	d.Set("extend_logging", logTrigger.ExtendLogging)
	d.Set("extend_logging_min_severity", logTrigger.ExtendLoggingMinSeverity)
	d.Set("response_body", logTrigger.ResponseBody)
	d.Set("response_code", logTrigger.ResponseCode)
	d.Set("log_to_syslog", logTrigger.LogToSyslog)
	d.Set("syslog_ip_address", logTrigger.SyslogIPAddress)
	d.Set("syslog_port", logTrigger.SyslogPort)
	d.Set("log_to_cef", logTrigger.LogToCEF)
	d.Set("cef_ip_address", logTrigger.CEFIPAddress)
	d.Set("cef_port", logTrigger.CEFPort)

	return nil
}

func GetLogTrigger(ctx context.Context, c *api.Client, id string) (models.LogTrigger, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
	{
		getLogTrigger(id: "`+id+`") {
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
	`, "getLogTrigger")

	if err != nil {
		return models.LogTrigger{}, err
	}

	logTrigger, err := utils.UnmarshalAs[models.LogTrigger](res)
	if err != nil {
		return models.LogTrigger{}, fmt.Errorf("failed to cinvert response to LogTrigger struct. Error: %w", err)
	}

	return logTrigger, nil
}
