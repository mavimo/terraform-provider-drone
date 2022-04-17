package drone

import (
	"context"
	"fmt"
	"regexp"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mavimo/terraform-provider-drone/drone/utils"
)

func resourceRepo() *schema.Resource {
	return &schema.Resource{
		Description: "Activate and configure a repository.",
		Schema: map[string]*schema.Schema{
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
			"trusted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Repository is trusted",
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
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "private",
				Description: "Repository visibility",
			},
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     ".drone.yml",
				Description: "Drone Configuration file",
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

	repository = &drone.RepoPatch{
		Config:     &config,
		Protected:  &protected,
		Trusted:    &trusted,
		Timeout:    &timeout,
		Visibility: &visibility,
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

	return diag.Diagnostics{}
}
