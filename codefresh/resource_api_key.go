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
		Manages an API Key tied to a user within an account or a service account within the current account.
		On the Codefresh SaaS platfrom this resource is only usable for service accounts.
		Management of API keys for users in other accounts requires admin priveleges and hence can only be done on Codefresh on-premises installations.
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
				Description:  "The ID of account in which the API key will be created. Required if user_id is set.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"user_id", "account_id"},
				ForceNew:     true,
			},
			"user_id": {
				Description:  "The ID of a user within the referenced `account_id` that will own the API key. Requires a Codefresh admin token and can be used only in Codefresh on-premises installations.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"user_id", "service_account_id"},
				RequiredWith: []string{"user_id", "account_id"},
				ForceNew:     true,
			},
			"service_account_id": {
				Description:  "The ID of the service account to create the API key for.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"user_id", "service_account_id"},
				ForceNew:     true,
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

	var (
		resp string
		err  error
	)

	if serviceAccountId := d.Get("service_account_id").(string); serviceAccountId != "" {
		resp, err = client.CreateApiKeyServiceUser(serviceAccountId, &apiKey)
	} else {
		accountID := d.Get("account_id").(string)
		userID := d.Get("user_id").(string)

		resp, err = client.CreateApiKey(userID, accountID, &apiKey)
	}

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

	var (
		apiKey *cfclient.ApiKey
		err    error
	)

	if serviceAccountId := d.Get("service_account_id").(string); serviceAccountId != "" {
		apiKey, err = client.GetAPIKeyServiceUser(keyID, serviceAccountId)
	} else {
		apiKey, err = client.GetAPIKey(keyID)
	}

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

	var err error

	if serviceAccountId := d.Get("service_account_id").(string); serviceAccountId != "" {
		err = client.UpdateAPIKeyServiceUser(&apiKey, serviceAccountId)
	} else {
		err = client.UpdateAPIKey(&apiKey)

	}

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

	var err error
	if serviceAccountId := d.Get("service_account_id").(string); serviceAccountId != "" {
		err = client.DeleteAPIKeyServiceUser(d.Id(), serviceAccountId)
	} else {
		err = client.DeleteAPIKey(d.Id())
	}

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
