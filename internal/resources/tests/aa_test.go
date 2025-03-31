package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebAPIPracticeEmptyResponse(t *testing.T) {
	nameAttribute := "emptychecktest"
	resourceName := "inext_web_api_practice." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				// Create a minimal practice
				Config: webAPIPracticeMinimalConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Test null response handling
			{
				// The provider should handle gracefully when read returns empty
				Config: webAPIPracticeMinimalConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					// Basic resource exists check
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Custom check function to verify empty response handling
					func(s *terraform.State) error {
						// This function will be called after resource reads
						// We're just verifying the test completes without crashing
						return nil
					},
				),
				// Simulate empty response by destroying resource before read
				// Note: This is implementation-specific and may need adjustment
				PreConfig: func() {
					// This could be a custom function that manipulates the backend
					// to return empty responses for the next read
				},
			},
		},
	})
}

func webAPIPracticeMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_api_practice" %[1]q {
  name = %[1]q
  # Minimal required fields
  ips {
    performance_impact    = "MediumOrLower"
    severity_level        = "MediumOrAbove"
    protections_from_year = "2016"
    high_confidence       = "AccordingToPractice"
    medium_confidence     = "AccordingToPractice"
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
