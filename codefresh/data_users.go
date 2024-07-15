package codefresh

import (
	"time"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves all users in the system. Requires a Codefresh admin token and applies only to Codefresh on-premises installations.",
		Read:        dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: *UserSchema(),
				},
			},
		},
	}
}

func dataSourceUsersRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	users, err := client.GetAllUsers()
	if err != nil {
		return err
	}

	err = mapDataUsersToResource(*users, d)
	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())

	return nil
}

func mapDataUsersToResource(users []cfclient.User, d *schema.ResourceData) error {

	var res = make([]map[string]interface{}, len(users))
	for i, user := range users {
		m := make(map[string]interface{})
		m["user_name"] = user.UserName
		m["email"] = user.Email
		m["status"] = user.Status
		if user.Personal != nil {
			m["personal"] = flattenPersonal(user.Personal)
		}
		m["short_profile"] = []map[string]interface{}{
			{"user_name": user.ShortProfile.UserName}}
		m["roles"] = user.Roles
		m["logins"] = flattenLogins(&user.Logins)
		m["user_id"] = user.ID

		res[i] = m
	}

	d.Set("users", res)

	return nil
}
