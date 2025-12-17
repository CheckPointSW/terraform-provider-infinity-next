package models

type SchemaValidationInput struct {
	ID        string `json:"id,omitempty"`
	OASSchema string `json:"OasSchema,omitempty"`
}

type AdvancedSettingInput struct {
	ID                 string `json:"id,omitempty"`
	IllegalHttpMethods string `json:"illegalHttpMethods,omitempty"`
	BodySize           int    `json:"bodySize,omitempty"`
	URLSize            int    `json:"urlSize,omitempty"`
	HeaderSize         int    `json:"headerSize,omitempty"`
	MaxObjectDepth     int    `json:"maxObjectDepth,omitempty"`
}

type APIAttacksInput struct {
	ID              string               `json:"id,omitempty"`
	MinimumSeverity string               `json:"minimumSeverity,omitempty"`
	AdvancedSetting AdvancedSettingInput `json:"advancedSetting,omitempty"`
}

type IPSInput struct {
	ID                  string `json:"id,omitempty"`
	PerformanceImpact   string `json:"performanceImpact,omitempty"`
	SeverityLevel       string `json:"severityLevel,omitempty"`
	ProtectionsFromYear string `json:"protectionsFromYear,omitempty"`
	HighConfidence      string `json:"highConfidence,omitempty"`
	MediumConfidence    string `json:"mediumConfidence,omitempty"`
	LowConfidence       string `json:"lowConfidence,omitempty"`
}

type WebAPIFileSecurityInput struct {
	ID                        string `json:"id,omitempty"`
	SeverityLevel             string `json:"severityLevel,omitempty"`
	HighConfidence            string `json:"highConfidence,omitempty"`
	MediumConfidence          string `json:"mediumConfidence,omitempty"`
	LowConfidence             string `json:"lowConfidence,omitempty"`
	AllowFileSizeLimit        string `json:"allowFileSizeLimit,omitempty"`
	FileSizeLimit             *int   `json:"fileSizeLimit,omitempty"`
	FileSizeLimitUnit         string `json:"fileSizeLimitUnit,omitempty"`
	FilesWithoutName          string `json:"filesWithoutName,omitempty"`
	RequiredArchiveExtraction bool   `json:"requiredArchiveExtraction,omitempty"`
	ArchiveFileSizeLimit      *int   `json:"archiveFileSizeLimit,omitempty"`
	ArchiveFileSizeLimitUnit  string `json:"archiveFileSizeLimitUnit,omitempty"`
	AllowArchiveWithinArchive string `json:"allowArchiveWithinArchive,omitempty"`
	AllowAnUnopenedArchive    string `json:"allowAnUnopenedArchive,omitempty"`
	AllowFileType             bool   `json:"allowFileType,omitempty"`
	RequiredThreatEmulation   bool   `json:"requiredThreatEmulation,omitempty"`
}

// CreateWebAPIPracticeInput represents the api input for creating a web API practice
type CreateWebAPIPracticeInput struct {
	Name             string                  `json:"name"`
	Visibility       string                  `json:"visibility,omitempty"`
	IPS              IPSInput                `json:"IPS,omitempty"`
	APIAttacks       APIAttacksInput         `json:"APIAttacks,omitempty"`
	SchemaValidation SchemaValidationInput   `json:"SchemaValidation,omitempty"`
	FileSecurity     WebAPIFileSecurityInput `json:"FileSecurity,omitempty"`
}
