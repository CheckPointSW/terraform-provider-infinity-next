package models

// FileWrapper represents the OASSchema field of the SchemaValidation field of the WebAPIPractice returned from the API
type FileWrapper struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

// SchemaValidation represents the SchemaValidation field of the WebAPIPractice returned from the API
type SchemaValidation struct {
	ID        string      `json:"id"`
	OASSchema FileWrapper `json:"OasSchema"`
}

// AdvancedSetting represents the AdvancedSetting field of the APIAttacks field of the WebAPIPractice returned from the API
type AdvancedSetting struct {
	ID                 string `json:"id"`
	IllegalHttpMethods string `json:"illegalHttpMethods,omitempty"`
	BodySize           int    `json:"bodySize,omitempty"`
	URLSize            int    `json:"urlSize,omitempty"`
	HeaderSize         int    `json:"headerSize,omitempty"`
	MaxObjectDepth     int    `json:"maxObjectDepth,omitempty"`
}

// APIAttacks represents the APIAttacks field of the WebAPIPractice returned from the API
type APIAttacks struct {
	ID              string          `json:"id"`
	MinimumSeverity string          `json:"minimumSeverity"`
	AdvancedSetting AdvancedSetting `json:"advancedSetting"`
}

// IPS represents the IPS field of the WebAPIPractice returned from the API
type IPS struct {
	ID                  string `json:"id"`
	PerformanceImpact   string `json:"performanceImpact"`
	SeverityLevel       string `json:"severityLevel"`
	ProtectionsFromYear string `json:"protectionsFromYear"`
	HighConfidence      string `json:"highConfidence"`
	MediumConfidence    string `json:"mediumConfidence"`
	LowConfidence       string `json:"lowConfidence"`
}

// WebAPIPractice represents the response from the API after creating the web API practice
type WebAPIPractice struct {
	ID               string           `json:"id"`
	IPS              IPS              `json:"IPS"`
	Name             string           `json:"name"`
	Category         string           `json:"category"`
	PracticeType     string           `json:"practiceType"`
	APIAttacks       APIAttacks       `json:"APIAttacks"`
	Default          bool             `json:"default"`
	SchemaValidation SchemaValidation `json:"SchemaValidation"`
}
