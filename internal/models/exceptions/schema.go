package models

type SchemaMatchExpression struct {
	Operator string                  `json:"operator,omitempty"`
	Operands []SchemaMatchExpression `json:"operand,omitempty"`
	Key      string                  `json:"key,omitempty"`
	Value    []string                `json:"value,omitempty"`
}

type SchemaExceptionObject struct {
	ID       string                  `json:"id,omitempty"`
	Match    []SchemaMatchExpression `json:"match"`
	Action   string                  `json:"action"`
	ActionID string                  `json:"action_id,omitempty"`
	Comment  string                  `json:"comment,omitempty"`
}
