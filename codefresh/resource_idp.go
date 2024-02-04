package codefresh

import (
	"log"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


func resourceIdp() *schema.Resource {
	return &schema.Resource{
		Description: "Identity providers used in Codefresh for user authentication.",
		Create:      resourceIDPCreate,
		Read:        resourceIDPRead,
		Update:      resourceIDPUpdate,
		Delete:      resourceIDPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Description: "The display name for the IDP.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the IDP, will be generated if not set",
				Type:        schema.TypeString,
				Optional: 	 true,
			},
			"client_type": {
				Description: "Type of the IDP",
				Type:        schema.TypeString,
				Computed: true,
			},
			"redirect_url": {
				Description: "API Callback url for the identity provider",
				Type:        schema.TypeString,
				Computed: true,
			},
			"redirect_ui_url": {
				Description: "UI Callback url for the identity provider",
				Type:        schema.TypeString,
				Computed: true,
			},
			"login_url": {
				Type:        schema.TypeString,
				Computed: true,
			},
			"github": {
				Description: "Settings for GitHub IDP",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Required: true,
							Sensitive: true,
						},
						"authentication_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "https://github.com/login/oauth/authorize",
						},
						"token_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "https://github.com/login/oauth/access_token",
						},
						"user_profile_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "https://api.github.com/user",
						},
						"api_host": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "api.github.com",
						},
						"api_path_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "/",
						},
					},
				},
			},
		},
	}
}

func resourceIDPCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	resp, err := client.CreateIDP(mapResourceToIDP(d))

	if err != nil {
		log.Printf("[DEBUG] Error while creating idp. Error = %v", err)
		return err
	}

	d.SetId(resp.ID)
	return resourceIDPRead(d, meta)
}

func resourceIDPRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	idpID := d.Id()

	cfClientIDP, err := client.GetIdpByID(idpID)

	if err != nil {
		log.Printf("[DEBUG] Error while getting IDP. Error = %v", err)
		return err
	}

	err = mapIDPToResource(*cfClientIDP, d)

	if err != nil {
		log.Printf("[DEBUG] Error while getting mapping response to IDP object. Error = %v", err)
		return err
	}

	return nil
}

func resourceIDPDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIDPUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func mapIDPToResource(cfClientIDP cfclient.IDP, d *schema.ResourceData) error {
	d.SetId(cfClientIDP.ID)
	d.Set("display_name", cfClientIDP.DisplayName)
	d.Set("name", cfClientIDP.ClientName)
	d.Set("redirect_url", cfClientIDP.RedirectUrl)
	d.Set("redirect_ui_url", cfClientIDP.RedirectUiUrl)
	d.Set("login_url", cfClientIDP.LoginUrl)

	if cfClientIDP.ClientType == "github" {
		d.Set("github.0.client_id", cfClientIDP.ClientId)
		d.Set("github.0.client_secret", cfClientIDP.ClientSecret)
		d.Set("github.0.authentication_url", cfClientIDP.AuthURL)
		d.Set("github.0.token_url", cfClientIDP.TokenURL)
		d.Set("github.0.user_profile_url", cfClientIDP.UserProfileURL)
		d.Set("github.0.api_host", cfClientIDP.ApiHost)
		d.Set("github.0.api_path_prefix", cfClientIDP.ApiPathPrefix)	
	}

	return nil
}

func mapResourceToIDP(d *schema.ResourceData) *cfclient.IDP {
	cfClientIDP := &cfclient.IDP{
		ID:               d.Id(),
		DisplayName: 	  d.Get("display_name").(string),
		ClientName:       d.Get("name").(string),
		RedirectUrl:      d.Get("redirect_url").(string),
		RedirectUiUrl:    d.Get("redirect_ui_url").(string),
		LoginUrl:         d.Get("login_url").(string),
	}

	// client_type - Set dynamically
	if _, ok := d.GetOk("github"); ok {
		cfClientIDP.ClientType = "github"
		cfClientIDP.ClientId = d.Get("github.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("github.0.client_secret").(string)
		cfClientIDP.AuthURL = d.Get("github.0.authentication_url").(string)
		cfClientIDP.TokenURL = d.Get("github.0.token_url").(string)
		cfClientIDP.UserProfileURL = d.Get("github.0.user_profile_url").(string)
		cfClientIDP.ApiHost = d.Get("github.0.api_host").(string)
		cfClientIDP.ApiPathPrefix = d.Get("github.0.api_path_prefix").(string)
	}

	return cfClientIDP
}



