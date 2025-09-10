package ratelimitpractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/rate-limit-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func mapToRateLimitRule(ruleAsMap map[string]any) models.RateLimitPracticeRuleInput {
	return models.RateLimitPracticeRuleInput{
		URI:     ruleAsMap["uri"].(string),
		Scope:   ruleAsMap["scope"].(string),
		Limit:   ruleAsMap["limit"].(int),
		Comment: ruleAsMap["comment"].(string),
		Action:  ruleAsMap["action"].(string),
	}
}

func CreateRateLimitPracticeInputFromResourceData(d *schema.ResourceData) (models.CreateRateLimitPracticeInput, error) {
	var res models.CreateRateLimitPracticeInput

	res.Name = d.Get("name").(string)
	res.Visibility = d.Get("visibility").(string)
	rules := utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "rule"), mapToRateLimitRule)
	if len(rules) > 0 {
		res.Rules = rules
	}

	return res, nil
}

func NewRateLimitPracticePractice(ctx context.Context, c *api.Client, input models.CreateRateLimitPracticeInput) (models.RateLimitPractice, error) {
	vars := map[string]any{"ownerId": nil, "mainMode": nil, "subPracticeModes": nil, "practiceInput": input}

	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newRateLimitPractice($ownerId: ID, $mainMode: PracticeMode, $subPracticeModes: [PracticeModeInput], $practiceInput: RateLimitPracticeInput)
					{
						newRateLimitPractice(ownerId: $ownerId, subPracticeModes: $subPracticeModes, mainMode: $mainMode, practiceInput: $practiceInput) {
							id
							name
							practiceType
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
				`, "newRateLimitPractice", vars)

	if err != nil {
		return models.RateLimitPractice{}, fmt.Errorf("failed to create new RateLimitPractice: %w", err)
	}

	practice, err := utils.UnmarshalAs[models.RateLimitPractice](res)
	if err != nil {
		return models.RateLimitPractice{}, fmt.Errorf("failed to convert response to RateLimitPractice struct %w", err)
	}

	return practice, err
}
