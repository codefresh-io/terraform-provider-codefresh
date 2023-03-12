package codefresh

import (
	"fmt"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"os"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if url := os.Getenv(ENV_CODEFRESH_API_URL); url != "" {
						return url, nil
					}
					return DEFAULT_CODEFRESH_API_URL, nil
				},
				Description: fmt.Sprintf("The Codefresh API URL. Defaults to `%s`. Can also be set using the `%s` environment variable.", DEFAULT_CODEFRESH_API_URL, ENV_CODEFRESH_API_URL),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: fmt.Sprintf("The Codefresh API token. Can also be set using the `%s` environment variable.", ENV_CODEFRESH_API_KEY),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"codefresh_account":         dataSourceAccount(),
			"codefresh_context":         dataSourceContext(),
			"codefresh_current_account": dataSourceCurrentAccount(),
			"codefresh_idps":            dataSourceIdps(),
			"codefresh_step_types":      dataSourceStepTypes(),
			"codefresh_team":            dataSourceTeam(),
			"codefresh_user":            dataSourceUser(),
			"codefresh_users":           dataSourceUsers(),
			"codefresh_registry":        dataSourceRegistry(),
			"codefresh_pipelines":       dataSourcePipelines(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_account":               resourceAccount(),
			"codefresh_account_admins":        resourceAccountAdmins(),
			"codefresh_api_key":               resourceApiKey(),
			"codefresh_context":               resourceContext(),
			"codefresh_registry":              resourceRegistry(),
			"codefresh_idp_accounts":          resourceIDPAccounts(),
			"codefresh_permission":            resourcePermission(),
			"codefresh_pipeline":              resourcePipeline(),
			"codefresh_pipeline_cron_trigger": resourcePipelineCronTrigger(),
			"codefresh_project":               resourceProject(),
			"codefresh_step_types":            resourceStepTypes(),
			"codefresh_user":                  resourceUser(),
			"codefresh_team":                  resourceTeam(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	apiURL := d.Get("api_url").(string)
	token := d.Get("token").(string)
	if token == "" {
		token = os.Getenv(ENV_CODEFRESH_API_KEY)
	}
	return cfClient.NewClient(apiURL, token, ""), nil
}
