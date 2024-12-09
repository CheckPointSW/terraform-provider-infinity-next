package models

import (
	"encoding/base64"
	"fmt"
	"mime"
)

const (
	SourceIdentifierValueIDSeparator = ";;;"
	//FileDataFilenameFormat           = "%s;"
	FileDataFormat = "data:%s;base64,%s"
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

type FileSchema struct {
	FilenameID      string `json:"filename_id,omitempty"`
	Filename        string `json:"filename,omitempty"`
	CertificateType string `json:"certificate_type,omitempty"`
	DataID          string `json:"data_id,omitempty"`
	Data            string `json:"data"`
	Type            string `json:"type,omitempty"`
	EnableID        string `json:"enable_id,omitempty"`
	Enable          bool   `json:"enable,omitempty"`
}

type FileSchemas []FileSchema

func NewFileSchemaEncode(filename, fileData, mTLSType, certificateType string, fileEnable bool) FileSchema {
	b64Data := base64.StdEncoding.EncodeToString([]byte(fileData))
	data := fmt.Sprintf(FileDataFormat, mime.TypeByExtension(certificateType), b64Data)
	//filenameFmt := fmt.Sprintf(FileDataFilenameFormat, filepath.Base(filename))

	return FileSchema{
		Filename: filename,
		Data:     data,
		Type:     mTLSType,
		Enable:   fileEnable,
	}
}
