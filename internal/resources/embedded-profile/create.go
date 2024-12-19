package embeddedprofile

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/embedded-profile"
)

const (
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

func CreateEmbeddedProfileInputFromResourceData(d *schema.ResourceData) (models.CreateEmbeddedProfileInput, error) {
	var res models.CreateEmbeddedProfileInput

	res.Name = d.Get("name").(string)
	res.UpgradeMode = d.Get("upgrade_mode").(string)
	if res.UpgradeMode == UpgradeModeScheduled {
		upgradeTime := handleScheduledUpgradeMode(d)
		res.UpgradeTime = &upgradeTime
	}

	res.OnlyDefinedApplications = d.Get("defined_applications_only").(bool)
	res.Authentication.MaxNumberOfAgents = d.Get("max_number_of_agents").(int)
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

func handleScheduledUpgradeMode(d *schema.ResourceData) models.ScheduleTimeInput {
	var res models.ScheduleTimeInput
	res.ScheduleType = d.Get("upgrade_time_schedule_type").(string)
	res.Time = d.Get("upgrade_time_hour").(string)
	res.Duration = d.Get("upgrade_time_duration").(int)
	if v, ok := d.GetOk("upgrade_time_week_days"); ok {
		weekDays := v.(*schema.Set).List()
		res.WeekDays = make([]string, 0, len(weekDays))
		for _, weekDayInterface := range weekDays {
			weekDay := weekDayInterface.(string)
			res.WeekDays = append(res.WeekDays, weekDay)
		}
	}

	if v, ok := d.GetOk("upgrade_time_days"); ok {
		days := v.(*schema.Set).List()
		res.Days = make([]int, 0, len(days))
		for _, dayInterface := range days {
			day := dayInterface.(int)
			res.Days = append(res.Days, day)
		}
	}

	return res
}

func NewEmbeddedProfile(ctx context.Context, c *api.Client, input models.CreateEmbeddedProfileInput) (models.EmbeddedProfile, error) {
	vars := map[string]any{"profileInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newEmbeddedProfile($profileInput: EmbeddedProfileInput)
					{	
						newEmbeddedProfile (profileInput: $profileInput) {
							id
							name
							profileType
							authentication {
								token
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
								... on ScheduleDaysInMonth {
									days
								}
							}
							onlyDefinedApplications
						}
					}
				`, "newEmbeddedProfile", vars)

	if err != nil {
		return models.EmbeddedProfile{}, fmt.Errorf("failed to create new EmbeddedProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.EmbeddedProfile](res)
	if err != nil {
		return models.EmbeddedProfile{}, fmt.Errorf("failed to convert response to EmbeddedProfile struct. Error: %w", err)
	}

	return profile, nil
}
