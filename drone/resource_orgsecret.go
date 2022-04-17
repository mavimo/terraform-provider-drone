package drone

import (
	"context"
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mavimo/terraform-provider-drone/drone/utils"
)

func resourceOrgSecret() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a org secret.",
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Description: "Organization name",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Secret name",
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "Secret value",
				Required:    true,
				Sensitive:   true,
				ForceNew:    false,
			},
			"allow_on_pull_request": {
				Type:        schema.TypeBool,
				Description: "Allow retrieving the secret on pull requests.",
				Optional:    true,
				ForceNew:    false,
			},
			"allow_push_on_pull_request": {
				Type:        schema.TypeBool,
				Description: "Allow pushing on pull requests",
				Optional:    true,
				ForceNew:    false,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceOrgSecretCreate,
		ReadContext:   resourceOrgSecretRead,
		UpdateContext: resourceOrgSecretUpdate,
		DeleteContext: resourceOrgSecretDelete,
		Exists:        resourceOrgSecretExists,
	}
}

func resourceOrgSecretCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace := data.Get("namespace").(string)
	secret, err := client.OrgSecretCreate(namespace, createOrgSecret(data))
	if err != nil {
		return diag.FromErr(err)
	}

	readOrgSecret(data, namespace, secret)

	return diag.Diagnostics{}
}

func resourceOrgSecretRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.OrgSecret(namespace, name)
	if err != nil {
		return diag.Errorf("failed to read Drone Org Secret: %s/%s", namespace, name)
	}

	readOrgSecret(data, namespace, secret)
	return nil
}

func resourceOrgSecretUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace, _, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.OrgSecretUpdate(namespace, createOrgSecret(data))
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("value", data.Get("value").(string))
	readOrgSecret(data, namespace, secret)

	return diag.Diagnostics{}
}

func resourceOrgSecretDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.OrgSecretDelete(namespace, name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceOrgSecretExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "secret_password")
	if err != nil {
		return false, err
	}

	secret, err := client.OrgSecret(namespace, name)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Org Secret: %s/%s", namespace, name)
	}

	exists := (secret.Name == name) && (err == nil)
	return exists, err
}

func createOrgSecret(data *schema.ResourceData) (secret *drone.Secret) {
	return &drone.Secret{
		Name:            data.Get("name").(string),
		Data:            data.Get("value").(string),
		PullRequest:     data.Get("allow_on_pull_request").(bool),
		PullRequestPush: data.Get("allow_push_on_pull_request").(bool),
	}
}

func readOrgSecret(data *schema.ResourceData, namespace string, secret *drone.Secret) {
	data.SetId(fmt.Sprintf("%s/%s", namespace, secret.Name))
	data.Set("name", secret.Name)
	data.Set("allow_on_pull_request", secret.PullRequest)
	data.Set("allow_push_on_pull_request", secret.PullRequestPush)
}
