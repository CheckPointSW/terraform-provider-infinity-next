package models

type UpdateSchemaValidationInput struct {
	ID        string `json:"id,omitempty"`
	OASSchema string `json:"OasSchema"`
}

type UpdateAPIAttacksInput struct {
	ID              string               `json:"id,omitempty"`
	MinimumSeverity string               `json:"minimumSeverity,omitempty"`
	AdvancedSetting AdvancedSettingInput `json:"advancedSetting,omitempty"`
}

type UpdateIPSInput struct {
	ID                  string `json:"id,omitempty"`
	PerformanceImpact   string `json:"performanceImpact,omitempty"`
	SeverityLevel       string `json:"severityLevel,omitempty"`
	ProtectionsFromYear string `json:"protectionsFromYear,omitempty"`
	HighConfidence      string `json:"highConfidence,omitempty"`
	MediumConfidence    string `json:"mediumConfidence,omitempty"`
	LowConfidence       string `json:"lowConfidence,omitempty"`
}

type UpdateWebAPIFileSecurityInput struct {
	ID                        string `json:"id,omitempty"`
	SeverityLevel             string `json:"severityLevel,omitempty"`
	HighConfidence            string `json:"highConfidence,omitempty"`
	MediumConfidence          string `json:"mediumConfidence,omitempty"`
	LowConfidence             string `json:"lowConfidence,omitempty"`
	AllowFileSizeLimit        string `json:"allowFileSizeLimit,omitempty"`
	FileSizeLimit             *int   `json:"fileSizeLimit,omitempty"`
	FileSizeLimitUnit         string `json:"fileSizeLimitUnit,omitempty"`
	FilesWithoutName          string `json:"filesWithoutName,omitempty"`
	RequiredArchiveExtraction *bool  `json:"requiredArchiveExtraction,omitempty"`
	ArchiveFileSizeLimit      *int   `json:"archiveFileSizeLimit,omitempty"`
	ArchiveFileSizeLimitUnit  string `json:"archiveFileSizeLimitUnit,omitempty"`
	AllowArchiveWithinArchive string `json:"allowArchiveWithinArchive,omitempty"`
	AllowAnUnopenedArchive    string `json:"allowAnUnopenedArchive,omitempty"`
	AllowFileType             *bool  `json:"allowFileType,omitempty"`
	RequiredThreatEmulation   *bool  `json:"requiredThreatEmulation,omitempty"`
}

type UpdatePracticeInput struct {
	Name             string                        `json:"name,omitempty"`
	Visibility       string                        `json:"visibility,omitempty"`
	IPS              UpdateIPSInput                `json:"IPS,omitempty"`
	APIAttacks       UpdateAPIAttacksInput         `json:"APIAttacks,omitempty"`
	SchemaValidation UpdateSchemaValidationInput   `json:"SchemaValidation,omitempty"`
	FileSecurity     UpdateWebAPIFileSecurityInput `json:"FileSecurity,omitempty"`
}

type DisplayObject struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	SubType      string `json:"subType,omitempty"`
	ObjectStatus string `json:"objectStatus,omitempty"`
}

type DisplayObjects []DisplayObject
