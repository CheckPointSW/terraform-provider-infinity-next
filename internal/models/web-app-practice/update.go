package models

type UpdateWebApplicationPracticeWebBotInput struct {
	ID                  string   `json:"id"`
	AddInjectURIs       []string `json:"addInjectURIs,omitempty"`
	RemoveInjectURIsIDs []string `json:"removeInjectURIs,omitempty"`
	AddValidURIs        []string `json:"addValidURIs,omitempty"`
	RemoveValidURIsIDs  []string `json:"removeValidURIs,omitempty"`
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

type UpdateWebApplicationPracticeInput struct {
	Name       string                                      `json:"name,omitempty"`
	IPS        UpdateWebApplicationPracticeIPSInput        `json:"IPS,omitempty"`
	WebAttacks UpdateWebApplicationPracticeWebAttacksInput `json:"WebAttacks,omitempty"`
	WebBot     UpdateWebApplicationPracticeWebBotInput     `json:"WebBot,omitempty"`
}
