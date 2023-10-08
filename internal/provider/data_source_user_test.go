package drone

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneUserDataSource(t *testing.T) {
	userName := "user-test"
	userEmail := "user-test@example.com"

	t.Run("check drone_user data source with an existing user", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_user" "test" {
				login  = "%s"
				active = true
				admin  = false
			}
		`, userName)

		configData := configResource + fmt.Sprintf(`
			data "drone_user" "test" {
				login = "%s"
			}
		`, userName)

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
						resource.TestCheckResourceAttr("data.drone_user.test", "active", "true"),
						resource.TestCheckResourceAttr("data.drone_user.test", "admin", "false"),
						resource.TestCheckResourceAttr("data.drone_user.test", "email", userEmail),
						resource.TestCheckResourceAttr("data.drone_user.test", "id", userName),
						resource.TestCheckResourceAttr("data.drone_user.test", "login", userName),
						resource.TestCheckResourceAttr("data.drone_user.test", "machine", "false"),
					),
				},
			},
		})
	})

	t.Run("check drone_user data source with a machine user", func(t *testing.T) {
		machineName := "machine-test"

		configResource := fmt.Sprintf(`
			resource "drone_user" "test" {
				login   = "%s"
				active  = true
				admin   = false
				machine = true
			}
		`, machineName)

		configData := configResource + fmt.Sprintf(`
			data "drone_user" "test" {
				login = "%s"
			}
		`, machineName)

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
						resource.TestCheckResourceAttr("data.drone_user.test", "active", "true"),
						resource.TestCheckResourceAttr("data.drone_user.test", "admin", "false"),
						resource.TestCheckResourceAttr("data.drone_user.test", "email", ""),
						resource.TestCheckResourceAttr("data.drone_user.test", "id", machineName),
						resource.TestCheckResourceAttr("data.drone_user.test", "login", machineName),
						resource.TestCheckResourceAttr("data.drone_user.test", "machine", "true"),
					),
				},
			},
		})
	})

	t.Run("check drone_user data source for non existing user", func(t *testing.T) {
		testMissingUserName := "missing"

		configData := fmt.Sprintf(`
			data "drone_user" "test" {
				login = "%s"
			}
		`, testMissingUserName)

		resource.Test(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config:      configData,
					ExpectError: regexp.MustCompile("client error 404"), // TODO: find a better error to catch messages
				},
			},
		})
	})
}
