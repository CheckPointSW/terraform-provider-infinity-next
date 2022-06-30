package tests

import (
	"fmt"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebUserResponseBasic(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_web_user_response." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: webUserResponseBasicConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":               nameAttribute,
						"mode":               "BlockPage",
						"http_response_code": "403",
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
				Config: webUserResponseUpdateCreateSourceIdentifiersConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":         nameAttribute,
						"mode":         "Redirect",
						"redirect_url": "http://localhost:1234/test",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWebUserResponseFull(t *testing.T) {
	nameAttribute := acctest.GenerateResourceName()
	resourceName := "inext_web_user_response." + nameAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
		Steps: []resource.TestStep{
			{
				Config: webUserResponseWithIdentifiersConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":               nameAttribute,
						"mode":               "BlockPage",
						"http_response_code": "403",
						"message_title":      "some message title",
						"message_body":       "some message body",
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
				Config: webUserResponseUpdateSourceIdentifiersConfig(nameAttribute),
				Check: resource.ComposeTestCheckFunc(
					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
						"name":         nameAttribute,
						"mode":         "Redirect",
						"redirect_url": "http://localhost:1234/test",
						"x_event_id":   "true",
					}),
						resource.TestCheckResourceAttrSet(resourceName, "id"),
					)...,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func webUserResponseBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_user_response" %[1]q {
	name                = %[1]q
	mode  				= "BlockPage"
	http_response_code  = 403
}
`, name)
}

func webUserResponseWithIdentifiersConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_user_response" %[1]q {
	name                = %[1]q
	mode  				= "BlockPage"
	http_response_code 	= 403
	message_title 		= "some message title"
	message_body 		= "some message body"
}
`, name)
}

func webUserResponseUpdateCreateSourceIdentifiersConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_user_response" %[1]q {
	name                = %[1]q
	mode  				= "Redirect"
	redirect_url 		= "http://localhost:1234/test"
}
`, name)
}

func webUserResponseUpdateSourceIdentifiersConfig(name string) string {
	return fmt.Sprintf(`
resource "inext_web_user_response" %[1]q {
	name                = %[1]q
	mode  				= "Redirect"
	redirect_url 		= "http://localhost:1234/test"
	x_event_id 			= true
}
`, name)
}
