package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebAppPracticeBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_web_app_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: webAppPracticeBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name": nameAttribute,
						"web_attacks.0.advanced_setting.0.csrf_protection":  "Disabled",
						"web_attacks.0.advanced_setting.0.max_object_depth": "40",
						"web_attacks.#":          "1",
						"web_bot.0.valid_uris.#": "0",
						"ips.0.severity_level":   "MediumOrAbove",
						"ips.0.%":                "7",
						"web_attacks.0.advanced_setting.0.body_size":        "1000000",
						"web_attacks.0.advanced_setting.0.url_size":         "32768",
						"ips.0.low_confidence":                              "Detect",
						"web_attacks.0.minimum_severity":                    "High",
						"web_attacks.0.advanced_setting.0.open_redirect":    "Disabled",
						"web_attacks.0.advanced_setting.0.%":                "9",
						"web_attacks.0.%":                                   "3",
						"web_bot.#":                                         "1",
						"web_attacks.0.advanced_setting.0.error_disclosure": "Disabled",
						"web_bot.0.%":                                       "5",
						"category":                                          "ThreatPrevention",
						"default":                                           "false",
						"ips.0.performance_impact":                          "MediumOrLower",
						"ips.0.protections_from_year":                       "2016",
						"ips.#":                                             "1",
						"web_attacks.0.advanced_setting.0.illegal_http_methods": "false",
						"web_attacks.0.advanced_setting.#":                      "1",
						"web_bot.0.inject_uris_ids.#":                           "0",
						"ips.0.medium_confidence":                               "Prevent",
						"web_bot.0.valid_uris_ids.#":                            "0",
						"web_attacks.0.advanced_setting.0.header_size":          "102400",
						"web_bot.0.inject_uris.#":                               "0",
						"ips.0.high_confidence":                                 "Prevent",
						"practice_type":                                         "WebApplication",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.advanced_setting.0.id"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: webAppPracticeUpdateBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                        nameAttribute,
						"category":                    "ThreatPrevention",
						"ips.0.protections_from_year": "2016",
						"ips.0.performance_impact":    "LowOrLower",
						"web_attacks.#":               "1",
						"web_bot.0.valid_uris.#":      "2",
						"web_attacks.0.advanced_setting.0.csrf_protection":  "Prevent",
						"web_attacks.0.minimum_severity":                    "High",
						"web_bot.0.inject_uris.#":                           "2",
						"web_bot.0.valid_uris_ids.#":                        "2",
						"web_attacks.0.%":                                   "3",
						"web_bot.0.%":                                       "5",
						"ips.0.high_confidence":                             "Detect",
						"web_attacks.0.advanced_setting.0.open_redirect":    "Disabled",
						"web_attacks.0.advanced_setting.#":                  "1",
						"ips.0.low_confidence":                              "Detect",
						"ips.0.%":                                           "7",
						"web_attacks.0.advanced_setting.0.header_size":      "1000",
						"web_attacks.0.advanced_setting.0.max_object_depth": "1000",
						"web_bot.0.inject_uris_ids.#":                       "2",
						"web_bot.#":                                         "1",
						"default":                                           "false",
						"%":                                                 "8",
						"web_attacks.0.advanced_setting.0.illegal_http_methods": "true",
						"web_attacks.0.advanced_setting.0.url_size":             "1000",
						"web_attacks.0.advanced_setting.0.body_size":            "1000",
						"web_attacks.0.advanced_setting.0.error_disclosure":     "AccordingToPractice",
						"ips.0.severity_level":                                  "LowOrAbove",
						"ips.0.medium_confidence":                               "Detect",
						"web_attacks.0.advanced_setting.0.%":                    "9",
						"ips.#":                                                 "1",
						"practice_type":                                         "WebApplication",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.valid_uris_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.valid_uris_ids.1"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.inject_uris_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.inject_uris_ids.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.valid_uris.*", "url1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.valid_uris.*", "url2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.inject_uris.*", "url1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.inject_uris.*", "url2"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.advanced_setting.0.id"),
					)...,
				),
			},
		},
	})
}

func TestAccWebAppPracticeFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_web_app_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: webAppPracticeFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                           nameAttribute,
						"category":                       "ThreatPrevention",
						"ips.0.protections_from_year":    "2016",
						"ips.0.performance_impact":       "LowOrLower",
						"web_attacks.#":                  "1",
						"web_bot.0.valid_uris.#":         "2",
						"web_attacks.0.minimum_severity": "High",
						"web_bot.0.inject_uris.#":        "2",
						"web_bot.0.valid_uris_ids.#":     "2",
						"web_attacks.0.%":                "3",
						"web_bot.0.%":                    "5",
						"ips.0.high_confidence":          "Detect",
						"ips.0.low_confidence":           "Detect",
						"ips.0.%":                        "7",
						"web_bot.0.inject_uris_ids.#":    "2",
						"web_bot.#":                      "1",
						"default":                        "false",
						"ips.0.severity_level":           "LowOrAbove",
						"ips.0.medium_confidence":        "Detect",
						"ips.#":                          "1",
						"practice_type":                  "WebApplication",
						"web_attacks.0.advanced_setting.0.csrf_protection":      "Prevent",
						"web_attacks.0.advanced_setting.0.max_object_depth":     "1000",
						"web_attacks.0.advanced_setting.0.body_size":            "1000",
						"web_attacks.0.advanced_setting.0.url_size":             "1000",
						"web_attacks.0.advanced_setting.0.open_redirect":        "Disabled",
						"web_attacks.0.advanced_setting.0.%":                    "9",
						"web_attacks.0.advanced_setting.0.error_disclosure":     "AccordingToPractice",
						"web_attacks.0.advanced_setting.0.illegal_http_methods": "true",
						"web_attacks.0.advanced_setting.#":                      "1",
						"web_attacks.0.advanced_setting.0.header_size":          "1000",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.valid_uris_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.valid_uris_ids.1"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.inject_uris_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.inject_uris_ids.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.valid_uris.*", "url1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.valid_uris.*", "url2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.inject_uris.*", "url1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.inject_uris.*", "url2"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.advanced_setting.0.id"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: webAppPracticeUpdateFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                           nameAttribute,
						"ips.0.high_confidence":          "Prevent",
						"web_bot.#":                      "1",
						"ips.0.performance_impact":       "MediumOrLower",
						"ips.0.protections_from_year":    "2020",
						"ips.0.severity_level":           "HighOrAbove",
						"ips.0.%":                        "7",
						"web_attacks.0.minimum_severity": "Critical",
						"web_bot.0.inject_uris_ids.#":    "2",
						"web_bot.0.valid_uris_ids.#":     "2",
						"web_attacks.0.%":                "3",
						"category":                       "ThreatPrevention",
						"default":                        "false",
						"ips.0.low_confidence":           "Detect",
						"practice_type":                  "WebApplication",
						"web_bot.0.valid_uris.#":         "2",
						"web_attacks.#":                  "1",
						"web_bot.0.%":                    "5",
						"ips.0.medium_confidence":        "Inactive",
						"web_bot.0.inject_uris.#":        "2",
						"ips.#":                          "1",
						"web_attacks.0.advanced_setting.0.csrf_protection":      "Learn",
						"web_attacks.0.advanced_setting.0.max_object_depth":     "1004",
						"web_attacks.0.advanced_setting.0.body_size":            "1001",
						"web_attacks.0.advanced_setting.0.url_size":             "1002",
						"web_attacks.0.advanced_setting.0.open_redirect":        "AccordingToPractice",
						"web_attacks.0.advanced_setting.0.%":                    "9",
						"web_attacks.0.advanced_setting.0.error_disclosure":     "Prevent",
						"web_attacks.0.advanced_setting.0.illegal_http_methods": "false",
						"web_attacks.0.advanced_setting.#":                      "1",
						"web_attacks.0.advanced_setting.0.header_size":          "1003",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.advanced_setting.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.valid_uris_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.valid_uris_ids.1"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.inject_uris_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "web_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "web_bot.0.inject_uris_ids.1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.valid_uris.*", "url3"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.valid_uris.*", "url4"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.inject_uris.*", "url3"),
						resource.TestCheckTypeSetElemAttr(resourceName, "web_bot.0.inject_uris.*", "url4"),
					)...,
				),
			},
		},
	})
}

func webAppPracticeBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_app_practice" %[1]q {
	name = %[1]q
}
`, name)
}

func webAppPracticeUpdateBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_app_practice" %[1]q {
	name                          = %[1]q
	ips {
		performance_impact    = "LowOrLower"
		severity_level        = "LowOrAbove"
		protections_from_year = "2016"
		high_confidence       = "Detect"
		medium_confidence     = "Detect"
		low_confidence        = "Detect"
	}
	web_attacks {
		minimum_severity = "High"
		advanced_setting {
			csrf_protection      = "Prevent"
			open_redirect        = "Disabled"
			error_disclosure     = "AccordingToPractice"
			body_size            = 1000
			url_size             = 1000
			header_size          = 1000
			max_object_depth     = 1000
			illegal_http_methods = true
		}
	}
	web_bot {
		inject_uris = ["url1", "url2"]
		valid_uris  = ["url1", "url2"]
	}
}
`, name)
}

func webAppPracticeFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_app_practice" %[1]q {
	name                          = %[1]q
	ips {
		performance_impact    = "LowOrLower"
		severity_level        = "LowOrAbove"
		protections_from_year = "2016"
		high_confidence       = "Detect"
		medium_confidence     = "Detect"
		low_confidence        = "Detect"
	}
	web_attacks {
		minimum_severity = "High"
		advanced_setting {
			csrf_protection      = "Prevent"
			open_redirect        = "Disabled"
			error_disclosure     = "AccordingToPractice"
			body_size            = 1000
			url_size             = 1000
			header_size          = 1000
			max_object_depth     = 1000
			illegal_http_methods = true
		}
	}
	web_bot {
		inject_uris = ["url1", "url2"]
		valid_uris  = ["url1", "url2"]
	}
}
`, name)
}

func webAppPracticeUpdateFullConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_app_practice" %[1]q {
	name                          = %[1]q
	ips {
		performance_impact    = "MediumOrLower"
		severity_level        = "HighOrAbove"
		protections_from_year = "2020"
		high_confidence       = "Prevent"
		medium_confidence     = "Inactive"
		low_confidence        = "Detect"
	}
	web_attacks {
		minimum_severity = "Critical"
		advanced_setting {
			csrf_protection      = "Learn"
			open_redirect        = "AccordingToPractice"
			error_disclosure     = "Prevent"
			body_size            = 1001
			url_size             = 1002
			header_size          = 1003
			max_object_depth     = 1004
			illegal_http_methods = false
		}
	}
	web_bot {
		inject_uris = ["url3", "url4"]
		valid_uris  = ["url3", "url4"]
	}
}
`, name)
}
