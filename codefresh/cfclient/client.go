package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/envutil"
)

// Client token, host, htpp.Client
type Client struct {
	Token        string
	TokenHeader  string
	Host         string
	HostV2       string
	featureFlags map[string]bool
	Client       *http.Client
}

// RequestOptions  path, method, etc
type RequestOptions struct {
	Path         string
	Method       string
	Body         []byte
	QS           map[string]string
	XAccessToken string
}

// HttpClient returns a client which can be configured to communicate on a server with custom timeout settings
func NewHttpClient(hostname string, hostnameV2 string, token string, tokenHeader string) *Client {
	if tokenHeader == "" {
		tokenHeader = "Authorization"
	}

	// Configurable HTTP transport with proper connection pooling and timeouts to prevent "connection reset by peer" errors,
	// default values are equivalent to default &http.Client{} settings
	transport := &http.Transport{
		// Limit maximum idle connections per host to prevent connection exhaustion
		MaxIdleConnsPerHost: envutil.GetEnvAsInt("CF_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST", 2),
		// Limit total idle connections
		MaxIdleConns: envutil.GetEnvAsInt("CF_HTTP_MAX_IDLE_CONNECTIONS", 100),
		// Close idle connections after specified seconds to prevent server-side timeouts
		IdleConnTimeout: time.Duration(envutil.GetEnvAsInt("CF_HTTP_IDLE_CONNECTION_TIMEOUT", 90)) * time.Second,
		// Timeout for TLS handshake in seconds
		TLSHandshakeTimeout: time.Duration(envutil.GetEnvAsInt("CF_HTTP_TLS_HANDSHAKE_TIMEOUT", 10)) * time.Second,
		// Timeout for expecting response headers in seconds, 0 - no limits
		ResponseHeaderTimeout: time.Duration(envutil.GetEnvAsInt("CF_HTTP_RESPONSE_HEADER_TIMEOUT", 0)) * time.Second,
		// Disable connection reuse for more stable connections
		DisableKeepAlives: envutil.GetEnvAsBool("CF_HTTP_DISABLE_KEEPALIVES", false),
	}

	httpClient := &http.Client{
		Transport: transport,
		// Overall request timeout in seconds, 0 - no limits
		Timeout: time.Duration(envutil.GetEnvAsInt("CF_HTTP_GLOBAL_TIMEOUT", 0)) * time.Second,
	}

	return &Client{
		Host:         hostname,
		HostV2:       hostnameV2,
		Token:        token,
		TokenHeader:  tokenHeader,
		Client:       httpClient,
		featureFlags: map[string]bool{},
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

	body, err := io.ReadAll(resp.Body)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v, %s", resp.Status, string(body))
	}
	return body, nil
}

func (client *Client) isFeatureFlagEnabled(flagName string) (bool, error) {

	if len(client.featureFlags) == 0 {
		currAcc, err := client.GetCurrentAccount()

		if err != nil {
			return false, err
		}

		client.featureFlags = currAcc.FeatureFlags
	}

	if val, ok := client.featureFlags[flagName]; ok {
		return val, nil
	}

	return false, nil
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
