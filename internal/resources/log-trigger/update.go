package logtrigger

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/log-trigger"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateLogTriggerInputFromResourceData(d *schema.ResourceData) (models.UpdateLogTriggerInput, error) {
	var ret models.UpdateLogTriggerInput
	ret.Name = d.Get("name").(string)
	ret.Verbosity = d.Get("verbosity").(string)
	ret.ComplianceWarnings = d.Get("compliance_warnings").(bool)
	ret.ComplianceViolations = d.Get("compliance_violations").(bool)
	ret.AccessControlAllowEvents = d.Get("access_control_allow_events").(bool)
	ret.AccessControlDropEvents = d.Get("access_control_drop_events").(bool)
	ret.ThreatPreventionDetectEvents = d.Get("threat_prevention_detect_events").(bool)
	ret.ThreatPreventionPreventEvents = d.Get("threat_prevention_prevent_events").(bool)
	ret.WebRequests = d.Get("web_requests").(bool)
	ret.WebURLPath = d.Get("web_url_path").(bool)
	ret.WebURLQuery = d.Get("web_url_query").(bool)
	ret.WebHeaders = d.Get("web_headers").(bool)
	ret.WebBody = d.Get("web_body").(bool)
	ret.LogToCloud = d.Get("log_to_cloud").(bool)
	ret.LogToAgent = d.Get("log_to_agent").(bool)
	ret.ExtendLogging = d.Get("extend_logging").(bool)
	ret.ResponseBody = d.Get("response_body").(bool)
	ret.ResponseCode = d.Get("response_code").(bool)
	ret.LogToSyslog = d.Get("log_to_syslog").(bool)
	ret.LogToCEF = d.Get("log_to_cef").(bool)

	if _, extendLoggingMinSeverity, hasChange := utils.MustGetChange[string](d, "extend_logging_min_severity"); hasChange {
		ret.ExtendLoggingMinSeverity = extendLoggingMinSeverity
	}

	if _, syslogIPAddress, hasChange := utils.MustGetChange[string](d, "syslog_ip_address"); hasChange {
		ret.SyslogIPAddress = syslogIPAddress
	}

	if _, syslogProtocol, hasChange := utils.MustGetChange[string](d, "syslog_protocol"); hasChange {
		ret.SyslogProtocol = syslogProtocol
	}

	if _, syslogPortNum, hasChange := utils.MustGetChange[int](d, "syslog_port"); hasChange {
		ret.SyslogPort = syslogPortNum
	}

	if _, cefIPAddress, hasChange := utils.MustGetChange[string](d, "cef_ip_address"); hasChange {
		ret.CEFIPAddress = cefIPAddress
	}

	if _, cefPortNum, hasChange := utils.MustGetChange[int](d, "cef_port"); hasChange {
		ret.CEFPort = cefPortNum
	}

	if _, cefProtocol, hasChange := utils.MustGetChange[string](d, "cef_protocol"); hasChange {
		ret.CEFProtocol = cefProtocol
	}

	return ret, nil

}

func UpdateLogTrigger(ctx context.Context, c *api.Client, id any, triggerInput models.UpdateLogTriggerInput) (bool, error) {
	vars := map[string]any{"triggerInput": triggerInput, "id": id}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation updateLogTrigger($triggerInput: LogTriggerInput, $id: ID)
					{	
					updateLogTrigger (triggerInput: $triggerInput, id: $id) 
					}
				`, "updateLogTrigger", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid id %#v should be of type bool", res)
	}

	return isUpdated, err
}
