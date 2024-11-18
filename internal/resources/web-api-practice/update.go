package webapipractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-api-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateWebAPIPracticeInputFromResourceData(d *schema.ResourceData) (models.UpdatePracticeInput, error) {
	var updateInput models.UpdatePracticeInput

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		updateInput.Name = newName
	}

	if oldIPSSlice, newIPSSlice, hasChange := utils.GetChangeWithParse(d, "ips", parseSchemaIPS); hasChange && len(newIPSSlice) > 0 {
		if len(oldIPSSlice) > 0 {
			newIPSSlice[0].ID = oldIPSSlice[0].ID
		}

		updateInput.IPS = newIPSSlice[0]
	}

	if oldAPIAttacks, newAPIAttacks, hasChange := utils.GetChangeWithParse(d, "api_attacks", parseSchemaAPIAttacks); hasChange && len(newAPIAttacks) > 0 {
		if len(oldAPIAttacks) > 0 {
			newAPIAttacks[0].ID = oldAPIAttacks[0].ID
		}

		updateInput.APIAttacks = newAPIAttacks[0]
	}

	if oldSchemaValidation, newSchemaValidation, hasChange := utils.GetChangeWithParse(d, "schema_validation", parseSchemaValidation); hasChange && len(newSchemaValidation) > 0 {
		if len(oldSchemaValidation) > 0 {
			newSchemaValidation[0].ID = oldSchemaValidation[0].ID
		}

		updateInput.SchemaValidation = newSchemaValidation[0]
	}

	return updateInput, nil
}

func UpdateWebAPIPractice(ctx context.Context, c *api.Client, id string, input models.UpdatePracticeInput) (bool, error) {
	vars := map[string]any{"practiceInput": input, "id": id, "ownerId": nil}

	res, err := c.MakeGraphQLRequest(ctx, `
		mutation updateWebAPIPractice($practiceInput: WebAPIPracticeUpdateInput, $id: ID!, $ownerId: ID)
		{
			updateWebAPIPractice(practiceInput: $practiceInput, id: $id, ownerId: $ownerId)
		}
	`, "updateWebAPIPractice", vars)

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateWebAPIPractice response %#v should be of type bool", res)
	}

	return value, err
}

func parseSchemaIPS(schemaIPS any) []models.UpdateIPSInput {
	input := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaIPS), mapToIPSInput)
	return utils.Map(input, utils.MustUnmarshalAs[models.UpdateIPSInput, models.IPSInput])
}

func parseSchemaAPIAttacks(schemaAPIAttacks any) []models.UpdateAPIAttacksInput {
	input := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaAPIAttacks), mapToAPIAttacksInput)
	return utils.Map(input, utils.MustUnmarshalAs[models.UpdateAPIAttacksInput, models.APIAttacksInput])
}

func parseSchemaValidation(validation any) []models.UpdateSchemaValidationInput {
	input := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](validation), mapToSchemaValidationInput)
	return utils.Map(input, utils.MustUnmarshalAs[models.UpdateSchemaValidationInput, models.SchemaValidationInput])
}
