package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUsersRead,
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

	client := meta.(*cfClient.Client)

	users, err := client.GetAllUsers()
	if err != nil {
		return err
	}

	err = mapDataUsersToResource(*users, d)

	d.SetId(time.Now().UTC().String())

	return nil
}

func mapDataUsersToResource(users []cfClient.User, d *schema.ResourceData) error {

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
