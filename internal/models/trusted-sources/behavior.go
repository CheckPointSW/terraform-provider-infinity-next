package models

import "fmt"

const (
	TrustedSourceIDSeparator = ";;;"
)

type TrustedSourceSource struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type TrustedSourceBehavior struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name"`
	NumOfSources       int                   `json:"numOfSources"`
	SourcesIdentifiers []TrustedSourceSource `json:"sourcesIdentifiers,omitempty"`
}

func (trustedSource *TrustedSourceSource) CreateSchemaID() string {
	return fmt.Sprintf("%s%s%s", trustedSource.Source, TrustedSourceIDSeparator, trustedSource.ID)
}
