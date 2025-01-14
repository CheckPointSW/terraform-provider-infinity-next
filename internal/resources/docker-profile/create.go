package dockerprofile

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/docker-profile"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateDockerProfileInputFromResourceData(d *schema.ResourceData) (models.CreateDockerProfileInput, error) {
	var res models.CreateDockerProfileInput

	res.Name = d.Get("name").(string)

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

func NewDockerProfile(ctx context.Context, c *api.Client, input models.CreateDockerProfileInput) (models.DockerProfile, error) {
	vars := map[string]any{"profileInput": input}
	res, err := c.MakeGraphQLRequest(ctx, `
				mutation newDockerProfile($profileInput: DockerProfileInput)
					{	
						newDockerProfile (profileInput: $profileInput) {
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
				`, "newDockerProfile", vars)

	if err != nil {
		return models.DockerProfile{}, fmt.Errorf("failed to create new DockerProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.DockerProfile](res)
	if err != nil {
		return models.DockerProfile{}, fmt.Errorf("failed to convert response to DockerProfile struct. Error: %w", err)
	}

	return profile, nil
}
