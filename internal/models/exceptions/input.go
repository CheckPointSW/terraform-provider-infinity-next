package models

type Action struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Match struct {
	Type     string   `json:"type"`
	Operator string   `json:"op"`
	Items    []Match  `json:"items,omitempty"`
	Key      string   `json:"key,omitempty"`
	Value    []string `json:"value,omitempty"`
}

type ExceptionObjectInput struct {
	ID      string   `json:"id,omitempty"`
	Match   string   `json:"match"`
	Actions []string `json:"actions"`
	Comment string   `json:"comment,omitempty"`
}

type ExceptionObjectInputs []ExceptionObjectInput

// CreateExceptionBehaviorInput represents the api input for creating an Exception behavior object
type CreateExceptionBehaviorInput struct {
	Name       string                `json:"name,omitempty"`
	Visibility string                `json:"visibility,omitempty"`
	Exceptions ExceptionObjectInputs `json:"exceptions,omitempty"`
}
