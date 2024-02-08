package cfclient

import (
	"errors"
	"fmt"
	"log"
	"net/url"
)

type IDP struct {
	ID            string   `json:"_id,omitempty"`
	IsGlobal	  bool     `json:"isGlobal,omitempty"` // This is not part of the schema, rather it is used to determine if the IDP is gloal or not and choose endpoints accordingly
	Access_token  string   `json:"access_token,omitempty"`
	Accounts      []string `json:"accounts,omitempty"`
	ClientName    string   `json:"clientName,omitempty"` // IDP name
	ClientType    string   `json:"clientType,omitempty"` // IDP type
	DisplayName   string   `json:"displayName,omitempty"`
	LoginUrl      string   `json:"loginUrl,omitempty"`      // Login url in Codefresh
	RedirectUiUrl string   `json:"redirectUiUrl,omitempty"` // Redicrect url Codefresh UI
	RedirectUrl   string   `json:"redirectUrl,omitempty"`
	ClientId      string   `json:"clientId,omitempty"`      // All providers (base)
	ClientSecret  string   `json:"clientSecret,omitempty"`  // All providers (base)
	ApiHost       string   `json:"apiHost,omitempty"`       // GitHub
	ApiPathPrefix string   `json:"apiPathPrefix,omitempty"` // Github
	// Bitbucket, Gitlab
	ApiURL string `json:"apiURL,omitempty"`
	// Azure, Okta, onelogin,saml
	AppId string `json:"appId,omitempty"`
	// Github, Gitlab
	AuthURL string `json:"authURL,omitempty"`
	// saml, okta, onelogin, auth0, azure, google, google-cloud-sr
	ClientHost string `json:"clientHost,omitempty"`
	// Azure
	CookieIv string `json:"cookieIv,omitempty"`
	// Azure
	CookieKey string `json:"cookieKey,omitempty"`
	// Azure
	IDPLoginUrl string `json:"IDPLoginUrl,omitempty"`
	// Bitbucket
	RefreshTokenURL string `json:"refreshTokenURL,omitempty"`
	// Multiple - computed
	Scopes []string `json:"scopes,omitempty"`
	// Azure
	Tenant      string `json:"tenant,omitempty"`
	TokenSecret string `json:"tokenSecret,omitempty"`
	// Okta, Bitbucket, GitHub, Keycloak
	TokenURL string `json:"tokenURL,omitempty"`
	// Github, Gitlab
	UserProfileURL string `json:"userProfileURL,omitempty"`
	// Okta
	SyncMirrorAccounts []string `json:"syncMirrorAccounts,omitempty"`
	// Google, Ldap
	AllowedGroupsForSync string `json:"allowedGroupsForSync,omitempty"`
	// Google
	Subject string `json:"subject,omitempty"`
	// Google
	KeyFile string `json:"keyfile,omitempty"`
	// Google
	SyncField string `json:"syncField,omitempty"`
}

// Return the appropriate API endpoint for platform and account scoped IDPs
func getAPIEndpoint(isGlobal bool) string {
	// If IDP is platform scoped
	if isGlobal {
		return "/admin/idp"
	} else {
		return "/idp/account"
	}
}

func (client *Client) CreateIDP(idp *IDP) (*IDP, error) {

	body, err := EncodeToJSON(idp)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   getAPIEndpoint(idp.IsGlobal),
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		log.Printf("[DEBUG] Call to API for IDP creation failed with Error = %v for Body %v", err, body)
		return nil, err
	}

	var respIDP IDP
	err = DecodeResponseInto(resp, &respIDP)
	if err != nil {
		return nil, err
	}

	return &respIDP, nil
}

// Currently on update the API returns a different structure for accounts than on read making the client crash on decode
// For now we are disabling response decode and in the resource will instead call the read function again
func (client *Client) UpdateIDP(idp *IDP) error {

	body, err := EncodeToJSON(idp)

	if err != nil {
		return err
	}
	opts := RequestOptions{
		Path:   getAPIEndpoint(idp.IsGlobal),
		Method: "PUT",
		Body:   body,
	}

	_, err = client.RequestAPI(&opts)

	if err != nil {
		log.Printf("[DEBUG] Call to API for IDP update failed with Error = %v for Body %v", err, body)
		return err
	}

	// var respIDP IDP
	// err = DecodeResponseInto(resp, &respIDP)
	// if err != nil {
	// 	return nil, err
	// }

	return nil
}

func (client *Client) DeleteIDP(id string) error {
	baseUrl := getAPIEndpoint(true)
	fullPath := fmt.Sprintf("%s/%s", baseUrl, url.PathEscape(id))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteIDPAccount(id string) error {

	body, err := EncodeToJSON(map[string]interface{}{"id": id})

	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   getAPIEndpoint(false),
		Method: "DELETE",
		Body: body,
	}

	_, err = client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}



// get all idps
func (client *Client) GetIDPs() (*[]IDP, error) {
	fullPath := "/admin/idp"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var idps []IDP

	err = DecodeResponseInto(resp, &idps)
	if err != nil {
		return nil, err
	}

	return &idps, nil
}

// get idp id by idp name
func (client *Client) GetIdpByName(idpName string) (*IDP, error) {

	idpList, err := client.GetIDPs()
	if err != nil {
		return nil, err
	}

	for _, idp := range *idpList {
		if idp.ClientName == idpName {
			return &idp, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("[ERROR] IDP with name %s isn't found.", idpName))
}

func (client *Client) GetIdpByID(idpID string) (*IDP, error) {

	idpList, err := client.GetIDPs()
	if err != nil {
		return nil, err
	}

	for _, idp := range *idpList {
		if idp.ID == idpID {
			return &idp, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("[ERROR] IDP with ID %s isn't found.", idpID))
}

// get account idps
func (client *Client) GetAccountIDPs() (*[]IDP, error) {
	fullPath := "/idp/account"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var idps []IDP

	err = DecodeResponseInto(resp, &idps)
	if err != nil {
		return nil, err
	}

	return &idps, nil
}

func (client *Client) GetAccountIdpByID(idpID string) (*IDP, error) {

	idpList, err := client.GetAccountIDPs()
	if err != nil {
		return nil, err
	}

	for _, idp := range *idpList {
		if idp.ID == idpID {
			return &idp, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("[ERROR] IDP with ID %s isn't found.", idpID))
}

// add account to idp
func (client *Client) AddAccountToIDP(accountId, idpId string) error {

	body := fmt.Sprintf(`{"accountId":"%s","IDPConfigId":"%s"}`, accountId, idpId)

	opts := RequestOptions{
		Path:   "/admin/idp/addAccount",
		Method: "POST",
		Body:   []byte(body),
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

// remove account form idp
// doesn't implemente
