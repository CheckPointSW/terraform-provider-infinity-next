package models

import (
	"encoding/base64"
	"fmt"
)

const (
	SourceIdentifierValueIDSeparator = ";;;"
	FileDataFormat                   = "data:%s;base64,%s"

	mTLSFileTypePEM = ".pem"
	mTLSFileTypeCRT = ".crt"
	mTLSFileTypeDER = ".der"
	mTLSFileTypeP12 = ".p12"
	mTLSFileTypePFX = ".pfx"
	mTLSFileTypeP7B = ".p7b"
	mTLSFileTypeP7C = ".p7c"
	mTLSFileTypeCER = ".cer"

	instructionsBlockTypeJSON = ".json"
	instructionsBlockTypeYAML = ".yaml"

	mimeTypePEM  = "application/octet-stream"
	mimeTypeDER  = "application/x-x509-ca-cert"
	mimeTypeP12  = "application/x-pkcs12"
	mimeTypeP7B  = "application/x-pkcs7-certificates"
	mimeTypeP7C  = "application/pkcs7-mime"
	mimeTypeJSON = "application/json"
	mimeTypeYAML = "application/octet-stream"
)

// SchemaPracticeMode represents a PracticeMode field of a practice field of a web API asset as it is saved in the state file
// this structure is aligned with the input schema (see web-api-asset.go file)
type SchemaPracticeMode struct {
	Mode        string `json:"mode"`
	SubPractice string `json:"sub_practice,omitempty"`
}

// SchemaPracticeWrapper represents a field of web API asset as it is saved in the state file
// this structure is aligned with the input schema (see web-api-asset.go file)
type SchemaPracticeWrapper struct {
	PracticeWrapperID string            `json:"practice_wrapper_id"`
	PracticeID        string            `json:"id"`
	MainMode          string            `json:"main_mode,omitempty"`
	SubPracticeModes  map[string]string `json:"sub_practices_modes,omitempty"`
	Triggers          []string          `json:"triggers,omitempty"`
}

// SchemaSourceIdentifier represents the SourceIdentifier field of a web APi asset as it is saved in the state file
// this structure is aligned with the input schema (see web-api-asset.go file)
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

// MTLSSchema represents a field of web API asset as it is saved in the state file
// this structure is aligned with the input schema (see web-api-asset.go file)
type MTLSSchema struct {
	FilenameID      string `json:"filename_id,omitempty"`
	Filename        string `json:"filename,omitempty"`
	CertificateType string `json:"certificate_type,omitempty"`
	DataID          string `json:"data_id,omitempty"`
	Data            string `json:"data"`
	Type            string `json:"type,omitempty"`
	EnableID        string `json:"enable_id,omitempty"`
	Enable          bool   `json:"enable,omitempty"`
}

type MTLSSchemas []MTLSSchema

// FileExtensionToMimeType returns the MIME type for a given file extension
// if the extension is not recognized, it returns "application/octet-stream" - a generic binary file MIME type
// the function is used to set the MIME type for the certificate type in the MTLSSchema
// the certificate types that are allowed displayed in the web-api-asset.go file
func FileExtensionToMimeType(extension string) string {
	switch extension {
	case mTLSFileTypePEM:
		return mimeTypePEM
	case mTLSFileTypeDER, mTLSFileTypeCER, mTLSFileTypeCRT:
		return mimeTypeDER
	case mTLSFileTypeP12, mTLSFileTypePFX:
		return mimeTypeP12
	case mTLSFileTypeP7B:
		return mimeTypeP7B
	case mTLSFileTypeP7C:
		return mimeTypeP7C
	case instructionsBlockTypeJSON:
		return mimeTypeJSON
	case instructionsBlockTypeYAML:
		return mimeTypeYAML
	default:
		return mimeTypePEM
	}
}

// MimeTypeToFileExtension returns the file extension for a given MIME type
// the function is used to set the certificate type in the MTLSSchema
func MimeTypeToFileExtension(mimeType string, isMTLS bool) string {
	if !isMTLS {
		switch mimeType {
		case mimeTypeJSON:
			return instructionsBlockTypeJSON
		case mimeTypeYAML:
			return instructionsBlockTypeYAML
		default:
			return instructionsBlockTypeYAML
		}
	}

	switch mimeType {
	case mimeTypePEM:
		return mTLSFileTypePEM
	case mimeTypeDER:
		return mTLSFileTypeDER
	case mimeTypeP12:
		return mTLSFileTypeP12
	case mimeTypeP7B:
		return mTLSFileTypeP7B
	case mimeTypeP7C:
		return mTLSFileTypeP7C
	default:
		return mTLSFileTypePEM
	}
}

func NewFileSchemaEncode(filename, fileData, mTLSType, certificateType string, fileEnable bool) MTLSSchema {
	b64Data := base64.StdEncoding.EncodeToString([]byte(fileData))
	data := fmt.Sprintf(FileDataFormat, FileExtensionToMimeType(certificateType), b64Data)

	return MTLSSchema{
		Filename: filename,
		Data:     data,
		Type:     mTLSType,
		Enable:   fileEnable,
	}
}

// BlockSchema represents a field of web application asset as it is saved in the state file
// this structure is aligned with the input schema (see web-app-asset.go file)
type BlockSchema struct {
	FilenameID   string `json:"filename_id,omitempty"`
	Filename     string `json:"filename,omitempty"`
	FilenameType string `json:"filename_type,omitempty"`
	DataID       string `json:"data_id,omitempty"`
	Data         string `json:"data"`
	Type         string `json:"type,omitempty"`
	EnableID     string `json:"enable_id,omitempty"`
	Enable       bool   `json:"enable,omitempty"`
}

type BlockSchemas []BlockSchema

func NewFileSchemaEncodeBlocks(filename, fileData, fileType, blockType string, fileEnable bool) BlockSchema {
	b64Data := base64.StdEncoding.EncodeToString([]byte(fileData))
	data := fmt.Sprintf(FileDataFormat, FileExtensionToMimeType(fileType), b64Data)
	return BlockSchema{
		Filename:     filename,
		Data:         data,
		FilenameType: fileType,
		Type:         blockType,
		Enable:       fileEnable,
	}
}

// CustomHeaderSchema represents a field of web application asset as it is saved in the state file
// this structure is aligned with the input schema (see web-app-asset.go file)
type CustomHeaderSchema struct {
	HeaderID string `json:"header_id,omitempty"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type CustomHeadersSchemas []CustomHeaderSchema
