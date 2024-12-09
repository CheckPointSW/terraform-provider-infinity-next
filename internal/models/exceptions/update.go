package models

type AddExceptionObjectInput struct {
	Match   string   `json:"match"`
	Actions []string `json:"actions"`
	Comment string   `json:"comment,omitempty"`
}

type UpdateExceptionObjectActionInput struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

type UpdateExceptionsObjectInputs []UpdateExceptionObjectActionInput

type ExceptionObjectActionUpdate struct {
	ID            string                       `json:"id,omitempty"`
	Match         string                       `json:"match,omitempty"`
	AddActions    []string                     `json:"addActions,omitempty"`
	RemoveActions []string                     `json:"removeActions,omitempty"`
	UpdateActions UpdateExceptionsObjectInputs `json:"updateActions,omitempty"`
	Comment       string                       `json:"comment,omitempty"`
}

type ExceptionObjectActionsUpdate []ExceptionObjectActionUpdate

type UpdateExceptionBehaviorInput struct {
	Name             string                       `json:"name,omitempty"`
	Visibility       string                       `json:"visibility,omitempty"`
	AddExceptions    []AddExceptionObjectInput    `json:"addExceptions,omitempty"`
	RemoveExceptions []string                     `json:"removeExceptions,omitempty"`
	UpdateExceptions ExceptionObjectActionsUpdate `json:"updateExceptions,omitempty"`
}
