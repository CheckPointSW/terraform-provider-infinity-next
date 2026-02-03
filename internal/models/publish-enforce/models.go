package models

// PublishChangesResult represents the result of a publish operation
type PublishChangesResult struct {
	IsValid  bool
	Errors   []ValidationMessage
	Warnings []ValidationMessage
}

// ValidationMessage represents a validation error or warning message
type ValidationMessage struct {
	Message string `json:"message"`
}

// EnforcePolicyResult represents the result of an enforce operation
type EnforcePolicyResult struct {
	ID string
}

// AsyncPublishResult represents the result of an async publish operation
type AsyncPublishResult struct {
	ID string
}
