package main

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return codefresh.Provider()
		},
	})
}
