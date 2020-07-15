package codefresh

import (
	"time"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUsersRead,

		Schema: map[string]*schema.Schema{

			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	users, err := client.ListUsers()
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
			{"user_name": user.ShortProfile.UserName},}
		m["roles"] = user.Roles
		m["logins"] = flattenLogins(&user.Logins)
		m["id"] = user.ID

		res[i] = m
	}

	d.Set("users", res)

	return nil
}

func flattenPersonal(personal *cfClient.Personal) []map[string]interface{} {
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

func flattenLogins(logins *[]cfClient.Login) []map[string]interface{} {

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