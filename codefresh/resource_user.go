package codefresh

import (
	"log"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUsersCreate,
		Read:   resourceUsersRead,
		Update: resourceUsersUpdate,
		Delete: resourceUsersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"accounts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"personal": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems:    1,
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
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// "credentials": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	MaxItems:    1,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"permissions": {
						// 				Type:     schema.TypeList,
						// 				Optional: true,
						// 				Elem: &schema.Schema{
						// 					Type: schema.TypeString,
						// 				},
						// 			},
						// 		},
						// 	},
						// },
						"idp_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sso": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceUsersCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	user := mapResourceToUser(d)

	resp, err := client.AddPendingUser(user)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	if d.Get("activate").(bool) {
		client.ActivateUser(d.Id())
	}

	return nil
}

func resourceUsersRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

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

	client := meta.(*cfClient.Client)

	accountList := d.Get("accounts").(*schema.Set).List()

	userId := d.Id()

	accounts, err := client.GetAccountsList(convertStringArr(accountList))
	if err != nil {
		return err
	}

	err = client.UpdateUserAccounts(userId, *accounts)
	if err != nil {
		return err
	}

	return nil
}

func resourceUsersDelete(d *schema.ResourceData, meta interface{}) error {

	// To research
	// it's impossible sometimes to delete user - limit of runtimes or collaborators should be increased.

	client := meta.(*cfClient.Client)

	userName := d.Get("user_name").(string)

	err := client.DeleteUser(userName)
	if err != nil {
		return err
	}

	return nil
}

func mapUserToResource(user cfClient.User, d *schema.ResourceData) error {

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
	d.Set("roles", user.Roles)
	d.Set("login", flattenUserLogins(&user.Logins))

	return nil
}

func flattenUserAccounts(accounts []cfClient.Account) []string {

	var accountList []string

	for _, account := range accounts {
		accountList = append(accountList, account.ID)
	}

	return accountList
}

func flattenUserLogins(logins *[]cfClient.Login) []map[string]interface{} {

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

func mapResourceToUser(d *schema.ResourceData) *cfClient.NewUser {

	roles := d.Get("roles").(*schema.Set).List()
	accounts := d.Get("accounts").(*schema.Set).List()

	user := &cfClient.NewUser{
		ID:       d.Id(),
		UserName: d.Get("user_name").(string),
		Email:    d.Get("email").(string),
		Roles:    convertStringArr(roles),
		Account:  convertStringArr(accounts),
	}

	if _, ok := d.GetOk("personal"); ok {
		user.Personal = &cfClient.Personal{
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
			login := cfClient.Login{
				// Credentials: cfClient.Credentials{
				// 	Permissions: loginData.Get("credentials.permissions").([]string),
				// },
				IDP: cfClient.IDP{
					ID:         idpID,
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

	// 	permissions := convertStringArr(d.Get(fmt.Sprintf("login.%v.credentials.0.permissions", idx)).([]interface{}))
	// 	login := cfClient.Login{
	// 		Credentials: cfClient.Credentials{
	// 			Permissions: permissions,
	// 		},
	// 		Idp: d.Get(fmt.Sprintf("login.%v.idp_id", idx)).(string),
	// 		Sso: d.Get(fmt.Sprintf("login.%v.sso", idx)).(bool),
	// 	}
	// 	user.Logins = append(user.Logins, login)
	// }

	return user
}
