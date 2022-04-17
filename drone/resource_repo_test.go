package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mavimo/terraform-provider-drone/drone/utils"
)

func testRepoConfigBasic(user, repo string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
    }
    `, user, repo)
}

func TestRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRepoConfigBasic(testDroneUser, "hook-test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"repository",
						fmt.Sprintf("%s/hook-test", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"visibility",
						"private",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"timeout",
						"60",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"protected",
						"false",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"trusted",
						"false",
					),
				),
			},
		},
	})
}

func testRepoDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_repo" {
			continue
		}

		owner, repo, err := utils.ParseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		repositories, err := client.RepoList()

		for _, repository := range repositories {
			if (repository.Namespace == owner) && (repository.Name == repo) {
				err = client.RepoDisable(owner, repo)
				if err != nil {
					return fmt.Errorf("Repo still exists: %s/%s", owner, repo)
				}
			}
		}
	}

	return nil
}
