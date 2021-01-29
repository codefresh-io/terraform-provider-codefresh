package codefresh

import (
	"fmt"
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionCreate,
		Read:   resourcePermissionRead,
		Update: resourcePermissionUpdate,
		Delete: resourcePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"team": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": {
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
