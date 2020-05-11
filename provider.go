package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_project":      resourceProject(),
			"codefresh_pipeline":     resourcePipeline(),
			"codefresh_cron_event":   resourceCronEvent(),
			"codefresh_cron_trigger": resourceCronTrigger(),
			"codefresh_user":         resourceUser(),
			"codefresh_environment":  resourceEnvironment(),
		},
	}
}
