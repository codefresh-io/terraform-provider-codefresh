package idp

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	SupportedIdps = []string{GitHub, GitLab, Okta, Google, Auth0, Auth0, OneLogin, Keycloak, SAML, LDAP}
	IdpSchema     = map[string]*schema.Schema{
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
			Description: "Login url using the IDP to Codefresh",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"github": {
			Description:  "Settings for GitHub IDP",
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: SupportedIdps,
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
			ExactlyOneOf: SupportedIdps,
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
			ExactlyOneOf: SupportedIdps,
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
				},
			},
		},
		"google": {
			Description:  "Settings for Google IDP",
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: SupportedIdps,
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
					"admin_email": {
						Type:        schema.TypeString,
						Description: "Email of a user with admin permissions on google, relevant only for synchronization",
						Optional:    true,
					},
					"json_keyfile": {
						Type:        schema.TypeString,
						Description: "JSON keyfile for google service account used for synchronization",
						Optional:    true,
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
			ExactlyOneOf: SupportedIdps,
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
			ExactlyOneOf: SupportedIdps,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"client_secret": {
						Type:        schema.TypeString,
						Description: "Client secret from Azure",
						Required:    true,
						Sensitive:   true,
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
			ExactlyOneOf: SupportedIdps,
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
						// DiffSuppressFunc: surpressObfuscatedFields(),
					},
				},
			},
		},
		"keycloak": {
			Description:  "Settings for Keycloak IDP",
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: SupportedIdps,
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
			ExactlyOneOf: SupportedIdps,
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
					"admin_email": {
						Type:        schema.TypeString,
						Description: "Valid for GSuite only: Email of a user with admin permissions on google, relevant only for synchronization",
						Optional:    true,
					},
					"access_token": {
						Type:        schema.TypeString,
						Description: "Valid for Okta only: The Okta API token generated in Okta, used to sync groups and their users from Okta to Codefresh",
						Optional:    true,
					},
				},
			},
		},
		"ldap": {
			Description:  "Settings for Keycloak IDP",
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: SupportedIdps,
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
)
