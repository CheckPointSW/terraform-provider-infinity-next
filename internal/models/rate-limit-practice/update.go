package models

type AddRateLimitPracticeRule struct {
	URI     string `json:"URI"`
	Scope   string `json:"scope"`
	Limit   int    `json:"limit"`
	Comment string `json:"comment"`
	Action  string `json:"action"`
}

type UpdateRateLimitPracticeRule struct {
	ID      string `json:"id"`
	URI     string `json:"URI"`
	Scope   string `json:"scope"`
	Limit   int    `json:"limit"`
	Comment string `json:"comment"`
	Action  string `json:"action"`
}

type UpdateRateLimitPracticeInput struct {
	Name        string                        `json:"name,omitempty"`
	Visibility  string                        `json:"visibility,omitempty"`
	AddRules    []AddRateLimitPracticeRule    `json:"addRules,omitempty"`
	RemoveRules []string                      `json:"removeRules,omitempty"`
	UpdateRules []UpdateRateLimitPracticeRule `json:"updateRules,omitempty"`
}

type DisplayObject struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	SubType      string `json:"subType,omitempty"`
	ObjectStatus string `json:"objectStatus,omitempty"`
}

type DisplayObjects []DisplayObject
