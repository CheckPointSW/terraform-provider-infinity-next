package models

const (
	SourceIdentifierValueIDSeparator = ";;;"
)

// SchemaPracticeMode represents a PracticeMode field of a practice field of a
// web application asset as it is saved in the state file
// this structure is aligned with the input schema (see web-api-asset.go file)
type SchemaPracticeMode struct {
	Mode        string `json:"mode"`
	SubPractice string `json:"sub_practice,omitempty"`
}

// SchemaPracticeWrapper represents a field of web application asset as it is saved in the state file
// this structure is aligned with the input schema (see web-api-asset.go file)
type SchemaPracticeWrapper struct {
	PracticeWrapperID string            `json:"practice_wrapper_id"`
	PracticeID        string            `json:"id"`
	MainMode          string            `json:"main_mode,omitempty"`
	SubPracticeModes  map[string]string `json:"sub_practices_modes,omitempty"`
	Triggers          []string          `json:"triggers,omitempty"`
}

// SchemaSourceIdentifier represents the SourceIdentifier field of a web application asset as it is saved in the state file
// this structure is aligned with the input schema (see web-app-asset.go file)
type SchemaSourceIdentifier struct {
	ID               string   `json:"id,omitempty"`
	SourceIdentifier string   `json:"identifier"`
	Values           []string `json:"values"`
	ValuesIDs        []string `json:"values_ids"`
}

type SchemaTag struct {
	ID    string `json:"id,omitempty"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
