package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client token, host, htpp.Client
type Client struct {
	Token       string
	TokenHeader string
	Host        string
	HostV2      string
	Client      *http.Client
}

// RequestOptions  path, method, etc
type RequestOptions struct {
	Path         string
	Method       string
	Body         []byte
	QS           map[string]string
	XAccessToken string
}

// NewClient returns a new client configured to communicate on a server with the
// given hostname and to send an Authorization Header with the value of
// token
func NewClient(hostname string, hostnameV2 string, token string, tokenHeader string) *Client {
	if tokenHeader == "" {
		tokenHeader = "Authorization"
	}
	return &Client{
		Host:        hostname,
		HostV2:      hostnameV2,
		Token:       token,
		TokenHeader: tokenHeader,
		Client:      &http.Client{},
	}

}

// RequestAPI http request to Codefresh API
func (client *Client) RequestAPI(opt *RequestOptions) ([]byte, error) {
	finalURL := fmt.Sprintf("%s%s", client.Host, opt.Path)
	if opt.QS != nil {
		finalURL += ToQS(opt.QS)
	}
	request, err := http.NewRequest(opt.Method, finalURL, bytes.NewBuffer(opt.Body))
	if err != nil {
		return nil, err
	}

	tokenHeader := client.TokenHeader
	if tokenHeader == "" {
		tokenHeader = "Authorization"
	}
	request.Header.Set(tokenHeader, client.Token)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Client.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	// todo: maybe other 2**?
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("%v, %s", resp.Status, string(body))
	}
	return body, nil
}

func (client *Client) RequestApiXAccessToken(opt *RequestOptions) ([]byte, error) {
	finalURL := fmt.Sprintf("%s%s", client.Host, opt.Path)
	if opt.QS != nil {
		finalURL += ToQS(opt.QS)
	}
	request, err := http.NewRequest(opt.Method, finalURL, bytes.NewBuffer(opt.Body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("x-access-token", opt.XAccessToken)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Client.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v, %s", resp.Status, string(body))
	}
	return body, nil
}

// ToQS add extra parameters to path
func ToQS(qs map[string]string) string {
	var arr = []string{}
	for k, v := range qs {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	return "?" + strings.Join(arr, "&")
}

// DecodeResponseInto json Unmarshall
func DecodeResponseInto(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

// EncodeToJSON json Marshal
func EncodeToJSON(object interface{}) ([]byte, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return body, nil
}
