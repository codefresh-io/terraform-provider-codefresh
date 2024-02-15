package main

import (
	"os"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	debugMode := (os.Getenv(codefresh.ENV_CODEFRESH_PLUGIN_DEBUG) != "")
	providerAddr := os.Getenv(codefresh.ENV_CODEFRESH_PLUGIN_ADDR)
	if providerAddr == "" {
		providerAddr = codefresh.DEFAULT_CODEFRESH_PLUGIN_ADDR
	}
	plugin.Serve(&plugin.ServeOpts{
		ProviderAddr: providerAddr, // Required for debug attaching
		ProviderFunc: codefresh.Provider,
		Debug:        debugMode,
	})
}
