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
	ID                        string `json:"id"`
	SeverityLevel             string `json:"severity_level,omitempty"`
	HighConfidence            string `json:"high_confidence,omitempty"`
	MediumConfidence          string `json:"medium_confidence,omitempty"`
	LowConfidence             string `json:"low_confidence,omitempty"`
	AllowFileSizeLimit        string `json:"allow_file_size_limit,omitempty"`
	FileSizeLimit             int    `json:"file_size_limit,omitempty"`
	FileSizeLimitUnit         string `json:"file_size_limit_unit,omitempty"`
	FilesWithoutName          string `json:"files_without_name,omitempty"`
	RequiredArchiveExtraction bool   `json:"required_archive_extraction,omitempty"`
	ArchiveFileSizeLimit      int    `json:"archive_file_size_limit,omitempty"`
	ArchiveFileSizeLimitUnit  string `json:"archive_file_size_limit_unit,omitempty"`
	AllowArchiveWithinArchive string `json:"allow_archive_within_archive,omitempty"`
	AllowAnUnopenedArchive    string `json:"allow_an_unopened_archive,omitempty"`
	AllowFileType             bool   `json:"allow_file_type,omitempty"`
	RequiredThreatEmulation   bool   `json:"required_threat_emulation,omitempty"`
}

func (schemaIDs IDs) ToIndicatorsMap() map[string]string {
	ret := make(map[string]string)
	for _, id := range schemaIDs {
		uriAndID := strings.Split(id, WebApplicationURIIDSeparator)
		ret[uriAndID[0]] = uriAndID[1]
	}

	return ret
}
