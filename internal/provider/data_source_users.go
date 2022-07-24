package drone

import (
	"context"
	"sort"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving all Drone users",
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"logins": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List with the login name for all the active users in the Drone instance",
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	users, err := client.UserList()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to retrieve users",
			Detail:   err.Error(),
		})

		return diags
	}

	id := make([]string, 0)
	logins := make([]string, 0)

	for _, user := range users {
		id = append(id, user.Login)
		logins = append(logins, user.Login)
	}

	sort.Strings(logins)
	d.Set("logins", logins)

	d.SetId(utils.BuildChecksumID(id))

	return diags
}
