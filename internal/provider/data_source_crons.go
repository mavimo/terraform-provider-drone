package drone

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func dataSourceCrons() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving all Drone crons for a repo",
		ReadContext: dataSourceCronsRead,
		Schema: map[string]*schema.Schema{
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Repository name",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			"names": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List with cron name for a specific repository",
			},
		},
	}
}

func dataSourceCronsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	owner, repo, err := utils.ParseRepo(d.Get("repository").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	crons, err := client.CronList(owner, repo)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to read Drone Crons for repo: %s/%ss", owner, repo),
			Detail:   err.Error(),
		})

		return diags
	}

	ids := make([]string, 0)
	slugs := make([]string, 0)

	for _, cron := range crons {
		ids = append(ids, cron.Name)
		slugs = append(slugs, cron.Name)
	}

	sort.Strings(slugs)
	d.Set("names", slugs)
	d.Set("repository", fmt.Sprintf("%s/%s", owner, repo))

	d.SetId(utils.BuildChecksumID(ids))

	return diags
}
