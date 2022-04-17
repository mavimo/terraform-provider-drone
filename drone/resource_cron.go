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

func resourceCron() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a repository cron job.",
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
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicate if the cron should be disabled.",
			},
			"event": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The event for this cron job. Only allowed value is `push`.",
				ValidateFunc: validation.StringInSlice([]string{
					drone.EventPush,
				}, false),
			},
			"branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The branch to run this cron job on.",
				Default:     "master",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Cron job name",
				Required:    true,
				ForceNew:    true,
			},
			"expr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cron interval. Allowed values are `@daily`, `@weekly`, `@monthly`, and `@yearly`.",
				Default:     "@monthly",
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"@hourly",
					"@daily",
					"@weekly",
					"@monthly",
					"@yearly",
				}, false),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCronCreate,
		ReadContext:   resourceCronRead,
		UpdateContext: resourceCronUpdate,
		DeleteContext: resourceCronDelete,
		Exists:        resourceCronExists,
	}
}

func resourceCronCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	cron, err := client.CronCreate(owner, repo, createCron(data))

	if err != nil {
		return diag.FromErr(err)
	}

	return readCron(data, cron, owner, repo, err)
}

func resourceCronRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	cronName := data.Get("name").(string)
	cron, err := client.Cron(owner, repo, cronName)
	if err != nil {
		return diag.Errorf("failed to read Drone Cron: %s/%s/%s", owner, repo, cronName)
	}

	return readCron(data, cron, owner, repo, err)
}

func resourceCronUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	cronName := data.Get("name").(string)
	cron, err := client.CronUpdate(owner, repo, cronName, updateCron(data))

	return readCron(data, cron, owner, repo, err)
}

func resourceCronDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	cronName := data.Get("name").(string)
	err = client.CronDelete(owner, repo, cronName)

	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceCronExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, err := utils.ParseRepo(data.Get("repository").(string))

	if err != nil {
		return false, err
	}
	cronName := data.Get("name").(string)
	cron, err := client.Cron(owner, repo, cronName)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Cron: %s/%s/%s", owner, repo, cronName)
	}

	exists := cron.Name == cronName
	return exists, err
}

func createCron(data *schema.ResourceData) (repository *drone.Cron) {
	branch := data.Get("branch").(string)
	disabled := data.Get("disabled").(bool)
	expr := data.Get("expr").(string)
	name := data.Get("name").(string)
	event := data.Get("event").(string)

	cron := &drone.Cron{
		Disabled: disabled,
		Branch:   branch,
		Expr:     expr,
		Event:    event,
		Name:     name,
	}

	return cron
}

func updateCron(data *schema.ResourceData) (repository *drone.CronPatch) {
	branch := data.Get("branch").(string)
	disabled := data.Get("disabled").(bool)
	event := data.Get("event").(string)

	cron := &drone.CronPatch{
		Disabled: utils.Bool(disabled),
		Branch:   &branch,
		Event:    &event,
	}
	return cron
}

func readCron(data *schema.ResourceData, cron *drone.Cron, namespace string, repo string, err error) diag.Diagnostics {
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%d", cron.ID))

	data.Set("repository", fmt.Sprintf("%s/%s", namespace, repo))
	data.Set("branch", cron.Branch)
	data.Set("disabled", cron.Disabled)
	data.Set("expr", cron.Expr)
	data.Set("name", cron.Name)

	return diag.Diagnostics{}
}
