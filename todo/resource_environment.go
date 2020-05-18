package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type environmentMetadata struct {
	Name      string `json:"name,omitempty"`
	AccountID string `json:"accountId,omitempty"`
	ID        string `json:"id,omitempty"`
}

type environmentFilter struct {
	Cluster   string `json:"cluster,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type environmentEndpoint struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type environmentSpec struct {
	Type      string                `json:"type,omitempty"`
	Filters   []environmentFilter   `json:"filters,omitempty"`
	Endpoints []environmentEndpoint `json:"endpoints,omitempty"`
}

type environment struct {
	Metadata environmentMetadata `json:"metadata,omitempty"`
	Spec     environmentSpec     `json:"spec,omitempty"`
	Version  string              `json:"version,omitempty"`
	Kind     string              `json:"kind,omitempty"`
}

func (e *environment) getID() string {
	return e.Metadata.ID
}

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEnvironmentImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, _ interface{}) error {
	return createCodefreshObject(
		fmt.Sprintf("%v/environments-v2", getCfUrl()),
		"POST",
		d,
		mapResourceToEnvironment,
		readEnvironment,
	)
}

func resourceEnvironmentRead(d *schema.ResourceData, _ interface{}) error {
	return readCodefreshObject(
		d,
		getEnvironmentFromCodefresh,
		mapEnvironmentToResource)
}

func resourceEnvironmentUpdate(d *schema.ResourceData, _ interface{}) error {
	url := fmt.Sprintf("%v/environments-v2/%v", getCfUrl(), d.Id())
	return updateCodefreshObject(
		d,
		url,
		"PUT",
		mapResourceToEnvironment,
		readEnvironment,
		resourceEnvironmentRead)
}

func resourceEnvironmentDelete(d *schema.ResourceData, _ interface{}) error {
	url := fmt.Sprintf("%v/environments-v2/%v", getCfUrl(), d.Id())
	return deleteCodefreshObject(url)
}

func resourceEnvironmentImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
		getEnvironmentFromCodefresh,
		mapEnvironmentToResource)
}

func readEnvironment(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	environment := &environment{}
	err := json.Unmarshal(b, environment)
	if err != nil {
		return nil, err
	}
	return environment, nil
}

func getEnvironmentFromCodefresh(d *schema.ResourceData) (codefreshObject, error) {
	id := d.Id()
	url := fmt.Sprintf("%v/environments-v2/%v", getCfUrl(), id)
	return getFromCodefresh(d, url, readEnvironment)
}

func mapEnvironmentToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	environment := cfObject.(*environment)
	d.SetId(environment.Metadata.ID)

	err := d.Set("name", environment.Metadata.Name)
	if err != nil {
		return err
	}

	if len(environment.Spec.Filters) != 1 {
		return fmt.Errorf("expected exactly one filter, but found %v", len(environment.Spec.Filters))
	}

	err = d.Set("namespace", environment.Spec.Filters[0].Namespace)
	if err != nil {
		return err
	}

	err = d.Set("cluster", environment.Spec.Filters[0].Cluster)
	if err != nil {
		return err
	}

	err = d.Set("account_id", environment.Metadata.AccountID)
	if err != nil {
		return err
	}

	if len(environment.Spec.Endpoints) > 0 {
		err = d.Set("url", environment.Spec.Endpoints[0].URL)
		if err != nil {
			return err
		}
	}

	return nil
}

func mapResourceToEnvironment(d *schema.ResourceData) codefreshObject {
	environment := &environment{
		Metadata: environmentMetadata{
			Name:      d.Get("name").(string),
			AccountID: d.Get("account_id").(string),
		},
		Spec: environmentSpec{
			Type: "kubernetes",
			Filters: []environmentFilter{
				{
					Cluster:   d.Get("cluster").(string),
					Namespace: d.Get("namespace").(string),
				},
			},
		},
		Version: "1.0",
		Kind:    "environment",
	}

	if d.Get("url").(string) != "" {
		environment.Spec.Endpoints = []environmentEndpoint{
			{
				Name: "default",
				URL:  d.Get("url").(string),
			},
		}
	}
	return environment
}
