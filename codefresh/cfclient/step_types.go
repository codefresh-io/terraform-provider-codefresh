package cfclient

import (
	"fmt"
	"log"
	"net/url"

	"github.com/iancoleman/orderedmap"
)

type StepTypesVersions struct {
	Name     string
	Versions []StepTypesVersion
}
type StepTypesVersion struct {
	VersionNumber string
	StepTypes     StepTypes
}

type StepTypes struct {
	Version  string                 `json:"version,omitempty"`
	Kind     string                 `json:"kind,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Spec     SpecStepTypes          `json:"spec,omitempty"`
}

type SpecStepTypes struct {
	Arguments     string                 `json:"arguments,omitempty"`
	Delimiters    map[string]interface{} `json:"delimiters,omitempty"`
	Returns       string                 `json:"returns,omitempty"`
	Steps         *orderedmap.OrderedMap `json:"steps,omitempty"`
	StepsTemplate string                 `json:"stepsTemplate,omitempty"`
}

func (stepTypes *StepTypes) GetID() string {
	return stepTypes.Metadata["name"].(string)
}

func (client *Client) GetStepTypesVersions(name string) ([]string, error) {
	fullPath := fmt.Sprintf("/step-types/%s/versions", url.PathEscape(name))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}
	var respStepTypesVersions []string
	err = DecodeResponseInto(resp, &respStepTypesVersions)
	if err != nil {
		return nil, err
	}
	return respStepTypesVersions, nil
}

func (client *Client) GetStepTypes(identifier string) (*StepTypes, error) {
	fullPath := fmt.Sprintf("/step-types/%s", url.PathEscape(identifier))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}
	var respStepTypes StepTypes
	err = DecodeResponseInto(resp, &respStepTypes)
	if err != nil {
		return nil, err
	}

	return &respStepTypes, nil

}

func (client *Client) CreateStepTypes(stepTypes *StepTypes) (*StepTypes, error) {

	body, err := EncodeToJSON(stepTypes)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/step-types",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respStepTypes StepTypes
	err = DecodeResponseInto(resp, &respStepTypes)
	if err != nil {
		log.Printf("[DEBUG] Error while decoding step types. Error = %v, Response: %v", err, respStepTypes)
		return nil, err
	}
	return &respStepTypes, nil

}

func (client *Client) UpdateStepTypes(stepTypes *StepTypes) (*StepTypes, error) {

	body, err := EncodeToJSON(stepTypes)

	if err != nil {
		return nil, err
	}

	fullPath := fmt.Sprintf("/step-types/%s", url.PathEscape(stepTypes.Metadata["name"].(string)+":"+stepTypes.Metadata["version"].(string)))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respStepTypes StepTypes
	err = DecodeResponseInto(resp, &respStepTypes)
	if err != nil {
		return nil, err
	}

	return &respStepTypes, nil

}

func (client *Client) DeleteStepTypes(name string) error {

	fullPath := fmt.Sprintf("/step-types/%s", url.PathEscape(name))
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
