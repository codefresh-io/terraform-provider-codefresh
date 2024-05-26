package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountGitopsSettings() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves gitops settings for the active account",
		Read:        dataSourceAccountGitopsSettingsRead,
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_provider_api_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_config_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAccountGitopsSettingsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	var accountGitopsInfo *cfclient.GitopsActiveAccountInfo

	accountGitopsInfo, err := client.GetActiveGitopsAccountInfo()

	if err != nil {
		return err
	}

	return mapDataAccountGitopsSettingsToResource(accountGitopsInfo, d)
}

func mapDataAccountGitopsSettingsToResource(account *cfclient.GitopsActiveAccountInfo, d *schema.ResourceData) error {

	if account == nil || account.ID == "" {
		return fmt.Errorf("data.codefresh_account - failed to mapDataAccountToResource")
	}
	d.SetId(account.ID)
	d.Set("_id", account.ID)
	d.Set("name", account.AccountName)
	d.Set("admins", account.Admins)
	d.Set("git_provider", account.GitProvider)
	d.Set("git_provider_api_url", account.GitApiUrl)
	d.Set("shared_config_repo", account.SharedConfigRepo)

	return nil
}
