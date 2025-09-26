package drone

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserSelf() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving the currently authenticated Drone user",
		ReadContext: dataSourceUserSelfRead,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Login name",
			},
			"machine": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the user a machine?",
			},
		},
	}
}

func dataSourceUserSelfRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	user, err := client.Self()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read currently authenticated Drone user",
			Detail:   err.Error(),
		})

		return diags
	}

	d.Set("active", user.Active)
	d.Set("admin", user.Admin)
	d.Set("email", user.Email)
	d.Set("login", user.Login)
	d.Set("machine", user.Machine)

	d.SetId(user.Login)

	return diags
}
