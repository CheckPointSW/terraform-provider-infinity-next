package models

import "fmt"

const (
	WebApplicationURIIDSeparator = ";;;"
)

type WebApplicationAdvancedSetting struct {
	ID                 string `json:"id"`
	CSRFProtection     string `json:"CSRFProtection,omitempty"`
	OpenRedirect       string `json:"openRedirect,omitempty"`
	ErrorDisclosure    string `json:"errorDisclosure,omitempty"`
	IllegalHttpMethods string `json:"illegalHttpMethods,omitempty"`
	BodySize           int    `json:"bodySize,omitempty"`
	URLSize            int    `json:"urlSize,omitempty"`
	HeaderSize         int    `json:"headerSize,omitempty"`
	MaxObjectDepth     int    `json:"maxObjectDepth,omitempty"`
}

type WebApplicationWebAttacks struct {
	ID              string                        `json:"id"`
	MinimumSeverity string                        `json:"minimumSeverity"`
	AdvancedSetting WebApplicationAdvancedSetting `json:"advancedSetting"`
}

type WebApplicationURI struct {
	ID  string `json:"id"`
	URI string `json:"URI"`
}

type WebApplicationWebBot struct {
	ID         string              `json:"id"`
	InjectURIs []WebApplicationURI `json:"injectURIs"`
	ValidURIs  []WebApplicationURI `json:"validURIs"`
}

type WebApplicationIPS struct {
	ID                  string `json:"id"`
	PerformanceImpact   string `json:"performanceImpact"`
	SeverityLevel       string `json:"severityLevel"`
	ProtectionsFromYear string `json:"protectionsFromYear"`
	HighConfidence      string `json:"highConfidence"`
	MediumConfidence    string `json:"mediumConfidence"`
	LowConfidence       string `json:"lowConfidence"`
}

type FileSecurity struct {
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

type WebApplicationPractice struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Category     string                   `json:"category"`
	PracticeType string                   `json:"practiceType"`
	Visibility   string                   `json:"visibility"`
	IPS          WebApplicationIPS        `json:"IPS"`
	WebBot       WebApplicationWebBot     `json:"WebBot"`
	WebAttacks   WebApplicationWebAttacks `json:"WebAttacks"`
	FileSecurity FileSecurity             `json:"FileSecurity"`
	Default      bool                     `json:"default"`
}

func (uri *WebApplicationURI) CreateSchemaID() string {
	return fmt.Sprintf("%s%s%s", uri.URI, WebApplicationURIIDSeparator, uri.ID)
}
