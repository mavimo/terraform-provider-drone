package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mavimo/terraform-provider-drone/drone/utils"
)

func testCronConfigBasic(user, repo, name string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
    }

    resource "drone_cron" "cron" {
      repository = "${drone_repo.repo.repository}"
      name       = "%s"
			expr       = "@monthly"
			event      = "push"
    }
    `,
		user,
		repo,
		name,
	)
}

func TestCron(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCronDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCronConfigBasic(
					testDroneUser,
					"hook-test",
					"cron_job",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_cron.cron",
						"repository",
						fmt.Sprintf("%s/hook-test", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_cron.cron",
						"name",
						"cron_job",
					),
				),
			},
		},
	})
}

func testCronDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_cron" {
			continue
		}

		owner, repo, err := utils.ParseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		err = client.CronDelete(owner, repo, resource.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf(
				"Cron job still exists: %s/%s:%s",
				owner,
				repo,
				resource.Primary.Attributes["name"],
			)
		}
	}

	return nil
}
