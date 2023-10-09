package drone

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneUsersDataSource(t *testing.T) {
	machineName := "machine-test"
	userName := "user-test"

	t.Run("check drone_user data source with an existing user", func(t *testing.T) {
		configResource := fmt.Sprintf(`
			resource "drone_user" "user" {
				login  = "%s"
				active = true
				admin  = false
			}
			resource "drone_user" "machine" {
				login   = "%s"
				active  = true
				admin   = false
				machine = true
			}
		`, userName, machineName)

		configData := configResource + `
			data "drone_users" "test" {
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
						resource.TestCheckResourceAttrSet("data.drone_users.test", "id"),
						// We check that are correctly sorted
						resource.TestCheckResourceAttr("data.drone_users.test", "logins.0", machineName),
						resource.TestCheckResourceAttr("data.drone_users.test", "logins.1", "terraform"),
						resource.TestCheckResourceAttr("data.drone_users.test", "logins.2", userName),
						// We check that we get only the users we created (plus the admin)
						resource.TestCheckResourceAttr("data.drone_users.test", "logins.#", "3"),
					),
				},
			},
		})
	})
}
