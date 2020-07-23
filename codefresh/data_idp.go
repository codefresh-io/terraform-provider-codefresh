package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIdps() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIdpRead,
		Schema: IdpSchema(),
	}
}

//IdpSchema - 
func IdpSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"client_type": {
			Type:     schema.TypeString,
			Optional: true,
		},			
		"client_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"access_token": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"client_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"client_secret": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"app_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cookie_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cookie_iv": {
			Type:     schema.TypeString,
			Computed: true,
		},								
		"tenant": {
			Type:     schema.TypeString,
			Computed: true,
		},										
		"scopes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"accounts": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceIdpRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	idps, err := client.GetIDPs()
	if err != nil {
		return err
	}

	_id, _idOk := d.GetOk("_id")
	displayName, displayNameOk := d.GetOk("displayName")
	clientType, clientTypeOk := d.GetOk("clientType")

	if !(_idOk || displayNameOk || clientTypeOk) {
		return fmt.Errorf("[EROOR] Idp data_source - no parameters specified")
	}
	for _, idp := range *idps {
		if _idOk && _id.(string) != idp.ID {
			continue
		}
		if displayNameOk && displayName.(string) != idp.DisplayName {
			continue
		}
		if clientTypeOk && clientType.(string) != idp.ClientType {
			continue
		}
		err = mapDataIdpToResource(idp, d)
		if err != nil {
			return err
		}
	}

	if d.Id() == "" {
		return fmt.Errorf("[EROOR] Idp wasn't found")
	}

	return nil
}

func mapDataIdpToResource(idp cfClient.IDP, d *schema.ResourceData) error {
	
	d.SetId(idp.ID)

	d.Set("access_token", idp.Access_token)//    string   `json:"access_token,omitempty"`
	
	
	d.Set("accounts", flattenStringArr(idp.Accounts)) //
	//d.Set("apiHost", idp.ApiHost) //         string   `json:"apiHost,omitempty"`
	//d.Set("apiPathPrefix", idp.ApiPathPrefix) //   string   `json:"apiPathPrefix,omitempty"`
	//d.Set("apiURL", idp.ApiURL) //          string   `json:"apiURL,omitempty"`
	//d.Set("appId", idp.AppId) //          string   `json:"appId,omitempty"`
	//d.Set("authURL", idp.AuthURL) //        string   `json:"authURL,omitempty"`
	d.Set("client_host", idp.ClientHost) //     string   `json:"clientHost,omitempty"`
	d.Set("client_id", idp.ClientId) //       string   `json:"clientId,omitempty"`
	d.Set("client_name", idp.ClientName) //     string   `json:"clientName,omitempty"`
	d.Set("client_secret",idp.ClientSecret) //    string   `json:"clientSecret,omitempty"`
	d.Set("client_type", idp.ClientType) //      string   `json:"clientType,omitempty"`
	d.Set("cookie_iv", idp.CookieIv) //        string   `json:"cookieIv,omitempty"`
	d.Set("cookie_key", idp.CookieKey) //       string   `json:"cookieKey,omitempty"`
	d.Set("display_name", idp.DisplayName) //     string   `json:"displayName,omitempty"`
	d.Set("_id", idp.ID) //              string   `json:"_id,omitempty"`
	//d.Set("IDPLoginUrl", idp.IDPLoginUrl) //     string   `json:"IDPLoginUrl,omitempty"`
	//d.Set("loginUrl", idp.LoginUrl) //        string   `json:"loginUrl,omitempty"`
	//d.Set("redirectUiUrl", idp.RedirectUiUrl) //   string   `json:"redirectUiUrl,omitempty"`
	//d.Set("redirectUrl", idp.RedirectUrl) //     string   `json:"redirectUrl,omitempty"`
	//d.Set("refreshTokenURL", idp.RefreshTokenURL) // string   `json:"refreshTokenURL,omitempty"`
	d.Set("scopes", flattenStringArr(idp.Scopes)) //          []string `json:"scopes,omitempty"`
	d.Set("tenant", idp.Tenant) //          string   `json:"tenant,omitempty"`
	//d.Set("tokenSecret", idp.TokenSecret) //     string   `json:"tokenSecret,omitempty"`
	//d.Set("tokenURL", idp.TokenURL) //        string   `json:"tokenURL,omitempty"`
	//d.Set("userProfileURL", idp.UserProfileURL) //  string   `json:"userProfileURL,omitempty"`

	return nil
}