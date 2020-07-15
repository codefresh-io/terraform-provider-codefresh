package codefresh

import (
	"errors"
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiKeyCreate,
		Read:   resourceApiKeyRead,
		Update: resourceApiKeyUpdate,
		Delete: resourceApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": {
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
	client := meta.(*cfClient.Client)

	apiKey := *mapResourceToApiKey(d)
	accountID := d.Get("account_id").(string)

	resp, err := client.CreateApiKey(accountID, &apiKey)
	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	err = d.Set("token", resp)
	if err != nil {
		return err
	}

	client.Token = resp

	apiKeys, err := client.GetApiKeysList()
	if err != nil {
		return nil
	}

	var keyID string
	for _, key := range apiKeys {
		if key.Name == apiKey.Name {
			keyID = key.ID
		}
	}

	if keyID == "" {
		return errors.New("[ERROR] Key ID is not found.")
	}

	d.SetId(keyID)

	return nil
}

func resourceApiKeyRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	keyID := d.Id()
	if keyID == "" {
		d.SetId("")
		return nil
	}

	token := d.Get("token").(string)

	if token == "" {
		return errors.New("[ERROR] Can't read API Key. Token is empty.")
	}

	client.Token = token

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
	client := meta.(*cfClient.Client)

	apiKey := *mapResourceToApiKey(d)

	token := d.Get("token").(string)
	if token == "" {
		return errors.New("[ERROR] Can't read API Key. Token is empty.")
	}

	client.Token = token

	err := client.UpdateAPIKey(&apiKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

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

func mapApiKeyToResource(apiKey *cfClient.ApiKey, d *schema.ResourceData) error {

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

func mapResourceToApiKey(d *schema.ResourceData) *cfClient.ApiKey {
	scopes := d.Get("scopes").(*schema.Set).List()
	apiKey := &cfClient.ApiKey{
		ID:     d.Id(),
		Name:   d.Get("name").(string),
		Scopes: convertStringArr(scopes),
	}
	return apiKey
}
