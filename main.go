package main

import (
	"github.com/Lucretius/terraform-provider-drone/drone"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return drone.Provider()
		},
	})
}
