package models

type WebApplicationPracticeAdvancedSettingInput struct {
	ID                 string `json:"id,omitempty"`
	CSRFProtection     string `json:"CSRFProtection,omitempty"`
	OpenRedirect       string `json:"openRedirect,omitempty"`
	ErrorDisclosure    string `json:"errorDisclosure,omitempty"`
	IllegalHttpMethods string `json:"illegalHttpMethods,omitempty"`
	BodySize           int    `json:"bodySize,omitempty"`
	URLSize            int    `json:"urlSize,omitempty"`
	HeaderSize         int    `json:"headerSize,omitempty"`
	MaxObjectDepth     int    `json:"maxObjectDepth,omitempty"`
}

type WebApplicationPracticeWebAttacksInput struct {
	ID              string                                     `json:"id,omitempty"`
	MinimumSeverity string                                     `json:"minimumSeverity,omitempty"`
	AdvancedSetting WebApplicationPracticeAdvancedSettingInput `json:"advancedSetting,omitempty"`
}

type WebApplicationPracticeWebBotInput struct {
	ID         string   `json:"id,omitempty"`
	InjectURIs []string `json:"injectURIs,omitempty"`
	ValidURIs  []string `json:"validURIs,omitempty"`
}

type WebApplicationPractcieIPSInput struct {
	ID                  string `json:"id,omitempty"`
	PerformanceImpact   string `json:"performanceImpact,omitempty"`
	SeverityLevel       string `json:"severityLevel,omitempty"`
	ProtectionsFromYear string `json:"protectionsFromYear,omitempty"`
	HighConfidence      string `json:"highConfidence,omitempty"`
	MediumConfidence    string `json:"mediumConfidence,omitempty"`
	LowConfidence       string `json:"lowConfidence,omitempty"`
}

type CreateWebApplicationPracticeInput struct {
	Name       string                                `json:"name"`
	Visibility string                                `json:"visibility"`
	IPS        WebApplicationPractcieIPSInput        `json:"IPS,omitempty"`
	WebBot     WebApplicationPracticeWebBotInput     `json:"WebBot,omitempty"`
	WebAttacks WebApplicationPracticeWebAttacksInput `json:"WebAttacks,omitempty"`
}
