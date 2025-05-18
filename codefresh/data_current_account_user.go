package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCurrentAccountUser() *schema.Resource {
	return &schema.Resource{
		Description: "Returns a user the current Codefresh account by name or email.",
		Read:        dataSourceCurrentAccountUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"name", "email"},
				Optional:     true,
			},
			"email": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"name", "email"},
				Optional:     true,
			},
		},
	}
}

func dataSourceCurrentAccountUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)
	var currentAccount *cfclient.CurrentAccount
	var err error

	currentAccount, err = client.GetCurrentAccount()

	if err != nil {
		return err
	}

	if currentAccount == nil {
		return fmt.Errorf("data.codefresh_current_account - failed to get current_account")
	}

	var (
		userAttributeName  string
		userAttributeValue string
	)

	if _email, _emailOk := d.GetOk("email"); _emailOk {
		userAttributeName = "email"
		userAttributeValue = _email.(string)
	} else if _name, _nameOk := d.GetOk("name"); _nameOk {
		userAttributeName = "name"
		userAttributeValue = _name.(string)
	} else {
		return fmt.Errorf("data.codefresh_current_account_user - must specify name or email")
	}

	return mapDataCurrentAccountUserToResource(currentAccount, d, userAttributeName, userAttributeValue)

}

func mapDataCurrentAccountUserToResource(currentAccount *cfclient.CurrentAccount, d *schema.ResourceData, userAttributeName string, userAttributeValue string) error {

	if currentAccount == nil || currentAccount.ID == "" {
		return fmt.Errorf("data.codefresh_current_account - failed to mapDataCurrentAccountUserToResource no id for current account set")
	}

	isFound := false

	for _, user := range currentAccount.Users {
		if (userAttributeName == "name" && user.UserName == userAttributeValue) || (userAttributeName == "email" && user.Email == userAttributeValue) {
			isFound = true
			d.SetId(user.ID)
			err := d.Set("name", user.UserName)

			if err != nil {
				return err
			}


			err = d.Set("email", user.Email)

			if err != nil {
				return err
			}

			break
		}
	}

	if !isFound {
		return fmt.Errorf("data.codefresh_current_account_user - cannot find user with %s %s", userAttributeName, userAttributeValue)
	}

	return nil
}
