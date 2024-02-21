package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountIdp() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves an account level identity provider",
		Read:        dataSourceAccountIdpRead,
		Schema:      AccountIdpSchema(),
	}
}

// IdpSchema -
func AccountIdpSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"_id", "client_name"},
		},
		"client_name": {
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"_id", "client_name"},
		},
		"display_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"client_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"redirect_url": {
			Description: "API Callback url for the identity provider",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"redirect_ui_url": {
			Description: "UI Callback url for the identity provider",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"login_url": {
			Description: "Login url using the IDP to Codefresh",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func dataSourceAccountIdpRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	idps, err := client.GetAccountIDPs()
	if err != nil {
		return err
	}

	_id, _idOk := d.GetOk("_id")
	clientName, clientNameOk := d.GetOk("client_name")

	for _, idp := range *idps {
		if clientNameOk && clientName.(string) != idp.ClientName {
			continue
		}
		if _idOk && _id.(string) != idp.ID {
			continue
		}

		err = mapDataAccountIdpToResource(idp, d)
		if err != nil {
			return err
		}
	}

	if d.Id() == "" {
		return fmt.Errorf("[EROOR] Idp wasn't found")
	}

	return nil
}

func mapDataAccountIdpToResource(cfClientIDP cfclient.IDP, d *schema.ResourceData) error {

	d.SetId(cfClientIDP.ID)
	d.Set("client_name", cfClientIDP.ClientName)
	d.Set("client_type", cfClientIDP.ClientType)
	d.Set("display_name", cfClientIDP.DisplayName)
	d.Set("redirect_url", cfClientIDP.RedirectUrl)
	d.Set("redirect_ui_url", cfClientIDP.RedirectUiUrl)
	d.Set("login_url", cfClientIDP.LoginUrl)

	return nil
}
