package client

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"log"
	"net/http"
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

func (client *Client) GetAPIKey(keyID string) (*ApiKey, error) {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/%s", keyID),
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

func (client *Client) DeleteAPIKey(keyID string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/%s", keyID),
		Method: "DELETE",
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		fmt.Println(string(resp))
		return err
	}

	return nil
}

func (client *Client) UpdateAPIKey(key *ApiKey) error {

	keyID := key.ID
	if keyID == "" {
		return errors.New("[ERROR] Key ID is empty")
	}

	body, err := EncodeToJSON(key)
	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   fmt.Sprintf("/auth/key/%s", keyID),
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

func (client *Client) CreateApiKey(accountId string, apiKey *ApiKey) (string, error) {

	// Check collaborataros
	account, err := client.GetAccountByID(accountId)
	if err != nil {
		return "", err
	}
	if account.Limits == nil {
		log.Fatal("[ERROR] Collaborators are not set")
	}

	// add user
	userPrefix := acctest.RandString(10)
	userName := "tfuser" + userPrefix
	userEmail := userName + "@codefresh.io"

	user, err := client.AddNewUserToAccount(accountId, userName, userEmail)
	if err != nil {
		return "", err
	}
	userID := user.ID

	// activate
	_, err = client.ActivateUser(userID)
	if err != nil {
		return "", err
	}

	// set user as account admin
	err = client.SetUserAsAccountAdmin(accountId, userID)
	if err != nil {
		return "", nil
	}

	// login as user
	xAccessToken, err := client.GetXAccessToken(userID)
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

func (client *Client) GetXAccessToken(userID string) (string, error) {

	url := fmt.Sprintf("%s/admin/user/loginAsUser?userId=%s", client.Host, userID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", client.Token)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "cf-access-token" {
			return cookie.Value, nil
		}
	}

	return "", nil
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
