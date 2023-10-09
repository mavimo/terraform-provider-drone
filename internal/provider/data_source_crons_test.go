package drone

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneCronsDataSource(t *testing.T) {
	t.Run("check drone_crons data source with default values", func(t *testing.T) {
		orgName := "terraform"
		repoName := "repo-test-1"
		cronName1 := "cron-test-1"
		cronName2 := "cron-test-2"
		cronName3 := "cron-test-3"

		configResource := fmt.Sprintf(`
			resource "drone_repo" "test" {
				repository = "%s/%s"
			}

			resource "drone_cron" "test1" {
				repository = "${drone_repo.test.repository}"
				name       = "%s"
				event      = "push"
				expr       = "@daily"
			}

			resource "drone_cron" "test2" {
				repository = "${drone_repo.test.repository}"
				name       = "%s"
				event      = "push"
				expr       = "@weekly"
			}

			resource "drone_cron" "test3" {
				repository = "${drone_repo.test.repository}"
				name       = "%s"
				event      = "push"
				expr       = "@monthly"
			}
		`, orgName, repoName, cronName1, cronName2, cronName3)

		configData := configResource + fmt.Sprintf(`
			data "drone_crons" "test" {
				repository = "%s/%s"
			}
		`, orgName, repoName)

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
						resource.TestCheckResourceAttrSet("data.drone_crons.test", "id"),
						// We check that are correctly sorted
						resource.TestCheckResourceAttr("data.drone_crons.test", "names.0", cronName1),
						resource.TestCheckResourceAttr("data.drone_crons.test", "names.1", cronName2),
						resource.TestCheckResourceAttr("data.drone_crons.test", "names.2", cronName3),
						// We check that we get only the repos we created
						resource.TestCheckResourceAttr("data.drone_crons.test", "names.#", "3"),
					),
				},
			},
		})
	})

	t.Run("check drone_crons data source with non existing repository", func(t *testing.T) {
		orgName := "terraform"
		testMissingRepoName := "missing-repo"

		configData := fmt.Sprintf(`
			data "drone_crons" "test" {
				repository = "%s/%s"
			}
		`, orgName, testMissingRepoName)

		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configData,
					ExpectError: regexp.MustCompile(
						fmt.Sprintf("Error: failed to read Drone Crons for repo: %s/%s", orgName, testMissingRepoName),
					),
				},
			},
		})
	})
}
