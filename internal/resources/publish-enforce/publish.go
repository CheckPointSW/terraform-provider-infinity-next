package publishenforce

import (
	"context"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/publish-enforce"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ExecutePublish triggers a publish operation (same as `inext publish`)
func ExecutePublish(ctx context.Context, c *api.Client) error {
	result, err := PublishChanges(ctx, c)
	if err != nil {
		return err
	}

	if !result.IsValid {
		errorMsgs := make([]string, len(result.Errors))
		for i, e := range result.Errors {
			errorMsgs[i] = e.Message
		}
		return fmt.Errorf("publish failed with errors: %s", strings.Join(errorMsgs, ", "))
	}

	return nil
}

// ShouldPublishFromResourceData reads the publish value from ResourceData
// Returns true if publish should be executed
func ShouldPublishFromResourceData(d *schema.ResourceData) bool {
	return d.Get("publish").(bool)
}

// PublishChanges triggers a publish operation for the session
func PublishChanges(ctx context.Context, c *api.Client) (*models.PublishChangesResult, error) {
	query := `
		mutation publishChanges {
			publishChanges {
				isValid
				errors {
					message
				}
				warnings {
					message
				}
			}
		}
	`

	response, err := c.MakeGraphQLRequest(ctx, query, "publishChanges")
	if err != nil {
		return nil, fmt.Errorf("failed to execute publishChanges mutation: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("received nil response from publishChanges mutation")
	}

	responseMap, ok := response.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", response)
	}

	result := &models.PublishChangesResult{}

	if isValid, ok := responseMap["isValid"].(bool); ok {
		result.IsValid = isValid
	}

	if errors, ok := responseMap["errors"].([]any); ok {
		for _, e := range errors {
			if errMap, ok := e.(map[string]any); ok {
				if msg, ok := errMap["message"].(string); ok {
					result.Errors = append(result.Errors, models.ValidationMessage{Message: msg})
				}
			}
		}
	}

	if warnings, ok := responseMap["warnings"].([]any); ok {
		for _, w := range warnings {
			if warnMap, ok := w.(map[string]any); ok {
				if msg, ok := warnMap["message"].(string); ok {
					result.Warnings = append(result.Warnings, models.ValidationMessage{Message: msg})
				}
			}
		}
	}

	return result, nil
}
