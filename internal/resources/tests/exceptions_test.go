package tests

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/acctest"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// )

// func TestAccExceptionBasic(t *testing.T) {
// 	nameAttribute := acctest.GenerateResourceName()
// 	resourceName := "inext_exceptions." + nameAttribute
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheck(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: exceptionsBasicConfig(nameAttribute),
// 				Check: resource.ComposeTestCheckFunc(
// 					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
// 						"name":        nameAttribute,
// 						"%":           "3",
// 						"exception.#": "0",
// 					}),
// 						resource.TestCheckResourceAttrSet(resourceName, "id"))...,
// 				),
// 			},
// 			{
// 				Config: exceptionsBasicConfigUpdateAddExceptionMatchSingleField(nameAttribute),
// 				Check: resource.ComposeTestCheckFunc(
// 					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
// 						"name":                          nameAttribute,
// 						"exception.0.action":            "drop",
// 						"exception.0.match.0.value.#":   "1",
// 						"exception.0.match.0.operator":  "equals",
// 						"exception.0.comment":           "",
// 						"%":                             "3",
// 						"exception.0.match.0.operand.#": "0",
// 						"exception.#":                   "1",
// 						"exception.0.match.0.value.0":   "www.google.com",
// 						"exception.0.match.0.key":       "hostName",
// 						"exception.0.match.0.%":         "4",
// 						"exception.0.match.#":           "1",
// 						"exception.0.%":                 "5",
// 					}),
// 						resource.TestCheckResourceAttrSet(resourceName, "id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
// 					)...,
// 				),
// 			},
// 			{
// 				Config: exceptionsBasicConfigUpdateChangeExceptionFields(nameAttribute),
// 				Check: resource.ComposeTestCheckFunc(
// 					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
// 						"name":                                    nameAttribute,
// 						"exception.0.action":                      "skip",
// 						"%":                                       "3",
// 						"exception.0.comment":                     "test comment",
// 						"exception.0.match.0.value.#":             "0",
// 						"exception.0.match.0.operand.0.value.#":   "1",
// 						"exception.0.match.0.operand.0.%":         "4",
// 						"exception.0.match.0.operator":            "and",
// 						"exception.0.match.0.%":                   "4",
// 						"exception.0.match.0.key":                 "",
// 						"exception.0.match.0.operand.2.value.#":   "1",
// 						"exception.#":                             "1",
// 						"exception.0.match.0.operand.1.operator":  "equals",
// 						"exception.0.match.0.operand.1.value.#":   "1",
// 						"exception.0.match.0.operand.2.%":         "4",
// 						"exception.0.match.0.operand.2.operand.#": "0",
// 						"exception.0.match.0.operand.0.operand.#": "0",
// 						"exception.0.match.0.operand.0.operator":  "equals",
// 						"exception.0.match.0.operand.1.%":         "4",
// 						"exception.0.match.#":                     "1",
// 						"exception.0.match.0.operand.1.operand.#": "0",
// 						"exception.0.match.0.operand.2.operator":  "equals",
// 						"exception.0.match.0.operand.#":           "3",
// 						"exception.0.%":                           "5",
// 					}),
// 						resource.TestCheckResourceAttrSet(resourceName, "id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
// 					)...,
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func TestAccExceptionWithExceptionBlock(t *testing.T) {
// 	nameAttribute := acctest.GenerateResourceName()
// 	resourceName := "inext_exceptions." + nameAttribute
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheck(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      acctest.CheckResourceDestroyed([]string{resourceName}),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: exceptionsConfigWithException(nameAttribute),
// 				Check: resource.ComposeTestCheckFunc(
// 					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
// 						"name":                                    nameAttribute,
// 						"exception.0.action":                      "skip",
// 						"%":                                       "3",
// 						"exception.0.comment":                     "test comment",
// 						"exception.0.match.0.value.#":             "0",
// 						"exception.0.match.0.operand.0.value.#":   "1",
// 						"exception.0.match.0.operand.0.%":         "4",
// 						"exception.0.match.0.operator":            "and",
// 						"exception.0.match.0.%":                   "4",
// 						"exception.0.match.0.key":                 "",
// 						"exception.0.match.0.operand.2.value.#":   "1",
// 						"exception.#":                             "1",
// 						"exception.0.match.0.operand.1.operator":  "equals",
// 						"exception.0.match.0.operand.1.value.#":   "1",
// 						"exception.0.match.0.operand.2.%":         "4",
// 						"exception.0.match.0.operand.2.operand.#": "0",
// 						"exception.0.match.0.operand.0.operand.#": "0",
// 						"exception.0.match.0.operand.0.operator":  "equals",
// 						"exception.0.match.0.operand.1.%":         "4",
// 						"exception.0.match.#":                     "1",
// 						"exception.0.match.0.operand.1.operand.#": "0",
// 						"exception.0.match.0.operand.2.operator":  "equals",
// 						"exception.0.match.0.operand.#":           "3",
// 						"exception.0.%":                           "5",
// 					}),
// 						resource.TestCheckResourceAttrSet(resourceName, "id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
// 					)...),
// 			},
// 			{
// 				Config: exceptionsConfigWithExceptionAddException(nameAttribute),
// 				Check: resource.ComposeTestCheckFunc(
// 					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
// 						"name":                                    nameAttribute,
// 						"exception.0.match.0.operand.1.key":       "sourceIdentifier",
// 						"exception.1.match.0.operand.1.%":         "4",
// 						"exception.0.match.0.operand.0.operand.#": "0",
// 						"exception.1.match.0.operand.1.key":       "sourceIdentifier",
// 						"exception.1.match.0.operand.1.value.0":   "1.1.1.1/24",
// 						"exception.1.match.0.operand.2.operand.#": "0",
// 						"exception.1.match.0.value.#":             "0",
// 						"exception.0.match.0.operand.2.value.0":   "/logout",
// 						"exception.1.match.0.operand.0.%":         "4",
// 						"exception.1.match.#":                     "1",
// 						"exception.0.match.0.operand.0.value.#":   "1",
// 						"%":                                       "3",
// 						"exception.0.match.0.operand.0.key":       "hostName",
// 						"exception.0.match.0.operand.1.operand.#": "0",
// 						"exception.0.match.0.operand.#":           "3",
// 						"exception.1.match.0.operand.2.value.1":   "/login2",
// 						"exception.0.match.0.operand.0.%":         "4",
// 						"exception.1.match.0.operand.0.key":       "hostName",
// 						"exception.1.action":                      "skip",
// 						"exception.1.match.0.operand.0.operator":  "not-equals",
// 						"exception.1.match.0.operand.#":           "3",
// 						"exception.1.match.0.%":                   "4",
// 						"exception.0.match.0.value.#":             "0",
// 						"exception.0.match.0.operand.0.value.0":   "www.facebook.com",
// 						"exception.0.match.0.operand.1.%":         "4",
// 						"exception.0.match.#":                     "1",
// 						"exception.1.match.0.operand.1.operator":  "equals",
// 						"exception.1.match.0.operand.1.value.#":   "1",
// 						"exception.1.match.0.operand.2.operator":  "in",
// 						"exception.1.match.0.operand.2.value.#":   "2",
// 						"exception.#":                             "2",
// 						"exception.0.action":                      "drop",
// 						"exception.1.match.0.operand.2.key":       "url",
// 						"exception.1.match.0.operand.2.value.0":   "/login",
// 						"exception.1.%":                           "5",
// 						"exception.0.match.0.operand.2.key":       "url",
// 						"exception.1.match.0.operand.0.value.#":   "1",
// 						"exception.0.match.0.key":                 "",
// 						"exception.0.match.0.operand.0.operator":  "equals",
// 						"exception.0.match.0.operand.2.operand.#": "0",
// 						"exception.0.match.0.operand.2.%":         "4",
// 						"exception.1.match.0.operand.0.operand.#": "0",
// 						"exception.1.match.0.key":                 "",
// 						"exception.0.match.0.operand.1.value.0":   "2.2.2.2/24",
// 						"exception.0.match.0.operand.1.operator":  "equals",
// 						"exception.0.match.0.operand.2.operator":  "equals",
// 						"exception.0.match.0.%":                   "4",
// 						"exception.1.comment":                     "test comment",
// 						"exception.0.comment":                     "test comment",
// 						"exception.0.match.0.operator":            "and",
// 						"exception.0.match.0.operand.1.value.#":   "1",
// 						"exception.0.match.0.operand.2.value.#":   "1",
// 						"exception.1.match.0.operand.0.value.0":   "www.google.com",
// 						"exception.1.match.0.operand.2.%":         "4",
// 						"exception.1.match.0.operator":            "or",
// 						"exception.0.%":                           "5",
// 						"exception.1.match.0.operand.1.operand.#": "0",
// 					}),
// 						resource.TestCheckResourceAttrSet(resourceName, "id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.1.id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.1.action_id"),
// 					)...),
// 			},
// 			{
// 				Config: exceptionsConfigWithExceptionRemoveException(nameAttribute),
// 				Check: resource.ComposeTestCheckFunc(
// 					append(acctest.ComposeTestCheckResourceAttrsFromMap(resourceName, map[string]string{
// 						"name":                                    nameAttribute,
// 						"exception.0.action":                      "drop",
// 						"exception.0.match.0.operand.1.operator":  "equals",
// 						"exception.0.match.0.operand.1.value.0":   "2.2.2.2/24",
// 						"exception.0.comment":                     "test comment",
// 						"exception.0.match.0.operand.2.key":       "url",
// 						"exception.#":                             "1",
// 						"exception.0.match.0.operand.0.value.#":   "1",
// 						"exception.0.match.0.operand.2.operand.#": "0",
// 						"exception.0.match.#":                     "1",
// 						"exception.0.match.0.operand.2.%":         "4",
// 						"exception.0.match.0.value.#":             "0",
// 						"exception.0.match.0.operand.0.key":       "hostName",
// 						"exception.0.match.0.operand.0.%":         "4",
// 						"exception.0.match.0.operand.1.value.#":   "1",
// 						"exception.0.match.0.operand.1.key":       "sourceIdentifier",
// 						"exception.0.match.0.operand.2.value.0":   "/logout",
// 						"exception.0.%":                           "5",
// 						"exception.0.match.0.operator":            "and",
// 						"exception.0.match.0.%":                   "4",
// 						"%":                                       "3",
// 						"exception.0.match.0.operand.0.operand.#": "0",
// 						"exception.0.match.0.operand.0.value.0":   "www.facebook.com",
// 						"exception.0.match.0.operand.1.operand.#": "0",
// 						"exception.0.match.0.operand.#":           "3",
// 						"exception.0.match.0.operand.0.operator":  "equals",
// 						"exception.0.match.0.operand.1.%":         "4",
// 						"exception.0.match.0.operand.2.value.#":   "1",
// 						"exception.0.match.0.operand.2.operator":  "equals",
// 					}),
// 						resource.TestCheckResourceAttrSet(resourceName, "id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.id"),
// 						resource.TestCheckResourceAttrSet(resourceName, "exception.0.action_id"),
// 					)...,
// 				),
// 			},
// 			{
// 				ResourceName: resourceName,
// 				ImportState:  true,
// 			},
// 		},
// 	})
// }

// func exceptionsBasicConfig(name string) string {
// 	return fmt.Sprintf(`
// resource "inext_exceptions" %[1]q {
// 	name                = %[1]q
// }
// `, name)
// }

// func exceptionsBasicConfigUpdateAddExceptionMatchSingleField(name string) string {
// 	return fmt.Sprintf(`
// resource "inext_exceptions" %[1]q {
// 	name                = %[1]q
// 	exception {
// 		match {
// 		  key = "hostName"
// 		  value = ["www.google.com"]
// 		}
// 		action  = "drop"
// 	}
// }
// `, name)
// }

// func exceptionsBasicConfigUpdateChangeExceptionFields(name string) string {
// 	return fmt.Sprintf(`
// resource "inext_exceptions" %[1]q {
// 	name                = %[1]q
// 	exception {
// 		match {
// 		  operator = "and"
// 		  operand {
// 			  key = "hostName"
// 			  value = ["www.facebook.com"]
// 		  }
// 		  operand {
// 			  key = "url"
// 			  value = ["/login"]
// 		  }
// 		  operand {
// 			  key = "sourceIdentifier"
// 			  value = ["1.1.1.1/24"]
// 		  }
// 		}
// 		action  = "skip"
// 		comment = "test comment"
// 	}
// }
// `, name)
// }

// func exceptionsConfigWithException(name string) string {
// 	return fmt.Sprintf(`
// resource "inext_exceptions" %[1]q {
// 	name                = %[1]q
// 	exception {
// 		match {
// 			operator = "and"
// 		  	operand {
// 				  key = "hostName"
// 				  value = ["www.google.com"]
// 		  	}
// 		  	operand {
// 				  key = "url"
// 				  value = ["/login"]
// 		  	}
// 		  	operand {
// 				  key = "sourceIdentifier"
// 				  value = ["1.1.1.1/24"]
// 		  	}
// 		}
// 		action  = "skip"
// 		comment = "test comment"
// 	}
// }
// `, name)
// }

// func exceptionsConfigWithExceptionAddException(name string) string {
// 	return fmt.Sprintf(`
// resource "inext_exceptions" %[1]q {
// 	name                = %[1]q
// 	exception {
// 		match {
// 			operator = "or"
// 		  	operand {
// 				  operator = "not-equals"
// 				  key = "hostName"
// 				  value = ["www.google.com"]
// 		  	}
// 		  	operand {
// 				  operator = "in"
// 				  key = "url"
// 				  value = ["/login", "/login2"]
// 		  	}
// 		  	operand {
// 				  key = "sourceIdentifier"
// 				  value = ["1.1.1.1/24"]
// 		  	}
// 		}
// 		action  = "skip"
// 		comment = "test comment"
// 	}
// 	exception {
// 		match {
// 			operator = "and"
// 		  	operand {
// 				  key = "hostName"
// 				  value = ["www.facebook.com"]
// 		  	}
// 		  	operand {
// 				  key = "url"
// 				  value = ["/logout"]
// 		  	}
// 		  	operand {
// 				  key = "sourceIdentifier"
// 				  value = ["2.2.2.2/24"]
// 		  	}
// 		}
// 		action  = "drop"
// 		comment = "test comment"
// 	}
// }
// `, name)
// }

// func exceptionsConfigWithExceptionRemoveException(name string) string {
// 	return fmt.Sprintf(`
// resource "inext_exceptions" %[1]q {
// 	name                = %[1]q
// 	exception {
// 		match {
// 			operator = "and"
// 		  	operand {
// 				  key = "hostName"
// 				  value = ["www.facebook.com"]
// 		  	}
// 		  	operand {
// 				  key = "url"
// 				  value = ["/logout"]
// 		  	}
// 		  	operand {
// 				  key = "sourceIdentifier"
// 				  value = ["2.2.2.2/24"]
// 		  	}
// 		}
// 		action  = "drop"
// 		comment = "test comment"
// 	}
// }
// `, name)
// }
