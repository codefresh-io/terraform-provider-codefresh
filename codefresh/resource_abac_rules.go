package codefresh

import (
	"context"
	"fmt"
	"log"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validSetValues = []string{"REFRESH", "SYNC", "TERMINATE_SYNC", "VIEW_POD_LOGS", "APP_ROLLBACK"}

func resourceGitopsAbacRule() *schema.Resource {
	return &schema.Resource{
		Description: "Gitops ABAC Rules are used to setup access control and allow to define which teams have access to which resources based on tags and attributes.",
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
The type of resources the ABAC rules applies to. Possible values:
	* gitopsApplications
				`,
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"gitopsApplications",
				}, false),
			},
			"teams": {
				Description: "The IDs of the teams the ABAC rules apply to.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Description: `
The effective tags of the resource to apply the permission to. There are two special tags:
	* untagged: Apply to all resources without tags.
	* * (asterisk): Apply to all resources with any tag.
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
			"attribute": {
				Description: "Resource attribute that need to be validated",
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
	client := meta.(*cfclient.Client)

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

	client := meta.(*cfclient.Client)

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
	client := meta.(*cfclient.Client)

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
	client := meta.(*cfclient.Client)

	_, err := client.DeleteAbacRule(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func flattenAttributes(attributes []cfclient.EntityAbacAttribute) []map[string]interface{} {
	var res = make([]map[string]interface{}, len(attributes))
	for i, attribute := range attributes {
		m := make(map[string]interface{})
		m["name"] = attribute.Name
		m["key"] = attribute.Key
		m["value"] = attribute.Value
		res[i] = m
	}
	return res
}

func mapGitopsAbacRuleToResource(abacRule *cfclient.GitopsAbacRule, d *schema.ResourceData) error {

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

	if len(abacRule.Attributes) > 0 {
		err = d.Set("attribute", flattenAttributes(abacRule.Attributes))
		if err != nil {
			return err
		}
	}

	return nil
}

func mapResourceToGitopsAbacRule(d *schema.ResourceData) *cfclient.GitopsAbacRule {

	tagsI := d.Get("tags").(*schema.Set).List()
	var tags []string
	if len(tagsI) > 0 {
		tags = datautil.ConvertStringArr(tagsI)
	} else {
		tags = []string{"*", "untagged"}
	}

	abacRule := &cfclient.GitopsAbacRule{
		ID:         d.Id(),
		EntityType: d.Get("entity_type").(string),
		Teams:      datautil.ConvertStringArr(d.Get("teams").(*schema.Set).List()),
		Tags:       tags,
		Actions:    datautil.ConvertStringArr(d.Get("actions").(*schema.Set).List()),
		Attributes: []cfclient.EntityAbacAttribute{},
	}

	attributes := d.Get("attribute").([]interface{})
	for idx := range attributes {
		attr := cfclient.EntityAbacAttribute{
			Name:  d.Get(fmt.Sprintf("attribute.%v.name", idx)).(string),
			Key:   d.Get(fmt.Sprintf("attribute.%v.key", idx)).(string),
			Value: d.Get(fmt.Sprintf("attribute.%v.value", idx)).(string),
		}
		abacRule.Attributes = append(abacRule.Attributes, attr)
	}
	return abacRule
}
