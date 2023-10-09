package drone

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneTemplateDataSource(t *testing.T) {
	testOrg := "terraform"
	templateName := "template-test.yaml"

	t.Run("check drone_repo data source with default values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_template" "test" {
				namespace = "%s"
				name = "%s"
				data = "---"
			}
		`, testOrg, templateName)

		configData := configResource + fmt.Sprintf(`
			data "drone_template" "test" {
				namespace = "%s"
				name = "%s"
			}
		`, testOrg, templateName)

		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configResource,
					Check:  resource.ComposeTestCheckFunc(),
				},
				{
					Config: configData,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.drone_template.test", "namespace", testOrg),
						resource.TestCheckResourceAttr("data.drone_template.test", "name", templateName),
						resource.TestCheckResourceAttr("data.drone_template.test", "data", "---"),
					),
				},
			},
		})
	})

	t.Run("check drone_template data source for non existing template", func(t *testing.T) {
		testMissingTemplateName := "missing.yaml"

		configResource := fmt.Sprintf(`
			resource "drone_template" "test" {
				namespace = "%s"
				name = "%s"
				data = "---"
			}
		`, testOrg, templateName)

		configData := configResource + fmt.Sprintf(`
			data "drone_template" "test" {
				namespace = "%s"
				name = "%s"
			}
		`, testOrg, testMissingTemplateName)

		resource.Test(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configData,
					ExpectError: regexp.MustCompile(
						fmt.Sprintf("Template %s/%s not found", testOrg, testMissingTemplateName),
					),
				},
			},
		})
	})
}
