package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccountUserAssociation() *schema.Resource {
	return &schema.Resource{
		Description: `
		Associates a user with the account which the provider is authenticated against. If the user is not present in the system, an invitation will be sent to the specified email address.
		`,
		Create: resourceAccountUserAssociationCreate,
		Read:   resourceAccountUserAssociationRead,
		Update: resourceAccountUserAssociationUpdate,
		Delete: resourceAccountUserAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Description: `
				The email of the user to associate with the specified account.
				If the user is not present in the system, an invitation will be sent to this email.
				This field can only be changed when 'status' is 'pending'.
				`,
				Type:     schema.TypeString,
				Required: true,
			},
			"admin": {
				Description: "Whether to make this user an account admin.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"username": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The username of the associated user.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The status of the association.",
			},
		},
		CustomizeDiff: customdiff.All(
			// The email field is immutable, except for users with status "pending".
			customdiff.ForceNewIf("email", func(_ context.Context, d *schema.ResourceDiff, _ any) bool {
				return d.Get("status").(string) != "pending" && d.HasChange("email")
			}),
		),
	}
}

func resourceAccountUserAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)
	currentAccount, err := client.GetCurrentAccount()
	if err != nil {
		return err
	}

	user, err := client.AddNewUserToAccount(currentAccount.ID, "", d.Get("email").(string))
	if err != nil {
		return err
	}

	d.SetId(user.ID)

	if d.Get("admin").(bool) {
		err = client.SetUserAsAccountAdmin(currentAccount.ID, d.Id())
		if err != nil {
			return err
		}
	}

	d.Set("status", user.Status)

	return nil
}

func resourceAccountUserAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)
	currentAccount, err := client.GetCurrentAccount()
	if err != nil {
		return err
	}

	userID := d.Id()
	if userID == "" {
		d.SetId("")
		return nil
	}

	for _, user := range currentAccount.Users {
		if user.ID == userID {
			d.Set("email", user.Email)
			d.Set("username", user.UserName)
			d.Set("status", user.Status)
			d.Set("admin", false) // avoid missing attributes after import
			for _, admin := range currentAccount.Admins {
				if admin.ID == userID {
					d.Set("admin", true)
				}
			}
		}
	}

	if d.Id() == "" {
		return fmt.Errorf("a user with ID %s was not found", userID)
	}

	return nil
}

func resourceAccountUserAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	currentAccount, err := client.GetCurrentAccount()
	if err != nil {
		return err
	}

	if d.HasChange("email") {
		user, err := client.UpdateUserDetails(currentAccount.ID, d.Id(), d.Get("username").(string), d.Get("email").(string))
		if err != nil {
			return err
		}
		if user.Email != d.Get("email").(string) {
			return fmt.Errorf("failed to update user email, despite successful API response")
		}
	}

	if d.HasChange("admin") {
		if d.Get("admin").(bool) {
			err = client.SetUserAsAccountAdmin(currentAccount.ID, d.Id())
			if err != nil {
				return err
			}
		} else {
			err = client.DeleteUserAsAccountAdmin(currentAccount.ID, d.Id())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceAccountUserAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	currentAccount, err := client.GetCurrentAccount()
	if err != nil {
		return err
	}

	err = client.DeleteUserFromAccount(currentAccount.ID, d.Id())
	if err != nil {
		return err
	}

	return nil
}
