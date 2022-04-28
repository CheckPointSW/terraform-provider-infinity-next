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

// KubernetesProfile represents the profile object as it is returned from mgmt
type KubernetesProfile struct {
	ID                      string                      `json:"id"`
	Name                    string                      `json:"name"`
	ProfileType             string                      `json:"profileType"`
	ProfileSubType          string                      `json:"profileSubType"`
	Authentication          ReusableTokenAuthentication `json:"authentication,omitempty"`
	AdditionalSettings      []KeyValue                  `json:"additionalSettings"`
	OnlyDefinedApplications bool                        `json:"onlyDefinedApplications,omitempty"`
}
