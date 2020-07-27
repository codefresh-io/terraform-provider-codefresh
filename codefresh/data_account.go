package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admins": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},	
		},
	}
}


func dataSourceAccountRead(d *schema.ResourceData, meta interface{}) error {

	
	client := meta.(*cfClient.Client)
	var account *cfClient.Account
	var err error

	if _id, _idOk := d.GetOk("_id"); _idOk {
		account, err = client.GetAccountByID(_id.(string))
	} else if name, nameOk := d.GetOk("name"); nameOk {
		account, err = client.GetAccountByName(name.(string))
	} else {
		return fmt.Errorf("data.codefresh_account - must specify _id or name")
	}
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("data.codefresh_account - cannot find account")
	}

    return mapDataAccountToResource(account, d)

}

func mapDataAccountToResource(account *cfClient.Account, d *schema.ResourceData) error {
	
	if account == nil || account.ID == "" {
		return fmt.Errorf("data.codefresh_account - failed to mapDataAccountToResource")
	}
	d.SetId(account.ID)

	d.Set("_id", account.ID)
	d.Set("name", account.Name)
	d.Set("admins", account.Admins)

	return nil
}

