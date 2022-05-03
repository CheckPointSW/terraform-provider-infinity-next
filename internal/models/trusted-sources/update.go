package models

type UpdateTrustedSourceBehaviorInput struct {
	Name                        string   `json:"name,omitempty"`
	NumOfSources                int      `json:"numOfSources,omitempty"`
	AddSourcesIdentifiers       []string `json:"addSourcesIdentifiers,omitempty"`
	RemoveSourcesIdentifiersIDs []string `json:"removeSourcesIdentifiers,omitempty"`
}
