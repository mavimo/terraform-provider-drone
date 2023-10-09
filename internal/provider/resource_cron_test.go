package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func TestAccDroneCronResource(t *testing.T) {
	cronName := "cron_job"
	orgName := "terraform"
	repoName := "repo-test-1"

	setupConfig := fmt.Sprintf(`
		resource "drone_repo" "repo" {
			repository = "%s/%s"
		}
	`, orgName, repoName)

	createConfig := setupConfig + fmt.Sprintf(`
		resource "drone_cron" "cron" {
			repository = "${drone_repo.repo.repository}"
			name       = "%s"
			expr       = "@monthly"
			event      = "push"
		}
	`, cronName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCronDestroy,
		Steps: []resource.TestStep{
			{
				Config: setupConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("drone_cron.cron", "repository", fmt.Sprintf("%s/%s", orgName, repoName)),
					resource.TestCheckResourceAttr("drone_cron.cron", "name", cronName),
					resource.TestCheckResourceAttr("drone_cron.cron", "branch", "master"),
					resource.TestCheckResourceAttr("drone_cron.cron", "disabled", "false"),
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
