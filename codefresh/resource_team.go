package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Update: resourceTeamUpdate,
		Delete: resourceTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	team := *mapResourceToTeam(d)

	resp, err := client.CreateTeam(&team)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceTeamRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	teamID := d.Id()
	if teamID == "" {
		d.SetId("")
		return nil
	}

	team, err := client.GetTeamByID(teamID)
	if err != nil {
		return err
	}

	err = mapTeamToResource(team, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	team := *mapResourceToTeam(d)

	// Rename
	err := client.RenameTeam(team.ID, team.Name)
	if err != nil {
		return err
	}

	// Update users
	existingTeam, err := client.GetTeamByID(team.ID)
	if err != nil {
		return nil
	}

	desiredUsers := d.Get("users").(*schema.Set).List()

	usersToAdd, usersToDelete := cfClient.GetUsersDiff(convertStringArr(desiredUsers), existingTeam.Users)

	for _, userId := range usersToDelete{
		err := client.DeleteUserFromTeam(team.ID, userId)
		if err != nil {
			return err
		}
	}

	for _, userId := range usersToAdd{
		err := client.AddUserToTeam(team.ID, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	err := client.DeleteTeam(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapTeamToResource(team *cfClient.Team, d *schema.ResourceData) error {

	err := d.Set("name", team.Name)
	if err != nil {
		return err
	}

	err = d.Set("type", team.Type)
	if err != nil {
		return err
	}

	err = d.Set("account_id", team.Account)
	if err != nil {
		return err
	}

	err = d.Set("tags", team.Tags)
	if err != nil {
		return err
	}

	err = d.Set("users", flattenTeamUsers(team.Users))
	if err != nil {
		return err
	}

	return nil
}

func flattenTeamUsers(users []cfClient.TeamUser) []string {
	res := []string{}
	for _, user := range users {
		res = append(res, user.ID)
	}
	return res
}

func mapResourceToTeam(d *schema.ResourceData) *cfClient.Team {
	tags := d.Get("tags").(*schema.Set).List()
	team := &cfClient.Team{
		ID:      d.Id(),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Account: d.Get("account_id").(string),
		Tags:    convertStringArr(tags),
	}

	if _, ok := d.GetOk("users"); ok {
		users := d.Get("users").(*schema.Set).List()
		for _, id := range users {
			user := cfClient.TeamUser{
				ID: id.(string),
			}
			team.Users = append(team.Users, user)
		}
	}

	return team
}
