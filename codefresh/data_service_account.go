package codefresh

import (
	"fmt"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves a Codefresh service account by its ID or name.",
		Read:        dataSourceServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Description:  "Service account name",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"_id", "name"},
			},
			"assign_admin_role": {
				Description: "Whether or not account admin role is assigned to the service account",
				Type:        schema.TypeBool,
				Optional:    true,
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

func dataSourceServiceAccountRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)
	var serviceAccount *cfClient.ServiceUser
	var err error

	if _id, _idOk := d.GetOk("_id"); _idOk {
		serviceAccount, err = client.GetServiceUserByID(_id.(string))
	} else if name, nameOk := d.GetOk("name"); nameOk {
		serviceAccount, err = client.GetServiceUserByName(name.(string))
	}

	if err != nil {
		return err
	}

	if serviceAccount == nil {
		return fmt.Errorf("data.codefresh_service_account - cannot find service account")
	}

	return mapDataServiceAccountToResource(serviceAccount, d)

}

func mapDataServiceAccountToResource(serviceAccount *cfClient.ServiceUser, d *schema.ResourceData) error {

	if serviceAccount == nil || serviceAccount.ID == "" {
		return fmt.Errorf("data.codefresh_service_account - failed to mapDataServiceAccountToResource")
	}

	d.SetId(serviceAccount.ID)
	d.Set("name", serviceAccount.Name)
	d.Set("assign_admin_role", serviceAccount.HasAdminRole())

	teamIds := []string{}

	for _, team := range serviceAccount.Teams {
		teamIds = append(teamIds, team.ID)
	}

	d.Set("assigned_teams", teamIds)

	return nil
}
