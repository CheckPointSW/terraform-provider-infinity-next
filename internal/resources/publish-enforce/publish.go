package publishenforce

import (
	"context"
	"fmt"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/publish-enforce"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PublishOptions contains optional parameters for the async publish operation
type PublishOptions struct {
	ProfileTypes        []string
	SkipNginxValidation bool
}

// ExecutePublish triggers an async publish operation and waits for completion (same as `inext publish`)
func ExecutePublish(ctx context.Context, c *api.Client, opts *PublishOptions) error {
	result, err := AsyncPublishChanges(ctx, c, opts)
	if err != nil {
		return err
	}

	if result.ID == "" {
		return fmt.Errorf("async publish returned empty task ID")
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
		return fmt.Errorf("publish task %s failed", result.ID)
	default:
		return fmt.Errorf("publish task %s done with unknown status %s", result.ID, taskStatus)
	}
}

// ShouldPublishFromResourceData reads the publish value from ResourceData
// Returns true if publish should be executed
func ShouldPublishFromResourceData(d *schema.ResourceData) bool {
	return d.Get("publish").(bool)
}

// GetPublishOptionsFromResourceData reads publish options from ResourceData
func GetPublishOptionsFromResourceData(d *schema.ResourceData) *PublishOptions {
	opts := &PublishOptions{}

	if v, ok := d.GetOk("profile_types"); ok {
		profileTypes := v.([]any)
		opts.ProfileTypes = utils.MustSliceAs[string](profileTypes)
	}

	if v, ok := d.GetOk("skip_nginx_validation"); ok {
		opts.SkipNginxValidation = v.(bool)
	}

	return opts
}

// AsyncPublishChanges triggers an async publish operation for the session
func AsyncPublishChanges(ctx context.Context, c *api.Client, opts *PublishOptions) (*models.AsyncPublishResult, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`mutation {asyncPublishChanges`)

	// Build arguments if any options are provided
	var args []string
	if opts != nil {
		if len(opts.ProfileTypes) > 0 {
			quotedTypes := make([]string, len(opts.ProfileTypes))
			for i, pt := range opts.ProfileTypes {
				quotedTypes[i] = fmt.Sprintf(`"%s"`, pt)
			}

			args = append(args, fmt.Sprintf("profileTypes: [%s]", strings.Join(quotedTypes, ", ")))
		}

		if opts.SkipNginxValidation {
			args = append(args, "skipNginxValidation: true")
		}
	}

	if len(args) > 0 {
		queryBuilder.WriteString("(")
		queryBuilder.WriteString(strings.Join(args, ", "))
		queryBuilder.WriteString(")")
	}

	queryBuilder.WriteString("}")

	query := queryBuilder.String()

	response, err := c.MakeGraphQLRequest(ctx, query, "asyncPublishChanges")
	if err != nil {
		return nil, fmt.Errorf("failed to execute asyncPublishChanges mutation: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("received nil response from asyncPublishChanges mutation")
	}

	result := &models.AsyncPublishResult{}

	// asyncPublishChanges returns an ID directly, not an object
	if id, ok := response.(string); ok {
		result.ID = id
	} else {
		return nil, fmt.Errorf("unexpected response type: %T, expected string ID", response)
	}

	return result, nil
}
