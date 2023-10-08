package drone

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneTemplatesDataSource(t *testing.T) {
	testOrg := "terraform"
	templateName1 := "template-test-1.yaml"
	templateName2 := "template-test-2.yaml"
	templateName3 := "template-test-3.yaml"

	t.Run("check drone_repo data source with default values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
		resource "drone_template" "test1" {
			namespace = "%s"
			name = "%s"
			data = "---"
		}
		resource "drone_template" "test2" {
			namespace = "%s"
			name = "%s"
			data = "---"
		}
		resource "drone_template" "test3" {
			namespace = "%s"
			name = "%s"
			data = "---"
		}
		`, testOrg, templateName1, testOrg, templateName2, testOrg, templateName3)

		configData := configResource + fmt.Sprintf(`
			data "drone_templates" "test" {
				namespace = "%s"
			}
		`, testOrg)

		resource.Test(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configResource,
					Check:  resource.ComposeTestCheckFunc(),
				},
				{
					Config: configData,
					Check: resource.ComposeTestCheckFunc(
						// We check that id is computed
						resource.TestCheckResourceAttrSet("data.drone_templates.test", "id"),
						// We check that are correctly sorted
						resource.TestCheckResourceAttr("data.drone_templates.test", "names.0", templateName1),
						resource.TestCheckResourceAttr("data.drone_templates.test", "names.1", templateName2),
						resource.TestCheckResourceAttr("data.drone_templates.test", "names.2", templateName3),
						// We check that we get only the templates we created
						resource.TestCheckResourceAttr("data.drone_templates.test", "names.#", "3"),
					),
				},
			},
		})
	})
}
