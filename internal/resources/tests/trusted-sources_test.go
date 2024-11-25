package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTrustedSourcesBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_trusted_sources." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: trustedSourcesBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":               nameAttribute,
						"min_num_of_sources": "1",
						"visibility":         "Shared",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
					)...,
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: trustedSourcesUpdateCreateSourceIdentifiersConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"visibility":                "Local",
						"min_num_of_sources":        "2",
						"sources_identifiers.#":     "3",
						"sources_identifiers_ids.#": "3",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier3"),
					)...,
				),
			},
		},
	})
}

func TestAccTrustedSourcesFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_trusted_sources." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: trustedSourcesWithIdentifiersConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"visibility":                "Shared",
						"min_num_of_sources":        "1",
						"sources_identifiers.#":     "3",
						"sources_identifiers_ids.#": "3",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier2"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier3"),
					)...,
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: trustedSourcesUpdateSourceIdentifiersConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":                      nameAttribute,
						"visibility":                "Local",
						"min_num_of_sources":        "2",
						"sources_identifiers.#":     "4",
						"sources_identifiers_ids.#": "4",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier3"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier4"),
						resource.TestCheckTypeSetElemAttr(resourceName, "sources_identifiers.*", "identifier5"),
					)...,
				),
			},
		},
	})
}

func trustedSourcesBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_trusted_sources" %[1]q {
	name                = %[1]q
	min_num_of_sources  = 1
}
`, name)
}

func trustedSourcesWithIdentifiersConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_trusted_sources" %[1]q {
	name                = %[1]q
	visibility          = "Shared"
	min_num_of_sources  = 1
	sources_identifiers = ["identifier1", "identifier2", "identifier3"]
}
`, name)
}

func trustedSourcesUpdateCreateSourceIdentifiersConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_trusted_sources" %[1]q {
	name                = %[1]q
	visibility          = "Local"
	min_num_of_sources  = 2
	sources_identifiers = ["identifier1", "identifier2", "identifier3"]
}
`, name)
}

func trustedSourcesUpdateSourceIdentifiersConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_trusted_sources" %[1]q {
	name                = %[1]q
	visibility          = "Local"
	min_num_of_sources  = 2
	sources_identifiers = ["identifier1", "identifier3", "identifier4", "identifier5"]
}
`, name)
}
