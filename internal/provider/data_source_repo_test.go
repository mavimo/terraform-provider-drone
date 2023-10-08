package drone

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneRepoDataSource(t *testing.T) {
	testDroneUser := "terraform"
	repoName := "hook-test"

	t.Run("check drone_repo data source with default values", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_repo" "repo" {
				repository = "%s/%s"
			}
		`, testDroneUser, repoName)

		configData := configResource + fmt.Sprintf(`
			data "drone_repo" "test" {
				repository = "%s/%s"
			}
		`, testDroneUser, repoName)

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
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_pulls", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_push", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "cancel_running", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "configuration", ".drone.yml"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_forks", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "ignore_pulls", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "protected", "false"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "repository", fmt.Sprintf("%s/%s", testDroneUser, repoName)),
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
		`, testDroneUser, repoName)

		configData := configResource + fmt.Sprintf(`
			data "drone_repo" "test" {
				repository = "%s/%s"
			}
		`, testDroneUser, repoName)

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
						resource.TestCheckResourceAttr("data.drone_repo.test", "repository", fmt.Sprintf("%s/%s", testDroneUser, repoName)),
						resource.TestCheckResourceAttr("data.drone_repo.test", "timeout", "120"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "trusted", "true"),
						resource.TestCheckResourceAttr("data.drone_repo.test", "visibility", "private"),
					),
				},
			},
		})
	})
}
