package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
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
			"api_url_v2": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if url := os.Getenv(ENV_CODEFRESH_API2_URL); url != "" {
						return url, nil
					}
					return DEFAULT_CODEFRESH_API2_URL, nil
				},
				Description: fmt.Sprintf("The Codefresh gitops API URL. Defaults to `%s`. Can also be set using the `%s` environment variable.", DEFAULT_CODEFRESH_API2_URL, ENV_CODEFRESH_API2_URL),
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
			"codefresh_project":         dataSourceProject(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_account":                  resourceAccount(),
			"codefresh_account_user_association": resourceAccountUserAssociation(),
			"codefresh_account_admins":           resourceAccountAdmins(),
			"codefresh_api_key":                  resourceApiKey(),
			"codefresh_context":                  resourceContext(),
			"codefresh_registry":                 resourceRegistry(),
			"codefresh_idp_accounts":             resourceIDPAccounts(),
			"codefresh_permission":               resourcePermission(),
			"codefresh_pipeline":                 resourcePipeline(),
			"codefresh_pipeline_cron_trigger":    resourcePipelineCronTrigger(),
			"codefresh_project":                  resourceProject(),
			"codefresh_step_types":               resourceStepTypes(),
			"codefresh_user":                     resourceUser(),
			"codefresh_team":                     resourceTeam(),
			"codefresh_abac_rules":               resourceGitopsAbacRule(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	apiURL := d.Get("api_url").(string)
	apiURLV2 := d.Get("api_url_v2").(string)
	token := d.Get("token").(string)
	if token == "" {
		token = os.Getenv(ENV_CODEFRESH_API_KEY)
	}
	return cfclient.NewClient(apiURL, apiURLV2, token, ""), nil
}
