package main

import (
	"github.com/alertmixer/terraform-provider-amixr/amixr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: amixr.Provider}

	plugin.Serve(opts)
}
