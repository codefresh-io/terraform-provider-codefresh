package client

import (
	"fmt"
)

type HermesTriggerEvent struct {
	Type   string            `json:"type,omitempty"`
	Kind   string            `json:"kind,omitempty"`
	Filter string            `json:"filter,omitempty"`
	Secret string            `json:"secret,omitempty"`
	Values map[string]string `json:"values,omitempty"`
}

func (client *Client) GetHermesTriggerEvent(event string) (*HermesTriggerEvent, error) {
	fullPath := fmt.Sprintf("/hermes/triggers/%s", UriEncodeEvent(event))

	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Trigger Event: \n%v", err)
	}

	var hermesTriggerEvent HermesTriggerEvent
	err = DecodeResponseInto(resp, &hermesTriggerEvent)
	if err != nil {
		return nil, err
	}

	return &hermesTriggerEvent, nil
}

func (client *Client) CreateHermesTriggerEvent(event *HermesTriggerEvent) (string, error) {

	body, err := EncodeToJSON(event)
	if err != nil {
		return "", err
	}

	fullPath := "/hermes/events/"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return "", fmt.Errorf("failed to create Trigger Event: \n%v", err)
	}

	var eventString string
	err = DecodeResponseInto(resp, &eventString)
	if err != nil {
		return "", err
	}

	return eventString, err
}
