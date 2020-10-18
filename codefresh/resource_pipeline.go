package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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
			"original_yaml_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,

			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
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
				Optional: true,
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
						"runtime_environment": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dind_storage": {
										Type:     schema.TypeString,
										Optional: true,
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

	return resourcePipelineRead(d, meta)
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipelineID := d.Id()

	if pipelineID == "" {
		d.SetId("")
		return nil
	}

	pipeline, err := client.GetPipeline(pipelineID)
	if err != nil {
		return err
	}

	err = mapPipelineToResource(*pipeline, d)
	if err != nil {
		return err
	}

	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipeline := *mapResourceToPipeline(d)
	pipeline.Metadata.ID = d.Id()

	_, err := client.UpdatePipeline(&pipeline)
	if err != nil {
		return err
	}

	return resourcePipelineRead(d, meta)
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	err := client.DeletePipeline(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapPipelineToResource(pipeline cfClient.Pipeline, d *schema.ResourceData) error {

	err := d.Set("name", pipeline.Metadata.Name)
	if err != nil {
		return err
	}

	err = d.Set("revision", pipeline.Metadata.Revision)
	if err != nil {
		return err
	}

	err = d.Set("project_id", pipeline.Metadata.ProjectId)
	if err != nil {
		return err
	}

	err = d.Set("spec", flattenSpec(pipeline.Spec))
	if err != nil {
		return err
	}

	err = d.Set("tags", pipeline.Metadata.Labels.Tags)
	if err != nil {
		return err
	}

	err = d.Set("original_yaml_string", pipeline.Metadata.OriginalYamlString)
	if err != nil {
		return err
	}

	return nil
}

func flattenSpec(spec cfClient.Spec) []interface{} {

	var res = make([]interface{}, 0)
	m := make(map[string]interface{})

	if len(spec.Triggers) > 0 {
		m["trigger"] = flattenTriggers(spec.Triggers)
	}

	if spec.SpecTemplate != nil {
		m["spec_template"] = flattenSpecTemplate(*spec.SpecTemplate)
	}

	if len(spec.Variables) != 0 {
		m["variables"] = convertVariables(spec.Variables)
	}

	if spec.RuntimeEnvironment != (cfClient.RuntimeEnvironment{}) {
		m["runtime_environment"] = flattenSpecRuntimeEnvironment(spec.RuntimeEnvironment)
	}

	m["concurrency"] = spec.Concurrency

	m["priority"] = spec.Priority

	res = append(res, m)
	return res
}

func flattenSpecTemplate(spec cfClient.SpecTemplate) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"location": spec.Location,
			"repo":     spec.Repo,
			"context":  spec.Context,
			"revision": spec.Revision,
			"path":     spec.Path,
		},
	}
}

func flattenSpecRuntimeEnvironment(spec cfClient.RuntimeEnvironment) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":         spec.Name,
			"memory":       spec.Memory,
			"cpu":          spec.CPU,
			"dind_storage": spec.DindStorage,
		},
	}
}

func flattenTriggers(triggers []cfClient.Trigger) []map[string]interface{} {
	var res = make([]map[string]interface{}, len(triggers))
	for i, trigger := range triggers {
		m := make(map[string]interface{})
		m["name"] = trigger.Name
		m["description"] = trigger.Description
		m["context"] = trigger.Context
		m["repo"] = trigger.Repo
		m["branch_regex"] = trigger.BranchRegex
		m["modified_files_glob"] = trigger.ModifiedFilesGlob
		m["disabled"] = trigger.Disabled
		m["provider"] = trigger.Provider
		m["type"] = trigger.Type
		m["events"] = trigger.Events
		m["variables"] = convertVariables(trigger.Variables)

		res[i] = m
	}
	return res
}

func mapResourceToPipeline(d *schema.ResourceData) *cfClient.Pipeline {

	tags := d.Get("tags").(*schema.Set).List()

	pipeline := &cfClient.Pipeline{
		Metadata: cfClient.Metadata{
			Name:      d.Get("name").(string),
			Revision:  d.Get("revision").(int),
			ProjectId: d.Get("project_id").(string),
			Labels: cfClient.Labels{
				Tags: convertStringArr(tags),
			},
			OriginalYamlString: strings.Replace(
				d.Get("original_yaml_string").(string),
				"\n",
				"\n",
				-1),
		},
		Spec: cfClient.Spec{
			Priority:    d.Get("spec.0.priority").(int),
			Concurrency: d.Get("spec.0.concurrency").(int),
		},
	}

	if _, ok := d.GetOk("spec.0.spec_template"); ok {
		pipeline.Spec.SpecTemplate = &cfClient.SpecTemplate{
			Location: d.Get("spec.0.spec_template.0.location").(string),
			Repo:     d.Get("spec.0.spec_template.0.repo").(string),
			Path:     d.Get("spec.0.spec_template.0.path").(string),
			Revision: d.Get("spec.0.spec_template.0.revision").(string),
			Context:  d.Get("spec.0.spec_template.0.context").(string),
		}
	}

	if _, ok := d.GetOk("spec.0.runtime_environment"); ok {
		pipeline.Spec.RuntimeEnvironment = cfClient.RuntimeEnvironment{
			Name:        d.Get("spec.0.runtime_environment.0.name").(string),
			Memory:      d.Get("spec.0.runtime_environment.0.memory").(string),
			CPU:         d.Get("spec.0.runtime_environment.0.cpu").(string),
			DindStorage: d.Get("spec.0.runtime_environment.0.dind_storage").(string),
		}
	}

	variables := d.Get("spec.0.variables").(map[string]interface{})
	pipeline.SetVariables(variables)

	triggers := d.Get("spec.0.trigger").([]interface{})
	for idx := range triggers {
		events := d.Get(fmt.Sprintf("spec.0.trigger.%v.events", idx)).([]interface{})

		codefreshTrigger := cfClient.Trigger{
			Name:              d.Get(fmt.Sprintf("spec.0.trigger.%v.name", idx)).(string),
			Description:       d.Get(fmt.Sprintf("spec.0.trigger.%v.description", idx)).(string),
			Type:              d.Get(fmt.Sprintf("spec.0.trigger.%v.type", idx)).(string),
			Repo:              d.Get(fmt.Sprintf("spec.0.trigger.%v.repo", idx)).(string),
			BranchRegex:       d.Get(fmt.Sprintf("spec.0.trigger.%v.branch_regex", idx)).(string),
			ModifiedFilesGlob: d.Get(fmt.Sprintf("spec.0.trigger.%v.modified_files_glob", idx)).(string),
			Provider:          d.Get(fmt.Sprintf("spec.0.trigger.%v.provider", idx)).(string),
			Disabled:          d.Get(fmt.Sprintf("spec.0.trigger.%v.disabled", idx)).(bool),
			Context:           d.Get(fmt.Sprintf("spec.0.trigger.%v.context", idx)).(string),
			Events:            convertStringArr(events),
		}
		variables := d.Get(fmt.Sprintf("spec.0.trigger.%v.variables", idx)).(map[string]interface{})
		codefreshTrigger.SetVariables(variables)

		pipeline.Spec.Triggers = append(pipeline.Spec.Triggers, codefreshTrigger)
	}
	return pipeline
}
