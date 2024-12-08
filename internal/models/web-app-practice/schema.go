package models

import "strings"

type WebApplicationPracticeWebBotSchema struct {
	ID            string   `json:"id"`
	InjectURIs    []string `json:"inject_uris,omitempty"`
	InjectURIsIDs IDs      `json:"inject_uris_ids,omitempty"`
	ValidURIs     []string `json:"valid_uris,omitempty"`
	ValidURIsIDs  IDs      `json:"valid_uris_ids,omitempty"`
}

type IDs []string

type WebApplicationPracticeIPSSchema struct {
	ID                  string `json:"id"`
	PerformanceImpact   string `json:"performance_impact,omitempty"`
	SeverityLevel       string `json:"severity_level,omitempty"`
	ProtectionsFromYear string `json:"protections_from_year,omitempty"`
	HighConfidence      string `json:"high_confidence,omitempty"`
	MediumConfidence    string `json:"medium_confidence,omitempty"`
	LowConfidence       string `json:"low_confidence,omitempty"`
}

type WebApplicationPracticeAdvancedSettingSchema struct {
	ID                 string `json:"id"`
	CSRFProtection     string `json:"csrf_protection,omitempty"`
	OpenRedirect       string `json:"open_redirect,omitempty"`
	ErrorDisclosure    string `json:"error_disclosure,omitempty"`
	IllegalHttpMethods bool   `json:"illegal_http_methods,omitempty"`
	BodySize           int    `json:"body_size,omitempty"`
	URLSize            int    `json:"url_size,omitempty"`
	HeaderSize         int    `json:"header_size,omitempty"`
	MaxObjectDepth     int    `json:"max_object_depth,omitempty"`
}

type WebApplicationPracticeWebAttacksSchema struct {
	ID              string                                        `json:"id"`
	MinimumSeverity string                                        `json:"minimum_severity,omitempty"`
	AdvancedSetting []WebApplicationPracticeAdvancedSettingSchema `json:"advanced_setting,omitempty"`
}

type FileSecuritySchema struct {
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

func (schemaIDs IDs) ToIndicatorsMap() map[string]string {
	ret := make(map[string]string)
	for _, id := range schemaIDs {
		uriAndID := strings.Split(id, WebApplicationURIIDSeparator)
		ret[uriAndID[0]] = uriAndID[1]
	}

	return ret
}
