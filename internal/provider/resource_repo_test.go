package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func testRepoConfigBasic(user, repo string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
    }
    `, user, repo)
}

func TestAccDroneRepoResource(t *testing.T) {
	repoName := "repo-test-1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRepoConfigBasic(testDroneUser, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"repository",
						fmt.Sprintf("%s/%s", testDroneUser, repoName),
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
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"cancel_pulls",
						"false",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"cancel_push",
						"false",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"cancel_running",
						"false",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"ignore_forks",
						"false",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"ignore_pulls",
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

		repositories, _ := client.RepoList()

		for _, repository := range repositories {
			if (repository.Namespace == owner) && (repository.Name == repo) {
				err = client.RepoDisable(owner, repo)
				if err != nil {
					return fmt.Errorf("repo still exists: %s/%s", owner, repo)
				}
			}
		}
	}

	return nil
}
