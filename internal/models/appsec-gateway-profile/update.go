package models

type UpdateUpgradeTimeInput struct {
	ScheduleType string   `json:"scheduleType,omitempty"`
	Time         string   `json:"time,omitempty"`
	WeekDays     []string `json:"weekDays,omitempty"`
	Duration     *int     `json:"duration,omitempty"`
	Days         []int    `json:"days,omitempty"`
}

type UpdateKeyValue struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateAuthenticationInput struct {
	MaxNumberOfAgents *int `json:"maxNumberOfAgents,omitempty"`
}

type UpdateCloudGuardAppSecGatewayProfileInput struct {
	UpgradeTime                          *UpdateUpgradeTimeInput    `json:"upgradeTime,omitempty"`
	Name                                 string                     `json:"name,omitempty"`
	ProfileSubType                       string                     `json:"profileSubType,omitempty"`
	UpgradeMode                          string                     `json:"upgradeMode,omitempty"`
	AddAdditionalSettings                []AddKeyValue              `json:"addAdditionalSettings,omitempty"`
	UpdateAdditionalSettings             []UpdateKeyValue           `json:"updateAdditionalSettings,omitempty"`
	RemoveAdditionalSettings             []string                   `json:"removeAdditionalSettings,omitempty"`
	AddReverseProxyAdditionalSettings    []AddKeyValue              `json:"addReverseProxyAdditionalSettings,omitempty"`
	UpdateReverseProxyAdditionalSettings []UpdateKeyValue           `json:"updateReverseProxyAdditionalSettings,omitempty"`
	RemoveReverseProxyAdditionalSettings []string                   `json:"removeReverseProxyAdditionalSettings,omitempty"`
	ReverseProxyUpstreamTimeout          *int                       `json:"reverseProxyUpstreamTimeout,omitempty"`
	Authentication                       *UpdateAuthenticationInput `json:"authentication,omitempty"`
	CertificateType                      string                     `json:"certificateType,omitempty"`
	FailOpenInspection                   *bool                      `json:"failOpenInspection,omitempty"`
}
