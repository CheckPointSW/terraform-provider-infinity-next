package acctest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	resourceNameLength = 6
)

var resourceNameChartSet = acctest.CharSetAlpha + strings.ToUpper(acctest.CharSetAlpha)

func GenerateResourceName() string {
	return acctest.RandStringFromCharSet(resourceNameLength, resourceNameChartSet)
}

func MustReadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func CheckResourceDestroyed(resourcesNames []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, resourceName := range resourcesNames {
			resourceType, _, found := strings.Cut(resourceName, ".")
			if !found {
				return fmt.Errorf("resource name %s missing '.'", resourceName)
			}

			if rs, ok := s.RootModule().Resources[resourceName]; ok {
				rd := schema.ResourceData{}
				rd.SetId(rs.Primary.ID)
				diags := Provider.ResourcesMap[resourceType].ReadContext(context.Background(), &rd, Provider.Meta())
				if diags.HasError() {
					for _, d := range diags {
						if !strings.Contains(d.Summary, api.ErrorNotFound.Error()) {
							return errors.New(d.Summary)
						}
					}
				}
			}
		}

		return nil
	}
}
