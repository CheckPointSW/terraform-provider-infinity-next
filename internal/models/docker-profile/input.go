package models

type ReusableTokenAuthenticationInput struct {
	MaxNumberOfAgents *int `json:"maxNumberOfAgents,omitempty"`
}

type KeyValueInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateDockerProfileInput struct {
	Name                    string                           `json:"name"`
	AdditionalSettings      []KeyValueInput                  `json:"additionalSettings"`
	OnlyDefinedApplications *bool                            `json:"onlyDefinedApplications,omitempty"`
	Authentication          ReusableTokenAuthenticationInput `json:"authentication,omitempty"`
}
