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
	taskStatus, err := waitForTaskCompletion(ctx, c, result.ID)
	if err != nil {
		return err
	}

	switch taskStatus {
	case taskStatusSucceeded:
		return nil
	case taskStatusFailed:
		return fmt.Errorf("enforce policy task %s failed", result.ID)
	default:
		return fmt.Errorf("enforce policy task %s done with unknown status %s", result.ID, taskStatus)
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

// waitForTaskCompletion polls the task status until completion or timeout
func waitForTaskCompletion(ctx context.Context, c *api.Client, taskID string) (string, error) {
	taskStatus := taskStatusInProgress

	errch := make(chan error, 1)
	statusch := make(chan string, 1)

	go func() {
		for taskStatus == taskStatusInProgress {
			status, err := getTaskStatus(ctx, c, taskID)
			if err != nil {
				errch <- err
				return
			}

			taskStatus = status
			if taskStatus != taskStatusInProgress {
				statusch <- taskStatus
				return
			}

			time.Sleep(pollInterval)
		}
	}()

	// Timeout for the polling routine to finish
	select {
	case err := <-errch:
		return "", err
	case status := <-statusch:
		return status, nil
	case <-time.After(enforceTimeout):
		return "", fmt.Errorf("enforce policy task did not finish after %v, quitting", enforceTimeout)
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// getTaskStatus queries the status of a task
func getTaskStatus(ctx context.Context, c *api.Client, taskID string) (string, error) {
	query := fmt.Sprintf(`query {getTask(id: "%s") {id status}}`, taskID)

	response, err := c.MakeGraphQLRequest(ctx, query, "getTask")
	if err != nil {
		return "", fmt.Errorf("failed to get task status: %w", err)
	}

	if response == nil {
		return "", fmt.Errorf("received nil response from getTask query")
	}

	responseMap, ok := response.(map[string]any)
	if !ok {
		return "", fmt.Errorf("unexpected response type: %T", response)
	}

	status, ok := responseMap["status"].(string)
	if !ok {
		return "", fmt.Errorf("status not found in response")
	}

	return status, nil
}
