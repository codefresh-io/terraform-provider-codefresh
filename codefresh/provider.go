package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type: schema.TypeString,
				// Required: true,
				Optional: true,
				Default:  "https://g.codefresh.io/api",
			},
			"token": {
				Type:     schema.TypeString,
				Required: true,
				// Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("CODEFRESH_API_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_project":  resourceProject(),
			"codefresh_pipeline": resourcePipeline(),
			// "codefresh_cron_event":   resourceCronEvent(),
			// "codefresh_cron_trigger": resourceCronTrigger(),
			// "codefresh_user":         resourceUser(),
			// "codefresh_environment":  resourceEnvironment(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIServer: d.Get("api_url").(string),
		Token:     d.Get("token").(string),
	}

	// if err := config.LoadAndValidate(); err != nil {
	// 	return nil, err
	// }

	return &config, nil
}
