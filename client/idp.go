package client

import (
	"errors"
	"fmt"
)

type IDP struct {
	Access_token    string   `json:"access_token,omitempty"`
	Accounts        []string `json:"accounts,omitempty"`
	ApiHost         string   `json:"apiHost,omitempty"`
	ApiPathPrefix   string   `json:"apiPathPrefix,omitempty"`
	ApiURL          string   `json:"apiURL,omitempty"`
	AppId           string   `json:"appId,omitempty"`
	AuthURL         string   `json:"authURL,omitempty"`
	ClientHost      string   `json:"clientHost,omitempty"`
	ClientId        string   `json:"clientId,omitempty"`
	ClientName      string   `json:"clientName,omitempty"`
	ClientSecret    string   `json:"clientSecret,omitempty"`
	ClientType      string   `json:"clientType,omitempty"`
	CookieIv        string   `json:"cookieIv,omitempty"`
	CookieKey       string   `json:"cookieKey,omitempty"`
	DisplayName     string   `json:"displayName,omitempty"`
	ID              string   `json:"id,omitempty"`
	IDPLoginUrl     string   `json:"IDPLoginUrl,omitempty"`
	LoginUrl        string   `json:"loginUrl,omitempty"`
	RedirectUiUrl   string   `json:"redirectUiUrl,omitempty"`
	RedirectUrl     string   `json:"redirectUrl,omitempty"`
	RefreshTokenURL string   `json:"refreshTokenURL,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
	Tenant          string   `json:"tenant,omitempty"`
	TokenSecret     string   `json:"tokenSecret,omitempty"`
	TokenURL        string   `json:"tokenURL,omitempty"`
	UserProfileURL  string   `json:"userProfileURL,omitempty"`
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
