package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExceptionBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_exceptions." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: exceptionsBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name": nameAttribute,
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
			},
			{
				Config: exceptionsBasicConfigUpdateAddExceptionMatchSingleField(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                       nameAttribute,
						"exception.0.match.hostName": "www.google.com",
						"exception.0.action":         "drop",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
					)...,
				),
			},
			{
				Config: exceptionsBasicConfigUpdateChangeExceptionFields(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                               nameAttribute,
						"exception.0.match.hostName":         "www.facebook.com",
						"exception.0.match.url":              "/login",
						"exception.0.match.sourceIdentifier": "1.1.1.1/24",
						"exception.0.action":                 "skip",
						"exception.0.comment":                "test comment",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
					)...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExceptionWithExceptionBlock(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_exceptions." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: exceptionsConfigWithException(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                               nameAttribute,
						"exception.#":                        "1",
						"exception.0.match.hostName":         "www.google.com",
						"exception.0.match.url":              "/login",
						"exception.0.match.sourceIdentifier": "1.1.1.1/24",
						"exception.0.action":                 "skip",
						"exception.0.comment":                "test comment",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
					)...),
			},
			{
				Config: exceptionsConfigWithExceptionAddException(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                               nameAttribute,
						"exception.#":                        "2",
						"exception.0.match.hostName":         "www.google.com",
						"exception.0.match.url":              "/login",
						"exception.0.match.sourceIdentifier": "1.1.1.1/24",
						"exception.0.action":                 "skip",
						"exception.0.comment":                "test comment",
						"exception.1.match.hostName":         "www.facebook.com",
						"exception.1.match.url":              "/logout",
						"exception.1.match.sourceIdentifier": "2.2.2.2/24",
						"exception.1.action":                 "drop",
						"exception.1.comment":                "test comment",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.1.id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.1.action_id"),
					)...),
			},
			{
				Config: exceptionsConfigWithExceptionRemoveException(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                               nameAttribute,
						"exception.#":                        "1",
						"exception.0.match.hostName":         "www.facebook.com",
						"exception.0.match.url":              "/logout",
						"exception.0.match.sourceIdentifier": "2.2.2.2/24",
						"exception.0.action":                 "drop",
						"exception.0.comment":                "test comment",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
					)...,
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
		},
	})
}

func exceptionsBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_exceptions" %[1]q {
	name                = %[1]q
}
`, name)
}

func exceptionsBasicConfigUpdateAddExceptionMatchSingleField(name string) string {
	return fmt.Sprintf(`
resource "inext_exceptions" %[1]q {
	name                = %[1]q
	exception {
		match = {
		  hostName         = "www.google.com"
		}
		action  = "drop"
	  }
}
`, name)
}

func exceptionsBasicConfigUpdateChangeExceptionFields(name string) string {
	return fmt.Sprintf(`
resource "inext_exceptions" %[1]q {
	name                = %[1]q
	exception {
		match = {
		  hostName         = "www.facebook.com"
		  url              = "/login"
		  sourceIdentifier = "1.1.1.1/24"
		}
		action  = "skip"
		comment = "test comment"
	  }
}
`, name)
}

func exceptionsConfigWithException(name string) string {
	return fmt.Sprintf(`
resource "inext_exceptions" %[1]q {
	name                = %[1]q
	exception {
		match = {
		  hostName         = "www.google.com"
		  url              = "/login"
		  sourceIdentifier = "1.1.1.1/24"
		}
		action  = "skip"
		comment = "test comment"
	  }
}
`, name)
}

func exceptionsConfigWithExceptionAddException(name string) string {
	return fmt.Sprintf(`
resource "inext_exceptions" %[1]q {
	name                = %[1]q
	exception {
		match = {
		  hostName         = "www.google.com"
		  url              = "/login"
		  sourceIdentifier = "1.1.1.1/24"
		}
		action  = "skip"
		comment = "test comment"
	}
	exception {
		match = {
		  hostName         = "www.facebook.com"
		  url              = "/logout"
		  sourceIdentifier = "2.2.2.2/24"
		}
		action  = "drop"
		comment = "test comment"
	}
}
`, name)
}

func exceptionsConfigWithExceptionRemoveException(name string) string {
	return fmt.Sprintf(`
resource "inext_exceptions" %[1]q {
	name                = %[1]q
	exception {
		match = {
		  hostName         = "www.facebook.com"
		  url              = "/logout"
		  sourceIdentifier = "2.2.2.2/24"
		}
		action  = "drop"
		comment = "test comment"
	}
}
`, name)
}
