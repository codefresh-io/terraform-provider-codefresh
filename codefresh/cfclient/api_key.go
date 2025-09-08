package cfclient

import (
	"errors"
	"fmt"
	"log"
)

type ApiKeySubject struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"ref,omitempty"`
}

type ApiKeyScopeSnapshot struct {
	Scopes []string `json:"scopes,omitempty"`
	ID     string   `json:"_id,omitempty"`
	Date   string   `json:"date,omitempty"`
	V      int      `json:"__v,omitempty"`
}

type ApiKey struct {
	Subject       ApiKeySubject       `json:"subject,omitempty"`
	ID            string              `json:"_id,omitempty"`
	Name          string              `json:"name"`
	Scopes        []string            `json:"scopes,omitempty"`
	TokenPrefix   string              `json:"tokenPrefix,omitempty"`
	ScopeSnapshot ApiKeyScopeSnapshot `json:"scopeSnapshot,omitempty"`
	Created       string              `json:"created,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	User        struct {
		UserName string `json:"userName,omitempty"`
		Email    string `json:"email,omitempty"`
	} `json:"user"`
}

func (client *Client) GetAPIKey(userID string, accountId string, keyID string) (*ApiKey, error) {

	xAccessToken, err := client.GetXAccessToken(userID, accountId)

	if err != nil {
		return nil, err
	}

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/%s", keyID),
		XAccessToken: xAccessToken,
		Method: "GET",
	}

	resp, err := client.RequestApiXAccessToken(&opts)

	if err != nil {
		return nil, err
	}

	var apiKey ApiKey

	err = DecodeResponseInto(resp, &apiKey)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (client *Client) DeleteAPIKey(userID string, accountId string, keyID string) error {
	// login as user

	xAccessToken, err := client.GetXAccessToken(userID, accountId)

	if err != nil {
		return err
	}
	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/%s", keyID),
		Method: "DELETE",
		XAccessToken: xAccessToken,
	}

	resp, err := client.RequestApiXAccessToken(&opts)
	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	return nil
}

func (client *Client) UpdateAPIKey(userID string, accountId string,key *ApiKey) error {

	keyID := key.ID
	if keyID == "" {
		return errors.New("[ERROR] Key ID is empty")
	}

	body, err := EncodeToJSON(key)
	if err != nil {
		return err
	}

	var xAccessToken string

	// login as user
	xAccessToken, err = client.GetXAccessToken(userID, accountId)

	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/%s", keyID),
		Method: "PATCH",
		XAccessToken: xAccessToken,
		Body:   body,
	}

	resp, err := client.RequestApiXAccessToken(&opts)

	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	return nil
}

// CreateApiKey - creates api key for account by switch to the user and call /api/auth/keys
func (client *Client) CreateApiKey(userID string, accountId string, apiKey *ApiKey) (string, error) {

	// Check collaborataros
	account, err := client.GetAccountByID(accountId)

	if err != nil {
		return "", err
	}
	if account.Limits == nil {
		log.Fatal("[ERROR] Collaborators are not set")
	}

	var xAccessToken string

	// login as user
	xAccessToken, err = client.GetXAccessToken(userID, accountId)
	if err != nil {
		return "", err
	}

	// generate token
	apiToken, err := client.GenerateToken(xAccessToken, apiKey)
	if err != nil {
		return "", err
	}

	return apiToken, nil
}

// GetXAccessToken
func (client *Client) GetXAccessToken(userID string, accountId string) (string, error) {

	fullPath := fmt.Sprintf("/admin/user/loginAsUser?userId=%s", userID)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return "", err
	}

	var userCfAccessToken string
	var asUserTokenResponse TokenResponse

	err = DecodeResponseInto(resp, &asUserTokenResponse)
	if err != nil {
		return "", err
	}

	userCfAccessToken = asUserTokenResponse.AccessToken

	if userCfAccessToken == "" {
		return "", fmt.Errorf("Failed to GetXAccessToken for userId = %s", userID)
	}

	// change account
	fullPath = fmt.Sprintf("/user/changeaccount/%s", accountId)
	opts = RequestOptions{
		Path:         fullPath,
		Method:       "POST",
		XAccessToken: userCfAccessToken,
	}

	resp, err = client.RequestApiXAccessToken(&opts)

	if err != nil {
		return "", err
	}

	var accCfAccessToken string
	var changeAccountTokenResponse TokenResponse

	err = DecodeResponseInto(resp, &changeAccountTokenResponse)
	if err != nil {
		return "", err
	}

	accCfAccessToken = changeAccountTokenResponse.AccessToken

	if accCfAccessToken == "" {
		return "", fmt.Errorf("Failed to GetXAccessToken for userId = %s after ChangeAcocunt to %s", userID, accountId)
	}

	return accCfAccessToken, nil
}

func (client *Client) GenerateToken(xToken string, apiKey *ApiKey) (string, error) {

	body, err := EncodeToJSON(apiKey)
	if err != nil {
		return "", err
	}

	opts := RequestOptions{
		Path:         "/auth/key",
		Method:       "POST",
		XAccessToken: xToken,
		Body:         body,
	}

	resp, err := client.RequestApiXAccessToken(&opts)

	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func (client *Client) GetApiKeysList() ([]ApiKey, error) {
	fullPath := "/auth/keys"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var apiKeys []ApiKey

	err = DecodeResponseInto(resp, &apiKeys)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (client *Client) GetAPIKeyServiceUser(keyID string, serviceUserId string) (*ApiKey, error) {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/service-user/%s/%s", serviceUserId, keyID),
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var apiKey ApiKey

	err = DecodeResponseInto(resp, &apiKey)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (client *Client) DeleteAPIKeyServiceUser(keyID string, serviceUserId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/service-user/%s/%s", serviceUserId, keyID),
		Method: "DELETE",
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	return nil
}

func (client *Client) UpdateAPIKeyServiceUser(key *ApiKey, serviceUserId string) error {

	keyID := key.ID
	if keyID == "" {
		return errors.New("[ERROR] Key ID is empty")
	}

	body, err := EncodeToJSON(key)
	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/service-user/%s/%s", serviceUserId, keyID),
		Method: "PATCH",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	return nil
}

func (client *Client) CreateApiKeyServiceUser(serviceUserId string, apiKey *ApiKey) (string, error) {

	body, err := EncodeToJSON(apiKey)
	if err != nil {
		return "", err
	}

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/service-user/%s", serviceUserId),
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return "", err
	}

	return string(resp), nil
}
