package drone

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDroneUserSelfDataSource(t *testing.T) {
	t.Run("check drone_user_self data source for a logged user", func(t *testing.T) {
		configData := `
			data "drone_user_self" "test" {
			}
		`

		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: configData,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.drone_user_self.test", "active", "true"),
						resource.TestCheckResourceAttr("data.drone_user_self.test", "admin", "true"),
						resource.TestCheckResourceAttr("data.drone_user_self.test", "email", "terraform@example.com"),
						resource.TestCheckResourceAttr("data.drone_user_self.test", "id", "terraform"),
						resource.TestCheckResourceAttr("data.drone_user_self.test", "login", "terraform"),
						resource.TestCheckResourceAttr("data.drone_user_self.test", "machine", "false"),
					),
				},
			},
		})
	})

	t.Run("check drone_user_self data source error for a non-logged user", func(t *testing.T) {
		configData := `
			provider "drone" {
				token = "non-existing-token"
				alias = "notoken"
			}

			data "drone_user_self" "test" {
				provider = drone.notoken
			}
		`

		resource.Test(t, resource.TestCase{
			// ProviderFactories: map[string]func() (*schema.Provider, error){
			// 	"drone": func() (*schema.Provider, error) {
			// 		p := Provider()
			// 		p.Configure()
			// 		return nil, nil
			// 	},
			// },
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config:      configData,
					ExpectError: regexp.MustCompile("client error 401"),
				},
			},
		})
	})
}
