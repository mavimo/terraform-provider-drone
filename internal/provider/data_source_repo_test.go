package drone

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneRepoDataSource(t *testing.T) {
	testOrg := "terraform"
	testRepoName := "repo-test-1"

	t.Run("check drone_repo data source with default values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_repo" "repo" {
				repository = "%s/%s"
			}
		`, testOrg, testRepoName)

		configData := configResource + fmt.Sprintf(`
			data "drone_repo" "test" {
				repository = "%s/%s"
			}
		`, testOrg, testRepoName)

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
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_pulls", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_push", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_running", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "configuration", ".drone.yml"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_forks", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_pulls", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "protected", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "repository", fmt.Sprintf("%s/%s", testOrg, testRepoName)),
						resource.TestCheckResourceAttr("data.drone_repo.test", "timeout", "60"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "trusted", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "visibility", "private"),
					),
				},
			},
		})
	})

	t.Run("check drone_repo data source with custom values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_repo" "repo" {
				repository     = "%s/%s"
				cancel_pulls   = true
				cancel_push    = true
				cancel_running = true
				configuration  = ".drone.jsonnet"
				ignore_forks   = true
				ignore_pulls   = true
				protected      = true
				timeout        = 120
				trusted        = true
			}
		`, testOrg, testRepoName)

		configData := configResource + fmt.Sprintf(`
			data "drone_repo" "test" {
				repository = "%s/%s"
			}
		`, testOrg, testRepoName)

		resource.Test(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configResource,
					Check:  resource.ComposeTestCheckFunc(),
				},
				{
					Config: configData,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_pulls", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_push", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_running", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "configuration", ".drone.jsonnet"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_forks", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_pulls", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "protected", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "repository", fmt.Sprintf("%s/%s", testOrg, testRepoName)),
						resource.TestCheckResourceAttr("data.drone_repo.test", "timeout", "120"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "trusted", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "visibility", "private"),
					),
				},
			},
		})
	})

	t.Run("check drone_repo data source for non existing repositories", func(t *testing.T) {
		testMissingRepoName := "missing"

		configData := fmt.Sprintf(`
			data "drone_repo" "test" {
				repository = "%s/%s"
			}
		`, testOrg, testMissingRepoName)

		resource.Test(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config:      configData,
					ExpectError: regexp.MustCompile(fmt.Sprintf("Failed to read repo %s/%s", testOrg, testMissingRepoName)),
				},
			},
		})
	})
}
