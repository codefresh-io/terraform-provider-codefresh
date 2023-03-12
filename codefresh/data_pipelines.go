package codefresh

import (
	"fmt"
	"time"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePipelines() *schema.Resource {
	return &schema.Resource{
		Description: "This resource retrives Pipeline based on the provided attributes. Currently `spec` is limited to the YAML spec.",
		Read:        dataSourcePipelinesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Description: "The query to use when retrieving pipelines. See: [API Documentation](https://g.codefresh.io/api/#operation/pipelines-get-names).",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the attribute to filter by. Must be a valid path to a field in the pipeline object. Accepts regex expressions.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"values": {
							Description: "The acceptable values of the filter.",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        schema.TypeString,
						},
					},
				},
			},
			"pipelines": {
				Type:     schema.TypeList,
				Computed: true,
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
	for i, p := range res {
		match := false
		if d.Get(fmt.Sprintf("filter.%d", i)).(*schema.Set).Len() == 0 {
			match = true
		}

		//for _, f := range d.Get("filter").(*schema.Set).List() {
		//	name := f.(map[string]interface{})["name"].(string)
		//	values := f.(map[string]interface{})["values"].([]string)

		//	attribute, ok := p[name]
		//	if !ok {
		//		return errors.New("Invalid filter name: " + name)
		//	}

		//	for _, v := range values {
		//		r, err := regexp.Compile(v)
		//		if err != nil {
		//			return errors.New("Not a valid regular expression: " + v)
		//		}

		//		switch attribute.(type) {
		//		case bool, string:
		//			match = match || r.MatchString(attribute.(string))
		//		case []string:
		//			for _, s := range attribute.([]string) {
		//				match = match || r.MatchString(s)
		//			}
		//		}
		//	}
		//}

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
