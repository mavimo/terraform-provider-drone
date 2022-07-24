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

func dataSourceRepo() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving a Drone repository",
		ReadContext: dataSourceRepoRead,
		Schema: map[string]*schema.Schema{
			"cancel_pulls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Automatically cancel pending pull request builds",
			},
			"cancel_push": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Automatically cancel pending push builds",
			},
			"cancel_running": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Automatically cancel running builds if newer commit pushed",
			},
			"configuration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Drone Configuration file",
			},
			"ignore_forks": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Disable build for pull requests",
			},
			"ignore_pulls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Disable build for forks",
			},
			"protected": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Repository is protected",
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Repository name",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Build execution timeout in minutes",
			},
			"trusted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Repository is trusted",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Repository visibility",
			},
		},
	}
}

func dataSourceRepoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Refresh repository list
	if _, err := client.RepoListSync(); err != nil {
		return diag.FromErr(err)
	}

	repository := d.Get("repository").(string)
	owner, name, err := utils.ParseRepo(repository)
	if err != nil {
		return diag.FromErr(err)
	}

	repo, err := client.Repo(owner, name)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to read repo %s", repository),
			Detail:   err.Error(),
		})

		return diags
	}

	d.Set("cancel_pulls", repo.CancelPulls)
	d.Set("cancel_push", repo.CancelPush)
	d.Set("cancel_running", repo.CancelRunning)
	d.Set("configuration", repo.Config)
	d.Set("ignore_forks", repo.IgnoreForks)
	d.Set("ignore_pulls", repo.IgnorePulls)
	d.Set("protected", repo.Protected)
	d.Set("timeout", repo.Timeout)
	d.Set("trusted", repo.Trusted)
	d.Set("visibility", repo.Visibility)

	d.SetId(repository)

	return diags
}
