package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDockerProfileBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_docker_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: dockerProfileBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name": nameAttribute,
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
				Config: dockerProfileUpdateBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"max_number_of_agents":      "100",
						"additional_settings_ids.#": "2",
						"%":                         "8",
						"profile_type":              "Docker",
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

func TestAccDockerProfileFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_docker_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: dockerProfileFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"max_number_of_agents":      "100",
						"additional_settings_ids.#": "2",
						"%":                         "8",
						"profile_type":              "Docker",
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
				Config: dockerProfileUpdateFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"additional_settings.Key6":  "Value6",
						"profile_type":              "Docker",
						"additional_settings.%":     "3",
						"max_number_of_agents":      "101",
						"additional_settings_ids.#": "3",
						"additional_settings.Key2":  "Value11",
						"additional_settings.Key5":  "Value5",
						"%":                         "8",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
					)...,
				),
			},
		},
	})
}

func dockerProfileBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_docker_profile" %[1]q {
	name = %[1]q
}
`, name)
}

func dockerProfileUpdateBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_docker_profile" %[1]q {
	name                          = %[1]q
	max_number_of_agents = 100
	defined_applications_only = false
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
}
`, name)
}

func dockerProfileFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_docker_profile" %[1]q {
	name                 = %[1]q
	max_number_of_agents = 100
	defined_applications_only = false
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
}
`, name)
}

func dockerProfileUpdateFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_docker_profile" %[1]q {
	name                 = %[1]q
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
