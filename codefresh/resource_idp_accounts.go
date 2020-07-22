package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIDPAccounts() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountIDPCreate,
		Read:   resourceAccountIDPRead,
		Update: resourceAccountIDPUpdate,
		Delete: resourceAccountIDPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"idp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accounts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAccountIDPCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	accounts := convertStringArr(d.Get("accounts").(*schema.Set).List())

	idpName := d.Get("idp").(string)

	idp, err := client.GetIdpByName(idpName)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		client.AddAccountToIDP(account, idp.ID)
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

	err = d.Set("idp", idp.ClientName)
	if err != nil {
		return err
	}

	err = d.Set("accounts", idp.Accounts)
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

	desiredAccounts := convertStringArr(d.Get("accounts").(*schema.Set).List())

	for _, account := range desiredAccounts {
		if ok := cfClient.FindInSlice(existingAccounts, account); !ok {
			client.AddAccountToIDP(account, idp.ID)
		}
	}

	return nil
}
