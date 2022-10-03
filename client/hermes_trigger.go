package client

import (
	"fmt"
)

type HermesTrigger struct {
	Event      string    `json:"event,omitempty"`
	PipelineID string    `json:"pipeline,omitempty"`
	EventData  EventData `json:"event-data,omitempty"`
}

type EventData struct {
	Uri     string `json:"uri"`
	Type    string `json:"type"`
	Kind    string `json:"kind"`
	Account string `json:"account"`
	Secret  string `json:"secret"`
}

func (client *Client) GetHermesTriggerByEventAndPipeline(event string, pipeline string) (*HermesTrigger, error) {

	fullPath := fmt.Sprintf("/hermes/triggers/event/%s", UriEncodeEvent(event))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var hermesTriggerList []HermesTrigger

	err = DecodeResponseInto(resp, &hermesTriggerList)
	if err != nil {
		return nil, err
	}

	var hermesTrigger HermesTrigger
	for _, trigger := range hermesTriggerList {
		if trigger.PipelineID == pipeline {
			hermesTrigger = trigger
		}
	}
	if hermesTrigger.Event == "" {
		return nil, fmt.Errorf("No trigger found for event %s and pipeline %s", event, pipeline)
	}

	return &hermesTrigger, nil
}

func (client *Client) CreateHermesTriggerByEventAndPipeline(event string, pipeline string) error {

	fullPath := fmt.Sprintf("/hermes/triggers/%s/%s", UriEncodeEvent(event), pipeline)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
	}

	_, err := client.RequestAPI(&opts)
	return err
}

func (client *Client) DeleteHermesTriggerByEventAndPipeline(event string, pipeline string) error {
	fullPath := fmt.Sprintf("/hermes/triggers/%s/%s", UriEncodeEvent(event), pipeline)
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
