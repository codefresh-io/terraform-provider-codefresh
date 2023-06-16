package codefresh

import (
	"fmt"
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
			"entityType": {
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
				Type:         schema.TypeSet,
				Required:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ValidateFunc: ValidateSubset,
			},
			"attributes": {
				Description: "Resource attributes that need to be validated",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type: schema.TypeString,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func ValidateSubset(v interface{}, k string) (warnings []string, errors []error) {
	actions, ok := v.(*schema.Set)
	if !ok {
		errors = append(errors, fmt.Errorf("expected TypeSet for actions"))
		return
	}

	validActions := []string{"REFRESH", "SYNC", "TERMINATE_SYNC", "VIEW_POD_LOGS", "APP_ROLLBACK"} // Allowed values

	for _, team := range actions.List() {
		teamStr, ok := team.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected string element in actions"))
			continue
		}

		if !contains(validActions, teamStr) {
			errors = append(errors, fmt.Errorf("team %s is not a valid action", teamStr))
		}
	}

	return
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

	newGitopsAbacRule, err := client.CreateAbacRule("", &abacRule)
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

	abacRule, err := client.GetAbacRuleByID("", abacRuleID)
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
	resp, err := client.CreateAbacRule("", &abacRule)
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

	_, err := client.DeleteAbacRule("", d.Id())
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

	err = d.Set("entityType", abacRule.EntityType)
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
		EntityType: d.Get("entityType").(cfClient.AbacEntityValues),
		Teams:      d.Get("teams").([]string),
		Tags:       tags,
		Actions:    d.Get("actions").([]string),
		Attributes: d.Get("attributes").([]cfClient.EntityAbacAttribute),
	}

	return abacRule
}
