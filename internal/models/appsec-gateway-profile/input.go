package models

type AuthenticationInput struct {
	MaxNumberOfAgents int `json:"maxNumberOfAgents"`
}

type KeyValueInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpgradeTimeInput struct {
	ScheduleType string   `json:"scheduleType,omitempty"`
	Time         string   `json:"time,omitempty"`
	WeekDays     []string `json:"weekDays,omitempty"`
	Duration     int      `json:"duration,omitempty"`
	Days         []int    `json:"days,omitempty"`
}

type CreateCloudGuardAppSecGatewayProfileInput struct {
	UpgradeTime                    *UpgradeTimeInput   `json:"upgradeTime,omitempty"`
	Name                           string              `json:"name"`
	ProfileSubType                 string              `json:"profileSubType"`
	UpgradeMode                    string              `json:"upgradeMode,omitempty"`
	AdditionalSettings             []KeyValueInput     `json:"additionalSettings"`
	ReverseProxyAdditionalSettings []KeyValueInput     `json:"reverseProxyAdditionalSettings,omitempty"`
	ReverseProxyUpstreamTimeout    int                 `json:"reverseProxyUpstreamTimeout,omitempty"`
	Authentication                 AuthenticationInput `json:"authentication,omitempty"`
	CertificateType                string              `json:"certificateType,omitempty"`
	FailOpenInspection             bool                `json:"failOpenInspection,omitempty"`
}
