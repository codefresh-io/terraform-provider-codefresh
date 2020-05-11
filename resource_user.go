package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const codefreshAccountID = "123"

// codefresh expects a different json request format when
// creating a user than it returns when reading a user
type userCreate struct {
	UserDetails string `json:"userDetails,omitempty"`
}

func (u *userCreate) getID() string {
	panic("this should never be called. there's no id for userCreate.")
}

type user struct {
	ID    string `json:"_id,omitempty"`
	Email string `json:"email,omitempty"`
}

func (u *user) getID() string {
	return u.ID
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserImport,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, _ interface{}) error {
	return createCodefreshObject(
		fmt.Sprintf("%v/accounts/%s/adduser", getCfUrl(), codefreshAccountID),
		"POST",
		d,
		mapResourceToUser,
		readUserCreate,
	)
}

func resourceUserRead(d *schema.ResourceData, _ interface{}) error {
	return readCodefreshObject(
		d,
		getUserFromCodefresh,
		mapUserToResource)
}

func resourceUserDelete(d *schema.ResourceData, _ interface{}) error {
	url := fmt.Sprintf("%v/accounts/%v/%v", getCfUrl(), codefreshAccountID, d.Id())
	return deleteCodefreshObject(url)
}

func resourceUserImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return importCodefreshObject(
		d,
		getUserFromCodefresh,
		mapUserToResource)
}

func readUserCreate(_ *schema.ResourceData, b []byte) (codefreshObject, error) {
	user := &user{}
	err := json.Unmarshal(b, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getUserFromCodefresh(d *schema.ResourceData) (codefreshObject, error) {
	url := fmt.Sprintf("%v/accounts/%v/users", getCfUrl(), codefreshAccountID)
	return getFromCodefresh(d, url, readUserFromAll)
}

func readUserFromAll(d *schema.ResourceData, b []byte) (codefreshObject, error) {
	users := &[]user{}
	err := json.Unmarshal(b, users)
	if err != nil {
		return nil, err
	}
	return getUserByID(d.Id(), *users)
}

// returns empty user if no user found
// can cause a nil pointer exception?
func getUserByID(Id string, users []user) (codefreshObject, error) {
	for _, user := range users {
		if user.ID == Id {
			return &user, nil
		}
	}
	return &user{}, fmt.Errorf("Could not find user %v", Id)
}

func mapUserToResource(cfObject codefreshObject, d *schema.ResourceData) error {
	user := cfObject.(*user)
	d.SetId(user.ID)

	err := d.Set("email", user.Email)
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToUser(d *schema.ResourceData) codefreshObject {
	user := &userCreate{
		UserDetails: d.Get("email").(string),
	}
	return user
}
