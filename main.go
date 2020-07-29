package main

import (
	"context"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
	"os"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func main() {
	debugMode := (os.Getenv("CODEFRESH_PLUGIN_DEBUG") != "")
	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/-/codefresh",
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
