package kubernetesprofile

import (
	"fmt"
	"log"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/kubernetes-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateKubernetesProfile(c *api.Client, id any, input models.KubernetesProfileUpdateInput) (bool, error) {
	vars := map[string]any{"profileInput": input, "id": id}
	res, err := c.MakeGraphQLRequest(`
				mutation updateKubernetesProfile($profileInput: KubernetesProfileUpdateInput, $id: ID!)
					{
						updateKubernetesProfile (profileInput: $profileInput, id: $id)
					}
				`, "updateKubernetesProfile", vars)

	if err != nil {
		return false, err
	}

	value, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateKubernetesProfile response %#v should be of type bool", res)
	}

	return value, err
}

func UpdateKubernetesProfileInputFromResourceData(d *schema.ResourceData) (models.KubernetesProfileUpdateInput, error) {
	var res models.KubernetesProfileUpdateInput

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		res.Name = newName
	}

	if _, newProfileSubType, hasChange := utils.MustGetChange[string](d, "profile_sub_type"); hasChange {
		res.ProfileSubType = newProfileSubType
	}

	if _, new, hasChange := utils.MustGetChange[bool](d, "defined_applications_only"); hasChange {
		res.OnlyDefinedApplications = new
	}

	res.Authentication.MaxNumberOfAgents = d.Get("max_number_of_agents").(int)
	res.AddAdditionalSettings, res.UpdateAdditionalSettings, res.RemoveAdditionalSettings =
		handleUpdateAdditionalSetting(d, "additional_settings", "additional_settings_ids")

	return res, nil
}

func handleUpdateAdditionalSetting(d *schema.ResourceData, settingsKey, setttingsIDsKey string) ([]models.KeyValueInput, []models.KeyValueUpdateInput, []string) {
	if oldSettingMap, newSettingMap, hasChange := utils.GetChangeWithParse(d, settingsKey, utils.MustValueAs[map[string]any]); hasChange {
		// get reverse proxy additional settings ids - each in the format: "<key><additonalSettingsIDSeperator><ID>"
		additionalSettingsIDsMap := make(map[string]string)
		additionalSettingsIDsInterface := d.Get(setttingsIDsKey).([]any)
		for _, intefaceUnparsedID := range additionalSettingsIDsInterface {
			// parse ID
			keyAndID := strings.Split(intefaceUnparsedID.(string), additonalSettingsIDSeperator)
			key, settingID := keyAndID[0], keyAndID[1]
			additionalSettingsIDsMap[key] = settingID
		}

		// get settings to add or update
		var updateAdditionalSettings []models.KeyValueUpdateInput
		var addAdditionalSettings []models.KeyValueInput
		for newKey, newVal := range newSettingMap {
			if _, ok := oldSettingMap[newKey]; !ok {
				addAdditionalSettings = append(addAdditionalSettings, models.KeyValueInput{
					Key:   newKey,
					Value: newVal.(string),
				})

				continue
			}

			// if oldVal == newVal no need to update
			oldVal := oldSettingMap[newKey].(string)
			if oldVal == newVal.(string) {
				continue
			}

			// else - we need to update the key-value pair
			// id should exist, if not, log warning and continue
			settingID, ok := additionalSettingsIDsMap[newKey]
			if !ok {
				log.Printf("[WARN] Key %s does not have an ID in state. Removing and re-adding it with new value", newKey)
				continue
			}

			// updating the value
			updateAdditionalSettings = append(updateAdditionalSettings, models.KeyValueUpdateInput{
				Key:   newKey,
				Value: newVal.(string),
				ID:    settingID,
			})
		}

		// get settings to remove
		var removeAdditionalSettings []string
		for oldKey := range oldSettingMap {
			if _, ok := newSettingMap[oldKey]; !ok {
				// id should exist, if not, log warning and continue
				settingID, ok := additionalSettingsIDsMap[oldKey]
				if !ok {
					log.Printf("[WARN] Key %s does not have an ID in state. Removing and re-adding it with new value", oldKey)
					continue
				}

				removeAdditionalSettings = append(removeAdditionalSettings, settingID)
			}
		}

		return addAdditionalSettings, updateAdditionalSettings, removeAdditionalSettings
	}

	return nil, nil, nil

}
