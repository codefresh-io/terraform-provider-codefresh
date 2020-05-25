package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	var project cfClient.Project
	project = *mapResourceToProject(d)

	resp, err := client.CreateProject(&project)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	projectName := d.Get("name").(string)

	project, err := client.GetProjectByName(projectName)
	if err != nil {
		return err
	}

	if project.ID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	var project cfClient.Project
	project = *mapResourceToProject(d)

	err := client.UpdateProject(&project)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	err := client.DeleteProject(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToProject(d *schema.ResourceData) *cfClient.Project {
	tags := d.Get("tags").(*schema.Set).List()
	project := &cfClient.Project{
		ID:          d.Id(),
		ProjectName: d.Get("name").(string),
		Tags:        convertStringArr(tags),
	}
	variables := d.Get("variables").(map[string]interface{})
	project.SetVariables(variables)
	return project
}
