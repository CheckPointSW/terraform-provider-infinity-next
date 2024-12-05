package webapppractice

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/web-app-practice"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateWebApplicationPracticeInputFromResourceData(d *schema.ResourceData) (models.UpdateWebApplicationPracticeInput, error) {
	var updateInput models.UpdateWebApplicationPracticeInput

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		updateInput.Name = newName
	}

	if _, newVisibility, hasChange := utils.MustGetChange[string](d, "visibility"); hasChange {
		updateInput.Visibility = newVisibility
	}

	if oldIPSSlice, newIPSSlice, hasChange := utils.GetChangeWithParse(d, "ips", parseSchemaIPS); hasChange && len(newIPSSlice) > 0 {
		if len(oldIPSSlice) > 0 {
			newIPSSlice[0].ID = oldIPSSlice[0].ID
		}

		updateInput.IPS = newIPSSlice[0]
	}

	if oldWebAttacks, newWebAttacks, hasChange := utils.GetChangeWithParse(d, "web_attacks", parseSchemaWebAttacks); hasChange && len(newWebAttacks) > 0 {
		if len(oldWebAttacks) > 0 {
			newWebAttacks[0].ID = oldWebAttacks[0].ID
		}

		updateInput.WebAttacks = newWebAttacks[0]
	}

	if oldSchemaWebBotSlice, newSchemaWebBotSlice, hasChange := utils.GetChangeWithParse(d, "web_bot", parseToSchemaWebBot); hasChange && len(newSchemaWebBotSlice) > 0 {
		newSchemaWebBot := newSchemaWebBotSlice[0]
		if len(oldSchemaWebBotSlice) > 0 {
			oldSchemaWebBot := oldSchemaWebBotSlice[0]
			addedInjectURIs, removedInjectURIs := utils.SlicesDiff(oldSchemaWebBot.InjectURIs, newSchemaWebBot.InjectURIs)
			oldInjectURIsToIDsMap := oldSchemaWebBot.InjectURIsIDs.ToIndicatorsMap()
			var removedInjectURIsIDs []string
			for _, removedInjectURI := range removedInjectURIs {
				if removedID, ok := oldInjectURIsToIDsMap[removedInjectURI]; ok {
					removedInjectURIsIDs = append(removedInjectURIsIDs, removedID)
				}
			}

			addedValidURIs, removedValidURIs := utils.SlicesDiff(oldSchemaWebBot.ValidURIs, newSchemaWebBot.ValidURIs)
			oldValidURIsToIDsMap := oldSchemaWebBot.ValidURIsIDs.ToIndicatorsMap()
			var removedValidURIsIDs []string
			for _, removedValidURI := range removedValidURIs {
				if removedID, ok := oldValidURIsToIDsMap[removedValidURI]; ok {
					removedValidURIsIDs = append(removedValidURIsIDs, removedID)
				}
			}

			updateInput.WebBot = models.UpdateWebApplicationPracticeWebBotInput{
				ID:                  oldSchemaWebBot.ID,
				AddInjectURIs:       addedInjectURIs,
				RemoveInjectURIsIDs: removedInjectURIsIDs,
				UpdateInjectURIs:    models.UpdateURIsInputs{},
				AddValidURIs:        addedValidURIs,
				RemoveValidURIsIDs:  removedValidURIsIDs,
				UpdateValidURIs:     models.UpdateURIsInputs{},
			}

		} else {
			updateInput.WebBot = models.UpdateWebApplicationPracticeWebBotInput{
				AddInjectURIs: newSchemaWebBot.InjectURIs,
				AddValidURIs:  newSchemaWebBot.ValidURIs,
			}
		}
	}

	//if oldFileSecurity, newFileSecurity, hasChange := utils.GetChangeWithParse(d, "file_security", parseSchemaFileSecurity); hasChange && len(newFileSecurity) > 0 {
	//	if len(oldFileSecurity) > 0 {
	//		newFileSecurity[0].ID = oldFileSecurity[0].ID
	//	}
	//
	//	updateInput.FileSecurity = newFileSecurity[0]
	//}

	return updateInput, nil
}

func parseSchemaIPS(schemaIPS any) []models.UpdateWebApplicationPracticeIPSInput {
	input := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaIPS), mapToIPSInput)
	return utils.Map(input, utils.MustUnmarshalAs[models.UpdateWebApplicationPracticeIPSInput, models.WebApplicationPracticeIPSInput])
}

func parseSchemaWebAttacks(schemaWebAttacks any) []models.UpdateWebApplicationPracticeWebAttacksInput {
	input := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaWebAttacks), mapToWebAttacksInput)
	return utils.Map(input, utils.MustUnmarshalAs[models.UpdateWebApplicationPracticeWebAttacksInput, models.WebApplicationPracticeWebAttacksInput])
}

func parseToSchemaWebBot(schemaWebBot any) []models.WebApplicationPracticeWebBotSchema {
	parseFunc := func(schemaWebBotMap map[string]any) models.WebApplicationPracticeWebBotSchema {
		var ret models.WebApplicationPracticeWebBotSchema
		ret.ID = schemaWebBotMap["id"].(string)
		ret.InjectURIs = utils.MustSchemaCollectionToSlice[string](schemaWebBotMap["inject_uris"])
		ret.InjectURIsIDs = utils.MustSchemaCollectionToSlice[string](schemaWebBotMap["inject_uris_ids"])
		ret.ValidURIs = utils.MustSchemaCollectionToSlice[string](schemaWebBotMap["valid_uris"])
		ret.ValidURIsIDs = utils.MustSchemaCollectionToSlice[string](schemaWebBotMap["valid_uris_ids"])
		return ret
	}

	return utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaWebBot), parseFunc)
}

func UpdateWebApplicationPractice(ctx context.Context, c *api.Client, id string, input models.UpdateWebApplicationPracticeInput) (bool, error) {
	vars := map[string]any{"practiceInput": input, "id": id, "ownerId": nil}

	res, err := c.MakeGraphQLRequest(ctx, `
		mutation updateWebApplicationPractice($practiceInput: WebApplicationPracticeUpdateInput, $id: ID!, $ownerId: ID)
		{
			updateWebApplicationPractice(practiceInput: $practiceInput, id: $id, ownerId: $ownerId)
		}
	`, "updateWebApplicationPractice", vars)

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateWebApplicationPractice response %#v should be of type bool", res)
	}

	return value, err
}

//func parseSchemaFileSecurity(schemaFileSecurity any) []models.UpdateFileSecurity {
//	input := utils.Map(utils.MustSchemaCollectionToSlice[map[string]any](schemaFileSecurity), mapToFileSecurityInput)
//	return utils.Map(input, utils.MustUnmarshalAs[models.UpdateFileSecurity, models.FileSecurityInput])
//}
