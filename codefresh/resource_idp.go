package codefresh

import (
	"fmt"
	"log"
	"errors"
	"context"
	"regexp"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
)

var supportedIdps = []string{"github","gitlab", "okta"}

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
		CustomizeDiff: customdiff.All(
			// Recreate idp if the type has changed - we cannot simply do ForceNew on client_type as it is computed
			customdiff.ForceNewIf("client_type", func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				clientTypeInState := d.Get("client_type").(string)
				attributesForIdpTypeInState := d.Get(clientTypeInState)
				// If there is a different type of idp in the state, the idp needs to be recreated
				if attributesForIdpTypeInState == nil {
					d.SetNewComputed("client_type")
					return true
				} else if len(attributesForIdpTypeInState.([]interface{})) < 1 {
					d.SetNewComputed("client_type")
					return true
				} else {
					return false
				}
			}),
			// If name has changed for an account scoped IDP the provider needs to ignore it as the API always generates the name
			customdiff.If(func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				bIsGlobal := d.Get("is_global").(bool)
				return !bIsGlobal
				}, 
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error{ 
				old, _ := d.GetChange("name")
				if err := d.SetNew("name",old); err != nil {
					return err
				}
				return nil
			}),
		),
		Schema: map[string]*schema.Schema{
			"is_global": {
				Type: schema.TypeBool,
				Description: "If set to true IDP will be created globally for the entire platform - this requires a platform admin token. If false the IDP will be created at the level of a single account which is derived from the API token used. Defaults to false",
				Optional: true,
				Default: false,
				ForceNew: true,
			},
			"display_name": {
				Description: "The display name for the IDP.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the IDP, will be generated if not set",
				Type:        schema.TypeString,
				Computed: 	 true,
				Optional: true,
			},
			"client_type": {
				Description: "Type of the IDP. If not set it is derived from idp specific config object (github, gitlab etc)",
				Type:        schema.TypeString,
				Computed: true,
				ForceNew: true,
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
			"config_hash": {
				Type:        schema.TypeString,
				Computed: true,
			},
			"github": {
				Description: "Settings for GitHub IDP",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems: 1,
				ExactlyOneOf: supportedIdps,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Description: "Client ID from Github",
							Required: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Description: "Client secret from GitHub",
							Required: true,
							Sensitive: true,
						},
						"client_secret_encrypted": {
							Type:     schema.TypeString,
							Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
							Optional: true,
							Computed: true,
						},
						"authentication_url": {
							Type:     schema.TypeString,
							Description: "Authentication url, Defaults to https://github.com/login/oauth/authorize",
							Optional: true,
							Default: "https://github.com/login/oauth/authorize",
						},
						"token_url": {
							Type:     schema.TypeString,
							Description: "GitHub token endpoint url, Defaults to https://github.com/login/oauth/access_token",
							Optional: true,
							Default: "https://github.com/login/oauth/access_token",
						},
						"user_profile_url": {
							Type:     schema.TypeString,
							Description: "GitHub user profile url, Defaults to https://api.github.com/user",
							Optional: true,
							Default: "https://api.github.com/user",
						},
						"api_host": {
							Type:     schema.TypeString,
							Description: "GitHub API host, Defaults to api.github.com",
							Optional: true,
							Default: "api.github.com",
						},
						"api_path_prefix": {
							Type:     schema.TypeString,
							Description: "GitHub API url path prefix, defaults to /",
							Optional: true,
							Default: "/",
						},
					},
				},
			},
			"gitlab": {
				Description: "Settings for GitLab IDP",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems: 1,
				ExactlyOneOf: supportedIdps,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Description: "Client ID from Gitlab",
							Required: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Description: "Client secret from Gitlab",
							Required: true,
							Sensitive: true,
						},
						"client_secret_encrypted": {
							Type:     schema.TypeString,
							Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
							Optional: true,
							Computed: true,
						},
						"authentication_url": {
							Type:     schema.TypeString,
							Description: "Authentication url, Defaults to https://gitlab.com",
							Optional: true,
							Default: "https://gitlab.com",
						},
						"user_profile_url": {
							Type:     schema.TypeString,
							Description: "User profile url, Defaults to https://gitlab.com/api/v4/user",
							Optional: true,
							Default: "https://gitlab.com/api/v4/user",
						},
						"api_url": {
							Type:     schema.TypeString,
							Description: "Base url for Gitlab API, Defaults to https://gitlab.com/api/v4/",
							Optional: true,
							Default: "https://gitlab.com/api/v4/",
						},
					},
				},
			},
			"okta": {
				Description: "Settings for Okta IDP",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems: 1,
				ExactlyOneOf: supportedIdps,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Description: "Client ID in Okta, must be unique across all identity providers in Codefresh",
							Required: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Description: "Client secret in Okta",
							Required: true,
							Sensitive: true,
						},
						"client_secret_encrypted": {
							Type:     schema.TypeString,
							Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
							Optional: true,
							Computed: true,
						},
						"client_host": {
							Type:     schema.TypeString,
							Description: "The OKTA organization URL, for example, https://<company>.okta.com",
							ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile(`^(https?:\\/\\/)(\\S+)(\\.okta(preview|-emea)?\\.com$)`), "must be a valid okta url"),
							Required: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Description: "The Codefresh application ID in your OKTA organization",
							Optional: true,
						},
						"app_id_encrypted": {
							Type:     schema.TypeString,
							Description: "Computed app id in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
							Optional: true,
							Computed: true,
						},
						"sync_mirror_accounts": {
							Type:     schema.TypeList,
							Description: "The names of the additional Codefresh accounts to be synced from Okta",
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
	isGlobal := d.Get("is_global").(bool)

	var cfClientIDP *cfclient.IDP
	var err error

	if isGlobal {
		cfClientIDP, err = client.GetIdpByID(idpID)
	} else
	{
		cfClientIDP, err = client.GetAccountIdpByID(idpID)
	}

	if err != nil {
		if err.Error() == fmt.Sprintf("[ERROR] IDP with ID %s isn't found.", d.Id()) {
			d.SetId("")
			return nil
		}
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
	client := meta.(*cfclient.Client)
	idpID := d.Id()
	isGlobal := d.Get("is_global").(bool)

	var cfClientIDP *cfclient.IDP
	var err error

	if isGlobal {
		cfClientIDP, err = client.GetIdpByID(idpID)

		if err != nil {
			log.Printf("[DEBUG] Error while getting IDP. Error = %v", err)
			return err
		}

		if len(cfClientIDP.Accounts) < 1 {
			return errors.New("It is not allowed to delete IDPs without any assigned accounts as they are considered global. Assign at least one account before deleting")
		}

		err = client.DeleteIDP(d.Id())

		if err != nil {
			log.Printf("[DEBUG] Error while deleting IDP. Error = %v", err)
			return err
		}
	} else
	{
		err = client.DeleteIDPAccount(d.Id())

		if err != nil {
			log.Printf("[DEBUG] Error while deleting account level IDP. Error = %v", err)
			return err
		}
	}

	return nil
}

func resourceIDPUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	err := client.UpdateIDP(mapResourceToIDP(d))

	if err != nil {
		log.Printf("[DEBUG] Error while updating idp. Error = %v", err)
		return err
	}

	return resourceIDPRead(d, meta)
}

