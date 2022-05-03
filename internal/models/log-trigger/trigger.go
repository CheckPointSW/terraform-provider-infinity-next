package models

type LogTrigger struct {
	ID                           string `json:"id"`
	Name                         string `json:"name"`
	Verbosity                    string `json:"verbosity"`
	AccessControlAllowEvents     bool   `json:"acAllow"`
	AccessControlDropEvents      bool   `json:"acDrop"`
	ThreaPreventionDetectEvents  bool   `json:"tpDetect"`
	ThreaPreventionPreventEvents bool   `json:"tpPrevent"`
	WebRequests                  bool   `json:"webRequests"`
	WebURLPath                   bool   `json:"webUrlPath"`
	WebURLQuery                  bool   `json:"webUrlQuery"`
	WebHeaders                   bool   `json:"webHeaders"`
	WebBody                      bool   `json:"webBody"`
	LogToCloud                   bool   `json:"logToCloud"`
	LogToAgent                   bool   `json:"logToAgent"`
	ExtendLogging                bool   `json:"extendLogging"`
	ExtendLoggingMinSeverity     string `json:"extendLoggingMinSeverity,omitempty"`
	ResponseBody                 bool   `json:"responseBody"`
	ResponseCode                 bool   `json:"responseCode"`
	LogToSyslog                  bool   `json:"logToSyslog"`
	SyslogIPAddress              string `json:"syslogIpAddress,omitempty"`
	SyslogPort                   int    `json:"SyslogPortNum,omitempty"`
	LogToCEF                     bool   `json:"LogToCef"`
	CEFIPAddress                 string `json:"CefIpAddress,omitempty"`
	CEFPort                      int    `json:"CefPort,omitempty"`
}
