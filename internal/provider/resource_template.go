package drone

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a template.",
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Description: "Organization name",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:             schema.TypeString,
				Description:      "Template name",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: resourceTemplateNameValidation,
			},
			"data": {
				Type:        schema.TypeString,
				Description: "Template content",
				Required:    true,
				ForceNew:    false,
			},
			"id": {
				Description: "The string representation of the organization secret.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceTemplateCreate,
		ReadContext:   resourceTemplateRead,
		UpdateContext: resourceTemplateUpdate,
		DeleteContext: resourceTemplateDelete,
		Exists:        resourceTemplateExists,
	}
}

func resourceTemplateCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace := data.Get("namespace").(string)
	template, err := client.TemplateCreate(namespace, createTemplate(data))
	if err != nil {
		return diag.FromErr(err)
	}

	readTemplate(data, namespace, template)

	return diag.Diagnostics{}
}

func resourceTemplateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "template")
	if err != nil {
		return diag.FromErr(err)
	}

	template, err := client.Template(namespace, name)
	if err != nil {
		return diag.Errorf("failed to read Template: %s/%s", namespace, name)
	}

	readTemplate(data, namespace, template)
	return nil
}

func resourceTemplateUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "template")
	if err != nil {
		return diag.FromErr(err)
	}

	template, err := client.TemplateUpdate(namespace, name, createTemplate(data))
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("data", data.Get("data").(string))
	readTemplate(data, namespace, template)

	return diag.Diagnostics{}
}

func resourceTemplateDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "template")
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.TemplateDelete(namespace, name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceTemplateExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	namespace, name, err := utils.ParseOrgId(data.Id(), "template")
	if err != nil {
		return false, err
	}

	template, err := client.Template(namespace, name)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone Template: %s/%s", namespace, name)
	}

	exists := (template.Name == name) && (err == nil)
	return exists, err
}

func createTemplate(data *schema.ResourceData) (secret *drone.Template) {
	return &drone.Template{
		Name: data.Get("name").(string),
		Data: data.Get("data").(string),
	}
}

func readTemplate(data *schema.ResourceData, namespace string, template *drone.Template) {
	data.SetId(fmt.Sprintf("%s/%s", namespace, template.Name))
	data.Set("name", template.Name)
	data.Set("data", template.Data)
}

func resourceTemplateNameValidation(name interface{}, path cty.Path) diag.Diagnostics {
	switch filepath.Ext(name.(string)) {
	case
		".yaml",
		".jsonnet",
		".json",
		".starlark":
		return diag.Diagnostics{}
	}
	return diag.Errorf("Template name %s do not have a valid extension", name)
}
