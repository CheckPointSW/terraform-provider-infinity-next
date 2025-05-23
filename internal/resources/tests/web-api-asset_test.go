package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebAPIAssetBasic(t *testing.T) {
	assetNameAttribute := acctest.GenerateResourceName()
	profileNameAttribute := acctest.GenerateResourceName()
	trustedSourcesNameAttribute := acctest.GenerateResourceName()
	practiceNameAttribute := acctest.GenerateResourceName()
	logTriggerNameAttribute := acctest.GenerateResourceName()
	exceptionsNameAttribute := acctest.GenerateResourceName()
	assetResourceName := "inext_web_api_asset." + assetNameAttribute
	profileResourceName := "inext_appsec_gateway_profile." + profileNameAttribute
	trustedSourcesResourceName := "inext_trusted_sources." + trustedSourcesNameAttribute
	practiceResourceName := "inext_web_api_practice." + practiceNameAttribute
	logTriggerResourceName := "inext_log_trigger." + logTriggerNameAttribute
	exceptionsResourceName := "inext_exceptions." + exceptionsNameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy: acctest.CheckResourceDestroyed([]string{assetResourceName, profileResourceName, trustedSourcesResourceName,
			practiceResourceName, logTriggerResourceName, exceptionsResourceName}),
		Steps: []resource.TestStep{
			{
				Config: webAPIAssetBasicConfig(assetNameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(assetResourceName, map[string]string{
						"name":            assetNameAttribute,
						"urls.0":          fmt.Sprintf("http://host/%s/path1", assetNameAttribute),
						"urls.#":          "1",
						"%":               "32",
						"urls_ids.#":      "1",
						"main_attributes": fmt.Sprintf("{\"applicationUrls\":\"http://host/%s/path1\"}", assetNameAttribute),
					}),
						resource.TestCheckResourceAttrSet(assetResourceName, "id"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName: assetResourceName,
				ImportState:  true,
			},
			{
				Config: webAPIAssetUpdateBasicConfig(assetNameAttribute, profileNameAttribute, trustedSourcesNameAttribute,
					practiceNameAttribute, logTriggerNameAttribute, exceptionsNameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(assetResourceName, map[string]string{
						"name":                                  assetNameAttribute,
						"%":                                     "32",
						"read_only":                             "false",
						"upstream_url":                          "some url 5",
						"urls.#":                                "2",
						"urls_ids.#":                            "2",
						"profiles.#":                            "1",
						"practice.#":                            "1",
						"practice.0.%":                          "5",
						"practice.0.triggers.#":                 "1",
						"practice.0.sub_practices_modes.IPS":    "AccordingToPractice",
						"practice.0.sub_practices_modes.WebBot": "AccordingToPractice",
						"practice.0.sub_practices_modes.Snort":  "Disabled",
						"practice.0.main_mode":                  "Prevent",
						"source_identifier.0.%":                 "4",
						"source_identifier.1.%":                 "4",
						"source_identifier.2.%":                 "4",
						"source_identifier.2.values.#":          "1",
						"source_identifier.#":                   "3",
						"source_identifier.2.values_ids.#":      "1",
						"source_identifier.1.values_ids.#":      "1",
						"source_identifier.1.values.#":          "1",
						"source_identifier.0.values.#":          "1",
						"source_identifier.0.values_ids.#":      "1",
						"proxy_setting.#":                       "3",
						"proxy_setting.0.%":                     "3",
						"proxy_setting.1.%":                     "3",
						"proxy_setting.2.%":                     "3",
						"class":                                 "workload",
						"category":                              "cloud",
						"group":                                 "",
						"order":                                 "",
						"kind":                                  "",
						"family":                                "Web API",
						"main_attributes":                       fmt.Sprintf("{\"applicationUrls\":\"http://host/%[1]s/path2;http://host/%[1]s/path3\"}", assetNameAttribute),
						"asset_type":                            "WebAPI",
						"intelligence_tags":                     "",
						"tags.#":                                "1",
						"tags.0.%":                              "3",

						"mtls.#":                           "1",
						"mtls.0.filename":                  "cert.pem",
						"mtls.0.certificate_type":          ".pem",
						"mtls.0.data":                      "cert data",
						"mtls.0.type":                      "client",
						"mtls.0.enable":                    "true",
						"additional_instructions_blocks.#": "1",
						"additional_instructions_blocks.0.filename":      "location.json",
						"additional_instructions_blocks.0.filename_type": ".json",
						"additional_instructions_blocks.0.data":          "location data",
						"additional_instructions_blocks.0.type":          "location_instructions",
						"additional_instructions_blocks.0.enable":        "true",
						"custom_headers.#":                               "1",
					}),
						resource.TestCheckResourceAttrSet(assetResourceName, "id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "practice.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.2.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.2.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "id"),
						resource.TestCheckTypeSetElemAttr(assetResourceName, "urls.*", fmt.Sprintf("http://host/%s/path2", assetNameAttribute)),
						resource.TestCheckTypeSetElemAttr(assetResourceName, "urls.*", fmt.Sprintf("http://host/%s/path3", assetNameAttribute)),
						resource.TestCheckResourceAttrSet(assetResourceName, "tags.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers.0.header_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers_id"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWebAPIAssetFull(t *testing.T) {
	assetNameAttribute := acctest.GenerateResourceName()
	profileNameAttribute := acctest.GenerateResourceName()
	trustedSourcesNameAttribute := acctest.GenerateResourceName()
	practiceNameAttribute := acctest.GenerateResourceName()
	logTriggerNameAttribute := acctest.GenerateResourceName()
	exceptionsNameAttribute := acctest.GenerateResourceName()
	anotherProfileNameAttribute := acctest.GenerateResourceName()
	anotherTrustedSourcesNameAttribute := acctest.GenerateResourceName()
	anotherLogTriggerNameAttribute := acctest.GenerateResourceName()
	anotherExceptionsNameAttribute := acctest.GenerateResourceName()
	assetResourceName := "inext_web_api_asset." + assetNameAttribute
	profileResourceName := "inext_appsec_gateway_profile." + profileNameAttribute
	trustedSourcesResourceName := "inext_trusted_sources." + trustedSourcesNameAttribute
	practiceResourceName := "inext_web_api_practice." + practiceNameAttribute
	logTriggerResourceName := "inext_log_trigger." + logTriggerNameAttribute
	exceptionsResourceName := "inext_exceptions." + exceptionsNameAttribute
	anotherProfileResourceName := "inext_appsec_gateway_profile." + anotherProfileNameAttribute
	anotherTrustedSourcesResourceName := "inext_trusted_sources." + anotherTrustedSourcesNameAttribute
	anotherLogTriggerResourceName := "inext_log_trigger." + anotherLogTriggerNameAttribute
	anotherExceptionsResourceName := "inext_exceptions." + anotherExceptionsNameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy: acctest.CheckResourceDestroyed([]string{assetResourceName, profileResourceName, trustedSourcesResourceName,
			practiceResourceName, logTriggerResourceName, exceptionsResourceName, anotherProfileResourceName, anotherTrustedSourcesResourceName,
			anotherLogTriggerResourceName, anotherExceptionsResourceName}),
		Steps: []resource.TestStep{
			{
				Config: webAPIAssetFullConfig(assetNameAttribute, profileNameAttribute, trustedSourcesNameAttribute,
					practiceNameAttribute, logTriggerNameAttribute, exceptionsNameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(assetResourceName, map[string]string{
						"name":                                  assetNameAttribute,
						"%":                                     "32",
						"read_only":                             "false",
						"upstream_url":                          "some url 5",
						"urls.#":                                "2",
						"urls_ids.#":                            "2",
						"profiles.#":                            "1",
						"practice.#":                            "1",
						"practice.0.%":                          "5",
						"practice.0.triggers.#":                 "1",
						"practice.0.sub_practices_modes.IPS":    "AccordingToPractice",
						"practice.0.sub_practices_modes.WebBot": "AccordingToPractice",
						"practice.0.sub_practices_modes.Snort":  "Disabled",
						"practice.0.main_mode":                  "Learn",
						"source_identifier.0.%":                 "4",
						"source_identifier.1.%":                 "4",
						"source_identifier.2.%":                 "4",
						"source_identifier.2.values.#":          "1",
						"source_identifier.#":                   "3",
						"source_identifier.2.values_ids.#":      "1",
						"source_identifier.1.values_ids.#":      "1",
						"source_identifier.1.values.#":          "1",
						"source_identifier.0.values.#":          "1",
						"source_identifier.0.values_ids.#":      "1",
						"proxy_setting.#":                       "3",
						"proxy_setting.0.%":                     "3",
						"proxy_setting.1.%":                     "3",
						"proxy_setting.2.%":                     "3",
						"class":                                 "workload",
						"category":                              "cloud",
						"group":                                 "",
						"order":                                 "",
						"kind":                                  "",
						"family":                                "Web API",
						"main_attributes":                       fmt.Sprintf("{\"applicationUrls\":\"http://host/%[1]s/path1;http://host/%[1]s/path2\"}", assetNameAttribute),
						"asset_type":                            "WebAPI",
						"intelligence_tags":                     "",
						"tags.#":                                "2",
						"tags.0.%":                              "3",
						"tags.1.%":                              "3",

						"mtls.#":                           "1",
						"mtls.0.filename":                  "cert.der",
						"mtls.0.certificate_type":          ".der",
						"mtls.0.data":                      "cert data",
						"mtls.0.type":                      "client",
						"mtls.0.enable":                    "true",
						"additional_instructions_blocks.#": "1",
						"additional_instructions_blocks.0.filename":      "location.json",
						"additional_instructions_blocks.0.filename_type": ".json",
						"additional_instructions_blocks.0.data":          "location data",
						"additional_instructions_blocks.0.type":          "location_instructions",
						"additional_instructions_blocks.0.enable":        "true",
						"redirect_to_https":                              "true",
						"access_log":                                     "true",
						"custom_headers.#":                               "1",
					}),
						resource.TestCheckResourceAttrSet(assetResourceName, "id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "practice.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.2.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.2.id"),
						resource.TestCheckTypeSetElemAttr(assetResourceName, "urls.*", fmt.Sprintf("http://host/%s/path1", assetNameAttribute)),
						resource.TestCheckTypeSetElemAttr(assetResourceName, "urls.*", fmt.Sprintf("http://host/%s/path2", assetNameAttribute)),
						resource.TestCheckResourceAttrSet(assetResourceName, "tags.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "tags.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "redirect_to_https_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "access_log_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers.0.header_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers_id"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      assetResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: webAPIAssetUpdateFullConfig(assetNameAttribute, profileNameAttribute, trustedSourcesNameAttribute,
					practiceNameAttribute, logTriggerNameAttribute, exceptionsNameAttribute, anotherProfileNameAttribute,
					anotherTrustedSourcesNameAttribute, anotherLogTriggerNameAttribute, anotherExceptionsNameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(assetResourceName, map[string]string{
						"name":                                  assetNameAttribute,
						"%":                                     "32",
						"read_only":                             "false",
						"upstream_url":                          "some url 10",
						"urls.#":                                "2",
						"urls_ids.#":                            "2",
						"profiles.#":                            "1",
						"practice.#":                            "1",
						"practice.0.%":                          "5",
						"practice.0.triggers.#":                 "1",
						"practice.0.sub_practices_modes.IPS":    "Learn",
						"practice.0.sub_practices_modes.WebBot": "Inactive",
						"practice.0.sub_practices_modes.Snort":  "AccordingToPractice",
						"practice.0.main_mode":                  "Prevent",
						"source_identifier.0.%":                 "4",
						"source_identifier.1.%":                 "4",
						"source_identifier.2.%":                 "4",
						"source_identifier.2.values.#":          "2",
						"source_identifier.#":                   "3",
						"source_identifier.2.values_ids.#":      "2",
						"source_identifier.1.values_ids.#":      "2",
						"source_identifier.1.values.#":          "2",
						"source_identifier.0.values.#":          "2",
						"source_identifier.0.values_ids.#":      "2",
						"proxy_setting.#":                       "3",
						"proxy_setting.0.%":                     "3",
						"proxy_setting.1.%":                     "3",
						"proxy_setting.2.%":                     "3",
						"class":                                 "workload",
						"category":                              "cloud",
						"group":                                 "",
						"order":                                 "",
						"kind":                                  "",
						"family":                                "Web API",
						"main_attributes":                       fmt.Sprintf("{\"applicationUrls\":\"http://host/%[1]s/path3;http://host/%[1]s/path4\"}", assetNameAttribute),
						"asset_type":                            "WebAPI",
						"intelligence_tags":                     "",
						"tags.#":                                "3",
						"tags.0.%":                              "3",
						"tags.1.%":                              "3",
						"tags.2.%":                              "3",

						"mtls.#":                                         "2",
						"mtls.0.filename":                                "newfile.crt",
						"mtls.0.certificate_type":                        ".der",
						"mtls.0.data":                                    "new cert data",
						"mtls.0.type":                                    "server",
						"mtls.0.enable":                                  "true",
						"mtls.1.filename":                                "newfile2.p12",
						"mtls.1.certificate_type":                        ".p12",
						"mtls.1.data":                                    "new cert data2",
						"mtls.1.type":                                    "client",
						"mtls.1.enable":                                  "false",
						"additional_instructions_blocks.#":               "2",
						"additional_instructions_blocks.0.type":          "location_instructions",
						"additional_instructions_blocks.0.enable":        "false",
						"additional_instructions_blocks.1.filename":      "server.json",
						"additional_instructions_blocks.1.filename_type": ".json",
						"additional_instructions_blocks.1.data":          "server data",
						"additional_instructions_blocks.1.type":          "server_instructions",
						"additional_instructions_blocks.1.enable":        "true",
						"redirect_to_https":                              "false",
						"access_log":                                     "false",
						"custom_headers.#":                               "2",
					}),
						resource.TestCheckResourceAttrSet(assetResourceName, "id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "practice.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.2.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "source_identifier.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "proxy_setting.2.id"),
						resource.TestCheckTypeSetElemAttr(assetResourceName, "urls.*", fmt.Sprintf("http://host/%s/path3", assetNameAttribute)),
						resource.TestCheckTypeSetElemAttr(assetResourceName, "urls.*", fmt.Sprintf("http://host/%s/path4", assetNameAttribute)),
						resource.TestCheckResourceAttrSet(assetResourceName, "tags.0.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "tags.1.id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "tags.2.id"),

						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.0.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.1.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.1.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "mtls.1.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.0.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.1.filename_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.1.data_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "additional_instructions_blocks.1.enable_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "redirect_to_https_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "access_log_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers.0.header_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers.1.header_id"),
						resource.TestCheckResourceAttrSet(assetResourceName, "custom_headers_id"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func webAPIAssetBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_api_asset" %[1]q {
	name = %[1]q
	urls = ["http://host/%[1]s/path1"]
}
`, name)
}

func webAPIAssetUpdateBasicConfig(assetName, profileName, trustedSourcesName,
	practiceName, logTriggerName, exceptionsName string) string {
	return fmt.Sprintf(`
resource "inext_web_api_asset" %[1]q {
	name = %[1]q
	urls = ["http://host/%[1]s/path3", "http://host/%[1]s/path2"]
	profiles        = [inext_appsec_gateway_profile.%[2]s.id]
  	upstream_url    = "some url 5"
	practice {
		main_mode = "Prevent"
		sub_practices_modes = {
		  IPS    = "AccordingToPractice"
		  WebBot = "AccordingToPractice"
		  Snort  = "Disabled"
		}
		id         = inext_web_api_practice.%[4]s.id
		triggers   = [inext_log_trigger.%[5]s.id]
	}
  	proxy_setting {
    	key   = "some key"
    	value = "some value"
  	}
  	proxy_setting {
    	key   = "another key"
    	value = "another value"
  	}
  	proxy_setting {
    	key   = "last key"
    	value = "last value"
  	}
  	source_identifier {
    	identifier = "SourceIP"
    	values     = ["value3"]
  	}
  	source_identifier {
    	identifier = "XForwardedFor"
    	values     = ["value2"]
  	}
  	source_identifier {
    	identifier = "HeaderKey"
    	values     = ["value1"]
  	}
	tags {
		key   = "tagkey1"
		value = "tagvalue1"
	}
	mtls {
		filename = "cert.pem"
		certificate_type = ".pem"
		data	 = "cert data"
		type = "client"
		enable = true
	}
	additional_instructions_blocks {
		filename = "location.json"
		filename_type = ".json"
		data	 = "location data"
		type = "location_instructions"
		enable = true
	}
	custom_headers {
		name   = "header1"
		value  = "value1"
	}
}

resource "inext_appsec_gateway_profile" %[2]q {
	name                          = %[2]q
	profile_sub_type              = "Aws"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	reverseproxy_upstream_timeout = 3600
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key7 = "Value7"
		Key8 = "Value8"
	}
	additional_settings = {
		Key5 = "Value5"
		Key6 = "Value6"
	}
}

resource "inext_trusted_sources" %[3]q {
	name                = %[3]q
	min_num_of_sources  = 10
	sources_identifiers = ["identifier4", "identifier2", "identifier3"]
}

resource "inext_web_api_practice" %[4]q {
	name = %[4]q
	ips {
	  performance_impact    = "MediumOrLower"
	  severity_level        = "LowOrAbove"
	  protections_from_year = "2020"
	  high_confidence       = "Prevent"
	  medium_confidence     = "Detect"
	  low_confidence        = "Inactive"
	}
	api_attacks {
	  minimum_severity = "Critical"
	  advanced_setting {
		body_size            = 1000
		url_size             = 1000
		header_size          = 1000
		max_object_depth     = 1000
		illegal_http_methods = true
	  }
	}
}

resource "inext_log_trigger" %[5]q {
	name                             = %[5]q
	verbosity                        = "Extended" # enum of ["Minimal", "Standard", "Extended"]
	access_control_allow_events      = true
	access_control_drop_events       = true
	cef_ip_address                   = "10.0.0.1"
	cef_port                         = 81
	extend_logging                   = true
	extend_logging_min_severity      = "Critical" # enum of ["High", "Critical"]
	log_to_agent                     = true
	log_to_cef                       = true
	log_to_cloud                     = true
	log_to_syslog                    = true
	response_body                    = true
	response_code                    = true
	syslog_ip_address                = "10.0.0.2"
	syslog_port                      = 82
	threat_prevention_detect_events  = true
	threat_prevention_prevent_events = true
	web_body                         = true
	web_headers                      = false
	web_requests                     = true
	web_url_path                     = true
	web_url_query                    = true
}

resource "inext_exceptions" %[6]q {
	name = %[6]q
	exception {
		match {
		  key = "hostName"
		  value = ["www.google.com"]
		}
		action  = "drop"
	}
}
`, assetName, profileName, trustedSourcesName, practiceName, logTriggerName, exceptionsName)
}

func webAPIAssetFullConfig(assetName, profileName,
	trustedSourcesName, practiceName, logTriggerName, exceptionsName string) string {
	return fmt.Sprintf(`
resource "inext_web_api_asset" %[1]q {
	name = %[1]q
	urls = ["http://host/%[1]s/path1", "http://host/%[1]s/path2"]
	profiles        = [inext_appsec_gateway_profile.%[2]s.id]
	upstream_url    = "some url 5"
	practice {
	  main_mode = "Learn"
	  sub_practices_modes = {
		IPS    = "AccordingToPractice"
		WebBot = "AccordingToPractice"
		Snort  = "Disabled"
	  }
	  id         = inext_web_api_practice.%[4]s.id
	  triggers   = [inext_log_trigger.%[5]s.id]
	}

	proxy_setting {
	  key   = "some key"
	  value = "some value"
	}
	proxy_setting {
	  key   = "another key"
	  value = "another value"
	}
	proxy_setting {
	  key   = "last key"
	  value = "last value"
	}
	source_identifier {
	  identifier = "SourceIP"
	  values     = ["value3"]
	}
	source_identifier {
	  identifier = "XForwardedFor"
	  values     = ["value2"]
	}
	source_identifier {
	  identifier = "HeaderKey"
	  values     = ["value1"]
	}
	tags {
	  key   = "tagkey1"
	  value = "tagvalue1"
	}
	tags {
	  key   = "tagkey2"
	  value = "tagvalue2"
	}
	mtls {
		filename = "cert.der"
		certificate_type = ".der"
		data	 = "cert data"
		type = "client"
		enable = true
	}
	additional_instructions_blocks {
		filename = "location.json"
		filename_type = ".json"
		data	 = "location data"
		type = "location_instructions"
		enable = true
	}
	redirect_to_https = "true"
	access_log = "true"
	custom_headers {
		name   = "first header"
		value  = "first value"
	}
}

resource "inext_appsec_gateway_profile" %[2]q {
	name                          = %[2]q
	profile_sub_type              = "Aws"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	reverseproxy_upstream_timeout = 3600
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key7 = "Value7"
		Key8 = "Value8"
	}
	additional_settings = {
		Key5 = "Value5"
		Key6 = "Value6"
	}
}

resource "inext_trusted_sources" %[3]q {
	name                = %[3]q
	min_num_of_sources  = 10
	sources_identifiers = ["identifier4", "identifier2", "identifier3"]
}

resource "inext_web_api_practice" %[4]q {
	name = %[4]q
	ips {
	  performance_impact    = "MediumOrLower"
	  severity_level        = "LowOrAbove"
	  protections_from_year = "2020"
	  high_confidence       = "Prevent"
	  medium_confidence     = "Detect"
	  low_confidence        = "Inactive"
	}
	api_attacks {
	  minimum_severity = "Critical"
	  advanced_setting {
		body_size            = 1000
		url_size             = 1000
		header_size          = 1000
		max_object_depth     = 1000
		illegal_http_methods = true
	  }
	}
}

resource "inext_log_trigger" %[5]q {
	name                             = %[5]q
	verbosity                        = "Extended" # enum of ["Minimal", "Standard", "Extended"]
	access_control_allow_events      = true
	access_control_drop_events       = true
	cef_ip_address                   = "10.0.0.1"
	cef_port                         = 81
	extend_logging                   = true
	extend_logging_min_severity      = "Critical" # enum of ["High", "Critical"]
	log_to_agent                     = true
	log_to_cef                       = true
	log_to_cloud                     = true
	log_to_syslog                    = true
	response_body                    = true
	response_code                    = true
	syslog_ip_address                = "10.0.0.2"
	syslog_port                      = 82
	threat_prevention_detect_events  = true
	threat_prevention_prevent_events = true
	web_body                         = true
	web_headers                      = false
	web_requests                     = true
	web_url_path                     = true
	web_url_query                    = true
}

resource "inext_exceptions" %[6]q {
	name = %[6]q
	exception {
		match {
			operator = "or"
		  	operand {
				  operator = "not-equals"
				  key = "hostName"
				  value = ["www.google.com"]
		  	}
		  	operand {
				  operator = "in"
				  key = "url"
				  value = ["/login", "/login2"]
		  	}
		  	operand {
				  key = "sourceIdentifier"
				  value = ["1.1.1.1/24"]
		  	}
		}
		action  = "skip"
		comment = "test comment"
	}
	exception {
		match {
			operator = "and"
		  	operand {
				  key = "hostName"
				  value = ["www.facebook.com"]
		  	}
		  	operand {
				  key = "url"
				  value = ["/logout"]
		  	}
		  	operand {
				  key = "sourceIdentifier"
				  value = ["2.2.2.2/24"]
		  	}
		}
		action  = "drop"
		comment = "test comment"
	}
}
`, assetName, profileName, trustedSourcesName, practiceName, logTriggerName, exceptionsName)
}

func webAPIAssetUpdateFullConfig(assetName, profileName,
	trustedSourcesName, practiceName, logTriggerName, exceptionsName,
	anotherProfileName, anotherTrustedSourcesName, anotherLogTriggerName, anotherExcpetionsName string) string {
	return fmt.Sprintf(`
resource "inext_web_api_asset" %[1]q {
	name = %[1]q
	urls = ["http://host/%[1]s/path3", "http://host/%[1]s/path4"]
	profiles        = [inext_appsec_gateway_profile.%[7]s.id]
	upstream_url    = "some url 10"
	practice {
	  main_mode = "Prevent"
	  sub_practices_modes = {
		IPS    = "Learn"
		WebBot = "Inactive"
		Snort  = "AccordingToPractice"
	  }
	  id         = inext_web_api_practice.%[4]s.id
	  triggers   = [inext_log_trigger.%[9]s.id]
	}

	proxy_setting {
	  key   = "some key"
	  value = "some value2"
	}
	proxy_setting {
	  key   = "another key3"
	  value = "another value3"
	}
	proxy_setting {
	  key   = "last key"
	  value = "last value"
	}
	source_identifier {
	  identifier = "SourceIP"
	  values     = ["value4", "value5"]
	}
	source_identifier {
	  identifier = "XForwardedFor"
	  values     = ["value6", "value7"]
	}
	source_identifier {
	  identifier = "Cookie"
	  values     = ["value8", "value9"]
	}
	tags {
	  key   = "tagkey1"
	  value = "tagvalue2"
	}
	tags {
	  key   = "tagkey2"
      value = "tagvalue1"
	}
	tags {
	  key   = "tagkey3"
	  value = "tagvalue3"
	}
	mtls {
		filename = "newfile.crt"
		certificate_type = ".der"
		data	 = "new cert data"
		type = "server"
		enable = true
	}
	mtls {
		filename = "newfile2.p12"
		certificate_type = ".p12"
		data	 = "new cert data2"
		type = "client"
		enable = false
	}
	additional_instructions_blocks {
		filename = "location.json"
		filename_type = ".json"
		data	 = "location data"
		type = "location_instructions"
		enable = false
	}
	additional_instructions_blocks {
		filename = "server.json"
		filename_type = ".json"
		data	 = "server data"
		type = "server_instructions"
		enable = true
	}
	redirect_to_https = "false"
	access_log = "false"
	custom_headers {
		name   = "second header"
		value  = "second value"
	}
	custom_headers {
		name   = "first header"
		value  = "new first value"
	}
}

resource "inext_appsec_gateway_profile" %[2]q {
	name                          = %[2]q
	profile_sub_type              = "Aws"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	reverseproxy_upstream_timeout = 3600
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key7 = "Value7"
		Key8 = "Value8"
	}
	additional_settings = {
		Key5 = "Value5"
		Key6 = "Value6"
	}
}

resource "inext_appsec_gateway_profile" %[7]q {
	name                          = %[7]q
	profile_sub_type              = "Aws"
	upgrade_mode                  = "Scheduled"
	upgrade_time_schedule_type    = "DaysInWeek"
	upgrade_time_hour             = "12:00"
	upgrade_time_duration         = 10
	upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
	reverseproxy_upstream_timeout = 3600
	max_number_of_agents = 100
	reverseproxy_additional_settings = {
		Key7 = "Value7"
		Key8 = "Value8"
	}
	additional_settings = {
		Key5 = "Value5"
		Key6 = "Value6"
	}
}

resource "inext_trusted_sources" %[3]q {
	name                = %[3]q
	min_num_of_sources  = 10
	sources_identifiers = ["identifier4", "identifier2", "identifier3"]
}

resource "inext_trusted_sources" %[8]q {
	name                = %[8]q
	min_num_of_sources  = 10
	sources_identifiers = ["identifier4", "identifier2", "identifier3"]
}

resource "inext_web_api_practice" %[4]q {
	name = %[4]q
	ips {
	  performance_impact    = "MediumOrLower"
	  severity_level        = "LowOrAbove"
	  protections_from_year = "2020"
	  high_confidence       = "Prevent"
	  medium_confidence     = "Detect"
	  low_confidence        = "Inactive"
	}
	api_attacks {
	  minimum_severity = "Critical"
	  advanced_setting {
		body_size            = 1000
		url_size             = 1000
		header_size          = 1000
		max_object_depth     = 1000
		illegal_http_methods = true
	  }
	}
}

resource "inext_log_trigger" %[5]q {
	name                             = %[5]q
	verbosity                        = "Extended" # enum of ["Minimal", "Standard", "Extended"]
	access_control_allow_events      = true
	access_control_drop_events       = true
	cef_ip_address                   = "10.0.0.1"
	cef_port                         = 81
	extend_logging                   = true
	extend_logging_min_severity      = "Critical" # enum of ["High", "Critical"]
	log_to_agent                     = true
	log_to_cef                       = true
	log_to_cloud                     = true
	log_to_syslog                    = true
	response_body                    = true
	response_code                    = true
	syslog_ip_address                = "10.0.0.2"
	syslog_port                      = 82
	threat_prevention_detect_events  = true
	threat_prevention_prevent_events = true
	web_body                         = true
	web_headers                      = false
	web_requests                     = true
	web_url_path                     = true
	web_url_query                    = true
}

resource "inext_log_trigger" %[9]q {
	name                             = %[9]q
	verbosity                        = "Extended" # enum of ["Minimal", "Standard", "Extended"]
	access_control_allow_events      = true
	access_control_drop_events       = true
	cef_ip_address                   = "10.0.0.1"
	cef_port                         = 81
	extend_logging                   = true
	extend_logging_min_severity      = "Critical" # enum of ["High", "Critical"]
	log_to_agent                     = true
	log_to_cef                       = true
	log_to_cloud                     = true
	log_to_syslog                    = true
	response_body                    = true
	response_code                    = true
	syslog_ip_address                = "10.0.0.2"
	syslog_port                      = 82
	threat_prevention_detect_events  = true
	threat_prevention_prevent_events = true
	web_body                         = true
	web_headers                      = false
	web_requests                     = true
	web_url_path                     = true
	web_url_query                    = true
}

resource "inext_exceptions" %[6]q {
	name = %[6]q
	exception {
		match {
			operator = "or"
		  	operand {
				  operator = "not-equals"
				  key = "hostName"
				  value = ["www.google.com"]
		  	}
		  	operand {
				  operator = "in"
				  key = "url"
				  value = ["/login", "/login2"]
		  	}
		  	operand {
				  key = "sourceIdentifier"
				  value = ["1.1.1.1/24"]
		  	}
		}
		action  = "skip"
		comment = "test comment"
	}
	exception {
		match {
			operator = "and"
		  	operand {
				  key = "hostName"
				  value = ["www.facebook.com"]
		  	}
		  	operand {
				  key = "url"
				  value = ["/logout"]
		  	}
		  	operand {
				  key = "sourceIdentifier"
				  value = ["2.2.2.2/24"]
		  	}
		}
		action  = "drop"
		comment = "test comment"
	}
}

resource "inext_exceptions" %[10]q {
	name = %[10]q
	exception {
		match {
			operator = "and"
		  	operand {
				  key = "hostName"
				  value = ["www.facebook.com"]
		  	}
		  	operand {
				  key = "url"
				  value = ["/logout"]
		  	}
		  	operand {
				  key = "sourceIdentifier"
				  value = ["2.2.2.2/24"]
		  	}
		}
		action  = "drop"
		comment = "test comment"
	}
}
`, assetName, profileName, trustedSourcesName, practiceName, logTriggerName, exceptionsName,
		anotherProfileName, anotherTrustedSourcesName, anotherLogTriggerName, anotherExcpetionsName)
}
