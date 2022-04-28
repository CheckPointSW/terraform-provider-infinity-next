package models

import "fmt"

const (
	WebApplicationURIIDSeperator = ";;;"
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

type WebApplicationPractice struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Category     string                   `json:"category"`
	PracticeType string                   `json:"practiceType"`
	IPS          WebApplicationIPS        `json:"IPS"`
	WebBot       WebApplicationWebBot     `json:"WebBot"`
	WebAttacks   WebApplicationWebAttacks `json:"WebAttacks"`
	Default      bool                     `json:"default"`
}

func (uri *WebApplicationURI) CreateSchemaID() string {
	return fmt.Sprintf("%s%s%s", uri.URI, WebApplicationURIIDSeperator, uri.ID)
}
