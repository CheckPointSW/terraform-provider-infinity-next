package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubernetesProfileBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_kubernetes_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: kubernetesProfileBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                 nameAttribute,
						"profile_sub_type":     "AccessControl",
						"max_number_of_agents": "10",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: kubernetesProfileUpdateBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"profile_sub_type":          "AccessControl",
						"max_number_of_agents":      "100",
						"additional_settings_ids.#": "2",
						"%":                         "9",
						"profile_type":              "Kubernetes",
						"additional_settings.%":     "2",
						"additional_settings.Key1":  "Value1",
						"additional_settings.Key2":  "Value2",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.1"),
					)...,
				),
			},
		},
	})
}

func TestAccKubernetesProfileFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_kubernetes_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: kubernetesProfileFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"max_number_of_agents":      "100",
						"additional_settings_ids.#": "2",
						"%":                         "9",
						"profile_type":              "Kubernetes",
						"profile_sub_type":          "AccessControl",
						"additional_settings.%":     "2",
						"additional_settings.Key1":  "Value1",
						"additional_settings.Key2":  "Value2",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.1"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: kubernetesProfileUpdateFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"additional_settings.Key6":  "Value6",
						"profile_type":              "Kubernetes",
						"profile_sub_type":          "AppSec",
						"additional_settings.%":     "3",
						"max_number_of_agents":      "101",
						"additional_settings_ids.#": "3",
						"additional_settings.Key2":  "Value11",
						"additional_settings.Key5":  "Value5",
						"%":                         "9",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
					)...,
				),
			},
		},
	})
}

func kubernetesProfileBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_kubernetes_profile" %[1]q {
	name = %[1]q
	profile_sub_type = "AccessControl"
	max_number_of_agents = 10
}
`, name)
}

func kubernetesProfileUpdateBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_kubernetes_profile" %[1]q {
	name             = %[1]q
	profile_sub_type = "AccessControl"
	max_number_of_agents = 100
	defined_applications_only = false
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
}
`, name)
}

func kubernetesProfileFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_kubernetes_profile" %[1]q {
	name             = %[1]q
	profile_sub_type = "AccessControl"
	max_number_of_agents = 100
	defined_applications_only = false
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
}
`, name)
}

func kubernetesProfileUpdateFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_kubernetes_profile" %[1]q {
	name             = %[1]q
	profile_sub_type = "AppSec"
	max_number_of_agents = 101
	defined_applications_only = true
	additional_settings = {
		Key2 = "Value11"
		Key5 = "Value5"
		Key6 = "Value6"
	}
}
`, name)
}
