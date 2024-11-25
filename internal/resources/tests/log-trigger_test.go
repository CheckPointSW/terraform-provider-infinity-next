package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLogTriggerBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_log_trigger." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: logTriggerBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                             nameAttribute,
						"threat_prevention_detect_events":  "true",
						"threat_prevention_prevent_events": "true",
						"log_to_cloud":                     "true",
						"web_url_path":                     "true",
						"web_url_query":                    "true",
						"extend_logging_min_severity":      "High",
						"extend_logging":                   "true",
						"verbosity":                        "Standard",
						"compliance_warnings":              "true",
						"compliance_violations":            "true",
						"syslog_protocol":                  "UDP",
						"cef_protocol":                     "UDP",
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
				Config: logTriggerUpdateBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                             nameAttribute,
						"threat_prevention_detect_events":  "false",
						"threat_prevention_prevent_events": "false",
						"log_to_cloud":                     "false",
						"log_to_agent":                     "true",
						"log_to_cef":                       "true",
						"log_to_syslog":                    "true",
						"web_url_path":                     "false",
						"web_url_query":                    "false",
						"web_body":                         "true",
						"web_headers":                      "false",
						"web_requests":                     "true",
						"extend_logging_min_severity":      "Critical",
						"extend_logging":                   "false",
						"access_control_allow_events":      "true",
						"access_control_drop_events":       "true",
						"cef_ip_address":                   "10.0.0.1",
						"cef_port":                         "81",
						"verbosity":                        "Extended",
						"response_body":                    "true",
						"response_code":                    "true",
						"syslog_ip_address":                "10.0.0.2",
						"syslog_port":                      "82",
						"compliance_warnings":              "false",
						"compliance_violations":            "false",
						"syslog_protocol":                  "TCP",
						"cef_protocol":                     "TCP",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
			},
		},
	})
}

func logTriggerBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_log_trigger" %[1]q {
	name = %[1]q
}
`, name)
}

func logTriggerUpdateBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_log_trigger" %[1]q {
	verbosity                      	 = "Extended"
	access_control_allow_events      = true
	access_control_drop_events       = true
	cef_ip_address                   = "10.0.0.1"
	cef_port                         = 81
	cef_protocol                     = "TCP"
	compliance_violations            = false
	compliance_warnings              = false
	extend_logging                   = false
	extend_logging_min_severity      = "Critical"
	log_to_agent                     = true
	log_to_cef                       = true
	log_to_cloud                     = false
	log_to_syslog                    = true
	name                             = %[1]q
	response_body                    = true
	response_code                    = true
	syslog_ip_address                = "10.0.0.2"
	syslog_protocol                  = "TCP"
	syslog_port                      = 82
	threat_prevention_detect_events  = false
	threat_prevention_prevent_events = false
	web_body                         = true
	web_headers                      = false
	web_requests                     = true
	web_url_path                     = false
	web_url_query                    = false
}
`, name)
}
