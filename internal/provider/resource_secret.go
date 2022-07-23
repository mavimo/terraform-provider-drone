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

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a repository secret.",
		Schema: map[string]*schema.Schema{
			"allow_on_pull_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Allow retrieving the secret on pull requests",
			},
			"allow_push_on_pull_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow retrieving the secret for push on pull requests",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Secret name",
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
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Secret value",
			},
			"id": {
				Description: "The string representation of the secret.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Exists:        resourceSecretExists,
	}
}

func resourceSecretCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.SecretCreate(owner, repo, createSecret(data))

	data.Set("value", data.Get("value").(string))

	return readSecret(ctx, data, owner, repo, secret, err)
}

func resourceSecretRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, name, err := utils.ParseId(data.Id(), "secret_password")

	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.Secret(owner, repo, name)
	if err != nil {
		return diag.Errorf("failed to read Drone Secret: %s/%s/%s", owner, repo, name)
	}

	return readSecret(ctx, data, owner, repo, secret, err)
}

func resourceSecretUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.SecretUpdate(owner, repo, createSecret(data))

	data.Set("value", data.Get("value").(string))

	return readSecret(ctx, data, owner, repo, secret, err)
}

func resourceSecretDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, name, err := utils.ParseId(data.Id(), "secret_password")

	if err != nil {
		return diag.FromErr(err)
	}

	err = client.SecretDelete(owner, repo, name)

	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceSecretExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, name, err := utils.ParseId(data.Id(), "secret_password")

	if err != nil {
		return false, err
	}

	secret, err := client.Secret(owner, repo, name)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Secret: %s/%s/%s", owner, repo, name)
	}

	exists := secret.Name == name

	return exists, err
}

func createSecret(data *schema.ResourceData) (secret *drone.Secret) {
	secret = &drone.Secret{
		Name:            data.Get("name").(string),
		Data:            data.Get("value").(string),
		PullRequest:     data.Get("allow_on_pull_request").(bool),
		PullRequestPush: data.Get("allow_push_on_pull_request").(bool),
	}

	return
}

func readSecret(ctx context.Context, data *schema.ResourceData, owner, repo string, secret *drone.Secret, err error) diag.Diagnostics {
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s/%s/%s", owner, repo, secret.Name))

	data.Set("repository", fmt.Sprintf("%s/%s", owner, repo))
	data.Set("name", secret.Name)
	data.Set("allow_on_pull_request", secret.PullRequest)
	data.Set("allow_push_on_pull_request", secret.PullRequestPush)

	return diag.Diagnostics{}
}
