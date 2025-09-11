package models

type RateLimitPracticeRuleInput struct {
	URI     string `json:"URI"`
	Scope   string `json:"scope"`
	Limit   int    `json:"limit"`
	Comment string `json:"comment"`
	Action  string `json:"action"`
}

type CreateRateLimitPracticeInput struct {
	Name       string                       `json:"name"`
	Visibility string                       `json:"visibility"`
	Rules      []RateLimitPracticeRuleInput `json:"rules"`
}
