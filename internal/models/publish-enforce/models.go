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

// TaskResult represents the full result of a getTask query
type TaskResult struct {
	ID       string
	Status   string
	TaskData *TaskData
}

// TaskData contains task-specific data returned from getTask
type TaskData struct {
	PublishData *TaskPublishData
}

// TaskPublishData holds the publish validation result nested inside TaskData
type TaskPublishData struct {
	IsValid bool
	Errors  []ValidationMessage
}
