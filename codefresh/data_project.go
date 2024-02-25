package codefresh

import (
	"fmt"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves a project by its ID or name.",
		Read:        dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)
	var project *cfClient.Project
	var err error

	if _id, _idOk := d.GetOk("_id"); _idOk {
		project, err = client.GetProjectByID(_id.(string))
	} else if name, nameOk := d.GetOk("name"); nameOk {
		project, err = client.GetProjectByName(name.(string))
	}

	if err != nil {
		return err
	}

	if project == nil {
		return fmt.Errorf("data.codefresh_project - cannot find project")
	}

	return mapDataProjectToResource(project, d)

}

func mapDataProjectToResource(project *cfClient.Project, d *schema.ResourceData) error {

	if project == nil || project.ID == "" {
		return fmt.Errorf("data.codefresh_project - failed to mapDataProjectToResource")
	}
	d.SetId(project.ID)

	d.Set("_id", project.ID)
	d.Set("tags", project.Tags)

	return nil
}
