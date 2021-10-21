package main

import (
	"context"
	"log"
	"os"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func main() {
	debugMode := (os.Getenv("CODEFRESH_PLUGIN_DEBUG") != "")
	// for terraform 0.13: export CODEFRESH_PLUGIN_ADDR="codefresh.io/app/codefresh"
	providerAddr := os.Getenv("CODEFRESH_PLUGIN_ADDR")
	if providerAddr == "" {
		providerAddr = "registry.terraform.io/-/codefresh"
	}
	if debugMode {
		err := plugin.Debug(context.Background(), providerAddr,
			&plugin.ServeOpts{
				ProviderFunc: codefresh.Provider,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: codefresh.Provider,
		},
		)
	}
}
