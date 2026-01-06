package publishenforce

import (
	"context"
	"fmt"
	"time"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/publish-enforce"
)

const (
	taskStatusInProgress = "InProgress"
	taskStatusSucceeded  = "Succeeded"
	taskStatusFailed     = "Failed"
	enforceTimeout       = 10 * time.Second
	pollInterval         = 200 * time.Millisecond
)

// ExecuteEnforce triggers an enforce operation and waits for completion (same as `inext enforce`)
func ExecuteEnforce(ctx context.Context, c *api.Client) error {
	result, err := EnforcePolicy(ctx, c)
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
func EnforcePolicy(ctx context.Context, c *api.Client) (*models.EnforcePolicyResult, error) {
	query := `mutation {enforcePolicy {id}}`

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
