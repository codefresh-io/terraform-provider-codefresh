package client

import (
	"fmt"
	"log"
	"net/url"
)

type ContextErrorResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ContextMetadata struct {
	Name string `json:"name,omitempty"`
}

type Context struct {
	Metadata ContextMetadata `json:"metadata,omitempty"`
	Spec     ContextSpec     `json:"spec,omitempty"`
	Version  string          `json:"version,omitempty"`
}

type ContextSpec struct {
	Type string                 `json:"type,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func (context *Context) GetID() string {
	return context.Metadata.Name
}

func (client *Client) GetContext(name string) (*Context, error) {
	fullPath := fmt.Sprintf("/contexts/%s?decrypt=true", url.PathEscape(name))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}
	var respContext Context
	err = DecodeResponseInto(resp, &respContext)
	if err != nil {
		return nil, err
	}

	return &respContext, nil

}

func (client *Client) CreateContext(context *Context) (*Context, error) {

	body, err := EncodeToJSON(context)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/contexts",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		log.Printf("[DEBUG] Call to API for context creation failed with Error = %v for Body %v", err, body)
		return nil, err
	}

	var respContext Context
	err = DecodeResponseInto(resp, &respContext)
	if err != nil {
		return nil, err
	}

	return &respContext, nil

}

func (client *Client) UpdateContext(context *Context) (*Context, error) {

	body, err := EncodeToJSON(context)

	if err != nil {
		return nil, err
	}

	fullPath := fmt.Sprintf("/contexts/%s", url.PathEscape(context.Metadata.Name))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respContext Context
	err = DecodeResponseInto(resp, &respContext)
	if err != nil {
		return nil, err
	}

	return &respContext, nil

}

func (client *Client) DeleteContext(name string) error {

	fullPath := fmt.Sprintf("/contexts/%s", url.PathEscape(name))
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
