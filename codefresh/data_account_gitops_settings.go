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
			"id": {
				Type:        schema.TypeString,
				Description: "Account Id",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account name for active account",
			},
			"git_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Git provider name",
			},
			"git_provider_api_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Git provider API url",
			},
			"shared_config_repository": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Shared config repository url",
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
		return fmt.Errorf("cannot get gitops settings as account wasn't properly retrived")
	}
	d.SetId(account.ID)

	err := d.Set("name", account.AccountName)

	if err != nil {
		return err
	}

	err = d.Set("admins", account.Admins)

	if err != nil {
		return err
	}

	err = d.Set("git_provider", account.GitProvider)

	if err != nil {
		return err
	}

	err = d.Set("git_provider_api_url", account.GitApiUrl)

	if err != nil {
		return err
	}

	err = d.Set("shared_config_repository", account.SharedConfigRepo)

	if err != nil {
		return err
	}

	return nil
}
