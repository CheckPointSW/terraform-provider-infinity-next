package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEmbeddedProfileBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_embedded_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: embeddedProfileBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                 nameAttribute,
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
				Config: embeddedProfileUpdateBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                       nameAttribute,
						"upgrade_mode":               "Scheduled",
						"upgrade_time_schedule_type": "DaysInWeek",
						"upgrade_time_hour":          "12:00",
						"upgrade_time_duration":      "10",
						"upgrade_time_week_days.#":   "3",
						"max_number_of_agents":       "100",
						"additional_settings_ids.#":  "2",
						"%":                          "13",
						"profile_type":               "Embedded",
						"additional_settings.%":      "2",
						"upgrade_time_week_days.1":   "Monday",
						"upgrade_time_week_days.2":   "Thursday",
						"upgrade_time_week_days.0":   "Friday",
						"additional_settings.Key1":   "Value1",
						"additional_settings.Key2":   "Value2",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Monday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Thursday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Friday"),
					)...,
				),
			},
		},
	})
}

func TestAccEmbeddedProfileFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_embedded_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: embeddedProfileFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                       nameAttribute,
						"upgrade_mode":               "Scheduled",
						"upgrade_time_schedule_type": "DaysInWeek",
						"upgrade_time_hour":          "12:00",
						"upgrade_time_duration":      "10",
						"upgrade_time_week_days.#":   "3",
						"max_number_of_agents":       "100",
						"additional_settings_ids.#":  "2",
						"%":                          "13",
						"profile_type":               "Embedded",
						"additional_settings.%":      "2",
						"upgrade_time_week_days.1":   "Monday",
						"upgrade_time_week_days.2":   "Thursday",
						"upgrade_time_week_days.0":   "Friday",
						"additional_settings.Key1":   "Value1",
						"additional_settings.Key2":   "Value2",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Monday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Thursday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Friday"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: embeddedProfileUpdateFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                       nameAttribute,
						"upgrade_mode":               "Scheduled",
						"upgrade_time_hour":          "13:00",
						"additional_settings.Key6":   "Value6",
						"profile_type":               "Embedded",
						"additional_settings.%":      "3",
						"upgrade_time_week_days.0":   "Monday",
						"max_number_of_agents":       "101",
						"upgrade_time_duration":      "12",
						"upgrade_time_schedule_type": "DaysInWeek",
						"upgrade_time_week_days.#":   "2",
						"additional_settings_ids.#":  "3",
						"additional_settings.Key2":   "Value11",
						"additional_settings.Key5":   "Value5",
						"%":                          "13",
						"upgrade_time_week_days.1":   "Sunday",
					}),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Monday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Sunday"),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
					)...,
				),
			},
		},
	})
}

func embeddedProfileBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_embedded_profile" %[1]q {
	name = %[1]q
	max_number_of_agents = 10
}
`, name)
}

func embeddedProfileUpdateBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_embedded_profile" %[1]q {
	name                          = %[1]q
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	max_number_of_agents = 100
	defined_applications_only = false
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
}
`, name)
}

func embeddedProfileFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_embedded_profile" %[1]q {
	name                          = %[1]q
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	max_number_of_agents = 100
	defined_applications_only = false
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
}
`, name)
}

func embeddedProfileUpdateFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_embedded_profile" %[1]q {
	name                          = %[1]q
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "13:00"
	upgrade_time_duration         = 12
	upgrade_time_week_days        = ["Monday", "Sunday"]
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
