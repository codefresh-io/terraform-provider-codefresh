package codefresh

import (
	"context"
	"fmt"
	"log"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePermission() *schema.Resource {
	return &schema.Resource{
		Description: "Permissions are used to set up access control and define which teams have access to which clusters and pipelines based on tags.",
		Create:      resourcePermissionCreate,
		Read:        resourcePermissionRead,
		Update:      resourcePermissionUpdate,
		Delete:      resourcePermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
	* project
	* runtime-environment
				`,
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"pipeline",
					"cluster",
					"project",
					"runtime-environment",
				}, false),
			},
			"related_resource": {
				Description: `
Specifies the resource to use when evaluating the tags. Possible values:
	* project
`,
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"project",
				}, false),
			},
			"action": {
				Description: `
Action to be allowed. Possible values:
	* create
	* read (For runtime-environment resource, 'read' means 'assign')
	* update
	* delete
	* run (Only valid for pipeline resource)
	* approve (Only valid for pipeline resource)
	* debug (Only valid for pipeline resource)
				`,
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"create",
					"read",
					"update",
					"delete",
					"run",
					"approve",
					"debug",
				}, false),
			},
			"tags": {
				Description: `
The tags for which to apply the permission. Supports two custom tags:
	* untagged:  Apply to all resources without tags
  * (asterisk): Apply to all resources with any tag
				`,
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CustomizeDiff: customdiff.All(
			resourcePermissionCustomDiff,
		),
	}
}

func resourcePermissionCustomDiff(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChanges("resource", "related_resource") {
		if diff.Get("related_resource").(string) != "" && diff.Get("resource").(string) != "pipeline" {
			return fmt.Errorf("related_resource is only valid when resource is 'pipeline'")
		}
	}
	if diff.HasChanges("resource", "action") {
		if contains([]string{"run", "approve", "debug"}, diff.Get("action").(string)) && diff.Get("resource").(string) != "pipeline" {
			return fmt.Errorf("action %v is only valid when resource is 'pipeline'", diff.Get("action").(string))
		}
	}
	return nil
}

func resourcePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	permission := *mapResourceToPermission(d)

	newPermission, err := client.CreatePermission(&permission)
	if err != nil {
		return err
	}
	if newPermission == nil {
		return fmt.Errorf("resourcePermissionCreate - failed to create permission, empty response")
	}

	d.SetId(newPermission.ID)

	return resourcePermissionRead(d, meta)
}

func resourcePermissionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

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
	client := meta.(*cfclient.Client)
	permission := *mapResourceToPermission(d)

	// In case team, action or relatedResource or resource have changed - a new permission needs to be created (but without recreating the terraform resource as destruction of resources is alarming for end users)
	if d.HasChanges("team", "action", "related_resource", "resource") {
		deleteErr := resourcePermissionDelete(d, meta)

		if deleteErr != nil {
			log.Printf("[WARN] failed to delete permission %v: %v", permission, deleteErr)
		}

		resp, err := client.CreatePermission(&permission)

		if err != nil {
			return err
		}

		d.SetId(resp.ID)
		// Only tags can be updated
	} else if d.HasChange("tags") {
		err := client.UpdatePermissionTags(&permission)
		if err != nil {
			return err
		}
	}

	return resourcePermissionRead(d, meta)
}

func resourcePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	err := client.DeletePermission(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapPermissionToResource(permission *cfclient.Permission, d *schema.ResourceData) error {

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

	err = d.Set("related_resource", permission.RelatedResource)
	if err != nil {
		return err
	}

	err = d.Set("tags", permission.Tags)
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToPermission(d *schema.ResourceData) *cfclient.Permission {

	tagsI := d.Get("tags").(*schema.Set).List()
	var tags []string
	if len(tagsI) > 0 {
		tags = datautil.ConvertStringArr(tagsI)
	} else {
		tags = []string{"*", "untagged"}
	}
	permission := &cfclient.Permission{
		ID:              d.Id(),
		Team:            d.Get("team").(string),
		Action:          d.Get("action").(string),
		Resource:        d.Get("resource").(string),
		RelatedResource: d.Get("related_resource").(string),
		Tags:            tags,
	}

	return permission
}
