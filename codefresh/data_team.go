package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves a team by its ID or name.",
		Read:        dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceTeamRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	var team *cfclient.Team
	var err error

	if _id, _idOk := d.GetOk("_id"); _idOk {
		team, err = client.GetTeamByID(_id.(string))
	} else if name, nameOk := d.GetOk("name"); nameOk {
		// accountID, accountOk := d.GetOk("account_id");
		team, err = client.GetTeamByName(name.(string))
	}

	if err != nil {
		return err
	}

	if team == nil {
		return fmt.Errorf("data.codefresh_team - cannot find team")
	}

	return mapDataTeamToResource(team, d)

}

func mapDataTeamToResource(team *cfclient.Team, d *schema.ResourceData) error {

	if team == nil || team.ID == "" {
		return fmt.Errorf("data.codefresh_team - failed to mapDataTeamToResource")
	}
	d.SetId(team.ID)

	err := d.Set("_id", team.ID)

	if err != nil {
		return err
	}

	err = d.Set("account_id", team.Account)

	if err != nil {
		return err
	}

	err = d.Set("type", team.Type)

	if err != nil {
		return err
	}

	var users []string
	for _, user := range team.Users {
		users = append(users, user.ID)
	}

	err = d.Set("users", users)

	if err != nil {
		return err
	}

	err = d.Set("tags", team.Tags)

	if err != nil {
		return err
	}

	return nil
}
