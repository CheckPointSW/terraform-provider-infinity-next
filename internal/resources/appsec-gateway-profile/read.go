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
	additonalSettingsIDSeparator = ";;;"
)

func ReadCloudGuardAppSecGatewayProfileToResourceData(profile models.CloudGuardAppSecGatewayProfile, d *schema.ResourceData) error {
	d.SetId(profile.ID)
	d.Set("name", profile.Name)
	d.Set("profile_sub_type", profile.ProfileSubType)
	d.Set("profile_type", profile.ProfileType)
	d.Set("upgrade_mode", profile.UpgradeMode)
	if profile.UpgradeMode == UpgradeModeScheduled {
		d.Set("upgrade_time_schedule_type", profile.UpgradeTime.ScheduleType)
		d.Set("upgrade_time_hour", profile.UpgradeTime.Time)
		d.Set("upgrade_time_duration", profile.UpgradeTime.Duration)
		d.Set("upgrade_time_week_days", profile.UpgradeTime.WeekDays)
		d.Set("upgrade_time_days", profile.UpgradeTime.Days)
	}

	d.Set("reverseproxy_upstream_timeout", profile.ReverseProxyUpstreamTimeout)
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

	reverseProxyAdditionalSettingsIDs := make([]string, 0, len(profile.ReverseProxyAdditionalSettings))
	reverseProxyAdditionalSettingsKVs := make(map[string]any, len(profile.ReverseProxyAdditionalSettings))
	for _, kv := range profile.ReverseProxyAdditionalSettings {
		reverseProxyAdditionalSettingsIDs = append(reverseProxyAdditionalSettingsIDs,
			fmt.Sprintf("%s%s%s", kv.Key, additonalSettingsIDSeparator, kv.ID))
		reverseProxyAdditionalSettingsKVs[kv.Key] = kv.Value
	}

	d.Set("reverseproxy_additional_settings", reverseProxyAdditionalSettingsKVs)
	d.Set("reverseproxy_additional_settings_ids", reverseProxyAdditionalSettingsIDs)

	d.Set("certificate_type", profile.CertificateType)
	d.Set("fail_open_inspection", profile.FailOpenInspection)

	return nil
}

func GetCloudGuardAppSecGatewayProfile(ctx context.Context, c *api.Client, id string) (models.CloudGuardAppSecGatewayProfile, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getCloudGuardAppSecGatewayProfile(id: "`+id+`") {
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
				usedBy {
					id
					name
					type	
					subType
					objectStatus
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
				reverseProxyUpstreamTimeout
				reverseProxyAdditionalSettings {
					id
					key
					value
				}
				certificateType
				failOpenInspection
			}
		}
	`, "getCloudGuardAppSecGatewayProfile")

	if err != nil {
		return models.CloudGuardAppSecGatewayProfile{}, fmt.Errorf("failed to get CloudGuardAppSecGatewayProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.CloudGuardAppSecGatewayProfile](res)
	if err != nil {
		return models.CloudGuardAppSecGatewayProfile{}, fmt.Errorf("failed to convert response to CloudGuardAppSecGatewayProfile struct. Error: %w", err)
	}

	return profile, nil
}
