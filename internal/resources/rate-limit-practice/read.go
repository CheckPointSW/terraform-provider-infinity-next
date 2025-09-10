package ratelimitpractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/rate-limit-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadRateLimitPracticeToResourceData(practice models.RateLimitPractice, d *schema.ResourceData) error {
	d.SetId(practice.ID)
	d.Set("name", practice.Name)
	d.Set("practice_type", practice.PracticeType)
	d.Set("visibility", practice.Visibility)
	d.Set("category", practice.Category)
	d.Set("default", practice.Default)
	schemaRules := practice.Rules.ToSchema()
	schemaRulesMaps, err := utils.UnmarshalAs[[]map[string]any](schemaRules)
	if err != nil {
		return fmt.Errorf("failed to convert rate limit rule structs %+v to map. Error: %w", schemaRules, err)
	}

	d.Set("rule", schemaRulesMaps)

	return nil
}

func GetRateLimitPractice(ctx context.Context, c *api.Client, id string, mustFind bool) (models.RateLimitPractice, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getRateLimitPractice(id: "`+id+`") {
				id
				name
				practiceType
				visibility
				category
				default
				rules {
					id
					URI
					scope
					limit
					comment
					action
				}
			}
		}
	`, "getRateLimitPractice")

	if err != nil {
		return models.RateLimitPractice{}, fmt.Errorf("failed to get RateLimitPractice: %w", err)
	}

	if res == nil {
		if mustFind {
			return models.RateLimitPractice{}, api.ErrorNotFound
		}

		return models.RateLimitPractice{}, nil
	}

	practice, err := utils.UnmarshalAs[models.RateLimitPractice](res)
	if err != nil {
		return models.RateLimitPractice{}, fmt.Errorf("failed to convert response to RateLimitPractice struct. Error: %w", err)
	}

	return practice, nil
}
