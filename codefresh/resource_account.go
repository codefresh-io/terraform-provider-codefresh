package codefresh

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Description: `
		By creating different accounts for different teams within the same company a customer can achieve complete segregation of assets between the teams. Requires Codefresh admin token and hence is relevant only for on premise installations of Codefresh.
		`,
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The display name for the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"features": {
				Description: `
Features toggles for this account. Default:

* OfflineLogging: true
* ssoManagement: true
* teamsManagement: true
* abac: true
* customKubernetesCluster: true
`,
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				Default: map[string]bool{
					"OfflineLogging":          true,
					"ssoManagement":           true,
					"teamsManagement":         true,
					"abac":                    true,
					"customKubernetesCluster": true,
				},
			},
			"limits": {
				Description: "Limits for this account.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collaborators": {
							Description: "The number of collaborators allowed for this account.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"data_retention_weeks": {
							Description: "Specifies the number of weeks for which to store the builds (default: `5`).",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5,
						},
					},
				},
			},
			"build": {
				Description: "Build limits for this account.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parallel": {
							Description: "The number of parallel builds allowed for this account.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"nodes": {
							Description: "The number of nodes allowed for this account (default: `1`).",
							Type:        schema.TypeInt,
							Default:     1,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	account := *mapResourceToAccount(d)

	resp, err := client.CreateAccount(&account)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

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

	client := meta.(*cfclient.Client)

	account := *mapResourceToAccount(d)

	_, err := client.UpdateAccount(&account)
	if err != nil {
		return err
	}

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	err := client.DeleteAccount(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapAccountToResource(account *cfclient.Account, d *schema.ResourceData) error {
	err := d.Set("name", account.Name)
	if err != nil {
		return err
	}

	err = d.Set("features", account.Features)
	if err != nil {
		return err
	}

	if account.Limits != nil { // On-prem API does not return 'limits' field
		err = d.Set("limits", []map[string]interface{}{flattenLimits(*account.Limits)})
		if err != nil {
			return err
		}
	}

	err = d.Set("build", []map[string]interface{}{flattenBuild(*account.Build)})
	if err != nil {
		return err
	}

	return nil
}

func flattenLimits(limits cfclient.Limits) map[string]interface{} {
	res := make(map[string]interface{})
	res["collaborators"] = limits.Collaborators.Limit
	res["data_retention_weeks"] = limits.DataRetention.Weeks
	return res
}

func flattenBuild(build cfclient.Build) map[string]interface{} {
	res := make(map[string]interface{})
	res["parallel"] = build.Parallel
	res["nodes"] = build.Nodes
	return res
}
func mapResourceToAccount(d *schema.ResourceData) *cfclient.Account {
	account := &cfclient.Account{
		ID:   d.Id(),
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("features"); ok {
		account.SetFeatures(d.Get("features").(map[string]interface{}))
	}
	if _, ok := d.GetOk("limits"); ok {
		account.Limits = &cfclient.Limits{
			Collaborators: cfclient.Collaborators{
				Limit: d.Get("limits.0.collaborators").(int),
			},
			DataRetention: cfclient.DataRetention{
				Weeks: d.Get("limits.0.data_retention_weeks").(int),
			},
		}
	}

	if _, ok := d.GetOk("build"); ok {
		account.Build = &cfclient.Build{
			Parallel: d.Get("build.0.parallel").(int),
			Nodes:    d.Get("build.0.nodes").(int),
		}
	}
	return account
}
