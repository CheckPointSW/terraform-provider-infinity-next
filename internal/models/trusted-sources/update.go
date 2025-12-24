package models

type TrustedSourceUpdateInput struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type UpdateTrustedSourceBehaviorInput struct {
	Name                        string                     `json:"name,omitempty"`
	Visibility                  string                     `json:"visibility,omitempty"`
	NumOfSources                int                        `json:"numOfSources"`
	AddSourcesIdentifiers       []string                   `json:"addSourcesIdentifiers,omitempty"`
	RemoveSourcesIdentifiersIDs []string                   `json:"removeSourcesIdentifiers,omitempty"`
	UpdateSourcesIdentifiers    []TrustedSourceUpdateInput `json:"updateSourcesIdentifiers,omitempty"`
}

type DisplayObject struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	SubType      string `json:"subType,omitempty"`
	ObjectStatus string `json:"objectStatus,omitempty"`
}

type DisplayObjects []DisplayObject
