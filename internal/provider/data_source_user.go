package drone

import (
	"context"
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving a Drone user",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the user active?",
			},
			"admin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the user an admin?",
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Login name",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"machine": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the user a machine?",
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	login := d.Get("login").(string)

	user, err := client.User(login)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to read Drone user with id: %s", login),
			Detail:   err.Error(),
		})

		return diags
	}

	d.Set("active", user.Active)
	d.Set("admin", user.Admin)
	d.Set("email", user.Email)
	d.Set("login", user.Login)
	d.Set("machine", user.Machine)

	d.SetId(login)

	return diags
}
