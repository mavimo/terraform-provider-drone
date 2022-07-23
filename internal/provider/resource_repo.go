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

func resourceRepo() *schema.Resource {
	return &schema.Resource{
		Description: "Activate and configure a repository.",
		Schema: map[string]*schema.Schema{
			"cancel_pulls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Automatically cancel pending pull request builds",
			},
			"cancel_push": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Automatically cancel pending push builds",
			},
			"cancel_running": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Automatically cancel running builds if newer commit pushed",
			},
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     ".drone.yml",
				Description: "Drone Configuration file",
			},
			"ignore_forks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable build for pull requests",
			},
			"ignore_pulls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable build for forks",
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Repository name",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			"protected": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Repository is protected",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Build execution timeout in minutes",
			},
			"trusted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Repository is trusted",
			},
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "private",
				Description: "Repository visibility",
			},
			"id": {
				Description: "The string representation of the repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceRepoCreate,
		ReadContext:   resourceRepoRead,
		UpdateContext: resourceRepoUpdate,
		DeleteContext: resourceRepoDelete,
		Exists:        resourceRepoExists,
	}
}

func resourceRepoCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	// Refresh repository list
	if _, err := client.RepoListSync(); err != nil {
		return diag.FromErr(err)
	}

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.Repo(owner, repo)

	if err != nil {
		return diag.FromErr(err)
	}
	repository, err := client.RepoUpdate(owner, repo, createRepo(ctx, data))

	if err != nil {
		return diag.FromErr(err)
	}
	if !resp.Active {
		_, err = client.RepoEnable(owner, repo)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readRepo(ctx, data, repository, err)
}

func resourceRepoRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	repository, err := client.Repo(owner, repo)
	if err != nil {
		return diag.Errorf("failed to read Drone Repo: %s/%s", owner, repo)
	}

	return readRepo(ctx, data, repository, err)
}

func resourceRepoUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	repository, err := client.RepoUpdate(owner, repo, createRepo(ctx, data))

	return readRepo(ctx, data, repository, err)
}

func resourceRepoDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	err = client.RepoDisable(owner, repo)

	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceRepoExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Id())

	if err != nil {
		return false, err
	}

	repository, err := client.Repo(owner, repo)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Repo: %s/%s", owner, repo)
	}

	exists := (repository.Namespace == owner) && (repository.Name == repo)

	return exists, err
}

func createRepo(ctx context.Context, data *schema.ResourceData) (repository *drone.RepoPatch) {
	trusted := data.Get("trusted").(bool)
	protected := data.Get("protected").(bool)
	timeout := int64(data.Get("timeout").(int))
	visibility := data.Get("visibility").(string)
	config := data.Get("configuration").(string)
	cancelPulls := data.Get("cancel_pulls").(bool)
	cancelPush := data.Get("cancel_push").(bool)
	cancelRunning := data.Get("cancel_running").(bool)
	ignoreForks := data.Get("ignore_forks").(bool)
	ignorePulls := data.Get("ignore_pulls").(bool)

	repository = &drone.RepoPatch{
		Config:        &config,
		Protected:     &protected,
		Trusted:       &trusted,
		Timeout:       &timeout,
		Visibility:    &visibility,
		IgnoreForks:   &ignoreForks,
		IgnorePulls:   &ignorePulls,
		CancelPulls:   &cancelPulls,
		CancelPush:    &cancelPush,
		CancelRunning: &cancelRunning,
	}

	return repository
}

func readRepo(ctx context.Context, data *schema.ResourceData, repository *drone.Repo, err error) diag.Diagnostics {
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s/%s", repository.Namespace, repository.Name))

	data.Set("repository", fmt.Sprintf("%s/%s", repository.Namespace, repository.Name))
	data.Set("trusted", repository.Trusted)
	data.Set("protected", repository.Protected)
	data.Set("timeout", repository.Timeout)
	data.Set("visibility", repository.Visibility)
	data.Set("configuration", repository.Config)
	data.Set("cancel_pulls", repository.CancelPulls)
	data.Set("cancel_push", repository.CancelPush)
	data.Set("cancel_running", repository.CancelRunning)
	data.Set("ignore_forks", repository.IgnoreForks)
	data.Set("ignore_pulls", repository.IgnorePulls)

	return diag.Diagnostics{}
}
