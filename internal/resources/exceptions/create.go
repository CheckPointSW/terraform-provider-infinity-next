package exceptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/exceptions"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ParseSchemaMatchToInput(matchFromSchema models.SchemaMatchExpression) models.Match {
	var ret models.Match
	// commit
	// this is a condition match
	if len(matchFromSchema.Operands) == 0 {
		ret.Type = "condition"
		ret.Operator = matchFromSchema.Operator

		// default condition operator is "equals"
		if ret.Operator == "" {
			ret.Operator = "equals"
		}

		ret.Key = matchFromSchema.Key
		ret.Value = matchFromSchema.Value

		return ret
	}

	// this is an operator match
	ret.Type = "operator"
	ret.Operator = matchFromSchema.Operator

	// default condition operator is "and"
	if ret.Operator == "" {
		ret.Operator = "and"
	}

	for _, operand := range matchFromSchema.Operands {
		ret.Items = append(ret.Items, ParseSchemaMatchToInput(operand))
	}

	return ret

}

func mapToSchemaMatchExpression(conditionMap map[string]any) models.SchemaMatchExpression {
	var ret models.SchemaMatchExpression
	ret.Operator = conditionMap["operator"].(string)
	ret.Key = conditionMap["key"].(string)
	ret.Value = utils.MustSchemaCollectionToSlice[string](conditionMap["value"])
	ret.Operands = utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](conditionMap["operand"]), mapToSchemaMatchExpression)

	return ret
}

func mapToExceptionObjectInput(exceptionMap map[string]any) models.ExceptionObjectInput {
	var ret models.ExceptionObjectInput
	ret.Comment = exceptionMap["comment"].(string)
	action := models.Action{Key: "action", Value: exceptionMap["action"].(string)}
	actionBytes, err := json.Marshal(action)
	if err != nil {
		fmt.Printf("[WARN] failed to marshal action struct: %+v", err)
	}

	if id, ok := exceptionMap["id"]; ok {
		ret.ID = id.(string)
	}

	ret.Actions = []string{string(actionBytes)}
	matchExpression := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](exceptionMap["match"]), mapToSchemaMatchExpression)
	if len(matchExpression) > 0 {
		inputMatch := ParseSchemaMatchToInput(matchExpression[0])
		matchBytes, err := json.Marshal(inputMatch)
		if err != nil {
			fmt.Printf("[WARN] failed to marshal MatchExpression struct: %+v", err)
		}

		ret.Match = string(matchBytes)
	}

	return ret
}

func CreateExceptionBehaviorInputFromResourceData(d *schema.ResourceData) (models.CreateExceptionBehaviorInput, error) {
	var res models.CreateExceptionBehaviorInput

	res.Name = d.Get("name").(string)
	res.Visibility = d.Get("visibility").(string)
	res.Exceptions = utils.Map(utils.MustResourceDataCollectionToSlice[map[string]any](d, "exception"), mapToExceptionObjectInput)

	return res, nil
}

func NewExceptionBehavior(ctx context.Context, c *api.Client, input models.CreateExceptionBehaviorInput) (models.ExceptionBehavior, error) {
	vars := map[string]any{"ownerId": nil, "practiceId": nil, "behaviorInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newExceptionBehavior($ownerId: ID, $practiceId: ID, $behaviorInput: ExceptionBehaviorInput)
					{
						newExceptionBehavior(ownerId: $ownerId, practiceId: $practiceId, behaviorInput: $behaviorInput) {
							id
							name
							visibility
							exceptions {
								id
								match
								actions {
									id
									action
								}
								comment
							}
						}
					}
				`, "newExceptionBehavior", vars)

	if err != nil {
		return models.ExceptionBehavior{}, fmt.Errorf("failed to create new ExceptionBehavior: %w", err)
	}

	createRes, err := utils.UnmarshalAs[models.ExceptionBehavior](res)
	if err != nil {
		return models.ExceptionBehavior{}, fmt.Errorf("failed to convert response to ExceptionBehavior struct. Error: %w", err)
	}

	return createRes, err
}
