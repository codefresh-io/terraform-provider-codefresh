package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "This data source retrieves a user by email. Requires a Codefresh admin token and applies only to Codefresh on-premises installations.",
		Read:        dataSourceUserRead,
		Schema:      *UserSchema(),
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	users, err := client.GetAllUsers()
	if err != nil {
		return err
	}

	email := d.Get("email").(string)

	for _, user := range *users {
		if user.Email == email {
			err = mapDataUserToResource(user, d)
			if err != nil {
				return err
			}
		}
	}

	if d.Id() == "" {
		return fmt.Errorf("[EROOR] User %s wasn't found", email)
	}

	return nil
}

func mapDataUserToResource(user cfclient.User, d *schema.ResourceData) error {

	d.SetId(user.ID)
	err := d.Set("user_id", user.ID)

	if err != nil {
		return err
	}

	err = d.Set("user_name", user.UserName)

	if err != nil {
		return err
	}

	err = d.Set("email", user.Email)

	if err != nil {
		return err
	}

	err = d.Set("status", user.Status)

	if err != nil {
		return err
	}

	if user.Personal != nil {
		err = d.Set("personal", flattenPersonal(user.Personal))

		if err != nil {
			return err
		}
	}
	err = d.Set("short_profile",
		[]map[string]interface{}{
			{"user_name": user.ShortProfile.UserName},
		})

	if err != nil {
		return err
	}

	err = d.Set("roles", user.Roles)

	if err != nil {
		return err
	}

	err = d.Set("logins", flattenLogins(&user.Logins))

	if err != nil {
		return err
	}

	return nil
}

func flattenPersonal(personal *cfclient.Personal) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"first_name":   personal.FirstName,
			"last_name":    personal.LastName,
			"company_name": personal.CompanyName,
			"phone_number": personal.PhoneNumber,
			"country":      personal.Country,
		},
	}
}

func flattenLogins(logins *[]cfclient.Login) []map[string]interface{} {

	var res = make([]map[string]interface{}, len(*logins))
	for i, login := range *logins {
		m := make(map[string]interface{})
		m["credentials"] = []map[string]interface{}{
			{"permissions": login.Credentials.Permissions}}

		m["idp"] = []map[string]interface{}{
			{
				"id":          login.IDP.ID,
				"client_type": login.IDP.ClientType,
			},
		}

		res[i] = m
	}

	return res
}

func UserSchema() *map[string]*schema.Schema {
	return &map[string]*schema.Schema{
		"user_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"user_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"personal": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"first_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"last_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"company_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"phone_number": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"country": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"short_profile": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"roles": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"logins": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"credentials": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"permissions": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
					"idp": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"client_type": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
