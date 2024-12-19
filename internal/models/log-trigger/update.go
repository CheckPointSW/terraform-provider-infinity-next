package models

type UpdateLogTriggerInput struct {
	Name                          string `json:"name"`
	Verbosity                     string `json:"verbosity"`
	ComplianceWarnings            bool   `json:"complianceWarnings"`
	ComplianceViolations          bool   `json:"complianceViolations"`
	AccessControlAllowEvents      bool   `json:"acAllow"`
	AccessControlDropEvents       bool   `json:"acDrop"`
	ThreatPreventionDetectEvents  bool   `json:"tpDetect"`
	ThreatPreventionPreventEvents bool   `json:"tpPrevent"`
	WebRequests                   bool   `json:"webRequests"`
	WebURLPath                    bool   `json:"webUrlPath"`
	WebURLQuery                   bool   `json:"webUrlQuery"`
	WebHeaders                    bool   `json:"webHeaders"`
	WebBody                       bool   `json:"webBody"`
	LogToCloud                    bool   `json:"logToCloud"`
	LogToAgent                    bool   `json:"logToAgent"`
	ExtendLogging                 bool   `json:"extendLogging"`
	ExtendLoggingMinSeverity      string `json:"extendLoggingMinSeverity,omitempty"`
	ResponseBody                  bool   `json:"responseBody"`
	ResponseCode                  bool   `json:"responseCode"`
	LogToSyslog                   bool   `json:"logToSyslog"`
	SyslogIPAddress               string `json:"syslogIpAddress,omitempty"`
	SyslogProtocol                string `json:"syslogProtocol,omitempty"`
	SyslogPort                    int    `json:"syslogPortNum,omitempty"`
	LogToCEF                      bool   `json:"logToCef"`
	CEFIPAddress                  string `json:"cefIpAddress,omitempty"`
	CEFPort                       int    `json:"cefPortNum,omitempty"`
	CEFProtocol                   string `json:"cefProtocol,omitempty"`
}
