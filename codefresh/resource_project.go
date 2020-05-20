package codefresh

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type project struct {
	ID          string     `json:"id,omitempty"`
	ProjectName string     `json:"projectName,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Variables   []variable `json:"variables,omitempty"`
}

func (p *project) getID() string {
	return p.ID
}

func (p *project) setVariables(variables map[string]interface{}) {
	for key, value := range variables {
		p.Variables = append(p.Variables, variable{Key: key, Value: value.(string)})
	}
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
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
	return createCodefreshObject(
		meta.(*Config),
		"/projects",
		"POST",
		d,
		mapResourceToProject,
		readProject,
	)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	return readCodefreshObject(
		d,
		meta.(*Config),
		getProjectFromCodefresh,
		mapProjectToResource)
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	path := fmt.Sprintf("/projects/%v", d.Id())
	return updateCodefreshObject(
		d,
		meta.(*Config),
		path,
		"PATCH",
		mapResourceToProject,
		readProject,
		resourceProjectRead)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	path := fmt.Sprintf("/projects/%v", d.Id())
	return deleteCodefreshObject(meta.(*Config), path)
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
		meta.(*Config),
		getProjectFromCodefresh,
		mapProjectToResource)
}

func readProject(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	project := &project{}
	err := json.Unmarshal(b, project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func getProjectFromCodefresh(d *schema.ResourceData, c *Config) (codefreshObject, error) {
	projectName := d.Id()
	path := fmt.Sprintf("/projects/%v", projectName)
	return getFromCodefresh(d, c, path, readProject)
}

func mapProjectToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	project := cfObject.(*project)
	d.SetId(project.ID)

	err := d.Set("name", project.ProjectName)
	if err != nil {
		return err
	}

	err = d.Set("tags", project.Tags)
	if err != nil {
		return err
	}

	err = d.Set("variables", convertVariables(project.Variables))
	if err != nil {
		return err
	}
	return nil
}

func mapResourceToProject(d *schema.ResourceData) codefreshObject {
	tags := d.Get("tags").(*schema.Set).List()
	project := &project{
		ProjectName: d.Get("name").(string),
		Tags:        convertStringArr(tags),
	}
	variables := d.Get("variables").(map[string]interface{})
	project.setVariables(variables)
	return project
}
