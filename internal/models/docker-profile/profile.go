package models

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

// DockerProfile represents the profile object as it is returned from mgmt
type DockerProfile struct {
	ID                      string                      `json:"id"`
	Name                    string                      `json:"name"`
	ProfileType             string                      `json:"profileType"`
	Authentication          ReusableTokenAuthentication `json:"authentication,omitempty"`
	AdditionalSettings      []KeyValue                  `json:"additionalSettings"`
	UsedBy                  DisplayObjects              `json:"usedBy"`
	OnlyDefinedApplications bool                        `json:"onlyDefinedApplications,omitempty"`
}
