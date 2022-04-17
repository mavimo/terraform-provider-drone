package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mavimo/terraform-provider-drone/drone/utils"
)

func testSecretConfigBasic(user, repo, name, value string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
    }

    resource "drone_secret" "secret" {
      repository = "${drone_repo.repo.repository}"
      name       = "%s"
      value      = "%s"
    }
    `,
		user,
		repo,
		name,
		value,
	)
}

func TestSecret(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSecretConfigBasic(
					testDroneUser,
					"hook-test",
					"password",
					"1234567890",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"repository",
						fmt.Sprintf("%s/hook-test", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"name",
						"password",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"value",
						"1234567890",
					),
				),
			},
		},
	})
}

func testSecretDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_secret" {
			continue
		}

		owner, repo, err := utils.ParseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		err = client.SecretDelete(owner, repo, resource.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf(
				"Secret still exists: %s/%s:%s",
				owner,
				repo,
				resource.Primary.Attributes["name"],
			)
		}
	}

	return nil
}
