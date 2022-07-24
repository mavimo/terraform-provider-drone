package drone

import (
	"context"
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving a Drone template",
		ReadContext: dataSourceTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template name",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization name",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template content",
			},
		},
	}
}

func dataSourceTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	template, err := client.Template(namespace, name)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Template %s/%s not found", namespace, name),
			Detail:   err.Error(),
		})

		return diags
	}

	d.Set("data", template.Data)
	d.SetId(fmt.Sprintf("%s/%s", namespace, name))

	return diags
}
