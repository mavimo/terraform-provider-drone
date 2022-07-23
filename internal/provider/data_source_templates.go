package drone

import (
	"context"
	"sort"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func dataSourceTemplates() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving all Drone templates in a namespace",
		ReadContext: dataSourceTemplatesRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"names": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceTemplatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(drone.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	namespace := d.Get("namespace").(string)
	templates, err := client.TemplateList(namespace)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to retrieve templates",
			Detail:   err.Error(),
		})

		return diags
	}

	id := make([]string, 0)
	names := make([]string, 0)

	for _, template := range templates {
		id = append(id, template.Name)
		names = append(names, template.Name)
	}

	sort.Strings(names)
	d.Set("names", names)

	d.SetId(utils.BuildChecksumID(id))

	return diags
}
