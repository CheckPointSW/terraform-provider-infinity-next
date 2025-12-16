package models

type UpdateLogTriggerInput struct {
	Name                          string `json:"name"`
	Verbosity                     string `json:"verbosity"`
	ComplianceWarnings            *bool  `json:"complianceWarnings,omitempty"`
	ComplianceViolations          *bool  `json:"complianceViolations,omitempty"`
	AccessControlAllowEvents      *bool  `json:"acAllow,omitempty"`
	AccessControlDropEvents       *bool  `json:"acDrop,omitempty"`
	ThreatPreventionDetectEvents  *bool  `json:"tpDetect,omitempty"`
	ThreatPreventionPreventEvents *bool  `json:"tpPrevent,omitempty"`
	WebRequests                   *bool  `json:"webRequests,omitempty"`
	WebURLPath                    *bool  `json:"webUrlPath,omitempty"`
	WebURLQuery                   *bool  `json:"webUrlQuery,omitempty"`
	WebHeaders                    *bool  `json:"webHeaders,omitempty"`
	WebBody                       *bool  `json:"webBody,omitempty"`
	LogToCloud                    *bool  `json:"logToCloud,omitempty"`
	LogToAgent                    *bool  `json:"logToAgent,omitempty"`
	ExtendLogging                 *bool  `json:"extendLogging,omitempty"`
	ExtendLoggingMinSeverity      string `json:"extendLoggingMinSeverity,omitempty"`
	ResponseBody                  *bool  `json:"responseBody,omitempty"`
	ResponseCode                  *bool  `json:"responseCode,omitempty"`
	LogToSyslog                   *bool  `json:"logToSyslog,omitempty"`
	SyslogIPAddress               string `json:"syslogIpAddress,omitempty"`
	SyslogProtocol                string `json:"syslogProtocol,omitempty"`
	SyslogPort                    *int   `json:"syslogPortNum,omitempty"`
	LogToCEF                      *bool  `json:"logToCef,omitempty"`
	CEFIPAddress                  string `json:"cefIpAddress,omitempty"`
	CEFPort                       *int   `json:"cefPortNum,omitempty"`
	CEFProtocol                   string `json:"cefProtocol,omitempty"`
}

type TriggerUsedBy struct {
	Container string   `json:"container"`
	Practices []string `json:"practices"`
}

type TriggersUsedBy []TriggerUsedBy
