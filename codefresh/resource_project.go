package codefresh

import (
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description: `
The top-level concept in Codefresh. You can create projects to group pipelines that are related.
In most cases a single project will be a single application (that itself contains many micro-services).
You are free to use projects as you see fit. For example, you could create a project for a specific Kubernetes cluster or a specific team/department.
		`,
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The display name for the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tags": {
				Description: "A list of tags to mark a project for easy management and access control.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"variables": {
				Description: "Project variables.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	project := *mapResourceToProject(d)

	resp, err := client.CreateProject(&project)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	projectID := d.Id()
	if projectID == "" {
		d.SetId("")
		return nil
	}

	project, err := client.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	err = mapProjectToResource(project, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	project := *mapResourceToProject(d)

	err := client.UpdateProject(&project)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)
	// Adding a Retry backoff to address eventual consistency for the API
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 2 * time.Second
	err := backoff.Retry(
		func() error {
			err := client.DeleteProject(d.Id())
			if err != nil {
				log.Printf("Unable to destroy Project due to error %v", err)
			}
			return err
		}, expBackoff)
	if err != nil {
		return err
	}

	return nil
}

func mapProjectToResource(project *cfclient.Project, d *schema.ResourceData) error {

	err := d.Set("name", project.ProjectName)
	if err != nil {
		return err
	}

	err = d.Set("tags", project.Tags)
	if err != nil {
		return err
	}

	vars, _ := datautil.ConvertVariables(project.Variables)

	err = d.Set("variables", vars)
	if err != nil {
		return err
	}
	return nil
}

func mapResourceToProject(d *schema.ResourceData) *cfclient.Project {
	tags := d.Get("tags").(*schema.Set).List()
	project := &cfclient.Project{
		ID:          d.Id(),
		ProjectName: d.Get("name").(string),
		Tags:        datautil.ConvertStringArr(tags),
	}
	variables := d.Get("variables").(map[string]interface{})
	project.SetVariables(variables)
	return project
}
