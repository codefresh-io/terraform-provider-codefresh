package codefresh

import (
	"fmt"
	"regexp"
	"time"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePipelines() *schema.Resource {
	return &schema.Resource{
		Description: "This resource retrives all pipelines belonging to the current user, which can be optionally filtered by the name.",
		Read:        dataSourcePipelinesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Description: "The name regular expression to filter pipelines by.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"pipelines": {
				Description: "The returned list of pipelines. Note that `spec` is currently limited to the YAML, because of the complexity of the object.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeString,
						},
						"is_public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePipelinesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipelines, err := client.GetPipelines()
	if err != nil {
		return err
	}

	err = mapDataPipelinesToResource(*pipelines, d)
	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())

	return nil
}

func mapDataPipelinesToResource(pipelines []cfClient.Pipeline, d *schema.ResourceData) error {
	var res = make([]map[string]interface{}, len(pipelines))
	for i, p := range pipelines {
		m := make(map[string]interface{})
		m["id"] = p.Metadata.ID
		m["name"] = p.Metadata.Name
		m["project"] = p.Metadata.Project
		m["tags"] = p.Metadata.Labels.Tags
		m["is_public"] = p.Metadata.IsPublic
		m["spec"] = p.Metadata.OriginalYamlString

		res[i] = m
	}

	filteredPipelines := make([]map[string]interface{}, 0)
	for _, p := range res {
		match := false

		name, ok := d.GetOk("name_regex")
		if !ok {
			match = true
		} else {
			r, err := regexp.Compile(name.(string))
			if err != nil {
				return fmt.Errorf("`name_regex` is not a valid regular expression, %s", err.Error())
			}
			match = r.MatchString(p["name"].(string))
		}

		if match {
			filteredPipelines = append(filteredPipelines, p)
		}
	}

	err := d.Set("pipelines", filteredPipelines)
	if err != nil {
		return err
	}

	return nil
}
