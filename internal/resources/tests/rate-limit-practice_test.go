package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRateLimitPracticeBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_rate_limit_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: rateLimitPracticeEmptyConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":     nameAttribute,
						"category": "RateLimit",
						"default":  "false",
						"rule.#":   "0",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
					)...),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: rateLimitPracticeBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":          nameAttribute,
						"category":      "RateLimit",
						"default":       "false",
						"rule.#":        "1",
						"rule.0.uri":    "/api/v1/test",
						"rule.0.scope":  "Minute",
						"rule.0.limit":  "100",
						"rule.0.action": "Detect",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
					)...),
			},
		},
	})
}

func TestAccRateLimitPracticeFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_rate_limit_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: rateLimitPracticeFullConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":           nameAttribute,
						"category":       "RateLimit",
						"default":        "false",
						"rule.#":         "2",
						"rule.0.uri":     "/api/v1/test",
						"rule.0.scope":   "Minute",
						"rule.0.limit":   "100",
						"rule.0.action":  "Detect",
						"rule.1.uri":     "/api/v2/full",
						"rule.1.scope":   "Second",
						"rule.1.limit":   "50",
						"rule.1.comment": "Full test rule",
						"rule.1.action":  "Prevent",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
						resource.TestCheckResourceAttrSet(resourceName, "rule.1.id"),
					)...),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: rateLimitPracticeUpdateConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":          nameAttribute,
						"category":      "RateLimit",
						"default":       "false",
						"rule.#":        "1",
						"rule.0.uri":    "/api/v1/updated",
						"rule.0.scope":  "Minute",
						"rule.0.limit":  "200",
						"rule.0.action": "AccordingToPractice",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
					)...),
			},
		},
	})
}

func rateLimitPracticeEmptyConfig(name string) string {
	return fmt.Sprintf(`resource "inext_rate_limit_practice" %[1]q {
		name = %[1]q
	}`, name)
}

func rateLimitPracticeBasicConfig(name string) string {
	return fmt.Sprintf(`resource "inext_rate_limit_practice" %[1]q {
		name = %[1]q
		rule {
			uri = "/api/v1/test"
			scope = "Minute"
			limit = 100
			action = "Detect"
		}
	}`, name)
}

func rateLimitPracticeFullConfig(name string) string {
	return fmt.Sprintf(`resource "inext_rate_limit_practice" %[1]q {
		name = %[1]q
		rule {
			uri = "/api/v1/test"
			scope = "Minute"
			limit = 100
			action = "Detect"
		}
		rule {
			uri = "/api/v2/full"
			scope = "Second"
			limit = 50
			comment = "Full test rule"
			action = "Prevent"
		}
	}`, name)
}

func rateLimitPracticeUpdateConfig(name string) string {
	return fmt.Sprintf(`resource "inext_rate_limit_practice" %[1]q {
		name = %[1]q
		rule {
			uri = "/api/v1/updated"
			scope = "Minute"
			limit = 200
			action = "AccordingToPractice"
		}
	}`, name)
}
