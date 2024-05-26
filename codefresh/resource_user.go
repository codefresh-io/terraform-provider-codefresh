package codefresh

import (
	"errors"
	"log"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "This resource is used to manage a Codefresh user.",
		Create:      resourceUsersCreate,
		Read:        resourceUsersRead,
		Update:      resourceUsersUpdate,
		Delete:      resourceUsersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Description: "The username of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "Password - for users without SSO.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"has_password": {
				Description: "Whether the user has a local password.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"email": {
				Description: "The email of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"activate": {
				Description: "Whether to activate the user or to leave it as `pending`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"accounts": {
				Description: "A list of accounts IDs to assign the user to.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"personal": {
				Description: "Personal information about the user.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Description: "The first name of the user.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"last_name": {
							Description: "The last name of the user.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"company_name": {
							Description: "The company name of the user.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"phone_number": {
							Description: "The phone number of the user.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"country": {
							Description: "The country of the user.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"short_profile": {
				Description: "The computed short profile of the user.",
				Type:        schema.TypeList,
				Computed:    true,
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
				Description: "The roles of the user.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Description: "The status of the user (e.g. `new`, `pending`).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"login": {
				Description: "Login settings for the user.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idp_id": {
							Description: "The IdP ID for the user's login.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"sso": {
							Description: "Whether to enforce SSO for the user.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceUsersCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	user := mapResourceToNewUser(d)

	resp, err := client.AddPendingUser(user)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	// Adding user to users teams
	for _, accountID := range user.Account {
		_ = client.AddUserToTeamByAdmin(resp.ID, accountID, "users")
	}

	if d.Get("activate").(bool) {
		client.ActivateUser(d.Id())
	}

	if d.Get("password") != "" {
		client.UpdateLocalUserPassword(d.Get("user_name").(string), d.Get("password").(string))
	}

	return resourceUsersRead(d, meta)
}

func resourceUsersRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	userId := d.Id()

	if userId == "" {
		d.SetId("")
		return nil
	}

	user, err := client.GetUserByID(userId)
	if err != nil {
		return err
	}

	err = mapUserToResource(*user, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceUsersUpdate(d *schema.ResourceData, meta interface{}) error {
	// only accounts list

	client := meta.(*cfclient.Client)

	accountList := d.Get("accounts").(*schema.Set).List()

	userId := d.Id()

	accounts, err := client.GetAccountsList(datautil.ConvertStringArr(accountList))
	if err != nil {
		return err
	}

	err = client.UpdateUserAccounts(userId, *accounts)
	if err != nil {
		return err
	}

	// Adding user to users teams
	for _, account := range *accounts {
		_ = client.AddUserToTeamByAdmin(userId, account.ID, "users")
	}

	// Update local password
	err = updateUserLocalPassword(d, client)

	if err != nil {
		return err
	}

	return resourceUsersRead(d, meta)
}

func resourceUsersDelete(d *schema.ResourceData, meta interface{}) error {

	// To research
	// it's impossible sometimes to delete user - limit of runtimes or collaborators should be increased.

	client := meta.(*cfclient.Client)

	userName := d.Get("user_name").(string)

	err := client.DeleteUser(userName)
	if err != nil {
		return err
	}

	return nil
}

func mapUserToResource(user cfclient.User, d *schema.ResourceData) error {

	d.Set("user_name", user.UserName)
	d.Set("email", user.Email)
	d.Set("accounts", flattenUserAccounts(user.Account))
	d.Set("status", user.Status)
	if user.Personal != nil {
		d.Set("personal", flattenPersonal(user.Personal))
	}
	d.Set("short_profile",
		[]map[string]interface{}{
			{"user_name": user.ShortProfile.UserName},
		})
	d.Set("has_password", user.PublicProfile.HasPassword)
	d.Set("roles", user.Roles)
	d.Set("login", flattenUserLogins(&user.Logins))

	return nil
}

func flattenUserAccounts(accounts []cfclient.Account) []string {

	var accountList []string

	for _, account := range accounts {
		accountList = append(accountList, account.ID)
	}

	return accountList
}

func flattenUserLogins(logins *[]cfclient.Login) []map[string]interface{} {

	var res = make([]map[string]interface{}, len(*logins))
	for i, login := range *logins {
		m := make(map[string]interface{})
		// m["credentials"] = []map[string]interface{}{
		// 	{"permissions": login.Credentials.Permissions},
		// }

		m["idp_id"] = login.IDP.ID
		m["sso"] = login.Sso

		res[i] = m
	}

	return res
}

func mapResourceToNewUser(d *schema.ResourceData) *cfclient.NewUser {

	roles := d.Get("roles").(*schema.Set).List()
	accounts := d.Get("accounts").(*schema.Set).List()

	user := &cfclient.NewUser{
		ID:       d.Id(),
		UserName: d.Get("user_name").(string),
		Email:    d.Get("email").(string),
		Roles:    datautil.ConvertStringArr(roles),
		Account:  datautil.ConvertStringArr(accounts),
	}

	if _, ok := d.GetOk("personal"); ok {
		user.Personal = &cfclient.Personal{
			FirstName:   d.Get("personal.0.first_name").(string),
			LastName:    d.Get("personal.0.last_name").(string),
			CompanyName: d.Get("personal.0.company_name").(string),
			PhoneNumber: d.Get("personal.0.phone_number").(string),
			Country:     d.Get("personal.0.country").(string),
		}
	}

	if logins, ok := d.GetOk("login"); ok {
		loginsList := logins.(*schema.Set).List()
		for _, loginDataI := range loginsList {
			if loginData, isMap := loginDataI.(map[string]interface{}); isMap {
				idpID := loginData["idp_id"].(string)
				login := cfclient.Login{
					// Credentials: cfclient.Credentials{
					// 	Permissions: loginData.Get("credentials.permissions").([]string),
					// },
					IDP: cfclient.IDP{
						ID: idpID,
					},
					Sso: loginData["sso"].(bool),
				}
				user.Logins = append(user.Logins, login)
				log.Printf("[DEBUG] login = %v", login)
			}
		}
	}
	// logins := d.Get("login").(*schema.Set)

	// for idx := range logins {

	// 	permissions := datautil.ConvertStringArr(d.Get(fmt.Sprintf("login.%v.credentials.0.permissions", idx)).([]interface{}))
	// 	login := cfclient.Login{
	// 		Credentials: cfclient.Credentials{
	// 			Permissions: permissions,
	// 		},
	// 		Idp: d.Get(fmt.Sprintf("login.%v.idp_id", idx)).(string),
	// 		Sso: d.Get(fmt.Sprintf("login.%v.sso", idx)).(bool),
	// 	}
	// 	user.Logins = append(user.Logins, login)
	// }

	return user
}

func updateUserLocalPassword(d *schema.ResourceData, client *cfclient.Client) error {

	if d.HasChange("password") {
		hasPassword := d.Get("has_password").(bool)

		if _, ok := d.GetOk("user_name"); !ok {
			return errors.New("cannot update password as username attribute is not set")
		}

		userName := d.Get("user_name").(string)

		if password := d.Get("password"); password != "" {
			err := client.UpdateLocalUserPassword(userName, password.(string))

			if err != nil {
				return err
			}
			// If password is not set but has_password returns true, it means that it was removed
		} else if hasPassword {
			err := client.DeleteLocalUserPassword(userName)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
