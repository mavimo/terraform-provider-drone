package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testUserConfig = `
resource "drone_user" "octocat" {
	login = "octocat"
	admin = false
	active = true
	machine = true
}
`

func TestAccDroneUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_user.octocat",
						"login",
						"octocat",
					),
					func(s *terraform.State) error {

						if s.Modules[0].Resources["drone_user.octocat"].Primary.Attributes["token"] == "" {
							return fmt.Errorf("token is empty")
						}
						return nil
					},
				),
			},
			{
				Config: testUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_user.octocat",
						"login",
						"octocat",
					),
					func(s *terraform.State) error {

						if s.Modules[0].Resources["drone_user.octocat"].Primary.Attributes["token"] == "" {
							return fmt.Errorf("token is empty")
						}
						return nil
					},
				),
			},
		},
	})
}

func testUserDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_user" {
			continue
		}

		err := client.UserDelete(resource.Primary.Attributes["login"])

		if err == nil {
			return fmt.Errorf("User still exists: %s", resource.Primary.Attributes["login"])
		}
	}

	return nil
}
