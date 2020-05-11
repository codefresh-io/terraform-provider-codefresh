package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type cronTrigger struct {
	Pipeline string `json:"pipeline,omitempty"`
	Event    string `json:"event,omitempty"`
}

func (t *cronTrigger) getID() string {
	return t.Pipeline
}

func resourceCronTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceCronTriggerCreate,
		Read:   resourceCronTriggerRead,
		Delete: resourceCronTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCronTriggerImport,
		},
		Schema: map[string]*schema.Schema{
			"pipeline": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCronTriggerCreate(d *schema.ResourceData, _ interface{}) error {
	return createCodefreshObject(
		fmt.Sprintf("%v/hermes/triggers/%v/%v", codefreshURL, urlEncode(d.Get("event").(string)), d.Get("pipeline")),
		"POST",
		d,
		mapResourceToCronTrigger,
		readCronTriggerString,
	)
}

func resourceCronTriggerRead(d *schema.ResourceData, _ interface{}) error {
	return readCodefreshObject(
		d,
		getCronTriggerFromCodefresh,
		mapCronTriggerToResource)
}

// TODO: I don't think this is actually deleting anything, I have an open ticket with Codefresh
// https://support.codefresh.io/hc/en-us/requests/3167?page=1
func resourceCronTriggerDelete(d *schema.ResourceData, _ interface{}) error {
	cfURL := fmt.Sprintf("%v/hermes/triggers/%v/%v", codefreshURL, urlEncode(d.Get("event").(string)), d.Get("pipeline"))
	return deleteCodefreshObject(cfURL)
}

func resourceCronTriggerImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
		getCronTriggerFromCodefresh,
		mapCronTriggerToResource)
}

func mapResourceToCronTrigger(d *schema.ResourceData) codefreshObject {
	return &cronTrigger{
		Pipeline: d.Get("pipeline").(string),
		Event:    d.Get("event").(string),
	}
}

// readCronTriggerString reads a simple string response as is returned from CREATE
func readCronTriggerString(originalObject *schema.ResourceData, _ []byte) (codefreshObject, error) {
	// create doesn't return the object, so we'll just wing it and return the object we sent for creation
	return mapResourceToCronTrigger(originalObject), nil
}

// readCronTrigger reads a JSON type response as is returned from GET
func readCronTrigger(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	var cronTriggerResp []*cronTrigger

	err := json.Unmarshal(b, &cronTriggerResp)
	if err != nil {
		return nil, err
	}

	if len(cronTriggerResp) > 1 {
		return nil, fmt.Errorf("More than 1 trigger per pipeline is not supported. Found %v triggers.", len(cronTriggerResp))
	}
	if len(cronTriggerResp) == 0 {
		return nil, nil
	}
	return cronTriggerResp[0], nil
}

func getCronTriggerFromCodefresh(d *schema.ResourceData) (codefreshObject, error) {
	pipeline := d.Id()
	// get the trigger
	cfURL := fmt.Sprintf("%v/hermes/triggers/pipeline/%v", codefreshURL, pipeline)
	return getFromCodefresh(d, cfURL, readCronTrigger)
}

func mapCronTriggerToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	cronTrigger := cfObject.(*cronTrigger)
	d.SetId(cronTrigger.getID())

	err := d.Set("pipeline", cronTrigger.Pipeline)
	if err != nil {
		return err
	}

	err = d.Set("event", cronTrigger.Event)
	if err != nil {
		return err
	}

	return nil
}
