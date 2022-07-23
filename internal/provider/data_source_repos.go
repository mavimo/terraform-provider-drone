package drone

import (
	"context"
	"sort"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func dataSourceRepos() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving all repositories to which the user has explicit access in the host system",
		ReadContext: dataSourceReposRead,
		Schema: map[string]*schema.Schema{
			"repositories": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of names from repositories enabled in Drone",
			},
		},
	}
}

func dataSourceReposRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	repos, err := client.RepoListSync()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to retrieve repositories",
			Detail:   err.Error(),
		})

		return diags
	}

	id := make([]string, 0)
	slugs := make([]string, 0)

	for _, repo := range repos {
		id = append(id, repo.Slug)
		slugs = append(slugs, repo.Slug)
	}

	sort.Strings(slugs)
	d.Set("repositories", slugs)

	d.SetId(utils.BuildChecksumID(id))

	return diags
}
