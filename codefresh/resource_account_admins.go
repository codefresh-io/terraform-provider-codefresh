package codefresh

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccountAdmins() *schema.Resource {
	return &schema.Resource{
		Description: `
		Use this resource to set a list of admins for any account. Requires Codefresh admin token and hence is relevant only for on premise installations of Codefresh.
		`,
		Create: resourceAccountAdminsCreate,
		Read:   resourceAccountAdminsRead,
		Update: resourceAccountAdminsUpdate,
		Delete: resourceAccountAdminsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account ID for which to set up the list of admins.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"users": {
				Description: "A list of users to set up as account admins.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAccountAdminsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	admins := d.Get("users").(*schema.Set).List()

	accountId := d.Get("account_id").(string)

	for _, admin := range datautil.ConvertStringArr(admins) {
		err := client.SetUserAsAccountAdmin(accountId, admin)
		if err != nil {
			return err
		}
	}

	// d.SetId(time.Now().UTC().String())
	d.SetId(accountId)

	return nil
}

func resourceAccountAdminsDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	admins := d.Get("users").(*schema.Set).List()

	accountId := d.Get("account_id").(string)

	for _, admin := range datautil.ConvertStringArr(admins) {
		err := client.DeleteUserAsAccountAdmin(accountId, admin)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceAccountAdminsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	accountId := d.Id()

	d.Set("account_id", accountId)

	account, err := client.GetAccountByID(accountId)
	if err != nil {
		return nil
	}
	err = d.Set("users", account.Admins)
	if err != nil {
		return err
	}

	return nil
}

func resourceAccountAdminsUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	accountId := d.Get("account_id").(string)
	desiredAdmins := d.Get("users").(*schema.Set).List()

	account, err := client.GetAccountByID(accountId)
	if err != nil {
		return err
	}

	adminsToAdd, AdminsToDelete := cfclient.GetAccountAdminsDiff(datautil.ConvertStringArr(desiredAdmins), account.Admins)

	for _, userId := range AdminsToDelete {
		err := client.DeleteUserAsAccountAdmin(accountId, userId)
		if err != nil {
			return err
		}
	}

	for _, userId := range adminsToAdd {
		err := client.SetUserAsAccountAdmin(accountId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}
