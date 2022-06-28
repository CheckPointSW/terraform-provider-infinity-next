package appsecgatewayprofile

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/appsec-gateway-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ProfileSubTypeAws    string = "Aws"
	ProfileSubTypeAzure  string = "Azure"
	ProfileSubTypeVMware string = "VMware"
	ProfileSubTypeHyperV string = "HyperV"

	UpgradeModeAutomatic string = "Automatic"
	UpgradeModeManual    string = "Manual"
	UpgradeModeScheduled string = "Scheduled"

	ScheduleTypeDaily       string = "Daily"
	ScheduleTypeDaysInWeek  string = "DaysInWeek"
	ScheduleTypeDaysInMonth string = "DaysInMonth"

	WeekDaySunday    string = "Sunday"
	WeekDayMonday    string = "Monday"
	WeekDayTuesday   string = "Tuesday"
	WeekDayWednesday string = "Wednesday"
	WeekDayThursday  string = "Thursday"
	WeekDayFriday    string = "Friday"
	WeekDaySaturday  string = "Saturday"
)

func CreateCloudGuardAppSecGatewayProfileInputFromResourceData(d *schema.ResourceData) (models.CreateCloudGuardAppSecGatewayProfileInput, error) {
	var res models.CreateCloudGuardAppSecGatewayProfileInput

	res.Name = d.Get("name").(string)
	res.UpgradeMode = d.Get("upgrade_mode").(string)
	res.ProfileSubType = d.Get("profile_sub_type").(string)
	if res.UpgradeMode == UpgradeModeScheduled {
		upgradeTime := handleScheduledUpgradeMode(d)
		res.UpgradeTime = &upgradeTime
	}

	res.ReverseProxyUpstreamTimeout = d.Get("reverseproxy_upstream_timeout").(int)
	res.Authentication.MaxNumberOfAgents = d.Get("max_number_of_agents").(int)
	res.ReverseProxyAdditionalSettings = mapToKeyValueInput(d, "reverseproxy_additional_settings")
	res.AdditionalSettings = mapToKeyValueInput(d, "additional_settings")

	return res, nil
}

func mapToKeyValueInput(d *schema.ResourceData, key string) []models.KeyValueInput {
	mapUserInput := d.Get(key).(map[string]any)
	res := make([]models.KeyValueInput, 0, len(mapUserInput))
	for key, val := range mapUserInput {
		res = append(res,
			models.KeyValueInput{
				Key:   key,
				Value: val.(string),
			})
	}

	return res
}

func handleScheduledUpgradeMode(d *schema.ResourceData) models.UpgradeTimeInput {
	var res models.UpgradeTimeInput
	res.ScheduleType = d.Get("upgrade_time_schedule_type").(string)
	res.Time = d.Get("upgrade_time_hour").(string)
	res.Duration = d.Get("upgrade_time_duration").(int)
	res.WeekDays = utils.MustResourceDataCollectionToSlice[string](d, "upgrade_time_week_days")

	return res
}

func NewAppSecGatewayProfile(ctx context.Context, c *api.Client, input models.CreateCloudGuardAppSecGatewayProfileInput) (models.CloudGuardAppSecGatewayProfile, error) {
	vars := map[string]any{"profileInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newCloudGuardAppSecGatewayProfile($profileInput: CloudGuardAppSecGatewayProfileInput)
					{	
						newCloudGuardAppSecGatewayProfile (profileInput: $profileInput) {
							id
							name
							profileType
							profileSubType
							authentication {
								token
								authenticationType
								maxNumberOfAgents
							}
							additionalSettings {
								id
								key
								value
							}
							upgradeMode
							upgradeTime {
								scheduleType
								duration
								time
								... on ScheduleDaysInWeek {
									weekDays
								}
							}
							reverseProxyUpstreamTimeout
							reverseProxyAdditionalSettings {
								id
								key
								value
							}
						}
					}
				`, "newCloudGuardAppSecGatewayProfile", vars)

	if err != nil {
		return models.CloudGuardAppSecGatewayProfile{}, fmt.Errorf("failed to create new CloudGuardAppSecGatewayProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.CloudGuardAppSecGatewayProfile](res)
	if err != nil {
		return models.CloudGuardAppSecGatewayProfile{}, fmt.Errorf("failed to convert graphQL response to CloudGuardAppSecGatewayProfile struct. Error: %w", err)
	}

	return profile, err
}
