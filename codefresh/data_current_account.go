package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCurrentAccount() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceCurrentAccountRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},						
					},
				},
			},	
		},
	}
}


func dataSourceCurrentAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)
	var currentAccount *cfClient.CurrentAccount
	var err error

	currentAccount, err = client.GetCurrentAccount()
	if err != nil {
		return err
	}

	if currentAccount == nil {
		return fmt.Errorf("data.codefresh_current_account - failed to get current_account")
	}

    return mapDataCurrentAccountToResource(currentAccount, d)

}

func mapDataCurrentAccountToResource(currentAccount *cfClient.CurrentAccount, d *schema.ResourceData) error {
	
	if currentAccount == nil || currentAccount.ID == "" {
		return fmt.Errorf("data.codefresh_current_account - failed to mapDataCurrentAccountToResource")
	}
	d.SetId(currentAccount.ID)

	d.Set("_id", currentAccount.ID)
	d.Set("name", currentAccount.Name)
	
	// users := make(map[string](map[string]interface{}))
    // for n, user := range currentAccount.Users {
	// 	users[n] = make(map[string]interface{})
	// 	users[n]["name"] = user.UserName
	// 	users[n]["email"] = user.Email
	// 	users[n]["id"] = user.ID
	// }

	// d.Set("users", []map[string](map[string]interface{}){users})
	users := make([](map[string]interface{}), len(currentAccount.Users))
	for n, user := range currentAccount.Users {
		users[n] = make(map[string]interface{})
		users[n]["name"] = user.UserName
		users[n]["email"] = user.Email
		users[n]["id"] = user.ID
	}
	
    d.Set("users", users)
	return nil
}
