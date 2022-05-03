package models

type CreateTrustedSourceBehaviorInput struct {
	Name               string   `json:"name,omitempty"`
	Visibility         string   `json:"visibility,omitempty"`
	NumOfSources       int      `json:"numOfSources,omitempty"`
	SourcesIdentifiers []string `json:"sourcesIdentifiers,omitempty"`
}
