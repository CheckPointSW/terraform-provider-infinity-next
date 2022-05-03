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

type UpdatePracticeInput struct {
	Name             string                      `json:"name,omitempty"`
	Visibility       string                      `json:"visibility,omitempty"`
	IPS              UpdateIPSInput              `json:"IPS,omitempty"`
	APIAttacks       UpdateAPIAttacksInput       `json:"APIAttacks,omitempty"`
	SchemaValidation UpdateSchemaValidationInput `json:"SchemaValidation,omitempty"`
}
