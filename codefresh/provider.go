package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if url := os.Getenv("CODEFRESH_API_URL"); url != "" {
						return url, nil
					}
					return "https://g.codefresh.io/api", nil
				},
			},
			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"codefresh_users":           dataSourceUsers(),
			"codefresh_user":            dataSourceUser(),
			"codefresh_idps":            dataSourceIdps(),
			"codefresh_account":         dataSourceAccount(),
			"codefresh_team":            dataSourceTeam(),
			"codefresh_current_account": dataSourceCurrentAccount(),
			"codefresh_context":         dataSourceContext(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_project":        resourceProject(),
			"codefresh_pipeline":       resourcePipeline(),
			"codefresh_context":        resourceContext(),
			"codefresh_team":           resourceTeam(),
			"codefresh_account":        resourceAccount(),
			"codefresh_api_key":        resourceApiKey(),
			"codefresh_idp_accounts":   resourceIDPAccounts(),
			"codefresh_account_admins": resourceAccountAdmins(),
			"codefresh_user":           resourceUser(),
			"codefresh_permission":     resourcePermission(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	apiURL := d.Get("api_url").(string)
	token := d.Get("token").(string)
	if token == "" {
		token = os.Getenv("CODEFRESH_API_KEY")
	}
	return cfClient.NewClient(apiURL, token, ""), nil
}
