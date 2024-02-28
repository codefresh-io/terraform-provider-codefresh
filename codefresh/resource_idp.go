package codefresh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/idp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdp() *schema.Resource {
	return &schema.Resource{
		Description: "Codefresh global level identity provider. Requires Codefresh admin token, hence is relevant only for on-prem deployments of Codefresh",
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
		Schema: idp.IdpSchema,
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

	if cfClientIDP.ClientType == idp.GitHub {
		attributes := []map[string]interface{}{{
			"client_id": cfClientIDP.ClientId,
			// Codefresh API Returns the client secret as an encrypted string on the server side
			// hence we need to keep in the state the original secret the user provides along with the encrypted computed secret
			// for Terraform to properly calculate the diff
			"client_secret":      d.Get("github.0.client_secret"),
			"authentication_url": cfClientIDP.AuthURL,
			"token_url":          cfClientIDP.TokenURL,
			"user_profile_url":   cfClientIDP.UserProfileURL,
			"api_host":           cfClientIDP.ApiHost,
			"api_path_prefix":    cfClientIDP.ApiPathPrefix,
		}}

		d.Set("github", attributes)
	}

	if cfClientIDP.ClientType == idp.GitLab {
		attributes := []map[string]interface{}{{
			"client_id":          cfClientIDP.ClientId,
			"client_secret":      d.Get("gitlab.0.client_secret"),
			"authentication_url": cfClientIDP.AuthURL,
			"user_profile_url":   cfClientIDP.UserProfileURL,
			"api_url":            cfClientIDP.ApiURL,
		}}

		d.Set("gitlab", attributes)
	}

	if cfClientIDP.ClientType == idp.Okta {
		attributes := []map[string]interface{}{{
			"client_id":            cfClientIDP.ClientId,
			"client_secret":        d.Get("okta.0.client_secret"),
			"client_host":          cfClientIDP.ClientHost,
			"app_id":               d.Get("okta.0.app_id"),
			"sync_mirror_accounts": cfClientIDP.SyncMirrorAccounts,
			"access_token":         d.Get("okta.0.access_token"),
		}}

		d.Set("okta", attributes)
	}

	if cfClientIDP.ClientType == idp.Google {
		attributes := []map[string]interface{}{{
			"client_id":               cfClientIDP.ClientId,
			"client_secret":           d.Get("google.0.client_secret"),
			"admin_email":             d.Get("google.0.admin_email"),
			"json_keyfile":            d.Get("google.0.json_keyfile"),
			"allowed_groups_for_sync": cfClientIDP.AllowedGroupsForSync,
			"sync_field":              cfClientIDP.SyncField,
		}}

		d.Set("google", attributes)
	}

	if cfClientIDP.ClientType == idp.Auth0 {
		attributes := []map[string]interface{}{{
			"client_id":     cfClientIDP.ClientId,
			"client_secret": d.Get("auth0.0.client_secret"),
			"domain":        cfClientIDP.ClientHost,
		}}

		d.Set("auth0", attributes)
	}

	if cfClientIDP.ClientType == idp.Azure {

		syncInterval, err := strconv.Atoi(cfClientIDP.SyncInterval)
		if err != nil {
			return err
		}

		attributes := []map[string]interface{}{{
			"app_id":                   cfClientIDP.ClientId,
			"client_secret":            d.Get("azure.0.client_secret"),
			"object_id":                cfClientIDP.AppId,
			"autosync_teams_and_users": cfClientIDP.AutoGroupSync,
			"sync_interval":            syncInterval,
			"tenant":                   cfClientIDP.Tenant,
		}}

		d.Set("azure", attributes)
	}

	if cfClientIDP.ClientType == idp.OneLogin {
		attributes := []map[string]interface{}{{
			"client_id":     cfClientIDP.ClientId,
			"client_secret": d.Get("onelogin.0.client_secret"),
			"domain":        cfClientIDP.ClientHost,
			"api_client_id": cfClientIDP.ApiClientId,

			"api_client_secret": cfClientIDP.ApiClientSecret,
			"app_id":            cfClientIDP.AppId,
		}}

		d.Set("onelogin", attributes)
	}

	if cfClientIDP.ClientType == idp.Keycloak {
		attributes := []map[string]interface{}{{
			"client_id":     cfClientIDP.ClientId,
			"client_secret": d.Get("keycloak.0.client_secret"),
			"host":          cfClientIDP.Host,
			"realm":         cfClientIDP.Realm,
		}}

		d.Set("keycloak", attributes)
	}

	if cfClientIDP.ClientType == idp.SAML {
		syncInterval, err := strconv.Atoi(cfClientIDP.SyncInterval)
		if err != nil {
			return err
		}
		attributes := []map[string]interface{}{{
			"endpoint":                  cfClientIDP.EntryPoint,
			"application_certificate":   d.Get("saml.0.application_certificate"),
			"provider":                  cfClientIDP.SamlProvider,
			"allowed_groups_for_sync":   cfClientIDP.AllowedGroupsForSync,
			"autosync_teams_and_users":  cfClientIDP.AutoGroupSync,
			"activate_users_after_sync": cfClientIDP.ActivateUserAfterSync,
			"sync_interval":             syncInterval,
			"app_id":                    cfClientIDP.AppId,
			"client_host":               cfClientIDP.ClientHost,
			"json_keyfile":              d.Get("saml.0.json_keyfile"),
			"admin_email":               d.Get("saml.0.admin_email"),
			"access_token":              d.Get("saml.0.access_token"),
		}}

		d.Set("saml", attributes)
	}

	if cfClientIDP.ClientType == idp.LDAP {
		attributes := []map[string]interface{}{{
			"url":                     cfClientIDP.Url,
			"password":                d.Get("ldap.0.password"),
			"distinguished_name":      cfClientIDP.DistinguishedName,
			"search_base":             cfClientIDP.SearchBase,
			"search_filter":           cfClientIDP.SearchFilter,
			"certificate":             d.Get("ldap.0.certificate"),
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

	if _, ok := d.GetOk(idp.GitHub); ok {
		cfClientIDP.ClientType = idp.GitHub
		cfClientIDP.ClientId = d.Get("github.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("github.0.client_secret").(string)
		cfClientIDP.AuthURL = d.Get("github.0.authentication_url").(string)
		cfClientIDP.TokenURL = d.Get("github.0.token_url").(string)
		cfClientIDP.UserProfileURL = d.Get("github.0.user_profile_url").(string)
		cfClientIDP.ApiHost = d.Get("github.0.api_host").(string)
		cfClientIDP.ApiPathPrefix = d.Get("github.0.api_path_prefix").(string)
	}

	if _, ok := d.GetOk(idp.GitLab); ok {
		cfClientIDP.ClientType = idp.GitLab
		cfClientIDP.ClientId = d.Get("gitlab.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("gitlab.0.client_secret").(string)
		cfClientIDP.AuthURL = d.Get("gitlab.0.authentication_url").(string)
		cfClientIDP.UserProfileURL = d.Get("gitlab.0.user_profile_url").(string)
		cfClientIDP.ApiURL = d.Get("gitlab.0.api_url").(string)
	}

	if _, ok := d.GetOk(idp.Okta); ok {
		cfClientIDP.ClientType = idp.Okta
		cfClientIDP.ClientId = d.Get("okta.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("okta.0.client_secret").(string)
		cfClientIDP.ClientHost = d.Get("okta.0.client_host").(string)
		cfClientIDP.AppId = d.Get("okta.0.app_id").(string)
		cfClientIDP.SyncMirrorAccounts = datautil.ConvertStringArr(d.Get("okta.0.sync_mirror_accounts").([]interface{}))
		cfClientIDP.Access_token = d.Get("okta.0.access_token").(string)
	}

	if _, ok := d.GetOk(idp.Google); ok {
		cfClientIDP.ClientType = idp.Google
		cfClientIDP.ClientId = d.Get("google.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("google.0.client_secret").(string)
		cfClientIDP.KeyFile = d.Get("google.0.json_keyfile").(string)
		cfClientIDP.Subject = d.Get("google.0.admin_email").(string)
		cfClientIDP.AllowedGroupsForSync = d.Get("google.0.allowed_groups_for_sync").(string)
		cfClientIDP.SyncField = d.Get("google.0.sync_field").(string)
	}

	if _, ok := d.GetOk(idp.Auth0); ok {
		cfClientIDP.ClientType = idp.Auth0
		cfClientIDP.ClientId = d.Get("auth0.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("auth0.0.client_secret").(string)
		cfClientIDP.ClientHost = d.Get("auth0.0.domain").(string)
	}

	if _, ok := d.GetOk(idp.Azure); ok {
		cfClientIDP.ClientType = idp.Azure
		cfClientIDP.ClientId = d.Get("azure.0.app_id").(string)
		cfClientIDP.ClientSecret = d.Get("azure.0.client_secret").(string)
		cfClientIDP.AppId = d.Get("azure.0.object_id").(string)
		cfClientIDP.Tenant = d.Get("azure.0.tenant").(string)
		cfClientIDP.AutoGroupSync = d.Get("azure.0.autosync_teams_and_users").(bool)
		cfClientIDP.SyncInterval = strconv.Itoa(d.Get("azure.0.sync_interval").(int))
	}

	if _, ok := d.GetOk(idp.OneLogin); ok {
		cfClientIDP.ClientType = idp.OneLogin
		cfClientIDP.ClientId = d.Get("onelogin.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("onelogin.0.client_secret").(string)
		cfClientIDP.ClientHost = d.Get("onelogin.0.domain").(string)
		cfClientIDP.AppId = d.Get("onelogin.0.app_id").(string)
		cfClientIDP.ApiClientId = d.Get("onelogin.0.api_client_id").(string)
		cfClientIDP.ApiClientSecret = d.Get("onelogin.0.api_client_secret").(string)
	}

	if _, ok := d.GetOk(idp.Keycloak); ok {
		cfClientIDP.ClientType = idp.Keycloak
		cfClientIDP.ClientId = d.Get("keycloak.0.client_id").(string)
		cfClientIDP.ClientSecret = d.Get("keycloak.0.client_secret").(string)
		cfClientIDP.Host = d.Get("keycloak.0.host").(string)
		cfClientIDP.Realm = d.Get("keycloak.0.realm").(string)
	}

	if _, ok := d.GetOk(idp.SAML); ok {
		cfClientIDP.ClientType = idp.SAML
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

	if _, ok := d.GetOk(idp.LDAP); ok {
		cfClientIDP.ClientType = idp.LDAP
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
