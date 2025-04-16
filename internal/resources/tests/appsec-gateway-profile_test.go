package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppsecGatewayProfileBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_appsec_gateway_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: appsecGatewayProfileBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                 nameAttribute,
						"profile_sub_type":     "Aws",
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
				Config: appsecGatewayProfileUpdateBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                                   nameAttribute,
						"profile_sub_type":                       "Aws",
						"upgrade_mode":                           "Scheduled",
						"upgrade_time_schedule_type":             "DaysInWeek",
						"upgrade_time_hour":                      "12:00",
						"upgrade_time_duration":                  "10",
						"upgrade_time_week_days.#":               "3",
						"reverseproxy_upstream_timeout":          "3600",
						"max_number_of_agents":                   "100",
						"reverseproxy_additional_settings_ids.#": "2",
						"additional_settings_ids.#":              "2",
						"%":                                      "17",
						"profile_type":                           "CloudGuardAppSecGateway",
						"additional_settings.%":                  "2",
						"upgrade_time_week_days.1":               "Monday",
						"upgrade_time_week_days.2":               "Thursday",
						"reverseproxy_additional_settings.%":     "2",
						"upgrade_time_week_days.0":               "Friday",
						"reverseproxy_additional_settings.Key4":  "Value4",
						"reverseproxy_additional_settings.Key3":  "Value5",
						"additional_settings.Key1":               "Value1",
						"additional_settings.Key2":               "Value2",
						"certificate_type":                       "Gateway",
						"failOpenInspection":                     "true",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.1"),
						resource.TestCheckResourceAttrSet(resourceName, "reverseproxy_additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "reverseproxy_additional_settings_ids.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Monday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Thursday"),
						resource.TestCheckTypeSetElemAttr(resourceName, "upgrade_time_week_days.*", "Friday"),
					)...,
				),
			},
		},
	})
}

func TestAccAppsecGatewayProfileFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_appsec_gateway_profile." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: appsecGatewayProfileFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                                   nameAttribute,
						"profile_sub_type":                       "Aws",
						"upgrade_mode":                           "Scheduled",
						"upgrade_time_schedule_type":             "DaysInWeek",
						"upgrade_time_hour":                      "12:00",
						"upgrade_time_duration":                  "10",
						"upgrade_time_week_days.#":               "3",
						"reverseproxy_upstream_timeout":          "3600",
						"max_number_of_agents":                   "100",
						"reverseproxy_additional_settings_ids.#": "2",
						"additional_settings_ids.#":              "2",
						"%":                                      "19",
						"profile_type":                           "CloudGuardAppSecGateway",
						"additional_settings.%":                  "2",
						"upgrade_time_week_days.1":               "Monday",
						"upgrade_time_week_days.2":               "Thursday",
						"reverseproxy_additional_settings.%":     "2",
						"upgrade_time_week_days.0":               "Friday",
						"reverseproxy_additional_settings.Key4":  "Value4",
						"reverseproxy_additional_settings.Key3":  "Value5",
						"additional_settings.Key1":               "Value1",
						"additional_settings.Key2":               "Value2",
						"certificate_type":                       "Vault",
						"failOpenInspection":                     "true",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "authentication_token"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "additional_settings_ids.1"),
						resource.TestCheckResourceAttrSet(resourceName, "reverseproxy_additional_settings_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "reverseproxy_additional_settings_ids.1"),
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
				Config: appsecGatewayProfileUpdateFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                                   nameAttribute,
						"reverseproxy_additional_settings.Key8":  "Value8",
						"reverseproxy_additional_settings.Key3":  "Value10",
						"reverseproxy_additional_settings.Key7":  "Value7",
						"upgrade_mode":                           "Scheduled",
						"upgrade_time_hour":                      "13:00",
						"additional_settings.Key6":               "Value6",
						"reverseproxy_additional_settings.%":     "3",
						"profile_type":                           "CloudGuardAppSecGateway",
						"additional_settings.%":                  "3",
						"upgrade_time_week_days.0":               "Monday",
						"max_number_of_agents":                   "100",
						"upgrade_time_duration":                  "12",
						"reverseproxy_upstream_timeout":          "3601",
						"profile_sub_type":                       "Azure",
						"reverseproxy_additional_settings_ids.#": "3",
						"upgrade_time_schedule_type":             "DaysInWeek",
						"upgrade_time_week_days.#":               "2",
						"additional_settings_ids.#":              "3",
						"additional_settings.Key2":               "Value11",
						"additional_settings.Key5":               "Value5",
						"%":                                      "17",
						"upgrade_time_week_days.1":               "Sunday",
						"certificate_type":                       "Vault",
						"failOpenInspection":                     "false",
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

func appsecGatewayProfileBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_appsec_gateway_profile" %[1]q {
	name = %[1]q
	profile_sub_type = "Aws"
	max_number_of_agents = 10
	certificate_type = "Vault"
}
`, name)
}

func appsecGatewayProfileUpdateBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_appsec_gateway_profile" %[1]q {
	name                          = %[1]q
	profile_sub_type              = "Aws"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	reverseproxy_upstream_timeout = 3600
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key3 = "Value5"
		Key4 = "Value4"
	}
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
	certificate_type = "Gateway"
    fail_open_inspection = true
}
`, name)
}

func appsecGatewayProfileFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_appsec_gateway_profile" %[1]q {
	name                          = %[1]q
	profile_sub_type              = "Aws"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	reverseproxy_upstream_timeout = 3600
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key3 = "Value5"
		Key4 = "Value4"
	}
	additional_settings = {
		Key1 = "Value1"
		Key2 = "Value2"
	}
	certificate_type = "Vault"
    fail_open_inspection = true
}
`, name)
}

func appsecGatewayProfileUpdateFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_appsec_gateway_profile" %[1]q {
	name                          = %[1]q
	profile_sub_type              = "Azure"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "13:00"
	upgrade_time_duration         = 12
	upgrade_time_week_days        = ["Monday", "Sunday"]
	reverseproxy_upstream_timeout = 3601
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key3 = "Value10"
		Key7 = "Value7"
		Key8 = "Value8"
	}
	additional_settings = {
		Key2 = "Value11"
		Key5 = "Value5"
		Key6 = "Value6"
	}
	certificate_type = "Vault"
    fail_open_inspection = false
}
`, name)
}
