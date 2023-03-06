package codefresh

import (
	"fmt"
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePermission() *schema.Resource {
	return &schema.Resource{
		Description: "Permission are used to setup access control and allow to define which teams have access to which clusters and pipelines based on tags.",
		Create:      resourcePermissionCreate,
		Read:        resourcePermissionRead,
		Update:      resourcePermissionUpdate,
		Delete:      resourcePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"_id": {
				Description: "The permission ID.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"team": {
				Description: "The Id of the team the permissions apply to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"resource": {
				Description: `
The type of resources the permission applies to. Possible values:
	* pipeline
	* cluster
				`,
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "cluster" && v != "pipeline" {
						errs = append(errs, fmt.Errorf("%q must be between \"pipeline\" or \"cluster\", got: %s", key, v))
					}
					return
				},
			},
			"action": {
				Description: `
Action to be allowed. Possible values:
	* create
	* read
	* update
	* delete
	* run (Only valid for pipeline resource)
	* approve (Only valid for pipeline resource)
	* debug (Only valid for pipeline resource)
				`,
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "create" && v != "read" && v != "update" && v != "delete" && v != "run" && v != "approve" && v != "debug" {
						errs = append(errs, fmt.Errorf("%q must be between one of create,read,update,delete,approve,debug got: %s", key, v))
					}
					return
				},
			},
			"tags": {
				Description: `
The effective tags to apply the permission. It supports 2 custom tags:
	* untagged is a “tag” which refers to all clusters that don't have any tag.
	* (the star character) means all tags.
				`,
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	permission := *mapResourceToPermission(d)

	newPermission, err := client.CreatePermission(&permission)
	if err != nil {
		return err
	}
	if newPermission == nil {
		return fmt.Errorf("resourcePermissionCreate - failed to create permission, empty responce")
	}

	d.SetId(newPermission.ID)

	return resourcePermissionRead(d, meta)
}

func resourcePermissionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	permissionID := d.Id()
	if permissionID == "" {
		d.SetId("")
		return nil
	}

	permission, err := client.GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	err = mapPermissionToResource(permission, d)
	if err != nil {
		return err
	}

	return nil
}

func resourcePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	permission := *mapResourceToPermission(d)
	permission.ID = ""
	resp, err := client.CreatePermission(&permission)
	if err != nil {
		return err
	}

	deleteErr := resourcePermissionDelete(d, meta)
	if deleteErr != nil {
		log.Printf("[WARN] failed to delete permission %v: %v", permission, deleteErr)
	}
	d.SetId(resp.ID)

	return resourcePermissionRead(d, meta)
}

func resourcePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	err := client.DeletePermission(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapPermissionToResource(permission *cfClient.Permission, d *schema.ResourceData) error {

	err := d.Set("_id", permission.ID)
	if err != nil {
		return err
	}

	err = d.Set("team", permission.Team)
	if err != nil {
		return err
	}

	err = d.Set("action", permission.Action)
	if err != nil {
		return err
	}

	err = d.Set("resource", permission.Resource)
	if err != nil {
		return err
	}

	err = d.Set("tags", permission.Tags)
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToPermission(d *schema.ResourceData) *cfClient.Permission {

	tagsI := d.Get("tags").(*schema.Set).List()
	var tags []string
	if len(tagsI) > 0 {
		tags = convertStringArr(tagsI)
	} else {
		tags = []string{"*", "untagged"}
	}
	permission := &cfClient.Permission{
		ID:       d.Id(),
		Team:     d.Get("team").(string),
		Action:   d.Get("action").(string),
		Resource: d.Get("resource").(string),
		Tags:     tags,
	}

	return permission
}
