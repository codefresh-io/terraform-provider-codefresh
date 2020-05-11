package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/url"
	"strings"
)

type cronEventResp struct {
	Type   string `json:"type,omitempty"`
	Kind   string `json:"kind,omitempty"`
	Secret string `json:"secret,omitempty"`
	URI    string `json:"uri,omitempty"`
}

func (r *cronEventResp) toCronEvent() *cronEvent {
	parts := strings.Split(r.URI, ":")
	return &cronEvent{
		Type:   parts[0],
		Kind:   parts[1],
		Secret: parts[4],
		Values: cronEventValues{
			Message:    parts[3],
			Expression: parts[2],
		},
	}
}

type cronEvent struct {
	Type   string          `json:"type,omitempty"`
	Kind   string          `json:"kind,omitempty"`
	Secret string          `json:"secret,omitempty"`
	Values cronEventValues `json:"values,omitempty"`
}

type cronEventValues struct {
	Expression string `json:"expression,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (e *cronEvent) getID() string {
	return fmt.Sprintf("%v:%v:%v:%v:%v", e.Type, e.Kind, e.Values.Expression, e.Values.Message, e.Secret)
}

func (e *cronEvent) urlEncode() string {
	return urlEncode(e.getID())
}

func urlEncode(str string) string {
	return url.QueryEscape(url.QueryEscape(str))
}

func resourceCronEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceCronEventCreate,
		Read:   resourceCronEventRead,
		Delete: resourceCronEventDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCronEventImport,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cron",
				ForceNew: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "codefresh",
				ForceNew: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCronEventCreate(d *schema.ResourceData, _ interface{}) error {
	return createCodefreshObject(
		fmt.Sprintf("%v/hermes/events", codefreshURL),
		"POST",
		d,
		mapResourceToCronEvent,
		readCronEventString,
	)
}

func resourceCronEventRead(d *schema.ResourceData, _ interface{}) error {
	return readCodefreshObject(
		d,
		getCronEventFromCodefresh,
		mapCronEventToResource)
}

// TODO: I don't think this is actually deleting anything, I have an open ticket with Codefresh
// https://support.codefresh.io/hc/en-us/requests/3167?page=1
func resourceCronEventDelete(d *schema.ResourceData, _ interface{}) error {
	cfURL := fmt.Sprintf("%v/hermes/events/%v", codefreshURL, urlEncode(d.Id()))
	return deleteCodefreshObject(cfURL)
}

func resourceCronEventImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
		getCronEventFromCodefresh,
		mapCronEventToResource)
}

func mapResourceToCronEvent(d *schema.ResourceData) codefreshObject {
	dSecret := d.Get("secret").(string)
	if dSecret == "" {
		dSecret = "!generate"
	}
	cronEvent := &cronEvent{
		Type:   d.Get("type").(string),
		Kind:   d.Get("kind").(string),
		Secret: dSecret,
		Values: cronEventValues{
			Expression: d.Get("expression").(string),
			Message:    d.Get("message").(string),
		},
	}
	return cronEvent
}

// readCronEventString reads a simple string response as is returned from CREATE
func readCronEventString(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	var respString string
	err := json.Unmarshal(b, &respString)
	if err != nil {
		return nil, err
	}
	cronEventResp := &cronEventResp{
		URI: respString,
	}
	return cronEventResp.toCronEvent(), nil
}

// readCronEvent reads a JSON type response as is returned from GET
func readCronEvent(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	cronEventResp := &cronEventResp{}

	err := json.Unmarshal(b, cronEventResp)
	if err != nil {
		return nil, err
	}
	return cronEventResp.toCronEvent(), nil
}

func getCronEventFromCodefresh(d *schema.ResourceData) (codefreshObject, error) {
	// get the event
	event := d.Id()
	cfURL := fmt.Sprintf("%v/hermes/events/%v", codefreshURL, urlEncode(event))
	return getFromCodefresh(d, cfURL, readCronEvent)
}

func mapCronEventToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	cronEvent := cfObject.(*cronEvent)
	d.SetId(cronEvent.getID())

	err := d.Set("type", cronEvent.Type)
	if err != nil {
		return err
	}

	err = d.Set("kind", cronEvent.Kind)
	if err != nil {
		return err
	}

	err = d.Set("secret", cronEvent.Secret)
	if err != nil {
		return err
	}

	err = d.Set("message", cronEvent.Values.Message)
	if err != nil {
		return err
	}

	err = d.Set("expression", cronEvent.Values.Expression)
	if err != nil {
		return err
	}
	return nil
}
