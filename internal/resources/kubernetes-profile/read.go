package kubernetesprofile

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/kubernetes-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	additonalSettingsIDSeparator = ";;;"
)

func ReadKubernetesProfileToResourceData(profile models.KubernetesProfile, d *schema.ResourceData) error {
	d.SetId(profile.ID)
	d.Set("name", profile.Name)
	d.Set("profile_type", profile.ProfileType)
	d.Set("profile_sub_type", profile.ProfileSubType)
	d.Set("defined_applications_only", profile.OnlyDefinedApplications)
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

func GetKubernetesProfile(ctx context.Context, c *api.Client, id string) (models.KubernetesProfile, error) {
	res, err := c.MakeGraphQLRequest(ctx, `
		{
			getKubernetesProfile(id: "`+id+`") {
				id
				name
				profileType
				profileSubType
				authentication {
					token
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
				onlyDefinedApplications
			}
		}
	`, "getKubernetesProfile")

	if err != nil {
		return models.KubernetesProfile{}, fmt.Errorf("failed to get KubernetesProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.KubernetesProfile](res)
	if err != nil {
		return models.KubernetesProfile{}, fmt.Errorf("failed to convert response to KubernetesProfile struct. Error: %w", err)
	}

	return profile, nil
}
