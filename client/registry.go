package client

import (
	"fmt"
	"log"
	"net/url"
)

type Registry struct {
	// common
	Id               string `json:"_id,omitempty"`
	Name             string `json:"name,omitempty"`
	Kind             string `json:"kind,omitempty"`
	Default          bool   `json:"default,omitempty"`
	Primary          bool   `json:"primary,omitempty"`
	BehindFirewall   bool   `json:"behindFirewall,omitempty"`
	FallbackRegistry string `json:"fallbackRegistry,omitempty"`
	RepositoryPrefix string `json:"repositoryPrefix,omitempty"`
	Provider         string `json:"provider,omitempty"`

	// mostly all
	Domain   string `json:"domain,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	// bintray
	Token string `json:"token,omitempty"`

	// ecr
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
	Region          string `json:"region,omitempty"`

	// gcr, gar
	Keyfile string `json:"keyfile,omitempty"`

	// acr
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	//DenyCompositeDomain bool   `json:"denyCompositeDomain,omitempty"`
}

func (registry *Registry) GetID() string {
	return registry.Id
}

// GetRegistry identifier is ObjectId or name
func (client *Client) GetRegistry(identifier string) (*Registry, error) {
	fullPath := fmt.Sprintf("/registries/%s", url.PathEscape(identifier))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}
	var respRegistry Registry
	err = DecodeResponseInto(resp, &respRegistry)
	if err != nil {
		return nil, err
	}

	return &respRegistry, nil

}

func (client *Client) CreateRegistry(registry *Registry) (*Registry, error) {

	body, err := EncodeToJSON(registry)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/registries",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		log.Printf("[DEBUG] Call to API for registry creation failed with Error = %v for Body %v", err, body)
		return nil, err
	}

	var respRegistry Registry
	err = DecodeResponseInto(resp, &respRegistry)
	if err != nil {
		return nil, err
	}

	return &respRegistry, nil

}

func (client *Client) UpdateRegistry(registry *Registry) (*Registry, error) {

	body, err := EncodeToJSON(registry)

	if err != nil {
		return nil, err
	}

	fullPath := fmt.Sprintf("/registries/%s", url.PathEscape(registry.Id))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PATCH",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respRegistry Registry
	err = DecodeResponseInto(resp, &respRegistry)
	if err != nil {
		return nil, err
	}

	return &respRegistry, nil

}

func (client *Client) DeleteRegistry(name string) error {

	fullPath := fmt.Sprintf("/registries/%s", url.PathEscape(name))
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
