package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// GraphQLRequest GraphQL query
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLClient GraphQL client
type GraphQLClient struct {
	Token       string
	TokenHeader string
	Host        string
	Client      *http.Client
}

// NewGqlClient returns a new graphql client configured to communicate on a server with the
// given hostname and to send an Authorization Header with the value of token
func NewGqlClient(url, defaultUrl, apiKey string) *GraphQLClient {
	tokenHeader := "Authorization"
	hostname := url
	if hostname == "" {
		hostname = defaultUrl
	}
	token := apiKey
	return &GraphQLClient{
		Host:        hostname,
		Token:       token,
		TokenHeader: tokenHeader,
		Client:      &http.Client{},
	}
}

func (client *GraphQLClient) SendRequest(request GraphQLRequest) ([]byte, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", client.Host, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set(client.TokenHeader, client.Token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(resp.Status + " " + string(bodyBytes))
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeGraphQLResponseInto(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}
