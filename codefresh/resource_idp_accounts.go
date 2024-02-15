package codefresh

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIDPAccounts() *schema.Resource {
	return &schema.Resource{
		Description: `
This resource adds the list of provided account IDs to the IDP.  
Because of the current Codefresh API limitation it's impossible to remove account from IDP, thus deletion is not supported.
		`,
		Create: resourceIDPAccountsCreate,
		Read:   resourceIDPAccountsRead,
		Update: resourceIDPAccountsUpdate,
		Delete: resourceIDPAccountsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"idp_id": {
				Description: "The IdP name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_ids": {
				Description: "A list of account IDs to add to the IdP.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceIDPAccountsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	accountIds := datautil.ConvertStringArr(d.Get("account_ids").(*schema.Set).List())

	idpID := d.Get("idp_id").(string)

	idp, err := client.GetIdpByID(idpID)
	if err != nil {
		return err
	}

	for _, accountID := range accountIds {
		client.AddAccountToIDP(accountID, idp.ID)
	}

	d.SetId(idp.ID)

	return nil
}

func resourceIDPAccountsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	idpID := d.Id()
	if idpID == "" {
		d.SetId("")
		return nil
	}

	idp, err := client.GetIdpByID(idpID)
	if err != nil {
		return err
	}

	err = d.Set("idp_id", idp.ID)
	if err != nil {
		return err
	}

	err = d.Set("account_ids", idp.Accounts)
	if err != nil {
		return err
	}

	return nil
}

func resourceIDPAccountsDelete(_ *schema.ResourceData, _ interface{}) error {
	// todo
	// warning message
	return nil
}

func resourceIDPAccountsUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	idpID := d.Id()

	idp, err := client.GetIdpByID(idpID)
	if err != nil {
		return err
	}

	existingAccounts := idp.Accounts

	desiredAccounts := datautil.ConvertStringArr(d.Get("account_ids").(*schema.Set).List())

	for _, account := range desiredAccounts {
		if ok := cfclient.FindInSlice(existingAccounts, account); !ok {
			client.AddAccountToIDP(account, idp.ID)
		}
	}

	return nil
}
