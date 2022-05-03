package acctest

import (
	"context"
	"sync"
	"testing"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	ProviderName = "inext"
)

// ProviderFactories is a static map containing only the main provider instance
//
// Use other ProviderFactories functions, such as FactoriesAlternate,
// for tests requiring special provider configurations.
var ProviderFactories map[string]func() (*schema.Provider, error)

// Provider is the "main" provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// PreCheck(t) must be called before using this provider instance.
var Provider *schema.Provider

// testAccProviderConfigure ensures Provider is only configured once
//
// The PreCheck(t) function is invoked for every test and this prevents
// extraneous reconfiguration to the same values each time. However, this does
// not prevent reconfiguration that may happen should the address of
// Provider be errantly reused in ProviderFactories.
var testAccProviderConfigure sync.Once

func init() {
	Provider = provider.Provider()

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	ProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderName: func() (*schema.Provider, error) { return provider.Provider(), nil },
	}
}

// PreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call AWS APIs are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func PreCheck(t *testing.T) {
	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderConfigure.Do(func() {
		if err := Provider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); err != nil {
			t.Fatal(err)
		}
	})
}

func ComposeTestCheckResourceAttrsFromMap(resourceName string, attributesWithValues map[string]string) []resource.TestCheckFunc {
	funcs := make([]resource.TestCheckFunc, 0, len(attributesWithValues))
	for attributeName, attributeValue := range attributesWithValues {
		funcs = append(funcs, resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue))
	}

	return funcs
}
