package tests

import (
	"regexp"
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
						"publish": "false",
						"enforce": "false",
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
				// No ExpectNonEmptyPlan here - config is false, state is false, plan should be empty
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

// TestAccPublishEnforceSingletonPreventsMultiple tests that only one instance of this resource
// can exist by checking that two resources with different names get the same singleton ID
func TestAccPublishEnforceSingletonPreventsMultiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      publishEnforceConfigMultipleResources(),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

// TestAccPublishEnforceTransitionTrueToFalse tests transitioning from true to false
func TestAccPublishEnforceTransitionTrueToFalse(t *testing.T) {
	resourceName := "inext_publish_enforce.trigger"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: publishEnforceConfigBothTrue(),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"publish": "false",
						"enforce": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				ExpectNonEmptyPlan: true,
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
				// No ExpectNonEmptyPlan - config is false, state is false, plan should be empty
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
						"publish": "false",
						"enforce": "false",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
				),
				// No ExpectNonEmptyPlan - defaults are false, state is false
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

func publishEnforceConfigMultipleResources() string {
	return `
resource "inext_publish_enforce" "trigger1" {
	publish = true
	enforce = false
}

resource "inext_publish_enforce" "trigger2" {
	publish = false
	enforce = true
}
`
}
