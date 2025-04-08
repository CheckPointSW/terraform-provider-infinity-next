package models

// FileWrapper represents the OASSchema field of the SchemaValidation field of the WebAPIPractice returned from the API
type FileWrapper struct {
	Data        string `json:"data"`
	Name        string `json:"name"`
	Size        uint64 `json:"size"`
	IsFileExist bool   `json:"isFileExist"`
}

// SchemaValidation represents the SchemaValidation field of the WebAPIPractice returned from the API
type SchemaValidation struct {
	ID        string      `json:"id"`
	OASSchema FileWrapper `json:"OasSchema"`
}

// AdvancedSetting represents the AdvancedSetting field of the APIAttacks field of the WebAPIPractice returned from the API
type AdvancedSetting struct {
	ID                 string `json:"id"`
	IllegalHttpMethods string `json:"illegalHttpMethods,omitempty"`
	BodySize           int    `json:"bodySize,omitempty"`
	URLSize            int    `json:"urlSize,omitempty"`
	HeaderSize         int    `json:"headerSize,omitempty"`
	MaxObjectDepth     int    `json:"maxObjectDepth,omitempty"`
}

// APIAttacks represents the APIAttacks field of the WebAPIPractice returned from the API
type APIAttacks struct {
	ID              string          `json:"id"`
	MinimumSeverity string          `json:"minimumSeverity"`
	AdvancedSetting AdvancedSetting `json:"advancedSetting"`
}

// IPS represents the IPS field of the WebAPIPractice returned from the API
type IPS struct {
	ID                  string `json:"id"`
	PerformanceImpact   string `json:"performanceImpact"`
	SeverityLevel       string `json:"severityLevel"`
	ProtectionsFromYear string `json:"protectionsFromYear"`
	HighConfidence      string `json:"highConfidence"`
	MediumConfidence    string `json:"mediumConfidence"`
	LowConfidence       string `json:"lowConfidence"`
}

type WebAPIFileSecurity struct {
	ID                        string `json:"id"`
	SeverityLevel             string `json:"severityLevel"`
	HighConfidence            string `json:"highConfidence"`
	MediumConfidence          string `json:"mediumConfidence"`
	LowConfidence             string `json:"lowConfidence"`
	AllowFileSizeLimit        string `json:"allowFileSizeLimit"`
	FileSizeLimit             int    `json:"fileSizeLimit"`
	FileSizeLimitUnit         string `json:"fileSizeLimitUnit"`
	FilesWithoutName          string `json:"filesWithoutName"`
	RequiredArchiveExtraction bool   `json:"requiredArchiveExtraction"`
	ArchiveFileSizeLimit      int    `json:"archiveFileSizeLimit"`
	ArchiveFileSizeLimitUnit  string `json:"archiveFileSizeLimitUnit"`
	AllowArchiveWithinArchive string `json:"allowArchiveWithinArchive"`
	AllowAnUnopenedArchive    string `json:"allowAnUnopenedArchive"`
	AllowFileType             bool   `json:"allowFileType"`
	RequiredThreatEmulation   bool   `json:"requiredThreatEmulation"`
}

// WebAPIPractice represents the response from the API after creating the web API practice
type WebAPIPractice struct {
	ID               string             `json:"id"`
	IPS              IPS                `json:"IPS"`
	Name             string             `json:"name"`
	Category         string             `json:"category"`
	PracticeType     string             `json:"practiceType"`
	Visibility       string             `json:"visibility"`
	APIAttacks       APIAttacks         `json:"APIAttacks"`
	Default          bool               `json:"default"`
	SchemaValidation SchemaValidation   `json:"SchemaValidation"`
	FileSecurity     WebAPIFileSecurity `json:"FileSecurity"`
}
