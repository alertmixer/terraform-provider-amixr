package main

import (
	"github.com/alertmixer/terraform-provider-amixr/amixr"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return amixr.Provider()
		},
	})
}
