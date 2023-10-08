package drone

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateConfigBasic(
					"foo",
					"bar.baz",
					"lorem: ipsum",
				),
				ExpectError: regexp.MustCompile("Template name bar.baz do not have a valid extension"),
			},
			{
				Config: testTemplateConfigBasic(
					"foo",
					"bar.jsonnet",
					"lorem: ipsum",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_template.template",
						"namespace",
						"foo",
					),
					resource.TestCheckResourceAttr(
						"drone_template.template",
						"name",
						"bar.jsonnet",
					),
					resource.TestCheckResourceAttr(
						"drone_template.template",
						"data",
						"lorem: ipsum",
					),
				),
			},
		},
	})
}

func testTemplateConfigBasic(namespace, name, data string) string {
	return fmt.Sprintf(`
	resource "drone_template" "template" {
		namespace = "%s"
		name      = "%s"
		data      = "%s"
	}
	`, namespace, name, data)
}

func testTemplateDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_template" {
			continue
		}

		namespace := resource.Primary.Attributes["namespace"]
		name := resource.Primary.Attributes["name"]

		err := client.TemplateDelete(namespace, name)
		if err == nil {
			return fmt.Errorf("namespace Template still exists: %s/%s", namespace, name)
		}
	}

	return nil
}
