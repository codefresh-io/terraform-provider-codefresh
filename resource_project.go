package main

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
	return p.ProjectName
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

func resourceProjectCreate(d *schema.ResourceData, _ interface{}) error {
	return createCodefreshObject(
		fmt.Sprintf("%v/projects", getCfUrl()),
		"POST",
		d,
		mapResourceToProject,
		readProject,
	)
}

func resourceProjectRead(d *schema.ResourceData, _ interface{}) error {
	return readCodefreshObject(
		d,
		getProjectFromCodefresh,
		mapProjectToResource)
}

func resourceProjectUpdate(d *schema.ResourceData, _ interface{}) error {
	url := fmt.Sprintf("%v/projects/%v", getCfUrl(), d.Id())
	return updateCodefreshObject(
		d,
		url,
		"PATCH",
		mapResourceToProject,
		readProject,
		resourceProjectRead)
}

func resourceProjectDelete(d *schema.ResourceData, _ interface{}) error {
	url := fmt.Sprintf("%v/projects/%v", getCfUrl(), d.Id())
	return deleteCodefreshObject(url)
}

func resourceProjectImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
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

func getProjectFromCodefresh(d *schema.ResourceData) (codefreshObject, error) {
	projectName := d.Id()
	url := fmt.Sprintf("%v/projects/name/%v", getCfUrl(), projectName)
	return getFromCodefresh(d, url, readProject)
}

func mapProjectToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	project := cfObject.(*project)
	d.SetId(project.ProjectName)

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
