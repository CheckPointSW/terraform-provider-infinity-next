package models

type KeyValueUpdateInput struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KubernetesProfileUpdateInput struct {
	Name                     string                           `json:"name,omitempty"`
	ProfileSubType           string                           `json:"profileSubType,omitempty"`
	AddAdditionalSettings    []KeyValueInput                  `json:"addAdditionalSettings,omitempty"`
	UpdateAdditionalSettings []KeyValueUpdateInput            `json:"updateAdditionalSettings,omitempty"`
	RemoveAdditionalSettings []string                         `json:"removeAdditionalSettings,omitempty"`
	OnlyDefinedApplications  *bool                            `json:"onlyDefinedApplications,omitempty"`
	Authentication           ReusableTokenAuthenticationInput `json:"authentication,omitempty"`
}
