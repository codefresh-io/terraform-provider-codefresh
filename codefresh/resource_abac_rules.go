package codefresh

import (
	"context"
	"fmt"
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validSetValues = []string{"REFRESH", "SYNC", "TERMINATE_SYNC", "VIEW_POD_LOGS", "APP_ROLLBACK"}

func resourceGitopsAbacRule() *schema.Resource {
	return &schema.Resource{
		Description: "Gitops Abac Rules are used to setup access control and allow to define which teams have access to which resources based on tags and attributes.",
		Create:      resourceGitopsAbacRuleCreate,
		Read:        resourceGitopsAbacRuleRead,
		Update:      resourceGitopsAbacRuleUpdate,
		Delete:      resourceGitopsAbacRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The abac rule ID.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"entity_type": {
				Description: `
The type of resources the abac rules applies to. Possible values:
	* gitopsApplications
				`,
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"gitopsApplications",
				}, false),
			},
			"teams": {
				Description: "The Ids of teams the abac rules apply to.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Description: `
The effective tags to apply the permission. It supports 2 custom tags:
	* untagged is a “tag” which refers to all resources that don't have any tag.
	* (the star character) means all tags.
				`,
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"actions": {
				Description: `
Action to be allowed. Possible values:
	* REFRESH
	* SYNC
	* TERMINATE_SYNC
	* VIEW_POD_LOGS
	* APP_ROLLBACK
				`,
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attributes": {
				Description: "Resource attributes that need to be validated",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			actions := diff.Get("actions").(*schema.Set).List()

			for _, action := range actions {
				actionStr := action.(string)
				if !contains(validSetValues, actionStr) {
					return fmt.Errorf("Invalid action value: %s", actionStr)
				}
			}

			return nil
		},
	}
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func resourceGitopsAbacRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	abacRule := *mapResourceToGitopsAbacRule(d)

	newGitopsAbacRule, err := client.CreateAbacRule(&abacRule)
	if err != nil {
		return err
	}
	if newGitopsAbacRule == nil {
		return fmt.Errorf("resourceGitopsAbacRuleCreate - failed to create abac rule, empty response")
	}

	d.SetId(newGitopsAbacRule.ID)

	return resourceGitopsAbacRuleRead(d, meta)
}

func resourceGitopsAbacRuleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	abacRuleID := d.Id()
	if abacRuleID == "" {
		d.SetId("")
		return nil
	}

	abacRule, err := client.GetAbacRuleByID(abacRuleID)
	if err != nil {
		return err
	}

	err = mapGitopsAbacRuleToResource(abacRule, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceGitopsAbacRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	abacRule := *mapResourceToGitopsAbacRule(d)
	resp, err := client.CreateAbacRule(&abacRule)
	if err != nil {
		return err
	}

	deleteErr := resourceGitopsAbacRuleDelete(d, meta)
	if deleteErr != nil {
		log.Printf("[WARN] failed to delete permission %v: %v", abacRule, deleteErr)
	}
	d.SetId(resp.ID)

	return resourceGitopsAbacRuleRead(d, meta)
}

func resourceGitopsAbacRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	_, err := client.DeleteAbacRule(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapGitopsAbacRuleToResource(abacRule *cfClient.GitopsAbacRule, d *schema.ResourceData) error {

	err := d.Set("id", abacRule.ID)
	if err != nil {
		return err
	}

	err = d.Set("entity_type", abacRule.EntityType)
	if err != nil {
		return err
	}

	err = d.Set("teams", abacRule.Teams)
	if err != nil {
		return err
	}

	err = d.Set("tags", abacRule.Tags)
	if err != nil {
		return err
	}

	err = d.Set("actions", abacRule.Actions)
	if err != nil {
		return err
	}

	err = d.Set("attributes", abacRule.Attributes)
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToGitopsAbacRule(d *schema.ResourceData) *cfClient.GitopsAbacRule {

	tagsI := d.Get("tags").(*schema.Set).List()
	var tags []string
	if len(tagsI) > 0 {
		tags = convertStringArr(tagsI)
	} else {
		tags = []string{"*", "untagged"}
	}
	abacRule := &cfClient.GitopsAbacRule{
		ID:         d.Id(),
		EntityType: d.Get("entity_type").(string),
		Teams:      convertStringArr(d.Get("teams").(*schema.Set).List()),
		Tags:       tags,
		Actions:    convertStringArr(d.Get("actions").(*schema.Set).List()),
		Attributes: d.Get("attributes").([]cfClient.EntityAbacAttribute),
	}

	return abacRule
}
