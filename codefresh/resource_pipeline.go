package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type labels struct {
	Tags []string `json:"tags,omitempty"`
}

type metadata struct {
	Name   string `json:"name,omitempty"`
	ID     string `json:"id,omitempty"`
	Labels labels `json:"labels,omitempty"`
}

type specTemplate struct {
	Location string `json:"location,omitempty"`
	Repo     string `json:"repo,omitempty"`
	Path     string `json:"path,omitempty"`
	Revision string `json:"revision,omitempty"`
	Context  string `json:"context,omitempty"`
}

type trigger struct {
	Name              string     `json:"name,omitempty"`
	Description       string     `json:"description,omitempty"`
	Type              string     `json:"type,omitempty"`
	Repo              string     `json:"repo,omitempty"`
	Events            []string   `json:"events,omitempty"`
	BranchRegex       string     `json:"branchRegex,omitempty"`
	ModifiedFilesGlob string     `json:"modifiedFilesGlob,omitempty"`
	Provider          string     `json:"provider,omitempty"`
	Disabled          bool       `json:"disabled,omitempty"`
	Context           string     `json:"context,omitempty"`
	Variables         []variable `json:"variables,omitempty"`
}

func (t *trigger) setVariables(variables map[string]interface{}) {
	for key, value := range variables {
		t.Variables = append(t.Variables, variable{Key: key, Value: value.(string)})
	}
}

type spec struct {
	Variables    []variable    `json:"variables,omitempty"`
	SpecTemplate *specTemplate `json:"specTemplate,omitempty"`
	Triggers     []trigger     `json:"triggers,omitempty"`
	Priority     int           `json:"priority,omitempty"`
	Concurrency  int           `json:"concurrency,omitempty"`
}

type pipeline struct {
	Metadata metadata `json:"metadata,omitempty"`
	Spec     *spec    `json:"spec,omitempty"`
}

func (p *pipeline) setVariables(variables map[string]interface{}) {
	for key, value := range variables {
		p.Spec.Variables = append(p.Spec.Variables, variable{Key: key, Value: value.(string)})
	}
}

func (p *pipeline) getID() string {
	return p.Metadata.ID
}

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePipelineImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "git",
						},
						"repo": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"revision": {
							Type:     schema.TypeString,
							Required: true,
						},
						"context": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "github",
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"concurrency": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0, // zero is unlimited
						},
					},
				},
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"trigger": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "git",
						},
						"repo": {
							Type:     schema.TypeString,
							Required: true,
						},
						"branch_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/.*/gi",
						},
						"modified_files_glob": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"events": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"provider": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "github",
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"context": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "github",
						},
						"variables": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, m interface{}) error {
	return createCodefreshObject(
		m.(*Config),
		"/pipelines",
		"POST",
		d,
		mapResourceToPipeline,
		readPipeline,
	)
}

func resourcePipelineRead(d *schema.ResourceData, m interface{}) error {
	return readCodefreshObject(
		d,
		m.(*Config),
		getPipelineFromCodefresh,
		mapPipelineToResource)
}

func resourcePipelineUpdate(d *schema.ResourceData, m interface{}) error {
	path := fmt.Sprintf("/pipelines/%v?disableRevisionCheck=true", d.Id())
	return updateCodefreshObject(
		d,
		m.(*Config),
		path,
		"PUT",
		mapResourceToPipeline,
		readPipeline,
		resourcePipelineRead)
}

func resourcePipelineDelete(d *schema.ResourceData, m interface{}) error {
	path := fmt.Sprintf("/pipelines/%v", d.Id())
	return deleteCodefreshObject(m.(*Config), path)
}

func resourcePipelineImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
		m.(*Config),
		getPipelineFromCodefresh,
		mapPipelineToResource)
}

