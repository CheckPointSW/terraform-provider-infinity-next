package models

type ReusableTokenAuthenticationInput struct {
	MaxNumberOfAgents int `json:"maxNumberOfAgents"`
}

type KeyValueInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ScheduleTimeInput struct {
	ScheduleType string   `json:"scheduleType,omitempty"`
	Time         string   `json:"time,omitempty"`
	WeekDays     []string `json:"weekDays,omitempty"`
	Duration     int      `json:"duration,omitempty"`
}

type CreateEmbeddedProfileInput struct {
	UpgradeTime             *ScheduleTimeInput               `json:"upgradeTime,omitempty"`
	Name                    string                           `json:"name"`
	UpgradeMode             string                           `json:"upgradeMode,omitempty"`
	AdditionalSettings      []KeyValueInput                  `json:"additionalSettings"`
	OnlyDefinedApplications bool                             `json:"onlyDefinedApplications,omitempty"`
	Authentication          ReusableTokenAuthenticationInput `json:"authentication,omitempty"`
}
