package drone

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneReposDataSource(t *testing.T) {
	testOrg := "terraform"
	repoName1 := "repo-test-1"
	repoName2 := "repo-test-2"
	repoName3 := "repo-test-3"

	t.Run("check drone_repo data source with default values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
		resource "drone_repo" "test1" {
			repository = "%s/%s"
		}
		resource "drone_repo" "test2" {
			repository = "%s/%s"
		}
		resource "drone_repo" "test3" {
			repository = "%s/%s"
			
		}
		`, testOrg, repoName1, testOrg, repoName2, testOrg, repoName3)

		configData := configResource + `
			data "drone_repos" "test" {
			}
		`

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
						// We check that id is computed
						resource.TestCheckResourceAttrSet("data.drone_repos.test", "id"),
						// We check that are correctly sorted
						resource.TestCheckResourceAttr("data.drone_repos.test", "repositories.0", fmt.Sprintf("%s/%s", testOrg, repoName1)),
						resource.TestCheckResourceAttr("data.drone_repos.test", "repositories.1", fmt.Sprintf("%s/%s", testOrg, repoName2)),
						resource.TestCheckResourceAttr("data.drone_repos.test", "repositories.2", fmt.Sprintf("%s/%s", testOrg, repoName3)),
						// We check that we get only the repos we created
						resource.TestCheckResourceAttr("data.drone_repos.test", "repositories.#", "3"),
					),
				},
			},
		})
	})
}
