package models

type SchemaRateLimitPracticeRule struct {
	ID      string `json:"id"`
	URI     string `json:"uri"`
	Scope   string `json:"scope"`
	Limit   int    `json:"limit"`
	Comment string `json:"comment"`
	Action  string `json:"action"`
}

func (rule *SchemaRateLimitPracticeRule) GetUpdateRateLimitPracticeRule(newRule SchemaRateLimitPracticeRule) (UpdateRateLimitPracticeRule, bool) {
	var ret UpdateRateLimitPracticeRule
	ret.ID = rule.ID
	ret.URI = rule.URI
	isUpdate := false

	if rule.Scope != newRule.Scope {
		ret.Scope = newRule.Scope
		isUpdate = true
	} else {
		ret.Scope = rule.Scope
	}

	if rule.Limit != newRule.Limit {
		ret.Limit = newRule.Limit
		isUpdate = true
	} else {
		ret.Limit = rule.Limit
	}

	if rule.Comment != newRule.Comment {
		ret.Comment = newRule.Comment
		isUpdate = true
	} else {
		ret.Comment = rule.Comment
	}

	if rule.Action != newRule.Action {
		ret.Action = newRule.Action
		isUpdate = true
	} else {
		ret.Action = rule.Action
	}

	return ret, isUpdate

}
