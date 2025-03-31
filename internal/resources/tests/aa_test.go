package tests

import (
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebAPIPracticeNullFileSecurityID(t *testing.T) {
	nameAttribute := "nullfsidcheck"
	resourceName := "inext_web_api_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				// Create a basic practice without explicit file_security block
				Config: webAPIPracticeBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					// First verify the resource was created successfully
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", nameAttribute),

					// Check that file_security block was populated with defaults
					resource.TestCheckResourceAttr(resourceName, "file_security.#", "1"),

					// Verify other critical attributes are set (these would fail if we had a crash)
					resource.TestCheckResourceAttrSet(resourceName, "ips.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_attacks.0.advanced_setting.0.id"),

					// This field should exist with a value - even if ID was null in API response
					resource.TestCheckResourceAttrSet(resourceName, "file_security.0.id"),

					// Verify default values are set for file_security
					resource.TestCheckResourceAttr(resourceName, "file_security.0.severity_level", "MediumOrAbove"),
					resource.TestCheckResourceAttr(resourceName, "file_security.0.high_confidence", "AccordingToPractice"),
					resource.TestCheckResourceAttr(resourceName, "file_security.0.medium_confidence", "AccordingToPractice"),
					resource.TestCheckResourceAttr(resourceName, "file_security.0.low_confidence", "Detect"),
				),
			},
			// Do an update to exercise the read path again
			{
				Config: webAPIPracticeUpdateBasicConfig(nameAttribute, schemaValidationFilename, schemaValidationData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "file_security.0.id"),
					resource.TestCheckResourceAttr(resourceName, "file_security.0.severity_level", "Critical"),
				),
			},
			// Verify import still works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
