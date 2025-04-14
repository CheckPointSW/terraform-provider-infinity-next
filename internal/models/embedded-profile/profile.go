package models

type ScheduleTime struct {
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

type ReusableTokenAuthentication struct {
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

// EmbeddedProfile represents the profile object as it is returned from mgmt
type EmbeddedProfile struct {
	ID                      string                      `json:"id"`
	Name                    string                      `json:"name"`
	ProfileType             string                      `json:"profileType"`
	UpgradeMode             string                      `json:"upgradeMode,omitempty"`
	Authentication          ReusableTokenAuthentication `json:"authentication,omitempty"`
	AdditionalSettings      []KeyValue                  `json:"additionalSettings"`
	UsedBy                  DisplayObjects              `json:"usedBy"`
	UpgradeTime             *ScheduleTime               `json:"upgradeTime,omitempty"`
	OnlyDefinedApplications bool                        `json:"onlyDefinedApplications,omitempty"`
}
