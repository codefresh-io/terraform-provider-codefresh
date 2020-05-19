package codefresh

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type: schema.TypeString,
				Optional: true,
				Default:  "https://g.codefresh.io/api",
			},
			"token": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("CODEFRESH_API_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_project":  resourceProject(),
			"codefresh_pipeline": resourcePipeline(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIServer: d.Get("api_url").(string),
		Token:     d.Get("token").(string),
	}

	return &config, nil
}
