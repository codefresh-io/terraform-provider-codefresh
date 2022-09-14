package client

import (
	"fmt"
)

type HermesTrigger struct {
	Event     string `json:"event,omitempty"`
	Pipeline  string `json:"pipeline,omitempty"`
	EventData string `json:"event-data,omitempty"`
}

type EventData struct {
	Uri     string `json:"uri"`
	Type    string `json:"type"`
	Kind    string `json:"kind"`
	Account string `json:"account"`
	Secret  string `json:"secret"`
}

func (client *Client) GetHermesTriggerByEvent(event string) (*HermesTrigger, error) {

	fullPath := fmt.Sprintf("/hermes/triggers/%s", event)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var hermesTrigger HermesTrigger

	err = DecodeResponseInto(resp, &hermesTrigger)
	if err != nil {
		return nil, err
	}

	return &hermesTrigger, nil
}

func (client *Client) CreateHermesTriggerByEventAndPipeline(event string, pipeline string) (*HermesTrigger, error) {

	fullPath := fmt.Sprintf("/hermes/triggers/%s/%s", event, pipeline)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respHermesTrigger HermesTrigger
	err = DecodeResponseInto(resp, &respHermesTrigger)
	if err != nil {
		return nil, err
	}

	return &respHermesTrigger, nil
}

func (client *Client) DeleteHermesTriggerByEventAndPipeline(event string, pipeline string) error {
	fullPath := fmt.Sprintf("/hermes/triggers/%s/%s", event, pipeline)
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
