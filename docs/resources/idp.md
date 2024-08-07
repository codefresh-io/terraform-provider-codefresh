---
page_title: "codefresh_idp Resource - terraform-provider-codefresh"
subcategory: ""
description: |-
  Codefresh global level identity provider. Requires a Codefresh admin token and applies only to Codefresh on-premises installations.
---

# codefresh_idp (Resource)

Codefresh global level identity provider. Requires a Codefresh admin token and applies only to Codefresh on-premises installations.

## Example usage
```hcl
resource "codefresh_idp" "auth0-test" {
  display_name = "tf-auth0-example"
  
  auth0 {
    client_id = "auht0-codefresh-example"
    client_secret = "mysecret"
    domain = "codefresh.auth0.com"
  }
}
```
```hcl
resource "codefresh_idp" "google-example" {
  display_name = "tf-google-example"

  google {
    client_id = "google-codefresh-example"
    client_secret = "mysecret99"
    admin_email = "admin@codefresh.io"
    sync_field = "myfield"
    json_keyfile = <<EOT
    {  
      "installed":{  
          "client_id":"clientid",
          "project_id":"projectname",
          "auth_uri":"https://accounts.google.com/o/oauth2/auth",
          "token_uri":"https://accounts.google.com/o/oauth2/token",
          "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
      }
    }
    EOT
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the IDP.

### Optional

- `auth0` (Block List, Max: 1) Settings for Auth0 IDP (see [below for nested schema](#nestedblock--auth0))
- `azure` (Block List, Max: 1) Settings for Azure IDP (see [below for nested schema](#nestedblock--azure))
- `github` (Block List, Max: 1) Settings for GitHub IDP (see [below for nested schema](#nestedblock--github))
- `gitlab` (Block List, Max: 1) Settings for GitLab IDP (see [below for nested schema](#nestedblock--gitlab))
- `google` (Block List, Max: 1) Settings for Google IDP (see [below for nested schema](#nestedblock--google))
- `keycloak` (Block List, Max: 1) Settings for Keycloak IDP (see [below for nested schema](#nestedblock--keycloak))
- `ldap` (Block List, Max: 1) Settings for Keycloak IDP (see [below for nested schema](#nestedblock--ldap))
- `name` (String) Name of the IDP, will be generated if not set
- `okta` (Block List, Max: 1) Settings for Okta IDP (see [below for nested schema](#nestedblock--okta))
- `onelogin` (Block List, Max: 1) Settings for onelogin IDP (see [below for nested schema](#nestedblock--onelogin))
- `saml` (Block List, Max: 1) Settings for SAML IDP (see [below for nested schema](#nestedblock--saml))

### Read-Only

- `client_type` (String) Type of the IDP. Derived from idp specific config object (github, gitlab etc)
- `id` (String) The ID of this resource.
- `login_url` (String) Login url using the IDP to Codefresh
- `redirect_ui_url` (String) UI Callback url for the identity provider
- `redirect_url` (String) API Callback url for the identity provider

<a id="nestedblock--auth0"></a>
### Nested Schema for `auth0`

Required:

- `client_id` (String) Client ID from Auth0
- `client_secret` (String, Sensitive) Client secret from Auth0
- `domain` (String) The domain of the Auth0 application


<a id="nestedblock--azure"></a>
### Nested Schema for `azure`

Required:

- `app_id` (String) The Application ID from your Enterprise Application Properties in Azure AD
- `client_secret` (String, Sensitive) Client secret from Azure

Optional:

- `autosync_teams_and_users` (Boolean) Set to true to sync user accounts in Azure AD to your Codefresh account
- `object_id` (String) The Object ID from your Enterprise Application Properties in Azure AD
- `sync_interval` (Number) Sync interval in hours for syncing user accounts in Azure AD to your Codefresh account. If not set the sync inteval will be 12 hours
- `tenant` (String) Azure tenant


<a id="nestedblock--github"></a>
### Nested Schema for `github`

Required:

- `client_id` (String) Client ID from Github
- `client_secret` (String, Sensitive) Client secret from GitHub

Optional:

- `api_host` (String) GitHub API host, Defaults to api.github.com
- `api_path_prefix` (String) GitHub API url path prefix, defaults to /
- `authentication_url` (String) Authentication url, Defaults to https://github.com/login/oauth/authorize
- `token_url` (String) GitHub token endpoint url, Defaults to https://github.com/login/oauth/access_token
- `user_profile_url` (String) GitHub user profile url, Defaults to https://api.github.com/user


<a id="nestedblock--gitlab"></a>
### Nested Schema for `gitlab`

Required:

- `client_id` (String) Client ID from Gitlab
- `client_secret` (String, Sensitive) Client secret from Gitlab

Optional:

- `api_url` (String) Base url for Gitlab API, Defaults to https://gitlab.com/api/v4/
- `authentication_url` (String) Authentication url, Defaults to https://gitlab.com
- `user_profile_url` (String) User profile url, Defaults to https://gitlab.com/api/v4/user


<a id="nestedblock--google"></a>
### Nested Schema for `google`

Required:

- `client_id` (String) Client ID in Google, must be unique across all identity providers in Codefresh
- `client_secret` (String, Sensitive) Client secret in Google

Optional:

- `admin_email` (String) Email of a user with admin permissions on google, relevant only for synchronization
- `allowed_groups_for_sync` (String) Comma separated list of groups to sync
- `json_keyfile` (String) JSON keyfile for google service account used for synchronization
- `sync_field` (String) Relevant for custom schema-based synchronization only. See Codefresh documentation


<a id="nestedblock--keycloak"></a>
### Nested Schema for `keycloak`

Required:

- `client_id` (String) Client ID from Keycloak
- `client_secret` (String, Sensitive) Client secret from Keycloak
- `host` (String) The Keycloak URL

Optional:

- `realm` (String) The Realm ID for Codefresh in Keycloak. Defaults to master


<a id="nestedblock--ldap"></a>
### Nested Schema for `ldap`

Required:

- `password` (String, Sensitive) The password of the user defined in Distinguished name that will be used to search other users
- `search_base` (String) The search-user scope in LDAP notation
- `url` (String) ldap server url

Optional:

- `allowed_groups_for_sync` (String) To sync only by specified groups - specify a comma separated list of groups, by default all groups will be synced
- `certificate` (String) For ldaps only: The security certificate of the LDAP server. Do not convert to base64 or any other encoding
- `distinguished_name` (String) The username to be used to search other users in LDAP notation (combination of cn, ou,dc)
- `search_base_for_sync` (String) Synchronize using a custom search base, by deafult seach_base is used
- `search_filter` (String) The attribute by which to search for the user on the LDAP server. By default, set to uid. For the Azure LDAP server, set this field to sAMAccountName


<a id="nestedblock--okta"></a>
### Nested Schema for `okta`

Required:

- `client_host` (String) The OKTA organization URL, for example, https://<company>.okta.com
- `client_id` (String) Client ID in Okta, must be unique across all identity providers in Codefresh
- `client_secret` (String, Sensitive) Client secret in Okta

Optional:

- `access_token` (String) The Okta API token generated in Okta, used to sync groups and their users from Okta to Codefresh
- `app_id` (String) The Codefresh application ID in your OKTA organization
- `sync_mirror_accounts` (List of String) The names of the additional Codefresh accounts to be synced from Okta


<a id="nestedblock--onelogin"></a>
### Nested Schema for `onelogin`

Required:

- `client_id` (String) Client ID from Onelogin
- `client_secret` (String, Sensitive) Client secret from Onelogin
- `domain` (String) The domain to be used for authentication

Optional:

- `api_client_id` (String) Client ID for onelogin API, only needed if syncing users and groups from Onelogin
- `api_client_secret` (String) Client secret for onelogin API, only needed if syncing users and groups from Onelogin
- `app_id` (String) The Codefresh application ID in your Onelogin


<a id="nestedblock--saml"></a>
### Nested Schema for `saml`

Required:

- `application_certificate` (String, Sensitive) The security certificate of your Identity Provider. Paste the value directly on the field. Do not convert to base64 or any other encoding by hand
- `endpoint` (String) The SSO endpoint of your Identity Provider

Optional:

- `access_token` (String) Valid for Okta only: The Okta API token generated in Okta, used to sync groups and their users from Okta to Codefresh
- `activate_users_after_sync` (Boolean) Valid for Okta only: If set to true, Codefresh will automatically invite and activate new users added during the automated sync, without waiting for the users to accept the invitations. Defaults to false
- `admin_email` (String) Valid for GSuite only: Email of a user with admin permissions on google, relevant only for synchronization
- `allowed_groups_for_sync` (String) Valid for GSuite only: Comma separated list of groups to sync
- `app_id` (String) Valid for Okta only: The Codefresh application ID in Okta
- `autosync_teams_and_users` (Boolean) Valid for Okta/GSuite: Set to true to sync user accounts and teams in okta/gsuite to your Codefresh account
- `client_host` (String) Valid for Okta only: OKTA organization URL, for example, https://<company>.okta.com
- `json_keyfile` (String) Valid for GSuite only: JSON keyfile for google service account used for synchronization
- `provider` (String) SAML provider. Currently supported values - GSuite, okta or empty string for generic provider. Defaults to empty string
- `sync_interval` (Number) Valid for Okta/GSuite: Sync interval in hours for syncing user accounts in okta/gsuite to your Codefresh account. If not set the sync inteval will be 12 hours

## Import

Please note that secret fields are not imported. 
<br>All secrets should be provided in the configuration and applied after the import for the state to be consistent.

```sh
terraform import codefresh_account_idp.test <id>
```
