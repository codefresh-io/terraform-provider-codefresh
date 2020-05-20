package codefresh

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
	"strings"
)

type codefreshObject interface {
	getID() string
}

type variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type errorResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func convertStringArr(ifaceArr []interface{}) []string {
	return convertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

func convertVariables(vars []variable) map[string]string {
	res := make(map[string]string, len(vars))
	for _, v := range vars {
		res[v.Key] = v.Value
	}
	return res
}

func createCodefreshObject(
	c *Config,
	path string,
	httpMethod string,
	d *schema.ResourceData,
	mapFunc func(data *schema.ResourceData) codefreshObject,
	readFunc func(*schema.ResourceData, []byte) (codefreshObject, error),
) error {
	body, err := getBody(d, mapFunc)
	if err != nil {
		return err
	}

	url := c.APIServer + path
	apiKey := c.Token

	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode == 500 {
		errorResp, err := readErrorResponse(bodyBytes)
		if err != nil {
			return fmt.Errorf("%v (%v)", resp.Status, resp.StatusCode)
		}
		return fmt.Errorf("%v (%v) %v", errorResp.Message, errorResp.Status, errorResp.Error)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%v, %s", resp.Status, string(bodyBytes))
	}

	cfObject, err := readFunc(d, bodyBytes)
	if err != nil {
		return err
	}

	d.SetId(cfObject.getID())
	return nil
}

func readErrorResponse(b []byte) (*errorResponse, error) {
	errorResponse := &errorResponse{}
	err := json.Unmarshal(b, errorResponse)
	if err != nil {
		return nil, err
	}
	return errorResponse, nil
}

func getBody(d *schema.ResourceData, mapFunc func(data *schema.ResourceData) codefreshObject) (*strings.Reader, error) {
	cfObject := mapFunc(d)
	postBytes, err := json.Marshal(cfObject)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(postBytes)), nil
}

func readCodefreshObject(
	d *schema.ResourceData,
	c *Config,
	getFromSource func(d *schema.ResourceData, c *Config) (codefreshObject, error),
	mapToResource func(codefreshObject, *schema.ResourceData) error) error {
	cfObject, err := getFromSource(d, c)
	if err != nil {
		return err
	}

	if cfObject == nil {
		// if object was not found, clear ID, this signals to terraform that the resource does not exist in destination
		d.SetId("")
		return nil
	}

	err = mapToResource(cfObject, d)
	if err != nil {
		return err
	}

	return nil
}

func getFromCodefresh(
	d *schema.ResourceData,
	c *Config,
	path string,
	readFunc func(*schema.ResourceData, []byte) (codefreshObject, error),
) (codefreshObject, error) {

	url := c.APIServer + path
	apiKey := c.Token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// if not found, this is not an error, we just need to return nil to signal the thing was not found
	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	return readFunc(d, bodyBytes)
}

func updateCodefreshObject(
	d *schema.ResourceData,
	c *Config,
	path string,
	httpMethod string,
	mapFunc func(data *schema.ResourceData) codefreshObject,
	readFunc func(*schema.ResourceData, []byte) (codefreshObject, error),
	resourceRead func(d *schema.ResourceData, m interface{}) error,
) error {
	body, err := getBody(d, mapFunc)
	if err != nil {
		return err
	}

	url := c.APIServer + path
	apiKey := c.Token

	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	_, err = readFunc(d, bodyBytes)
	if err != nil {
		return err
	}

	return resourceRead(d, c)
}

func deleteCodefreshObject(c *Config, path string) error {
	url := c.APIServer + path
	apiKey := c.Token

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", apiKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// it seems like Codefresh should return a 400 error when a pipeline is not found, but it's actually returning
	// a 500, so we won't be able to tell the difference between a not found and an issue, we'll just have to
	// assume not found
	if resp.StatusCode == 500 {
		return nil
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

func importCodefreshObject(
	d *schema.ResourceData,
	c *Config,
	getFromCodefresh func(d *schema.ResourceData, c *Config) (codefreshObject, error),
	mapToResource func(codefreshObject, *schema.ResourceData) error,
) ([]*schema.ResourceData, error) {
	cfObject, err := getFromCodefresh(d, c)
	if err != nil {
		return nil, err
	}

	if cfObject == nil {
		// not found
		d.SetId("")
		return nil, errors.New("Failed to find object with id " + d.Id())
	}

	err = mapToResource(cfObject, d)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
