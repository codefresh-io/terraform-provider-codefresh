package codefresh

import (
	"fmt"	
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
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "cluster" || v != "pipeline" {
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
					if v != "create" || v != "read" || v != "update" || v != "delete" || v != "approve"  {
					  errs = append(errs, fmt.Errorf("%q must be between one of create,read,update,delete,approve, got: %s", key, v))
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
				Default: []string{"*", "untagged"},
			},
		},
	}
}

func resourcePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	permission := *mapResourceToPermission(d)

	resp, err := client.CreatePermission(&permission)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

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


	// existingPermission, err := client.GetPermissionByID(permission.ID)
	// if err != nil {
	// 	return nil
	// }

	resp, err := client.CreatePermission(&permission)
	if err != nil {
		return err
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

	err = d.Set("account", permission.Account)
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
	
	tags := d.Get("tags").(*schema.Set).List()
	permission := &cfClient.Permission{
		ID:      d.Id(),
		Team:    d.Get("team").(string),
		Action:  d.Get("action").(string),
		Resource: d.Get("string").(string),
		//Account: d.Get("account_id").(string),
		Tags:    convertStringArr(tags),
	}

	return permission
}
