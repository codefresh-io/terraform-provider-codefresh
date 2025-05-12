package cfclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// GraphQLRequest GraphQL query
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func (client *Client) SendGqlRequest(request GraphQLRequest) ([]byte, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", client.HostV2, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, err
	}

	tokenHeader := client.TokenHeader
	if tokenHeader == "" {
		tokenHeader = "Authorization"
	}
	req.Header.Set(tokenHeader, client.Token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
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
