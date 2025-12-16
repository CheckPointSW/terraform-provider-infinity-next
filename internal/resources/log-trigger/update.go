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

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		ret.Name = newName
	}

	if _, newVerbosity, hasChange := utils.MustGetChange[string](d, "verbosity"); hasChange {
		ret.Verbosity = newVerbosity
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "compliance_warnings"); hasChange {
		ret.ComplianceWarnings = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "compliance_violations"); hasChange {
		ret.ComplianceViolations = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "access_control_allow_events"); hasChange {
		ret.AccessControlAllowEvents = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "access_control_drop_events"); hasChange {
		ret.AccessControlDropEvents = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "threat_prevention_detect_events"); hasChange {
		ret.ThreatPreventionDetectEvents = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "threat_prevention_prevent_events"); hasChange {
		ret.ThreatPreventionPreventEvents = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "web_requests"); hasChange {
		ret.WebRequests = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "web_url_path"); hasChange {
		ret.WebURLPath = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "web_url_query"); hasChange {
		ret.WebURLQuery = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "web_headers"); hasChange {
		ret.WebHeaders = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "web_body"); hasChange {
		ret.WebBody = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "log_to_cloud"); hasChange {
		ret.LogToCloud = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "log_to_agent"); hasChange {
		ret.LogToAgent = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "extend_logging"); hasChange {
		ret.ExtendLogging = newVal
	}

	if _, extendLoggingMinSeverity, hasChange := utils.MustGetChange[string](d, "extend_logging_min_severity"); hasChange {
		ret.ExtendLoggingMinSeverity = extendLoggingMinSeverity
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "response_body"); hasChange {
		ret.ResponseBody = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "response_code"); hasChange {
		ret.ResponseCode = newVal
	}

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "log_to_syslog"); hasChange {
		ret.LogToSyslog = newVal
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

	if _, newVal, hasChange := utils.MustGetChange[bool](d, "log_to_cef"); hasChange {
		ret.LogToCEF = newVal
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
