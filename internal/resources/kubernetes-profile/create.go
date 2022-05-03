package kubernetesprofile

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	models "github.com/CheckPointSW/terraform-provider-infinity-next/internal/models/kubernetes-profile"
)

func CreateKubernetesProfileInputFromResourceData(d *schema.ResourceData) (models.CreateKubernetesProfileInput, error) {
	var res models.CreateKubernetesProfileInput

	res.Name = d.Get("name").(string)
	res.ProfileSubType = d.Get("profile_sub_type").(string)
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

func NewKubernetesProfile(c *api.Client, input models.CreateKubernetesProfileInput) (models.KubernetesProfile, error) {
	vars := map[string]any{"profileInput": input}
	res, err := c.MakeGraphQLRequest(`
				mutation newKubernetesProfile($profileInput: KubernetesProfileInput)
					{	
						newKubernetesProfile (profileInput: $profileInput) {
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
							onlyDefinedApplications
						}
					}
				`, "newKubernetesProfile", vars)

	if err != nil {
		return models.KubernetesProfile{}, fmt.Errorf("failed to create new KubernetesProfile: %w", err)
	}

	profile, err := utils.UnmarshalAs[models.KubernetesProfile](res)
	if err != nil {
		return models.KubernetesProfile{}, fmt.Errorf("failed to convert response to KubernetesProfile struct. Error: %w", err)
	}

	return profile, nil
}
