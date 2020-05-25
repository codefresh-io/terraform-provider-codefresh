package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"spec_template": {
							Type:     schema.TypeList,
							Optional: true,
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
						"trigger": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
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
										Optional: true,
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
										Optional: true,
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
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipeline := *mapResourceToPipeline(d)

	resp, err := client.CreatePipeline(&pipeline)
	if err != nil {
		return err
	}

	d.SetId(resp.Metadata.ID)

	return nil
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipelineName := d.Id()

	pipeline, err := client.GetPipeline(pipelineName)
	if err != nil {
		return err
	}

	if pipeline.Metadata.ID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	var pipeline cfClient.Pipeline
	pipeline = *mapResourceToPipeline(d)
	pipeline.Metadata.ID = d.Id()

	_, err := client.UpdatePipeline(&pipeline)
	if err != nil {
		return err
	}

	return nil
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	err := client.DeletePipeline(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToPipeline(d *schema.ResourceData) *cfClient.Pipeline {

	tags := d.Get("tags").(*schema.Set).List()
	pipeline := &cfClient.Pipeline{
		Metadata: cfClient.Metadata{
			Name:      d.Get("name").(string),
			ProjectId: d.Get("project_id").(string),
			Labels: cfClient.Labels{
				Tags: convertStringArr(tags),
			},
		},
		Spec: cfClient.Spec{
			SpecTemplate: cfClient.SpecTemplate{
				Location: d.Get("spec.0.spec_template.0.location").(string),
				Repo:     d.Get("spec.0.spec_template.0.repo").(string),
				Path:     d.Get("spec.0.spec_template.0.path").(string),
				Revision: d.Get("spec.0.spec_template.0.revision").(string),
				Context:  d.Get("spec.0.spec_template.0.context").(string),
			},
			Priority:    d.Get("spec.0.priority").(int),
			Concurrency: d.Get("spec.0.concurrency").(int),
		},
	}
	variables := d.Get("spec.0.variables").(map[string]interface{})
	pipeline.SetVariables(variables)

	triggers := d.Get("spec.trigger").(map[string]interface{})
	for idx := range triggers {
		events := d.Get(fmt.Sprintf("spec.0.trigger.%v.events", idx)).([]interface{})

		codefreshTrigger := cfClient.Trigger{
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
		variables := d.Get(fmt.Sprintf("spec.0.trigger.%v.variables", idx)).(map[string]interface{})
		codefreshTrigger.SetVariables(variables)

		pipeline.Spec.Triggers = append(pipeline.Spec.Triggers, codefreshTrigger)
	}
	return pipeline
}
