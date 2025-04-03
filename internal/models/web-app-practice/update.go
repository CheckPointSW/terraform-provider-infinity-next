package models

type UpdateURIInput struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

type UpdateURIsInputs []UpdateURIInput

type UpdateWebApplicationPracticeWebBotInput struct {
	ID                  string           `json:"id"`
	AddInjectURIs       []string         `json:"addInjectURIs,omitempty"`
	RemoveInjectURIsIDs []string         `json:"removeInjectURIs,omitempty"`
	UpdateInjectURIs    UpdateURIsInputs `json:"updateInjectURIs,omitempty"`
	AddValidURIs        []string         `json:"addValidURIs,omitempty"`
	RemoveValidURIsIDs  []string         `json:"removeValidURIs,omitempty"`
	UpdateValidURIs     UpdateURIsInputs `json:"updateValidURIs,omitempty"`
}

type UpdateWebApplicationPracticeAdvancedSettingInput struct {
	ID                 string `json:"id"`
	CSRFProtection     string `json:"CSRFProtection,omitempty"`
	OpenRedirect       string `json:"openRedirect,omitempty"`
	ErrorDisclosure    string `json:"errorDisclosure,omitempty"`
	BodySize           int    `json:"bodySize,omitempty"`
	URLSize            int    `json:"urlSize,omitempty"`
	HeaderSize         int    `json:"headerSize,omitempty"`
	MaxObjectDepth     int    `json:"maxObjectDepth,omitempty"`
	IllegalHttpMethods string `json:"illegalHttpMethods,omitempty"`
}

type UpdateWebApplicationPracticeWebAttacksInput struct {
	ID              string                                           `json:"id"`
	MinimumSeverity string                                           `json:"minimumSeverity,omitempty"`
	AdvancedSetting UpdateWebApplicationPracticeAdvancedSettingInput `json:"advancedSetting,omitempty"`
}

type UpdateWebApplicationPracticeIPSInput struct {
	ID                  string `json:"id"`
	PerformanceImpact   string `json:"performanceImpact,omitempty"`
	SeverityLevel       string `json:"severityLevel,omitempty"`
	ProtectionsFromYear string `json:"protectionsFromYear,omitempty"`
	HighConfidence      string `json:"highConfidence,omitempty"`
	MediumConfidence    string `json:"mediumConfidence,omitempty"`
	LowConfidence       string `json:"lowConfidence,omitempty"`
}

type UpdateFileSecurity struct {
	ID                        string `json:"id,omitempty"`
	SeverityLevel             string `json:"severityLevel,omitempty"`
	HighConfidence            string `json:"highConfidence,omitempty"`
	MediumConfidence          string `json:"mediumConfidence,omitempty"`
	LowConfidence             string `json:"lowConfidence,omitempty"`
	AllowFileSizeLimit        string `json:"allowFileSizeLimit,omitempty"`
	FileSizeLimit             int    `json:"fileSizeLimit,omitempty"`
	FileSizeLimitUnit         string `json:"fileSizeLimitUnit,omitempty"`
	FilesWithoutName          string `json:"filesWithoutName,omitempty"`
	RequiredArchiveExtraction bool   `json:"requiredArchiveExtraction,omitempty"`
	ArchiveFileSizeLimit      int    `json:"archiveFileSizeLimit,omitempty"`
	ArchiveFileSizeLimitUnit  string `json:"archiveFileSizeLimitUnit,omitempty"`
	AllowArchiveWithinArchive string `json:"allowArchiveWithinArchive,omitempty"`
	AllowAnUnopenedArchive    string `json:"allowAnUnopenedArchive,omitempty"`
	AllowFileType             bool   `json:"allowFileType,omitempty"`
	RequiredThreatEmulation   bool   `json:"requiredThreatEmulation,omitempty"`
}

type UpdateWebApplicationPracticeInput struct {
	Name         string                                      `json:"name,omitempty"`
	Visibility   string                                      `json:"visibility,omitempty"`
	IPS          UpdateWebApplicationPracticeIPSInput        `json:"IPS,omitempty"`
	WebAttacks   UpdateWebApplicationPracticeWebAttacksInput `json:"WebAttacks,omitempty"`
	WebBot       UpdateWebApplicationPracticeWebBotInput     `json:"WebBot,omitempty"`
	FileSecurity UpdateFileSecurity                          `json:"FileSecurity,omitempty"`
}

type DisplayObject struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	SubType      string `json:"subType,omitempty"`
	ObjectStatus string `json:"objectStatus,omitempty"`
}

type DisplayObjects []DisplayObject
