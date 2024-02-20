package codefresh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var supportedIdps = []string{"github", "gitlab", "okta", "google", "auth0", "azure", "onelogin", "keycloak", "saml", "ldap"}
var idpSchema = map[string]*schema.Schema{
	"display_name": {
		Description: "The display name for the IDP.",
		Type:        schema.TypeString,
		Required:    true,
	},
	"name": {
		Description: "Name of the IDP, will be generated if not set",
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    true,
	},
	"client_type": {
		Description: "Type of the IDP. Derived from idp specific config object (github, gitlab etc)",
		Type:        schema.TypeString,
		Computed:    true,
		ForceNew:    true,
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
		Type:     schema.TypeString,
		Computed: true,
	},
	"config_hash": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"github": {
		Description:  "Settings for GitHub IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID from Github",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret from GitHub",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"authentication_url": {
					Type:        schema.TypeString,
					Description: "Authentication url, Defaults to https://github.com/login/oauth/authorize",
					Optional:    true,
					Default:     "https://github.com/login/oauth/authorize",
				},
				"token_url": {
					Type:        schema.TypeString,
					Description: "GitHub token endpoint url, Defaults to https://github.com/login/oauth/access_token",
					Optional:    true,
					Default:     "https://github.com/login/oauth/access_token",
				},
				"user_profile_url": {
					Type:        schema.TypeString,
					Description: "GitHub user profile url, Defaults to https://api.github.com/user",
					Optional:    true,
					Default:     "https://api.github.com/user",
				},
				"api_host": {
					Type:        schema.TypeString,
					Description: "GitHub API host, Defaults to api.github.com",
					Optional:    true,
					Default:     "api.github.com",
				},
				"api_path_prefix": {
					Type:        schema.TypeString,
					Description: "GitHub API url path prefix, defaults to /",
					Optional:    true,
					Default:     "/",
				},
			},
		},
	},
	"gitlab": {
		Description:  "Settings for GitLab IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID from Gitlab",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret from Gitlab",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"authentication_url": {
					Type:        schema.TypeString,
					Description: "Authentication url, Defaults to https://gitlab.com",
					Optional:    true,
					Default:     "https://gitlab.com",
				},
				"user_profile_url": {
					Type:        schema.TypeString,
					Description: "User profile url, Defaults to https://gitlab.com/api/v4/user",
					Optional:    true,
					Default:     "https://gitlab.com/api/v4/user",
				},
				"api_url": {
					Type:        schema.TypeString,
					Description: "Base url for Gitlab API, Defaults to https://gitlab.com/api/v4/",
					Optional:    true,
					Default:     "https://gitlab.com/api/v4/",
				},
			},
		},
	},
	"okta": {
		Description:  "Settings for Okta IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID in Okta, must be unique across all identity providers in Codefresh",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret in Okta",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"client_host": {
					Type:         schema.TypeString,
					Description:  "The OKTA organization URL, for example, https://<company>.okta.com",
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(https?:\/\/)(\S+)(\.okta(preview|-emea)?\.com$)`), "must be a valid okta url"),
					Required:     true,
				},
				"app_id": {
					Type:        schema.TypeString,
					Description: "The Codefresh application ID in your OKTA organization",
					Optional:    true,
				},
				"app_id_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed app id in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"sync_mirror_accounts": {
					Type:        schema.TypeList,
					Description: "The names of the additional Codefresh accounts to be synced from Okta",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"access_token": {
					Type:        schema.TypeString,
					Description: "The Okta API token generated in Okta, used to sync groups and their users from Okta to Codefresh",
					Optional:    true,
				},
				"access_token_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed access token in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
			},
		},
	},
	"google": {
		Description:  "Settings for Google IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID in Google, must be unique across all identity providers in Codefresh",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret in Google",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"admin_email": {
					Type:        schema.TypeString,
					Description: "Email of a user with admin permissions on google, relevant only for synchronization",
					Optional:    true,
				},
				"admin_email_encrypted": {
					Type:        schema.TypeString,
					Description: "Admin email in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"json_keyfile": {
					Type:        schema.TypeString,
					Description: "JSON keyfile for google service account used for synchronization",
					Optional:    true,
				},
				"json_keyfile_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed JSON keyfile in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"allowed_groups_for_sync": {
					Type:        schema.TypeString,
					Description: "Comma separated list of groups to sync",
					Optional:    true,
				},
				"sync_field": {
					Type:        schema.TypeString,
					Description: "Relevant for custom schema-based synchronization only. See Codefresh documentation",
					Optional:    true,
				},
			},
		},
	},
	"auth0": {
		Description:  "Settings for Auth0 IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID from Auth0",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret from Auth0",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"domain": {
					Type:        schema.TypeString,
					Description: "The domain of the Auth0 application",
					Required:    true,
				},
			},
		},
	},
	"azure": {
		Description:  "Settings for Azure IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret from Azure",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"app_id": {
					Type:        schema.TypeString,
					Description: "The Application ID from your Enterprise Application Properties in Azure AD",
					Required:    true,
				},
				"tenant": {
					Type:        schema.TypeString,
					Description: "Azure tenant",
					Optional:    true,
				},
				"object_id": {
					Type:        schema.TypeString,
					Description: "The Object ID from your Enterprise Application Properties in Azure AD",
					Optional:    true,
				},
				"autosync_teams_and_users": {
					Type:        schema.TypeBool,
					Description: "Set to true to sync user accounts in Azure AD to your Codefresh account",
					Optional:    true,
					Default:     false,
				},
				"sync_interval": {
					Type:        schema.TypeInt,
					Description: "Sync interval in hours for syncing user accounts in Azure AD to your Codefresh account. If not set the sync inteval will be 12 hours",
					Optional:    true,
				},
			},
		},
	},
	"onelogin": {
		Description:  "Settings for onelogin IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID from Onelogin",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret from Onelogin",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"domain": {
					Type:        schema.TypeString,
					Description: "The domain to be used for authentication",
					Required:    true,
				},
				"app_id": {
					Type:        schema.TypeString,
					Description: "The Codefresh application ID in your Onelogin",
					Optional:    true,
				},
				"api_client_id": {
					Type:        schema.TypeString,
					Description: "Client ID for onelogin API, only needed if syncing users and groups from Onelogin",
					Optional:    true,
				},
				"api_client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret for onelogin API, only needed if syncing users and groups from Onelogin",
					Optional:    true,
					// When onelogin IDP is created on account level, after the first apply the client secret is returned obfuscated
					//DiffSuppressFunc: surpressObfuscatedFields(),
				},
			},
		},
	},
	"keycloak": {
		Description:  "Settings for Keycloak IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "Client ID from Keycloak",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "Client secret from Keycloak",
					Required:    true,
					Sensitive:   true,
				},
				"client_secret_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed client secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"host": {
					Type:         schema.TypeString,
					Description:  "The Keycloak URL",
					Required:     true,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(https?:\/\/)(\S+)$`), "must be a valid url"),
				},
				"realm": {
					Type:        schema.TypeString,
					Description: "The Realm ID for Codefresh in Keycloak. Defaults to master",
					Optional:    true,
					Default:     "master",
				},
			},
		},
	},
	"saml": {
		Description:  "Settings for SAML IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"endpoint": {
					Type:        schema.TypeString,
					Description: "The SSO endpoint of your Identity Provider",
					Required:    true,
				},
				"application_certificate": {
					Type:        schema.TypeString,
					Description: "The security certificate of your Identity Provider. Paste the value directly on the field. Do not convert to base64 or any other encoding by hand",
					Required:    true,
					Sensitive:   true,
				},
				"application_certificate_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed certificate secret in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"provider": {
					Type:         schema.TypeString,
					Description:  "SAML provider. Currently supported values - GSuite, okta or empty string for generic provider. Defaults to empty string",
					Optional:     true,
					Default:      "",
					ValidateFunc: validation.StringInSlice([]string{"", "okta", "GSuite"}, false),
				},
				"allowed_groups_for_sync": {
					Type:        schema.TypeString,
					Description: "Valid for GSuite only: Comma separated list of groups to sync",
					Optional:    true,
				},
				"autosync_teams_and_users": {
					Type:        schema.TypeBool,
					Description: "Valid for Okta/GSuite: Set to true to sync user accounts and teams in okta/gsuite to your Codefresh account",
					Optional:    true,
					Default:     false,
				},
				"sync_interval": {
					Type:        schema.TypeInt,
					Description: "Valid for Okta/GSuite: Sync interval in hours for syncing user accounts in okta/gsuite to your Codefresh account. If not set the sync inteval will be 12 hours",
					Optional:    true,
				},
				"activate_users_after_sync": {
					Type:        schema.TypeBool,
					Description: "Valid for Okta only: If set to true, Codefresh will automatically invite and activate new users added during the automated sync, without waiting for the users to accept the invitations. Defaults to false",
					Optional:    true,
					Default:     false,
				},
				"app_id": {
					Type:        schema.TypeString,
					Description: "Valid for Okta only: The Codefresh application ID in Okta",
					Optional:    true,
				},
				"client_host": {
					Type:        schema.TypeString,
					Description: "Valid for Okta only: OKTA organization URL, for example, https://<company>.okta.com",
					Optional:    true,
				},
				"json_keyfile": {
					Type:        schema.TypeString,
					Description: "Valid for GSuite only: JSON keyfile for google service account used for synchronization",
					Optional:    true,
				},
				"json_keyfile_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed JSON keyfile in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"admin_email": {
					Type:        schema.TypeString,
					Description: "Valid for GSuite only: Email of a user with admin permissions on google, relevant only for synchronization",
					Optional:    true,
				},
				"admin_email_encrypted": {
					Type:        schema.TypeString,
					Description: "Admin email in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"access_token": {
					Type:        schema.TypeString,
					Description: "Valid for Okta only: The Okta API token generated in Okta, used to sync groups and their users from Okta to Codefresh",
					Optional:    true,
				},
				"access_token_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed access token in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
			},
		},
	},
	"ldap": {
		Description:  "Settings for Keycloak IDP",
		Type:         schema.TypeList,
		Optional:     true,
		MaxItems:     1,
		ExactlyOneOf: supportedIdps,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:         schema.TypeString,
					Description:  "ldap server url",
					Required:     true,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ldap(s?):\/\/`), "must be a valid ldap url (must start with ldap:// or ldaps://)"),
				},
				"password": {
					Type:        schema.TypeString,
					Description: "The password of the user defined in Distinguished name that will be used to search other users",
					Required:    true,
					Sensitive:   true,
				},
				"password_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed password in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"distinguished_name": {
					Type:        schema.TypeString,
					Description: "The username to be used to search other users in LDAP notation (combination of cn, ou,dc)",
					Optional:    true,
					Computed:    true,
				},
				"search_base": {
					Type:        schema.TypeString,
					Description: "The search-user scope in LDAP notation",
					Required:    true,
				},
				"search_filter": {
					Type:        schema.TypeString,
					Description: "The attribute by which to search for the user on the LDAP server. By default, set to uid. For the Azure LDAP server, set this field to sAMAccountName",
					Optional:    true,
				},
				"certificate": {
					Type:        schema.TypeString,
					Description: "For ldaps only: The security certificate of the LDAP server. Do not convert to base64 or any other encoding",
					Optional:    true,
				},
				"certificate_encrypted": {
					Type:        schema.TypeString,
					Description: "Computed certificate in encrypted form as returned from Codefresh API. Only Codefresh can decrypt this value",
					Optional:    true,
					Computed:    true,
				},
				"allowed_groups_for_sync": {
					Type:        schema.TypeString,
					Description: "To sync only by specified groups - specify a comma separated list of groups, by default all groups will be synced",
					Optional:    true,
				},
				"search_base_for_sync": {
					Type:        schema.TypeString,
					Description: "Synchronize using a custom search base, by deafult seach_base is used",
					Optional:    true,
				},
			},
		},
	},
}

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
		),
		Schema: idpSchema,
	}
}

func resourceIDPCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	id, err := client.CreateIDP(mapResourceToIDP(d), true)

	if err != nil {
		log.Printf("[DEBUG] Error while creating idp. Error = %v", err)
		return err
	}

	d.SetId(id)
	return resourceIDPRead(d, meta)
}

func resourceIDPRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	idpID := d.Id()

	var cfClientIDP *cfclient.IDP
	var err error

	cfClientIDP, err = client.GetIdpByID(idpID)

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

	var cfClientIDP *cfclient.IDP
	var err error

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

	return nil
}

func resourceIDPUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	err := client.UpdateIDP(mapResourceToIDP(d), true)

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
			"client_id": cfClientIDP.ClientId,
			// Codefresh API Returns the client secret as an encrypted string on the server side
			// hence we need to keep in the state the original secret the user provides along with the encrypted computed secret
			// for Terraform to properly calculate the diff
			"client_secret":           d.Get("github.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"authentication_url":      cfClientIDP.AuthURL,
			"token_url":               cfClientIDP.TokenURL,
			"user_profile_url":        cfClientIDP.UserProfileURL,
			"api_host":                cfClientIDP.ApiHost,
			"api_path_prefix":         cfClientIDP.ApiPathPrefix,
		}}

		d.Set("github", attributes)
	}

	if cfClientIDP.ClientType == "gitlab" {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("gitlab.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"authentication_url":      cfClientIDP.AuthURL,
			"user_profile_url":        cfClientIDP.UserProfileURL,
			"api_url":                 cfClientIDP.ApiURL,
		}}

		d.Set("gitlab", attributes)
	}

	if cfClientIDP.ClientType == "okta" {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("okta.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"client_host":             cfClientIDP.ClientHost,
			"app_id":                  d.Get("okta.0.app_id"),
			"app_id_encrypted":        cfClientIDP.AppId,
			"sync_mirror_accounts":    cfClientIDP.SyncMirrorAccounts,
			"access_token":            d.Get("okta.0.access_token"),
			"access_token_encrypted":  cfClientIDP.Access_token,
		}}

		d.Set("okta", attributes)
	}

	if cfClientIDP.ClientType == "google" {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("google.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"admin_email":             d.Get("google.0.admin_email"),
			"admin_email_encrypted":   cfClientIDP.Subject,
			"json_keyfile":            d.Get("google.0.json_keyfile"),
			"json_keyfile_encrypted":  cfClientIDP.KeyFile,
			"allowed_groups_for_sync": cfClientIDP.AllowedGroupsForSync,
			"sync_field":              cfClientIDP.SyncField,
		}}

		d.Set("google", attributes)
	}

	if cfClientIDP.ClientType == "auth0" {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("auth0.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"domain":                  cfClientIDP.ClientHost,
		}}

		d.Set("auth0", attributes)
	}

	if cfClientIDP.ClientType == "azure" {

		syncInterval, err := strconv.Atoi(cfClientIDP.SyncInterval)

		if err != nil {
			return err
		}

		attributes := []map[string]interface{}{{
			"app_id":                   cfClientIDP.ClientId,
			"client_secret":            d.Get("azure.0.client_secret"),
			"client_secret_encrypted":  cfClientIDP.ClientSecret,
			"object_id":                cfClientIDP.AppId,
			"autosync_teams_and_users": cfClientIDP.AutoGroupSync,
			"sync_interval":            syncInterval,
			"tenant":                   cfClientIDP.Tenant,
		}}

		d.Set("azure", attributes)
	}

	if cfClientIDP.ClientType == "onelogin" {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("onelogin.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"domain":                  cfClientIDP.ClientHost,
			"api_client_id":           cfClientIDP.ApiClientId,

			"api_client_secret": cfClientIDP.ApiClientSecret,
			"app_id":            cfClientIDP.AppId,
		}}

		d.Set("onelogin", attributes)
	}

	if cfClientIDP.ClientType == "keycloak" {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("keycloak.0.client_secret"),
			"client_secret_encrypted": cfClientIDP.ClientSecret,
			"host":                    cfClientIDP.Host,
			"realm":                   cfClientIDP.Realm,
		}}

		d.Set("keycloak", attributes)
	}

	if cfClientIDP.ClientType == "saml" {
		syncInterval, err := strconv.Atoi(cfClientIDP.SyncInterval)

		if err != nil {
			return err
		}
		attributes := []map[string]interface{}{{
			"endpoint":                          cfClientIDP.EntryPoint,
			"application_certificate":           d.Get("saml.0.application_certificate"),
			"application_certificate_encrypted": cfClientIDP.ApplicationCert,
			"provider":                          cfClientIDP.SamlProvider,
			"allowed_groups_for_sync":           cfClientIDP.AllowedGroupsForSync,
			"autosync_teams_and_users":          cfClientIDP.AutoGroupSync,
			"activate_users_after_sync":         cfClientIDP.ActivateUserAfterSync,
			"sync_interval":                     syncInterval,
			"app_id":                            cfClientIDP.AppId,
			"client_host":                       cfClientIDP.ClientHost,
			"json_keyfile":                      d.Get("saml.0.json_keyfile"),
			"json_keyfile_encrypted":            cfClientIDP.KeyFile,
			"admin_email":                       d.Get("saml.0.admin_email"),
			"admin_email_encrypted":             cfClientIDP.Subject,
			"access_token":                      d.Get("saml.0.access_token"),
			"access_token_encrypted":            cfClientIDP.Access_token,
		}}

		d.Set("saml", attributes)
	}

	if cfClientIDP.ClientType == "ldap" {
		attributes := []map[string]interface{}{{
			"url":                     cfClientIDP.Url,
			"password":                d.Get("ldap.0.password"),
			"password_encrypted":      cfClientIDP.Password,
			"distinguished_name":      cfClientIDP.DistinguishedName,
			"search_base":             cfClientIDP.SearchBase,
			"search_filter":           cfClientIDP.SearchFilter,
			"certificate":             d.Get("ldap.0.certificate"),
			"certificate_encrypted":   cfClientIDP.Certificate,
			"allowed_groups_for_sync": cfClientIDP.AllowedGroupsForSync,
			"search_base_for_sync":    cfClientIDP.SearchBaseForSync,
		}}

		d.Set("ldap", attributes)
	}

	return nil
}

func mapResourceToIDP(d *schema.ResourceData) *cfclient.IDP {

	cfClientIDP := &cfclient.IDP{
		ID:            d.Id(),
		DisplayName:   d.Get("display_name").(string),
		ClientName:    d.Get("name").(string),
		RedirectUrl:   d.Get("redirect_url").(string),
		RedirectUiUrl: d.Get("redirect_ui_url").(string),
		LoginUrl:      d.Get("login_url").(string),
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
		cfClientIDP.Access_token = d.Get("okta.0.access_token").(string)
	}

	if _, ok := d.GetOk("google"); ok {
		cfClientIDP.ClientType = "google"
		cfClientIDP.ClientId = d.Get("google.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("google.0.client_secret").(string)
		cfClientIDP.KeyFile = d.Get("google.0.json_keyfile").(string)
		cfClientIDP.Subject = d.Get("google.0.admin_email").(string)
		cfClientIDP.AllowedGroupsForSync = d.Get("google.0.allowed_groups_for_sync").(string)
		cfClientIDP.SyncField = d.Get("google.0.sync_field").(string)
	}

	if _, ok := d.GetOk("auth0"); ok {
		cfClientIDP.ClientType = "auth0"
		cfClientIDP.ClientId = d.Get("auth0.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("auth0.0.client_secret").(string)
		cfClientIDP.ClientHost = d.Get("auth0.0.domain").(string)
	}

	if _, ok := d.GetOk("azure"); ok {
		cfClientIDP.ClientType = "azure"
		cfClientIDP.ClientId = d.Get("azure.0.app_id").(string)
		cfClientIDP.ClientSecret = d.Get("azure.0.client_secret").(string)
		cfClientIDP.AppId = d.Get("azure.0.object_id").(string)
		cfClientIDP.Tenant = d.Get("azure.0.tenant").(string)
		cfClientIDP.AutoGroupSync = d.Get("azure.0.autosync_teams_and_users").(bool)
		cfClientIDP.SyncInterval = strconv.Itoa(d.Get("azure.0.sync_interval").(int))
	}

	if _, ok := d.GetOk("onelogin"); ok {
		cfClientIDP.ClientType = "onelogin"
		cfClientIDP.ClientId = d.Get("onelogin.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("onelogin.0.client_secret").(string)
		cfClientIDP.ClientHost = d.Get("onelogin.0.domain").(string)
		cfClientIDP.AppId = d.Get("onelogin.0.app_id").(string)
		cfClientIDP.ApiClientId = d.Get("onelogin.0.api_client_id").(string)
		cfClientIDP.ApiClientSecret = d.Get("onelogin.0.api_client_secret").(string)
	}

	if _, ok := d.GetOk("keycloak"); ok {
		cfClientIDP.ClientType = "keycloak"
		cfClientIDP.ClientId = d.Get("keycloak.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("keycloak.0.client_secret").(string)
		cfClientIDP.Host = d.Get("keycloak.0.host").(string)
		cfClientIDP.Realm = d.Get("keycloak.0.realm").(string)
	}

	if _, ok := d.GetOk("saml"); ok {
		cfClientIDP.ClientType = "saml"
		cfClientIDP.SamlProvider = d.Get("saml.0.provider").(string)
		cfClientIDP.EntryPoint = d.Get("saml.0.endpoint").(string)
		cfClientIDP.ApplicationCert = d.Get("saml.0.application_certificate").(string)
		cfClientIDP.AllowedGroupsForSync = d.Get("saml.0.allowed_groups_for_sync").(string)
		cfClientIDP.AutoGroupSync = d.Get("saml.0.autosync_teams_and_users").(bool)
		cfClientIDP.ActivateUserAfterSync = d.Get("saml.0.activate_users_after_sync").(bool)
		cfClientIDP.SyncInterval = strconv.Itoa(d.Get("saml.0.sync_interval").(int))
		cfClientIDP.AppId = d.Get("saml.0.app_id").(string)
		cfClientIDP.ClientHost = d.Get("saml.0.client_host").(string)
		cfClientIDP.KeyFile = d.Get("saml.0.json_keyfile").(string)
		cfClientIDP.Subject = d.Get("saml.0.admin_email").(string)
		cfClientIDP.Access_token = d.Get("saml.0.access_token").(string)
	}

	if _, ok := d.GetOk("ldap"); ok {
		cfClientIDP.ClientType = "ldap"
		cfClientIDP.Url = d.Get("ldap.0.url").(string)
		cfClientIDP.Password = d.Get("ldap.0.password").(string)
		cfClientIDP.DistinguishedName = d.Get("ldap.0.distinguished_name").(string)
		cfClientIDP.SearchBase = d.Get("ldap.0.search_base").(string)
		cfClientIDP.SearchFilter = d.Get("ldap.0.search_filter").(string)
		cfClientIDP.Certificate = d.Get("ldap.0.certificate").(string)
		cfClientIDP.AllowedGroupsForSync = d.Get("ldap.0.allowed_groups_for_sync").(string)
		cfClientIDP.SearchBaseForSync = d.Get("ldap.0.search_base_for_sync").(string)
	}

	return cfClientIDP
}