func mapIDPToResource(cfClientIDP cfclient.IDP, d *schema.ResourceData) error {
	d.SetId(cfClientIDP.ID)
	d.Set("display_name", cfClientIDP.DisplayName)
	d.Set("name", cfClientIDP.ClientName)
	d.Set("redirect_url", cfClientIDP.RedirectUrl)
	d.Set("redirect_ui_url", cfClientIDP.RedirectUiUrl)
	d.Set("login_url", cfClientIDP.LoginUrl)
	d.Set("client_type", cfClientIDP.ClientType)

	if cfClientIDP.ClientType == "github" {
		attributes := []map[string]interface{}{{
			"client_id":            		cfClientIDP.ClientId,
			// Codefresh API Returns the client secret as an encrypted string on the server side
			// hence we need to keep in the state the original secret the user provides along with the encrypted computed secret
			// for Terraform to properly calculate the diff
			"client_secret": 				d.Get("github.0.client_secret"),
			"client_secret_encrypted": 		cfClientIDP.ClientSecret,
			"authentication_url":   		cfClientIDP.AuthURL,
			"token_url": 					cfClientIDP.TokenURL,
			"user_profile_url": 			cfClientIDP.UserProfileURL,
			"api_host": 					cfClientIDP.ApiHost,
			"api_path_prefix":      		cfClientIDP.ApiPathPrefix,
		}}

		d.Set("github", attributes)
	}

	if cfClientIDP.ClientType == "gitlab" {
		attributes := []map[string]interface{}{{
			"client_id":            		cfClientIDP.ClientId,
			"client_secret": 				d.Get("gitlab.0.client_secret"),
			"client_secret_encrypted": 		cfClientIDP.ClientSecret,
			"authentication_url":   		cfClientIDP.AuthURL,
			"user_profile_url": 			cfClientIDP.UserProfileURL,
			"api_url": 						cfClientIDP.ApiURL,
		}}

		d.Set("gitlab", attributes)
	}

	if cfClientIDP.ClientType == "okta" {
		attributes := []map[string]interface{}{{
			"client_id":            		cfClientIDP.ClientId,
			"client_secret": 				d.Get("okta.0.client_secret"),
			"client_secret_encrypted": 		cfClientIDP.ClientSecret,
			"client_host":   				cfClientIDP.ClientHost,
			"app_id": 						d.Get("okta.0.app_id"),
			"app_id_encrypted":             cfClientIDP.AppId,
			"sync_mirror_accounts": 		cfClientIDP.SyncMirrorAccounts,
		}}

		d.Set("okta", attributes)
	}

	return nil
}

func mapResourceToIDP(d *schema.ResourceData) *cfclient.IDP {
	cfClientIDP := &cfclient.IDP{
		ID:               d.Id(),
		IsGlobal: 		  d.Get("is_global").(bool),
		DisplayName: 	  d.Get("display_name").(string),
		ClientName:       d.Get("name").(string),
		RedirectUrl:      d.Get("redirect_url").(string),
		RedirectUiUrl:    d.Get("redirect_ui_url").(string),
		LoginUrl:         d.Get("login_url").(string),
	}

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

	if _, ok := d.GetOk("gitlab"); ok {
		cfClientIDP.ClientType = "gitlab"
		cfClientIDP.ClientId = d.Get("gitlab.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("gitlab.0.client_secret").(string)
		cfClientIDP.AuthURL = d.Get("gitlab.0.authentication_url").(string)
		cfClientIDP.UserProfileURL = d.Get("gitlab.0.user_profile_url").(string)
		cfClientIDP.ApiURL = d.Get("gitlab.0.api_url").(string)
	}

	if _, ok := d.GetOk("okta"); ok {
		cfClientIDP.ClientType = "okta"
		cfClientIDP.ClientId = d.Get("okta.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("okta.0.client_secret").(string)
		cfClientIDP.ClientHost = d.Get("okta.0.client_host").(string)
		cfClientIDP.AppId = d.Get("okta.0.app_id").(string)
		cfClientIDP.SyncMirrorAccounts = datautil.ConvertStringArr(d.Get("okta.0.sync_mirror_accounts").([]interface{}))
	}

	return cfClientIDP
}
