package tests

import (
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
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigWithProfileIds(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish":       "false",
						"enforce":       "false",
						"profile_ids.#": "2",
						"profile_ids.0": "profile-id-1",
						"profile_ids.1": "profile-id-2",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
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

func publishEnforceConfigWithProfileIds() string {
	return `
resource "inext_publish_enforce" "trigger" {
	publish     = true
	enforce     = true
	profile_ids = ["profile-id-1", "profile-id-2"]
}
`
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
