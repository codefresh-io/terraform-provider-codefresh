package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admins": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"limits": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collaborators": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"data_retention_weeks": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"build": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parallel": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nodes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	account := *mapResourceToAccount(d)

	resp, err := client.CreateAccount(&account)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	accountID := d.Id()
	if accountID == "" {
		d.SetId("")
		return nil
	}

	team, err := client.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	err = mapAccountToResource(team, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	account := *mapResourceToAccount(d)

	_, err := client.UpdateAccount(&account)
	if err != nil {
		return err
	}

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	err := client.DeleteAccount(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapAccountToResource(account *cfClient.Account, d *schema.ResourceData) error {

	err := d.Set("name", account.Name)
	if err != nil {
		return err
	}

	err = d.Set("admins", account.Admins)
	if err != nil {
		return err
	}

	err = d.Set("limits", []map[string]interface{}{flattenLimits(*account.Limits)})
	if err != nil {
		return err
	}

	err = d.Set("build", []map[string]interface{}{flattenBuild(*account.Build)})
	if err != nil {
		return err
	}

	return nil
}

func flattenLimits(limits cfClient.Limits) map[string]interface{} {
	res := make(map[string]interface{})
	res["collaborators"] = limits.Collaborators.Limit
	res["data_retention_weeks"] = limits.DataRetention.Weeks
	return res
}

func flattenBuild(build cfClient.Build) map[string]interface{} {
	res := make(map[string]interface{})
	res["parallel"] = build.Parallel
	res["nodes"] = build.Nodes
	return res
}
func mapResourceToAccount(d *schema.ResourceData) *cfClient.Account {
	admins := d.Get("admins").(*schema.Set).List()

	account := &cfClient.Account{
		ID:     d.Id(),
		Name:   d.Get("name").(string),
		Admins: convertStringArr(admins),
	}

	if _, ok := d.GetOk("limits"); ok {
		account.Limits = &cfClient.Limits{
			Collaborators: cfClient.Collaborators{
				Limit: d.Get("limits.0.collaborators").(int),
			},
			DataRetention: cfClient.DataRetention{
				Weeks: d.Get("limits.0.data_retention_weeks").(int),
			},
		}
	}

	if _, ok := d.GetOk("build"); ok {
		account.Build = &cfClient.Build{
			Parallel: d.Get("build.0.parallel").(int),
			Nodes:    d.Get("build.0.nodes").(int),
		}
	}
	return account
}
