package drone

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneCronDataSource(t *testing.T) {
	testOrg := "terraform"
	testRepoName := "repo-test-1"
	cronName := "cron-1"

	t.Run("check cron_repo data source with default values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_repo" "repo" {
				repository = "%s/%s"
			}
			resource "drone_cron" "test" {
				repository = "%s/%s"
				name       = "%s"
				event      = "push"
				expr       = "@daily"
			}
		`, testOrg, testRepoName, testOrg, testRepoName, cronName)

		configData := configResource + fmt.Sprintf(`
			data "drone_cron" "test" {
				repository = "%s/%s"
				name       = "%s"
			}
		`, testOrg, testRepoName, cronName)

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
						resource.TestCheckResourceAttr("data.drone_cron.test", "repository", fmt.Sprintf("%s/%s", testOrg, testRepoName)),
						resource.TestCheckResourceAttr("data.drone_cron.test", "name", cronName),

						resource.TestCheckResourceAttr("data.drone_cron.test", "disabled", "false"),
						resource.TestCheckResourceAttr("data.drone_cron.test", "event", "push"),
						resource.TestCheckResourceAttr("data.drone_cron.test", "branch", "master"),
						resource.TestCheckResourceAttr("data.drone_cron.test", "expr", "@daily"),

						resource.TestCheckResourceAttrSet("data.drone_cron.test", "id"),
					),
				},
			},
		})
	})

	// t.Run("check drone_repo data source with custom values", func(t *testing.T) {
	// 	configResource := fmt.Sprintf(`
	// 		resource "drone_repo" "repo" {
	// 			repository     = "%s/%s"
	// 			cancel_pulls   = true
	// 			cancel_push    = true
	// 			cancel_running = true
	// 			configuration  = ".drone.jsonnet"
	// 			ignore_forks   = true
	// 			ignore_pulls   = true
	// 			protected      = true
	// 			timeout        = 120
	// 			trusted        = true
	// 		}
	// 	`, testOrg, testRepoName)

	// 	configData := configResource + fmt.Sprintf(`
	// 		data "drone_repo" "test" {
	// 			repository = "%s/%s"
	// 		}
	// 	`, testOrg, testRepoName)

	// 	resource.Test(t, resource.TestCase{
	// 		Providers: testProviders,
	// 		Steps: []resource.TestStep{
	// 			{
	// 				Config: configResource,
	// 				Check:  resource.ComposeTestCheckFunc(),
	// 			},
	// 			{
	// 				Config: configData,
	// 				Check: resource.ComposeTestCheckFunc(
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_pulls", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_push", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_running", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "configuration", ".drone.jsonnet"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_forks", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_pulls", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "protected", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "repository", fmt.Sprintf("%s/%s", testOrg, testRepoName)),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "timeout", "120"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "trusted", "true"),
	// 					resource.TestCheckResourceAttr("data.drone_repo.test", "visibility", "private"),
	// 				),
	// 			},
	// 		},
	// 	})
	// })

	t.Run("check drone_repo data source for non existing repositories", func(t *testing.T) {
		testOrg := "terraform"
		testRepo := "repo-test-1"
		testMissingCronoName := "missing"

		configData := fmt.Sprintf(`
			resource "drone_repo" "repo" {
				repository = "%s/%s"
			}

			data "drone_cron" "test" {
				repository = "%s/%s"
				name       = "%s"
			}
		`, testOrg, testRepo, testOrg, testRepo, testMissingCronoName)

		resource.Test(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configData,
					ExpectError: regexp.MustCompile(
						fmt.Sprintf("Error: failed to read Drone Cron: %s/%s/%s", testOrg, testRepo, testMissingCronoName),
					),
				},
			},
		})
	})
}
