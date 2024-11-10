package tests

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	webAPIPracticeTestdataPath     = "testdata/web-api-practice"
	schemaValidationFilename       = "oasschema"
	schemaValidationFilenameUpdate = "oasschemaupdate"
)

var (
	schemaValidationData       = acctest.MustReadFile(path.Join(webAPIPracticeTestdataPath, schemaValidationFilename))
	schemaValidationDataUpdate = strings.Replace(schemaValidationData, "Intelligence", "update", 1)
)

func TestAccWebAPIPracticeBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_web_api_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: webAPIPracticeBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                           nameAttribute,
						"schema_validation.#":            "1",
						"ips.0.protections_from_year":    "2016",
						"ips.#":                          "1",
						"api_attacks.0.minimum_severity": "High",
						"ips.0.low_confidence":           "Detect",
						"ips.0.performance_impact":       "MediumOrLower",
						"ips.0.%":                        "7",
						//"schema_validation.0.filename":   "",
						"schema_validation.0.oas_schema.name": "",
						//"schema_validation.0.oas_schema.size": "",
						"schema_validation.0.oas_schema.data": "",
						"schema_validation.0.oas_schema.%":    "",
						"api_attacks.#":                       "1",
						"ips.0.severity_level":                "MediumOrAbove",
						//"schema_validation.0.data":       "",
						"practice_type":           "WebAPI",
						"default":                 "false",
						"ips.0.medium_confidence": "Prevent",
						"schema_validation.0.%":   "2",
						"category":                "ThreatPrevention",
						"api_attacks.0.%":         "3",
						"ips.0.high_confidence":   "Prevent",
						//"api_attacks.0.advanced_setting.0.body_size":            "1000000",
						//"api_attacks.0.advanced_setting.0.url_size":             "32768",
						//"api_attacks.0.advanced_setting.0.header_size":          "102400",
						//"api_attacks.0.advanced_setting.0.%":                    "6",
						//"api_attacks.0.advanced_setting.0.max_object_depth":     "40",
						//"api_attacks.0.advanced_setting.0.illegal_http_methods": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "schema_validation.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.advanced_setting.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: webAPIPracticeUpdateBasicConfig(nameAttribute, schemaValidationFilename, schemaValidationData),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                                                  nameAttribute,
						"schema_validation.0.filename":                          schemaValidationFilename,
						"schema_validation.0.data":                              schemaValidationData,
						"api_attacks.0.minimum_severity":                        "Critical",
						"ips.0.high_confidence":                                 "Detect",
						"practice_type":                                         "WebAPI",
						"ips.0.medium_confidence":                               "Detect",
						"ips.0.performance_impact":                              "LowOrLower",
						"api_attacks.0.advanced_setting.0.header_size":          "1000",
						"api_attacks.0.advanced_setting.0.illegal_http_methods": "true",
						"api_attacks.0.advanced_setting.0.body_size":            "1000",
						"api_attacks.0.advanced_setting.0.url_size":             "1000",
						"api_attacks.0.advanced_setting.0.%":                    "6",
						"api_attacks.0.advanced_setting.0.max_object_depth":     "1000",
						"api_attacks.0.advanced_setting.#":                      "1",
						"schema_validation.#":                                   "1",
						"default":                                               "false",
						"category":                                              "ThreatPrevention",
						"ips.0.low_confidence":                                  "Detect",
						"ips.0.protections_from_year":                           "2016",
						"ips.0.%":                                               "7",
						"schema_validation.0.%":                                 "3",
						"api_attacks.#":                                         "1",
						"ips.0.severity_level":                                  "LowOrAbove",
						"ips.#":                                                 "1",
						"api_attacks.0.%":                                       "3",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "schema_validation.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.advanced_setting.0.id"),
					)...,
				),
			},
		},
	})
}

func TestAccWebAPIPracticeFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_web_api_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: webAPIPracticeFullConfig(nameAttribute, schemaValidationFilename, schemaValidationData),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name": nameAttribute,
						//"schema_validation.0.filename":                          schemaValidationFilename,
						"schema_validation.0.data":       schemaValidationData,
						"api_attacks.0.minimum_severity": "Critical",
						"ips.0.high_confidence":          "Detect",
						"practice_type":                  "WebAPI",
						//"api_attacks.0.advanced_setting.0.url_size":             "1000",
						//"api_attacks.0.advanced_setting.0.%":                    "6",
						"ips.0.medium_confidence":  "Detect",
						"ips.0.performance_impact": "LowOrLower",
						//"api_attacks.0.advanced_setting.0.header_size":          "1000",
						//"api_attacks.0.advanced_setting.0.illegal_http_methods": "true",
						//"api_attacks.0.advanced_setting.0.body_size":            "1000",
						"schema_validation.#": "1",
						//"api_attacks.0.advanced_setting.0.max_object_depth":     "1000",
						"default":                          "false",
						"api_attacks.0.advanced_setting.#": "1",
						"category":                         "ThreatPrevention",
						"ips.0.low_confidence":             "Detect",
						"ips.0.protections_from_year":      "2016",
						"ips.0.%":                          "7",
						"schema_validation.0.%":            "3",
						"api_attacks.#":                    "1",
						"ips.0.severity_level":             "LowOrAbove",
						"ips.#":                            "1",
						"api_attacks.0.%":                  "3",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "schema_validation.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.advanced_setting.0.id"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: webAPIPracticeUpdateFullConfig(nameAttribute, schemaValidationFilenameUpdate, schemaValidationDataUpdate),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                         nameAttribute,
						"schema_validation.0.filename": schemaValidationFilenameUpdate,
						"schema_validation.0.data":     schemaValidationDataUpdate,
						"api_attacks.#":                "1",
						"category":                     "ThreatPrevention",
						"default":                      "false",
						"ips.0.high_confidence":        "Prevent",
						"api_attacks.0.advanced_setting.0.body_size":            "1001",
						"schema_validation.0.%":                                 "3",
						"api_attacks.0.minimum_severity":                        "High",
						"ips.0.protections_from_year":                           "2020",
						"ips.0.severity_level":                                  "Critical",
						"api_attacks.0.%":                                       "3",
						"ips.0.medium_confidence":                               "Inactive",
						"practice_type":                                         "WebAPI",
						"api_attacks.0.advanced_setting.0.illegal_http_methods": "false",
						"api_attacks.0.advanced_setting.0.%":                    "6",
						"ips.0.low_confidence":                                  "Detect",
						"ips.0.performance_impact":                              "MediumOrLower",
						"ips.0.%":                                               "7",
						"ips.#":                                                 "1",
						"api_attacks.0.advanced_setting.0.header_size":          "1003",
						"api_attacks.0.advanced_setting.0.max_object_depth":     "1004",
						"schema_validation.#":                                   "1",
						"api_attacks.0.advanced_setting.0.url_size":             "1002",
						"api_attacks.0.advanced_setting.#":                      "1",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "schema_validation.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.advanced_setting.0.id"),
					)...,
				),
			},
		},
	})
}

func webAPIPracticeBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_api_practice" %[1]q {
	name = %[1]q
	ips {
		performance_impact    = "MediumOrLower"
		severity_level        = "MediumOrAbove"
		protections_from_year = "2016"
		high_confidence       = "Prevent"
		medium_confidence     = "Prevent"
		low_confidence        = "Detect"
	}
	api_attacks {
		minimum_severity = "High"
		advanced_setting {
			body_size            = 1000000
			url_size             = 32768
			header_size          = 102400
			max_object_depth     = 40
			illegal_http_methods = false
		}
	}
}
`, name)
}

func webAPIPracticeUpdateBasicConfig(name, filename, data string) string {
	return fmt.Sprintf(`
resource "inext_web_api_practice" %[1]q {
	name                          = %[1]q
	ips {
		performance_impact    = "LowOrLower"   
		severity_level        = "LowOrAbove"
		protections_from_year = "2016"      
		high_confidence       = "Detect"    
		medium_confidence     = "Detect"    
		low_confidence        = "Detect"
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
	schema_validation {
		oas_schema {
			name = %[2]q
			data = %[3]q
		}
	}
}
`, name, filename, data)
}

func webAPIPracticeFullConfig(name, filename, data string) string {
	return fmt.Sprintf(`
resource "inext_web_api_practice" %[1]q {
	name                          = %[1]q
	ips {
		performance_impact    = "LowOrLower"   
		severity_level        = "LowOrAbove"
		protections_from_year = "2016"      
		high_confidence       = "Detect"    
		medium_confidence     = "Detect"    
		low_confidence        = "Detect" 
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
	schema_validation {
		oas_schema {
			name = %[2]q
			data = %[3]q
		}
	}
}
`, name, filename, data)
}

func webAPIPracticeUpdateFullConfig(name, filename, data string) string {
	return fmt.Sprintf(`
resource "inext_web_api_practice" %[1]q {
	name                          = %[1]q
	ips {
		performance_impact    = "MediumOrLower"   
		severity_level        = "Critical"
		protections_from_year = "2020"      
		high_confidence       = "Prevent"    
		medium_confidence     = "Inactive"    
		low_confidence        = "Detect"
	}
	api_attacks {
		minimum_severity = "High"
		advanced_setting {
			body_size            = 1001
			url_size             = 1002
			header_size          = 1003
			max_object_depth     = 1004
			illegal_http_methods = false
		}
	}
	schema_validation {
		oas_schema {
			name = %[2]q
			data = %[3]q
		}
	}
}
`, name, filename, data)
}
