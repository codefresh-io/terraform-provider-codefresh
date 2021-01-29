package codefresh

import (
	"fmt"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStepTypesVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStepTypesVersionsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceStepTypesVersionsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)
	var versions []string
	var err error

	if name, nameOk := d.GetOk("name"); nameOk {
		if versions, err = client.GetStepTypesVersions(name.(string)); err == nil {
			d.SetId(name.(string))
			d.Set("name", name.(string))
			d.Set("versions", versions)
		}
		return err
	}
	return fmt.Errorf("data.codefresh_step_types_versions - must specify name")

}
