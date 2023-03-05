package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIDPAccounts() *schema.Resource {
	return &schema.Resource{
		Description: `
This resource adds the list of provided account IDs to the IDP.  
Because of the current Codefresh API limitation it's impossible to remove account from IDP, thus deletion is not supported.
		`,
		Create: resourceAccountIDPCreate,
		Read:   resourceAccountIDPRead,
		Update: resourceAccountIDPUpdate,
		Delete: resourceAccountIDPDelete,
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

func resourceAccountIDPCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	accountIds := convertStringArr(d.Get("account_ids").(*schema.Set).List())

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

func resourceAccountIDPRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

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

func resourceAccountIDPDelete(_ *schema.ResourceData, _ interface{}) error {
	// todo
	// warning message
	return nil
}

func resourceAccountIDPUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	idpID := d.Id()

	idp, err := client.GetIdpByID(idpID)
	if err != nil {
		return err
	}

	existingAccounts := idp.Accounts

	desiredAccounts := convertStringArr(d.Get("account_ids").(*schema.Set).List())

	for _, account := range desiredAccounts {
		if ok := cfClient.FindInSlice(existingAccounts, account); !ok {
			client.AddAccountToIDP(account, idp.ID)
		}
	}

	return nil
}