func readPipeline(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	pipeline := &pipeline{}
	err := json.Unmarshal(b, pipeline)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func getPipelineFromCodefresh(d *schema.ResourceData, c *Config) (codefreshObject, error) {
	pipelineName := d.Id()
	path := fmt.Sprintf("/pipelines/%v", pipelineName)
	return getFromCodefresh(d, c, path, readPipeline)
}

func mapPipelineToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	pipeline := cfObject.(*pipeline)
	d.SetId(pipeline.Metadata.ID)

	err := d.Set("name", pipeline.Metadata.Name)
	if err != nil {
		return err
	}

	err = d.Set("spec", flattenSpec(pipeline.Spec))
	if err != nil {
		return err
	}

	err = d.Set("variables", convertVariables(pipeline.Spec.Variables))
	if err != nil {
		return err
	}

	err = d.Set("tags", pipeline.Metadata.Labels.Tags)
	if err != nil {
		return err
	}

	err = d.Set("trigger", flattenTriggers(pipeline.Spec.Triggers))
	if err != nil {
		return err
	}

	return nil
}

func flattenSpec(spec *spec) []map[string]interface{} {
	if spec.SpecTemplate == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"location":    spec.SpecTemplate.Location,
			"repo":        spec.SpecTemplate.Repo,
			"context":     spec.SpecTemplate.Context,
			"revision":    spec.SpecTemplate.Revision,
			"path":        spec.SpecTemplate.Path,
			"priority":    spec.Priority,
			"concurrency": spec.Concurrency,
		},
	}
}

func flattenTriggers(triggers []trigger) []map[string]interface{} {
	var res []map[string]interface{}
	for _, trigger := range triggers {
		res = append(res, map[string]interface{}{
			"name":                trigger.Name,
			"description":         trigger.Description,
			"context":             trigger.Context,
			"repo":                trigger.Repo,
			"branch_regex":        trigger.BranchRegex,
			"modified_files_glob": trigger.ModifiedFilesGlob,
			"disabled":            trigger.Disabled,
			"provider":            trigger.Provider,
			"type":                trigger.Type,
			"events":              trigger.Events,
			"variables":           convertVariables(trigger.Variables),
		})
	}
	return res
}

func mapResourceToPipeline(d *schema.ResourceData) codefreshObject {
	tags := d.Get("tags").(*schema.Set).List()
	pipeline := &pipeline{
		Metadata: metadata{
			Name: d.Get("name").(string),
			Labels: labels{
				Tags: convertStringArr(tags),
			},
		},
		Spec: &spec{
			SpecTemplate: &specTemplate{
				Location: d.Get("spec.0.location").(string),
				Repo:     d.Get("spec.0.repo").(string),
				Path:     d.Get("spec.0.path").(string),
				Revision: d.Get("spec.0.revision").(string),
				Context:  d.Get("spec.0.context").(string),
			},
			Priority:    d.Get("spec.0.priority").(int),
			Concurrency: d.Get("spec.0.concurrency").(int),
		},
	}
	variables := d.Get("variables").(map[string]interface{})
	pipeline.setVariables(variables)

	triggers := d.Get("trigger").([]interface{})
	for idx := range triggers {
		events := d.Get(fmt.Sprintf("trigger.%v.events", idx)).([]interface{})

		codefreshTrigger := trigger{
			Name:              d.Get(fmt.Sprintf("trigger.%v.name", idx)).(string),
			Description:       d.Get(fmt.Sprintf("trigger.%v.description", idx)).(string),
			Type:              d.Get(fmt.Sprintf("trigger.%v.type", idx)).(string),
			Repo:              d.Get(fmt.Sprintf("trigger.%v.repo", idx)).(string),
			BranchRegex:       d.Get(fmt.Sprintf("trigger.%v.branch_regex", idx)).(string),
			ModifiedFilesGlob: d.Get(fmt.Sprintf("trigger.%v.modified_files_glob", idx)).(string),
			Provider:          d.Get(fmt.Sprintf("trigger.%v.provider", idx)).(string),
			Disabled:          d.Get(fmt.Sprintf("trigger.%v.disabled", idx)).(bool),
			Context:           d.Get(fmt.Sprintf("trigger.%v.context", idx)).(string),
			Events:            convertStringArr(events),
		}
		variables := d.Get(fmt.Sprintf("trigger.%v.variables", idx)).(map[string]interface{})
		codefreshTrigger.setVariables(variables)

		pipeline.Spec.Triggers = append(pipeline.Spec.Triggers, codefreshTrigger)
	}
	return pipeline
}
