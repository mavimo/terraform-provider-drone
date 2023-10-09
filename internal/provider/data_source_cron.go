package drone

import (
	"context"
	"fmt"
	"regexp"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func dataSourceCron() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving a Drone cron",
		ReadContext: dataSourceReadCron,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cron job name",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicate if the cron should be disabled.",
			},
			"event": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event for this cron job. Only allowed value is `push`.",
			},
			"branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The branch to run this cron job on.",
			},
			"expr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cron interval. Allowed values are `@daily`, `@weekly`, `@monthly`, and `@yearly`.",
			},
			"id": {
				Description: "The string representation of the cron.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceReadCron(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	owner, repo, err := utils.ParseRepo(d.Get("repository").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	cronName := d.Get("name").(string)
	cron, err := client.Cron(owner, repo, cronName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to read Drone Cron: %s/%s/%s", owner, repo, cronName),
			Detail:   err.Error(),
		})

		return diags
	}

	d.SetId(fmt.Sprintf("%d", cron.ID))

	d.Set("repository", fmt.Sprintf("%s/%s", owner, repo))
	d.Set("event", cron.Event)
	d.Set("branch", cron.Branch)
	d.Set("disabled", cron.Disabled)
	d.Set("expr", cron.Expr)
	d.Set("name", cron.Name)

	return diags
}
