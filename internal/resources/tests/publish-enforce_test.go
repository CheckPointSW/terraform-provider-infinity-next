package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccPublishEnforceBasic tests that publish and enforce work when set to true
func TestAccPublishEnforceBasic(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigBothTrue(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "0",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				// State resets to false after apply, so config with true will always show a diff
				ExpectNonEmptyPlan: true,
			},
			{
				Config: publishEnforceConfigBothTrue(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "0",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPublishEnforcePublishOnly tests that only publish works when set to true
func TestAccPublishEnforcePublishOnly(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigPublishOnly(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
						"enforce": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPublishEnforceEnforceOnly tests that only enforce works when set to true
func TestAccPublishEnforceEnforceOnly(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigEnforceOnly(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "0",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPublishEnforceWithProfileIds tests that enforce works with specific profile IDs
func TestAccPublishEnforceWithProfileIds(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	profileName1 := acctest.GenerateResourceName()
	profileName2 := acctest.GenerateResourceName()
	profileResourceName1 := "inext_appsec_gateway_profile." + profileName1
	profileResourceName2 := "inext_appsec_gateway_profile." + profileName2
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{profileResourceName1, profileResourceName2}),
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigWithProfileIds(profileName1, profileName2),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "2",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttrSet(resourceName, "profile_ids.0"),
						resource.TestCheckResourceAttrSet(resourceName, "profile_ids.1"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPublishEnforceWithEmptyProfileIds tests that enforce works with empty profile IDs (enforce all)
func TestAccPublishEnforceWithEmptyProfileIds(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigWithEmptyProfileIds(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "0",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPublishEnforceFalseNoOp tests that setting both to false does nothing (no-op)
func TestAccPublishEnforceFalseNoOp(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigBothFalse(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
						"enforce": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
			},
			{
				Config: publishEnforceConfigBothFalse(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
						"enforce": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
			},
		},
	})
}

// TestAccPublishEnforceRepeatedTrueTriggersEachTime tests that applying true multiple times
// triggers the operation each time (not just the first time)
func TestAccPublishEnforceRepeatedTrueTriggersEachTime(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigPublishOnly(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: publishEnforceConfigPublishOnly(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: publishEnforceConfigPublishOnly(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPublishEnforceDelete tests that delete works correctly (resource can be removed)
func TestAccPublishEnforceDelete(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigBothFalse(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: `# Empty config to trigger destroy`,
			},
		},
	})
}

// TestAccPublishEnforceDefaults tests that default values (false) work correctly
func TestAccPublishEnforceDefaults(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigDefaults(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "0",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
			},
		},
	})
}

// Config helper functions

func publishEnforceConfigDefaults() string {
	return `
resource "inext_publish_enforce" "trigger" {
}
`
}

func publishEnforceConfigBothTrue() string {
	return `
resource "inext_publish_enforce" "trigger" {
	publish = true
	enforce = true
}
`
}

func publishEnforceConfigPublishOnly() string {
	return `
resource "inext_publish_enforce" "trigger" {
	publish = true
	enforce = false
}
`
}

func publishEnforceConfigEnforceOnly() string {
	return `
resource "inext_publish_enforce" "trigger" {
	publish = false
	enforce = true
}
`
}

func publishEnforceConfigBothFalse() string {
	return `
resource "inext_publish_enforce" "trigger" {
	publish = false
	enforce = false
}
`
}

func publishEnforceConfigWithProfileIds(profileName1, profileName2 string) string {
	return fmt.Sprintf(`
resource "inext_appsec_gateway_profile" %[1]q {
	name             = %[1]q
	profile_sub_type = "Aws"
	max_number_of_agents = 10
}

resource "inext_appsec_gateway_profile" %[2]q {
	name             = %[2]q
	profile_sub_type = "Aws"
	max_number_of_agents = 10
}

resource "inext_publish_enforce" "trigger" {
	publish     = true
	enforce     = true
	profile_ids = [inext_appsec_gateway_profile.%[1]s.id, inext_appsec_gateway_profile.%[2]s.id]

	depends_on = [
		inext_appsec_gateway_profile.%[1]s,
		inext_appsec_gateway_profile.%[2]s,
	]
}
`, profileName1, profileName2)
}

func publishEnforceConfigWithEmptyProfileIds() string {
	return `
resource "inext_publish_enforce" "trigger" {
	publish     = true
	enforce     = true
	profile_ids = []
}
`
}
