package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "A service account is an identity that provides automated processes, applications, and services with the necessary permissions to interact securely with the Codefresh platform",
		Create:      resourceServiceAccountCreate,
		Read:        resourceServiceAccountRead,
		Update:      resourceServiceAccountUpdate,
		Delete:      resourceServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Service account display name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"assign_admin_role": {
				Description: "Whether or not to assign account admin role to the service account",
				Type: 	  schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"assigned_teams": {
				Description: "A list of team IDs the service account is be assigned to",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceServiceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	newSerivceAccount := *mapResourceToServiceAccount(d)

	resp, err := client.CreateServiceUser(&newSerivceAccount)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceServiceAccountRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	serviceAccountID := d.Id()

	if serviceAccountID == "" {
		d.SetId("")
		return nil
	}

	serviceAccount, err := client.GetServiceUserByID(serviceAccountID)

	if err != nil {
		return err
	}

	err = mapServiceAccountToResource(serviceAccount, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceServiceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	updateServiceAccount := *mapResourceToServiceAccount(d)


	_, err := client.UpdateServiceUser(&updateServiceAccount)

	if err != nil {
		return err
	}

	return nil
}

func resourceServiceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	err := client.DeleteServiceUser(d.Id())

	if err != nil {
		return err
	}

	return nil
}

func mapServiceAccountToResource(serviceAccount *cfclient.ServiceUser, d *schema.ResourceData) error {

	if serviceAccount == nil {
		return fmt.Errorf("mapServiceAccountToResource - cannot find service account")
	}
	err := d.Set("name", serviceAccount.Name)

	if err != nil {
		return err
	}

	teamIds := []string{}

	for _, team := range serviceAccount.Teams {
		teamIds = append(teamIds, team.ID)
	}

	err = d.Set("assigned_teams", teamIds)

	if err != nil {
		return err
	}

	err = d.Set("assign_admin_role", serviceAccount.HasAdminRole())

	if err != nil {
		return err
	}

	return nil
}

func flattenServiceAccountTeams(users []cfclient.TeamUser) []string {
	res := []string{}
	for _, user := range users {
		res = append(res, user.ID)
	}
	return res
}

func mapResourceToServiceAccount(d *schema.ResourceData) *cfclient.ServiceUserCreateUpdate {

	return &cfclient.ServiceUserCreateUpdate{
		ID:             d.Id(),
		Name:           d.Get("name").(string),
		TeamIDs:        datautil.ConvertStringArr(d.Get("assigned_teams").(*schema.Set).List()),
		AssignAdminRole: d.Get("assign_admin_role").(bool),
	}
}
