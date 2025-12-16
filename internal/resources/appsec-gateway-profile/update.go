package appsecgatewayprofile

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/appsec-gateway-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateAppSecGatewayProfile(ctx context.Context, c *api.Client, id any, input models.UpdateCloudGuardAppSecGatewayProfileInput) (bool, error) {
	vars := map[string]any{"profileInput": input, "id": id}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation updateCloudGuardAppSecGatewayProfile($profileInput: CloudGuardAppSecGatewayProfileUpdateInput, $id: ID!)
					{
						updateCloudGuardAppSecGatewayProfile (profileInput: $profileInput, id: $id)
					}
				`, "updateCloudGuardAppSecGatewayProfile", vars)

	if err != nil {
		return false, err
	}

	isUpdated, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("invalid updateCloudGuardAppSecGatewayProfile response %#v should be of type bool", res)
	}

	return isUpdated, err
}

func UpdateCloudGuardAppSecGatewayProfileInputFromResourceData(d *schema.ResourceData) (models.UpdateCloudGuardAppSecGatewayProfileInput, error) {
	var res models.UpdateCloudGuardAppSecGatewayProfileInput

	if _, newName, hasChange := utils.MustGetChange[string](d, "name"); hasChange {
		res.Name = newName
	}

	if _, newUpgradeMode, hasChange := utils.MustGetChange[string](d, "upgrade_mode"); hasChange {
		res.UpgradeMode = newUpgradeMode
	}

	var currUpgradeMode string
	if res.UpgradeMode != "" {
		currUpgradeMode = res.UpgradeMode
	} else {
		currUpgradeMode = d.Get("upgrade_mode").(string)
	}

	if currUpgradeMode == UpgradeModeScheduled {
		var upgradeTime models.UpdateUpgradeTimeInput
		if _, newScheduleType, hasChange := utils.MustGetChange[string](d, "upgrade_time_schedule_type"); hasChange {
			upgradeTime.ScheduleType = newScheduleType
		}

		if _, newUpgradeTimeHour, hasChange := utils.MustGetChange[string](d, "upgrade_time_hour"); hasChange {
			upgradeTime.Time = newUpgradeTimeHour
		}

		if _, newUpgradeTimeDuration, hasChange := utils.MustGetChange[int](d, "upgrade_time_duration"); hasChange {
			upgradeTime.Duration = newUpgradeTimeDuration
		}

		if _, newUpgradeTimeWeekDays, hasChange := utils.GetChangeWithParse(d, "upgrade_time_week_days", utils.MustSchemaCollectionToSlice[string]); hasChange {
			upgradeTime.WeekDays = newUpgradeTimeWeekDays
		}

		if _, newUpgradeTime, hasChange := utils.MustGetChange[models.UpdateUpgradeTimeInput](d, "upgrade_time"); hasChange {
			upgradeTime = newUpgradeTime
		}

		res.UpgradeTime = &upgradeTime
	}

	if _, newReverseProxyTimeout, hasChange := utils.MustGetChange[int](d, "reverseproxy_upstream_timeout"); hasChange {
		res.ReverseProxyUpstreamTimeout = newReverseProxyTimeout
	}

	if _, newMaxNumberOfAgents, hasChange := utils.MustGetChange[int](d, "max_number_of_agents"); hasChange {
		res.Authentication.MaxNumberOfAgents = newMaxNumberOfAgents
	}

	res.AddReverseProxyAdditionalSettings, res.UpdateReverseProxyAdditionalSettings, res.RemoveReverseProxyAdditionalSettings =
		handleUpdateAdditionalSetting(d, "reverseproxy_additional_settings", "reverseproxy_additional_settings_ids")

	res.AddAdditionalSettings, res.UpdateAdditionalSettings, res.RemoveAdditionalSettings =
		handleUpdateAdditionalSetting(d, "additional_settings", "additional_settings_ids")

	if _, newProfileSubType, hasChange := utils.MustGetChange[string](d, "profile_sub_type"); hasChange {
		res.ProfileSubType = newProfileSubType
	}

	if _, newCertificateType, hasChange := utils.MustGetChange[string](d, "certificate_type"); hasChange {
		res.CertificateType = newCertificateType
	}

	if _, newFailOpenInspection, hasChange := utils.MustGetChange[bool](d, "fail_open_inspection"); hasChange {
		res.FailOpenInspection = newFailOpenInspection
	}

	return res, nil
}

func handleUpdateAdditionalSetting(d *schema.ResourceData, settingsKey, setttingsIDsKey string) ([]models.AddKeyValue, []models.UpdateKeyValue, []string) {
	if oldAdditionalSetting, newAdditionalSetting, hasChange := utils.GetChangeWithParse(d, settingsKey, utils.MustValueAs[map[string]any]); hasChange {
		oldAdditionalSettingIDs := utils.MustResourceDataCollectionToSlice[string](d, setttingsIDsKey)
		oldAdditionalSettingIDsindicatorMap := make(map[string]string)
		for _, settingID := range oldAdditionalSettingIDs {
			keyAndID := strings.Split(settingID, additonalSettingsIDSeparator)
			key, settingID := keyAndID[0], keyAndID[1]
			oldAdditionalSettingIDsindicatorMap[key] = settingID
		}

		// get settings to add or update
		var updateAdditionalSettings []models.UpdateKeyValue
		var addAdditionalSettings []models.AddKeyValue
		for newKey, newVal := range newAdditionalSetting {
			// if value does not exist in oldSettings - add it
			if _, ok := oldAdditionalSetting[newKey]; !ok {
				addAdditionalSettings = append(addAdditionalSettings, models.AddKeyValue{
					Key:   newKey,
					Value: newVal.(string),
				})

				continue
			}

			// if oldVal == newVal no need to update
			oldVal := oldAdditionalSetting[newKey].(string)
			if oldVal == newVal.(string) {
				continue
			}

			// else - we need to update the key-value pair
			// id should exist, if not, log warning and continue
			settingID, ok := oldAdditionalSettingIDsindicatorMap[newKey]
			if !ok {
				log.Printf("[WARN] Key %s does not have an ID in state, re-adding it with new value", newKey)
				addAdditionalSettings = append(addAdditionalSettings, models.AddKeyValue{
					Key:   newKey,
					Value: newVal.(string),
				})
				continue
			}

			// updating the value
			updateAdditionalSettings = append(updateAdditionalSettings, models.UpdateKeyValue{
				Key:   newKey,
				Value: newVal.(string),
				ID:    settingID,
			})
		}

		// get settings to remove
		var removeAdditionalSettings []string
		for oldKey := range oldAdditionalSetting {
			if _, ok := newAdditionalSetting[oldKey]; !ok {
				// id should exist, if not, log warning and continue
				settingID, ok := oldAdditionalSettingIDsindicatorMap[oldKey]
				if !ok {
					log.Printf("[WARN] Key %s does not have an ID in state", oldKey)
					continue
				}

				removeAdditionalSettings = append(removeAdditionalSettings, settingID)
			}
		}

		return addAdditionalSettings, updateAdditionalSettings, removeAdditionalSettings
	}

	return nil, nil, nil

}
