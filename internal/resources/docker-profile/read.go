package dockerprofile

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/docker-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	additonalSettingsIDSeperator = ";;;"
)

func ReadDockerProfileToResourceData(profile models.DockerProfile, d *schema.ResourceData) error {
	d.SetId(profile.ID)
	d.Set("name", profile.Name)
	d.Set("profile_type", profile.ProfileType)
	d.Set("defined_applications_only", profile.OnlyDefinedApplications)
	d.Set("max_number_of_agents", profile.Authentication.MaxNumberOfAgents)
	d.Set("authentication_token", profile.Authentication.Token)

	additionalSettingsIDs := make([]string, 0, len(profile.AdditionalSettings))
	additionalSettingsKVs := make(map[string]any, len(profile.AdditionalSettings))
	for _, kv := range profile.AdditionalSettings {
		additionalSettingsIDs = append(additionalSettingsIDs,
			fmt.Sprintf("%s%s%s", kv.Key, additonalSettingsIDSeperator, kv.ID))
		additionalSettingsKVs[kv.Key] = kv.Value
	}

	d.Set("additional_settings", additionalSettingsKVs)
	d.Set("additional_settings_ids", additionalSettingsIDs)

	return nil
}

func GetDockerProfile(c *api.Client, id string) (models.DockerProfile, error) {
	res, err := c.MakeGraphQLRequest(`
		{
			getDockerProfile(id: "`+id+`") {
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
				onlyDefinedApplications
			}
		}
	`, "getDockerProfile")

	if err != nil {
		return models.DockerProfile{}, fmt.Errorf("failed to get DockerProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.DockerProfile](res)
	if err != nil {
		return models.DockerProfile{}, fmt.Errorf("failed to convert response to DockerProfile struct. Error: %w", err)
	}

	return profile, nil
}
