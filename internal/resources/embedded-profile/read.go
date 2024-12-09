package embeddedprofile

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/embedded-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	additonalSettingsIDSeparator = ";;;"
)

func ReadEmbeddedProfileToResourceData(profile models.EmbeddedProfile, d *schema.ResourceData) error {
	d.SetId(profile.ID)
	d.Set("name", profile.Name)
	d.Set("profile_type", profile.ProfileType)
	d.Set("upgrade_mode", profile.UpgradeMode)
	d.Set("defined_applications_only", profile.OnlyDefinedApplications)
	if profile.UpgradeMode == UpgradeModeScheduled {
		d.Set("upgrade_time_schedule_type", profile.UpgradeTime.ScheduleType)
		d.Set("upgrade_time_hour", profile.UpgradeTime.Time)
		d.Set("upgrade_time_duration", profile.UpgradeTime.Duration)
		d.Set("upgrade_time_week_days", profile.UpgradeTime.WeekDays)
		d.Set("upgrade_time_days", profile.UpgradeTime.Days)
	}

	d.Set("max_number_of_agents", profile.Authentication.MaxNumberOfAgents)
	d.Set("authentication_token", profile.Authentication.Token)

	additionalSettingsIDs := make([]string, 0, len(profile.AdditionalSettings))
	additionalSettingsKVs := make(map[string]any, len(profile.AdditionalSettings))
	for _, kv := range profile.AdditionalSettings {
		additionalSettingsIDs = append(additionalSettingsIDs,
			fmt.Sprintf("%s%s%s", kv.Key, additonalSettingsIDSeparator, kv.ID))
		additionalSettingsKVs[kv.Key] = kv.Value
	}

	d.Set("additional_settings", additionalSettingsKVs)
	d.Set("additional_settings_ids", additionalSettingsIDs)

	return nil
}

func GetEmbeddedProfile(ctx context.Context, c *api.Client, id string) (models.EmbeddedProfile, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getEmbeddedProfile(id: "`+id+`") {
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
	`, "getEmbeddedProfile")

	if err != nil {
		return models.EmbeddedProfile{}, fmt.Errorf("failed to get EmbeddedProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.EmbeddedProfile](res)
	if err != nil {
		return models.EmbeddedProfile{}, fmt.Errorf("failed to convert response to EmbeddedProfile struct. Error: %w", err)
	}

	return profile, nil
}
