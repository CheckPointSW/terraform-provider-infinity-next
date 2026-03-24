package publishenforce

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/publish-enforce"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	taskStatusInProgress = "InProgress"
	taskStatusSucceeded  = "Succeeded"
	taskStatusFailed     = "Failed"
	enforceTimeout       = 5 * time.Minute
	pollInterval         = 300 * time.Millisecond
	maxTransientErrors   = 5
)

// ShouldEnforceFromResourceData reads the enforce value from ResourceData
// Returns true if enforce should be executed
func ShouldEnforceFromResourceData(d *schema.ResourceData) bool {
	return d.Get("enforce").(bool)
}

// GetProfileIDsFromResourceData reads the profile_ids value from ResourceData
// Returns an empty slice if not set
func GetProfileIDsFromResourceData(d *schema.ResourceData) []string {
	_, newVal, _ := utils.MustGetChange[[]any](d, "profile_ids")
	if newVal == nil {
		return []string{}
	}

	return utils.MustSliceAs[string](newVal)
}

// ExecuteEnforce triggers an enforce operation and waits for completion (same as `inext enforce`)
// If profileIDs is empty, all profiles will be enforced; otherwise only the specified profiles
func ExecuteEnforce(ctx context.Context, c *api.Client, profileIDs []string) error {
	result, err := EnforcePolicy(ctx, c, profileIDs)
	if err != nil {
		return err
	}

	if result.ID == "" {
		return fmt.Errorf("enforce policy returned empty task ID")
	}

	// Poll for task completion
	taskResult, err := waitForTaskCompletion(ctx, c, result.ID)
	if err != nil {
		return err
	}

	switch taskResult.Status {
	case taskStatusSucceeded:
		return nil
	case taskStatusFailed:
		return fmt.Errorf("enforce policy task %s failed", result.ID)
	default:
		return fmt.Errorf("enforce policy task %s done with unknown status %s", result.ID, taskResult.Status)
	}
}

// EnforcePolicy triggers an enforce operation
// If profileIDs is empty, all profiles will be enforced; otherwise only the specified profiles
func EnforcePolicy(ctx context.Context, c *api.Client, profileIDs []string) (*models.EnforcePolicyResult, error) {
	var query string
	if len(profileIDs) == 0 {
		query = `mutation {enforcePolicy {id}}`
	} else {
		// Build the profilesIds array for the GraphQL query
		quotedIDs := make([]string, len(profileIDs))
		for i, id := range profileIDs {
			quotedIDs[i] = fmt.Sprintf(`"%s"`, id)
		}

		query = fmt.Sprintf(`mutation {enforcePolicy(profilesIds: [%s]) {id}}`, strings.Join(quotedIDs, ", "))
	}

	response, err := c.MakeGraphQLRequest(ctx, query, "enforcePolicy")
	if err != nil {
		return nil, fmt.Errorf("failed to execute enforcePolicy mutation: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("received nil response from enforcePolicy mutation")
	}

	responseMap, ok := response.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", response)
	}

	result := &models.EnforcePolicyResult{}

	if id, ok := responseMap["id"].(string); ok {
		result.ID = id
	}

	return result, nil
}

// waitForTaskCompletion polls the task until completion or timeout and returns the full task result
func waitForTaskCompletion(ctx context.Context, c *api.Client, taskID string) (*models.TaskResult, error) {
	errch := make(chan error, 1)
	resultch := make(chan *models.TaskResult, 1)

	go func() {
		consecutiveErrors := 0
		var lastErr error

		for {
			result, err := getTask(ctx, c, taskID)
			if err != nil {
				consecutiveErrors++
				lastErr = err
				// Retry on transient errors
				if consecutiveErrors >= maxTransientErrors {
					errch <- fmt.Errorf("failed after %d consecutive errors, last error: %w", consecutiveErrors, lastErr)
					return
				}

				time.Sleep(pollInterval * 3)
				continue
			}

			// Reset error counter on success
			consecutiveErrors = 0

			if result.Status != taskStatusInProgress {
				resultch <- result
				return
			}

			time.Sleep(pollInterval)
		}
	}()

	// Timeout for the polling routine to finish
	select {
	case err := <-errch:
		return &models.TaskResult{}, err
	case result := <-resultch:
		return result, nil
	case <-time.After(enforceTimeout):
		return &models.TaskResult{}, fmt.Errorf("enforce policy task did not finish after %v, quitting", enforceTimeout)
	case <-ctx.Done():
		return &models.TaskResult{}, ctx.Err()
	}
}

// getTask queries the full result of a task including publish validation data
func getTask(ctx context.Context, c *api.Client, taskID string) (*models.TaskResult, error) {
	query := fmt.Sprintf(`query {getTask(id: "%s") {id status taskData {publishData {isValid errors {message}}}}}`, taskID)

	response, err := c.MakeGraphQLRequest(ctx, query, "getTask")
	if err != nil {
		return &models.TaskResult{}, fmt.Errorf("failed to get task: %w", err)
	}

	if response == nil {
		return &models.TaskResult{}, fmt.Errorf("received nil response from getTask query")
	}

	responseMap, ok := response.(map[string]any)
	if !ok {
		return &models.TaskResult{}, fmt.Errorf("unexpected response type: %T", response)
	}

	status, ok := responseMap["status"].(string)
	if !ok {
		return &models.TaskResult{}, fmt.Errorf("status not found in response")
	}

	result := &models.TaskResult{
		Status: status,
	}

	if id, ok := responseMap["id"].(string); ok {
		result.ID = id
	}

	if taskData, ok := responseMap["taskData"].(map[string]any); ok {
		if publishData, ok := taskData["publishData"].(map[string]any); ok {
			pd := &models.TaskPublishData{}

			if isValid, ok := publishData["isValid"].(bool); ok {
				pd.IsValid = isValid
			}

			if errs, ok := publishData["errors"].([]any); ok {
				for _, e := range errs {
					if errMap, ok := e.(map[string]any); ok {
						if msg, ok := errMap["message"].(string); ok {
							pd.Errors = append(pd.Errors, models.ValidationMessage{Message: msg})
						}
					}
				}
			}

			result.TaskData = &models.TaskData{PublishData: pd}
		}
	}

	return result, nil
}
