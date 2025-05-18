package codefresh

import (
	"fmt"
	"log"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStepTypes() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve the published versions of step-types.",
		Read:        dataSourceStepTypesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"step_types_yaml": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceStepTypesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	var err error
	var versions []string
	stepTypesIdentifier := d.Get("name").(string)

	d.SetId(stepTypesIdentifier)
	if versions, err = client.GetStepTypesVersions(stepTypesIdentifier); err == nil {
		var stepVersions cfclient.StepTypesVersions
		stepVersions.Name = stepTypesIdentifier
		err = d.Set("versions", versions)

		if err != nil {
			return err
		}

		for _, version := range versions {
			stepTypes, err := client.GetStepTypes(stepTypesIdentifier + ":" + version)
			if err != nil {
				log.Printf("[DEBUG] Skipping version %v due to error %v", version, err)
			} else {
				stepVersion := cfclient.StepTypesVersion{
					VersionNumber: version,
					StepTypes:     *stepTypes,
				}
				stepVersions.Versions = append(stepVersions.Versions, stepVersion)
			}
		}
		return mapStepTypesVersionsToResource(stepVersions, d)
	}

	return fmt.Errorf("data.codefresh_step_types - was unable to retrieve the versions for step_type %s", stepTypesIdentifier)

}

// func mapDataSetTypesToResource(stepTypesVersions cfclient.StepTypesVersions, d *schema.ResourceData) error {
// 	err := d.Set("name", stepTypesVersions.Name)
// 	if err != nil {
// 		return err
// 	}
// 	err = d.Set("version", flattenVersions(stepTypesVersions.Name, stepTypesVersions.Versions))
// 	return err
// }
