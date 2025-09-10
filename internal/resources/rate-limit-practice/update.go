package ratelimitpractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/rate-limit-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseToSchemaRateLimitRule(schemaRuleAny any) []models.SchemaRateLimitPracticeRule {
	parseFunc := func(schemaRuleMap map[string]any) models.SchemaRateLimitPracticeRule {
		var ret models.SchemaRateLimitPracticeRule
		ret.ID = schemaRuleMap["id"].(string)
		ret.URI = schemaRuleMap["uri"].(string)
		ret.Scope = schemaRuleMap["scope"].(string)
		ret.Limit = schemaRuleMap["limit"].(int)
		ret.Comment = schemaRuleMap["comment"].(string)
		ret.Action = schemaRuleMap["action"].(string)
		return ret
	}

	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaRuleAny), parseFunc)
}

func UpdateRateLimitPracticeInputFromResourceData(d *schema.ResourceData) (models.UpdateRateLimitPracticeInput, error) {
	var updateInput models.UpdateRateLimitPracticeInput

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		updateInput.Name = newName
	}

	if _, newVisibility, hasChange := utils.MustGetChange[string](d, "visibility"); hasChange {
		updateInput.Visibility = newVisibility
	}

	if oldSchemaRulesSlice, newSchemaRulesSlice, hasChange := utils.GetChangeWithParse(d, "rule", parseToSchemaRateLimitRule); hasChange {
		var rulesIDsToRemove []string
		var rulesToAdd []models.AddRateLimitPracticeRule
		var rulesToUpdate []models.UpdateRateLimitPracticeRule

		// create map from uri to schema rule (since it's unique for each rule)
		oldRulesMap := make(map[string]models.SchemaRateLimitPracticeRule)
		for _, oldSchemaRule := range oldSchemaRulesSlice {
			oldRulesMap[oldSchemaRule.URI] = oldSchemaRule
		}

		newRulesMap := make(map[string]models.SchemaRateLimitPracticeRule)
		for _, newSchemaRule := range newSchemaRulesSlice {
			newRulesMap[newSchemaRule.URI] = newSchemaRule
		}

		// iterate over old rules to get rules to update/remove
		for _, oldSchemaRule := range oldSchemaRulesSlice {
			newRule, ok := newRulesMap[oldSchemaRule.URI]
			if !ok {
				// old rule should be deleted
				rulesIDsToRemove = append(rulesIDsToRemove, oldSchemaRule.ID)
				continue
			}

			// exists in both old and new, let's check if should be updated
			if updateRule, isUpdate := oldSchemaRule.GetUpdateRateLimitPracticeRule(newRule); isUpdate {
				rulesToUpdate = append(rulesToUpdate, updateRule)
			}
		}

		// iterate over new rules and get rules to add
		for _, newSchemaRule := range newSchemaRulesSlice {
			if _, ok := oldRulesMap[newSchemaRule.URI]; !ok {
				// new rule should be added
				rulesToAdd = append(rulesToAdd, models.AddRateLimitPracticeRule{
					URI:     newSchemaRule.URI,
					Scope:   newSchemaRule.Scope,
					Limit:   newSchemaRule.Limit,
					Comment: newSchemaRule.Comment,
					Action:  newSchemaRule.Action,
				})
			}
		}

		updateInput.RemoveRules = rulesIDsToRemove
		updateInput.AddRules = rulesToAdd
		updateInput.UpdateRules = rulesToUpdate
	}

	return updateInput, nil
}

func UpdateRateLimitPractice(ctx context.Context, c *api.Client, id string, input models.UpdateRateLimitPracticeInput) (bool, error) {
	vars := map[string]any{"practiceInput": input, "id": id, "ownerId": nil}

	res, err := c.MakeGraphQLRequest(ctx, `
		mutation updateRateLimitPractice($practiceInput: RateLimitPracticeUpdateInput, $id: ID!, $ownerId: ID)
		{
			updateRateLimitPractice(practiceInput: $practiceInput, id: $id, ownerId: $ownerId)
		}
	`, "updateRateLimitPractice", vars)

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateRateLimitPractice response %#v should be of type bool", res)
	}

	return value, err
}
