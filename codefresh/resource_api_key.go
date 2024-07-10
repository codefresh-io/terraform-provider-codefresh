package codefresh

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiKey() *schema.Resource {
	return &schema.Resource{
		Description: `
		Manages an API Key tied to an Account and a User.
		Requires Codefresh admin token and hence is relevant only for on premise installations of Codefresh.
		`,
		Create: resourceApiKeyCreate,
		Read:   resourceApiKeyRead,
		Update: resourceApiKeyUpdate,
		Delete: resourceApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The display name for the API key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_id": {
				Description: "The ID of account in which the API key will be created.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_id": {
				Description: "The ID of a user within the referenced `account_id` that will own the API key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token": {
				Description: "The resulting API key.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"scopes": {
				Description: `
A list of access scopes for the API key. The possible values:
	* agent
	* agents
	* audit
	* build
	* cluster
	* clusters
	* environments-v2
	* github-action
	* helm
	* kubernetes
	* pipeline
	* project
	* repos
	* runner-installation
	* step-type
	* step-types
	* view
	* workflow
				`,
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	apiKey := *mapResourceToApiKey(d)
	accountID := d.Get("account_id").(string)
	userID := d.Get("user_id").(string)

	resp, err := client.CreateApiKey(userID, accountID, &apiKey)
	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	err = d.Set("token", resp)
	if err != nil {
		return err
	}

	// Codefresh tokens are in the form xxxxxxxxxxxx.xxxxxxxxx the first half serves as the id
	d.SetId(strings.Split(resp, ".")[0])

	return nil
}

func resourceApiKeyRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	keyID := d.Id()
	if keyID == "" {
		d.SetId("")
		return nil
	}

	token := d.Get("token").(string)

	if token == "" {
		return errors.New("[ERROR] Can't read API Key. Token is empty.")
	}

	apiKey, err := client.GetAPIKey(keyID)
	if err != nil {
		return err
	}

	err = mapApiKeyToResource(apiKey, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceApiKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	apiKey := *mapResourceToApiKey(d)

	token := d.Get("token").(string)
	if token == "" {
		return errors.New("[ERROR] Can't read API Key. Token is empty.")
	}

	err := client.UpdateAPIKey(&apiKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	token := d.Get("token").(string)
	if token == "" {
		return errors.New("[ERROR] Can't read API Key. Token is empty.")
	}

	err := client.DeleteAPIKey(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapApiKeyToResource(apiKey *cfclient.ApiKey, d *schema.ResourceData) error {

	err := d.Set("name", apiKey.Name)
	if err != nil {
		return err
	}

	err = d.Set("scopes", apiKey.Scopes)
	if err != nil {
		return err
	}
	return nil
}

func mapResourceToApiKey(d *schema.ResourceData) *cfclient.ApiKey {
	scopes := d.Get("scopes").(*schema.Set).List()
	apiKey := &cfclient.ApiKey{
		ID:     d.Id(),
		Name:   d.Get("name").(string),
		Scopes: datautil.ConvertStringArr(scopes),
	}
	return apiKey
}
