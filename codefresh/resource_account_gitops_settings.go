package codefresh

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/gitops"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccountGitopsSettings() *schema.Resource {
	return &schema.Resource{
		Description: "Codefresh account gitops setting - such as git provider, API URL for the git provider and internal shared config repository",
		Read:        resourceAccountGitopsSettingsRead,
		Create:      resourceAccountGitopsSettingsUpdate,
		Update:      resourceAccountGitopsSettingsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// Delete not implemenented as gitops settings cannot be removed, only updated
		Delete: resourceAccountGitopsSettingsDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Description: "Account Id",
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Account name for active account",
				Computed:    true,
			},
			"git_provider": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("Git provider name - currently supported values are: %s", strings.Join(gitops.GetSupportedGitProvidersList(), " ,")),
				Required:     true,
				ValidateFunc: validation.StringInSlice(gitops.GetSupportedGitProvidersList(), false),
			},
			"git_provider_api_url": {
				Type:         schema.TypeString,
				Description:  "Git provider API url. If not provided can automatically be set for known SaaS git providers. For example - for github it will be set to https://api.github.com",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(https?:\/\/)(\S+)$`), "must be a valid url"),
				Optional:     true,
				// When an empty value for provider url is provided, check if the old one was set by getting the default and surpress diff in such case
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						defaultProviderUrl, err := gitops.GetDefaultAPIUrlForProvider(d.Get("git_provider").(string))

						if err != nil {
							return false
						}

						if *defaultProviderUrl == old {
							return true
						}
					}
					return false
				},
			},
			"shared_config_repository": {
				Type:         schema.TypeString,
				Description:  "Shared config repository url. Must be a valid git url which contains `.git`. May also includ path and branch references",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(https?:\/\/)(\S+)(.git)(\S*)$`), "must be a valid git url and must contain .git For example https://github.com/owner/repo.git or https://github.com/owner/repo.git/some/path?ref=branch-name"),
				Required:     true,
			},
		},
	}
}

func resourceAccountGitopsSettingsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	var accountGitopsInfo *cfclient.GitopsActiveAccountInfo

	accountGitopsInfo, err := client.GetActiveGitopsAccountInfo()

	if err != nil {
		return err
	}

	return mapAccountGitopsSettingsToResource(accountGitopsInfo, d)
}

func resourceAccountGitopsSettingsUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	var gitApiUrl string

	if _, ok := d.GetOk("git_provider_api_url"); !ok {
		url, err := gitops.GetDefaultAPIUrlForProvider(d.Get("git_provider").(string))

		if err != nil {
			return err
		}

		gitApiUrl = *url

	} else {
		gitApiUrl = d.Get("git_provider_api_url").(string)
	}

	err := client.UpdateActiveGitopsAccountSettings(d.Get("git_provider").(string), gitApiUrl, d.Get("shared_config_repository").(string))

	if err != nil {
		return err
	}

	return resourceAccountGitopsSettingsRead(d, meta)
}

// Settings cannot be deleted, only updated
func resourceAccountGitopsSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func mapAccountGitopsSettingsToResource(account *cfclient.GitopsActiveAccountInfo, d *schema.ResourceData) error {

	if account == nil || account.ID == "" {
		return fmt.Errorf("cannot get gitops settings as account wasn't properly retrived")
	}

	d.SetId(account.ID)
	d.Set("name", account.AccountName)
	d.Set("git_provider", account.GitProvider)
	d.Set("git_provider_api_url", account.GitApiUrl)
	d.Set("shared_config_repository", account.SharedConfigRepo)

	return nil
}
