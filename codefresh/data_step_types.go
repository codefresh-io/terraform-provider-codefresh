package codefresh

import (
	"fmt"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStepTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStepTypesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"step_types_yaml": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStepTypesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)
	var stepTypes *cfClient.StepTypes
	var err error
	identifier := d.Get("name").(string)
	version, versionOk := d.GetOk("version")

	if versionOk {
		identifier = identifier + ":" + version.(string)
	}
	stepTypes, err = client.GetStepTypes(identifier)
	if err != nil {
		return err
	}

	if stepTypes == nil {
		return fmt.Errorf("data.codefresh_step_types - cannot find step-types")
	}

	return mapDataSetTypesToResource(stepTypes, d)
}

func mapDataSetTypesToResource(stepTypes *cfClient.StepTypes, d *schema.ResourceData) error {

	if stepTypes == nil || stepTypes.Metadata["name"].(string) == "" {
		return fmt.Errorf("data.codefresh_step_types - failed to mapDataSetTypesToResource")
	}
	d.SetId(stepTypes.Metadata["name"].(string))

	d.Set("name", d.Id())

	stepTypesYaml, err := yaml.Marshal(stepTypes)
	if err != nil {
		return err
	}
	d.Set("step_types_yaml", string(stepTypesYaml))

	return nil
}
