package models

type TrustedSourceUpdateInput struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type UpdateTrustedSourceBehaviorInput struct {
	Name                        string                     `json:"name,omitempty"`
	Visibility                  string                     `json:"visibility,omitempty"`
	NumOfSources                int                        `json:"numOfSources,omitempty"`
	AddSourcesIdentifiers       []string                   `json:"addSourcesIdentifiers,omitempty"`
	RemoveSourcesIdentifiersIDs []string                   `json:"removeSourcesIdentifiers,omitempty"`
	UpdateSourcesIdentifiers    []TrustedSourceUpdateInput `json:"updateSourcesIdentifiers,omitempty"`
}
