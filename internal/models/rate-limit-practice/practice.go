package models

type RateLimitPracticeRule struct {
	ID      string `json:"id"`
	URI     string `json:"URI"`
	Scope   string `json:"scope"`
	Limit   int    `json:"limit"`
	Comment string `json:"comment"`
	Action  string `json:"action"`
}

type RateLimitPracticeRules []RateLimitPracticeRule

type RateLimitPractice struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	PracticeType string                 `json:"practiceType"`
	Visibility   string                 `json:"visibility,omitempty"`
	ObjectStatus string                 `json:"objectStatus"`
	Category     string                 `json:"category"`
	Default      bool                   `json:"default"`
	UsedBy       int                    `json:"usedBy"`
	Rules        RateLimitPracticeRules `json:"rules"`
}

func (rules RateLimitPracticeRules) ToSchema() []SchemaRateLimitPracticeRule {
	ret := make([]SchemaRateLimitPracticeRule, len(rules))
	for i, rule := range rules {
		ret[i] = SchemaRateLimitPracticeRule{
			ID:      rule.ID,
			URI:     rule.URI,
			Scope:   rule.Scope,
			Limit:   rule.Limit,
			Comment: rule.Comment,
			Action:  rule.Action,
		}
	}

	return ret
}
