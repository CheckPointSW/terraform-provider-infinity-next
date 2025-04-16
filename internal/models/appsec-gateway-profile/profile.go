package models

type UpgradeTime struct {
	ScheduleType string   `json:"scheduleType,omitempty"`
	Time         string   `json:"time,omitempty"`
	WeekDays     []string `json:"weekDays,omitempty"`
	Duration     int      `json:"duration,omitempty"`
	Days         []int    `json:"days,omitempty"`
}

type KeyValue struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Authentication struct {
	Token             string `json:"token"`
	MaxNumberOfAgents int    `json:"maxNumberOfAgents"`
}

type DisplayObject struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	SubType      string `json:"subType,omitempty"`
	ObjectStatus string `json:"objectStatus,omitempty"`
}

type DisplayObjects []DisplayObject

// CloudGuardAppSecGatewayProfile represents the profile object as it is returned from mgmt
type CloudGuardAppSecGatewayProfile struct {
	ID                             string         `json:"id"`
	Name                           string         `json:"name"`
	ProfileType                    string         `json:"profileType"`
	ProfileSubType                 string         `json:"profileSubType"`
	UpgradeMode                    string         `json:"upgradeMode,omitempty"`
	Authentication                 Authentication `json:"authentication,omitempty"`
	AdditionalSettings             []KeyValue     `json:"additionalSettings"`
	UsedBy                         DisplayObjects `json:"usedBy"`
	ReverseProxyAdditionalSettings []KeyValue     `json:"reverseProxyAdditionalSettings,omitempty"`
	UpgradeTime                    *UpgradeTime   `json:"upgradeTime,omitempty"`
	ReverseProxyUpstreamTimeout    int            `json:"reverseProxyUpstreamTimeout,omitempty"`
	CertificateType                string         `json:"certificateType,omitempty"`
	FailOpenInspection             bool           `json:"failOpenInspection"`
}
