package models

type AddExceptionObjectInput struct {
	Match   string   `json:"match"`
	Actions []string `json:"actions"`
	Comment string   `json:"comment,omitempty"`
}

type UpdateExceptionBehaviorInput struct {
	Name             string                    `json:"name,omitempty"`
	AddExceptions    []AddExceptionObjectInput `json:"addExceptions,omitempty"`
	RemoveExceptions []string                  `json:"removeExceptions,omitempty"`
}
