package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	FileDataFilenameFormat = "%s;$$:$$;"
	FileDataFormat         = "data:application/octet-stream;base64,%s"
)

var (
	ErrorInvalidDataFormat = errors.New("invalid data format")
)

type SchemaAdvancedSetting struct {
	ID                 string `json:"id"`
	IllegalHttpMethods bool   `json:"illegal_http_methods,omitempty"`
	BodySize           int    `json:"body_size,omitempty"`
	URLSize            int    `json:"url_size,omitempty"`
	HeaderSize         int    `json:"header_size,omitempty"`
	MaxObjectDepth     int    `json:"max_object_depth,omitempty"`
}

type SchemaAPIAttacks struct {
	ID              string                  `json:"id"`
	MinimumSeverity string                  `json:"minimum_severity"`
	AdvancedSetting []SchemaAdvancedSetting `json:"advanced_setting"`
}

type SchemaIPS struct {
	ID                  string `json:"id"`
	PerformanceImpact   string `json:"performance_impact"`
	SeverityLevel       string `json:"severity_level"`
	ProtectionsFromYear string `json:"protections_from_year"`
	HighConfidence      string `json:"high_confidence"`
	MediumConfidence    string `json:"medium_confidence"`
	LowConfidence       string `json:"low_confidence"`
}

type FileSchema struct {
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	Data     string `json:"data"`
}

type OASSchema struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

type SchemaValidationSchema struct {
	ID        string    `json:"id"`
	OASSchema OASSchema `json:"OasSchema"`
}

func NewFileSchemaEncode(filename, fileData string) FileSchema {
	b64Data := base64.StdEncoding.EncodeToString([]byte(fileData))
	data := fmt.Sprintf(FileDataFormat, b64Data)
	filenameFmt := fmt.Sprintf(FileDataFilenameFormat, filepath.Base(filename))

	return FileSchema{
		Filename: filename,
		Data:     filenameFmt + data,
	}
}

func NewFileSchemaDecode(filename, b64Data string) (*FileSchema, error) {
	if _, bEncodedData, found := strings.Cut(b64Data, "base64,"); found {
		bDecodedData, err := base64.StdEncoding.DecodeString(bEncodedData)
		if err != nil {
			return nil, fmt.Errorf("failed decoding base64 string %s: %w", bEncodedData, err)
		}

		return &FileSchema{
			Filename: filename,
			Data:     string(bDecodedData),
		}, nil
	}

	return nil, ErrorInvalidDataFormat
}
