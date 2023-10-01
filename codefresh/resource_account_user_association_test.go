package codefresh

import (
	"fmt"
	"testing"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCodefreshAccountUserAssociationGenerateUserEmail() string {
	randomEmailFormat := "terraform-test-user+%s@codefresh.io" // use + addressing
	return fmt.Sprintf(randomEmailFormat, acctest.RandString(10))
}

func testAccCodefreshActivateUser(s *terraform.State, email string) error {
	c := testAccProvider.Meta().(*cfClient.Client)
	currentAccount, err := c.GetCurrentAccount()
	if err != nil {
		return fmt.Errorf("failed to get current account: %s", err)
	}
	for _, user := range currentAccount.Users {
		if user.Email == email {
			_, err = c.ActivateUser(user.ID)
			if err != nil {
				return fmt.Errorf("failed to activate user: %s", err)
			}
		}
	}
	return nil
}

func TestAccCodefreshAccountUserAssociation_Activation(t *testing.T) {
	resourceName := "codefresh_account_user_association.test_user"

	testUserEmail := testAccCodefreshAccountUserAssociationGenerateUserEmail()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshAccountUserAssociationConfig(testUserEmail, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", testUserEmail),
					resource.TestCheckResourceAttr(resourceName, "admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
				),
			},
			{
				RefreshState: true,
				Check: func(s *terraform.State) error {
					return testAccCodefreshActivateUser(s, testUserEmail)
				},
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "new"),
				),
			},
			{
				// Test resource import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodefreshAccountUserAssociation_StatusPending_Email_ForceNew(t *testing.T) {
	resourceName := "codefresh_account_user_association.test_user"

	testUserEmail := testAccCodefreshAccountUserAssociationGenerateUserEmail()
	var resourceId string
	var err error

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshAccountUserAssociationConfig(testUserEmail, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", testUserEmail),
					resource.TestCheckResourceAttr(resourceName, "admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
				),
			},
			{
				RefreshState: true,
				Check: func(s *terraform.State) error {
					resourceId, err = testAccGetResourceId(s, resourceName)
					return err
				},
			},
			{
				// Test that an email change on a pending user does NOT force a new resource
				PreConfig: func() {
					testUserEmail = testAccCodefreshAccountUserAssociationGenerateUserEmail()
				},
				Config: testAccCodefreshAccountUserAssociationConfig(testUserEmail, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", testUserEmail),
					resource.TestCheckResourceAttr(resourceName, "admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
					func(s *terraform.State) error {
						newResourceId, err := testAccGetResourceId(s, resourceName)
						if err != nil {
							return err
						}
						if resourceId != newResourceId {
							return fmt.Errorf("did not expect email change on pending user to force a new resource")
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccCodefreshAccountUserAssociation_StatusNew_Email_ForceNew(t *testing.T) {
	resourceName := "codefresh_account_user_association.test_user"

	testUserEmail := testAccCodefreshAccountUserAssociationGenerateUserEmail()
	var resourceId string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshAccountUserAssociationConfig(testUserEmail, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", testUserEmail),
					resource.TestCheckResourceAttr(resourceName, "admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
				),
			},
			{
				RefreshState: true,
				Check: func(s *terraform.State) error {
					_, err := testAccGetResourceId(s, resourceName)
					return err
				},
			},
			{
				RefreshState: true,
				Check: func(s *terraform.State) error {
					return testAccCodefreshActivateUser(s, testUserEmail)
				},
			},
			{
				// Test that an email change on an activated user DOES force a new resource
				PreConfig: func() {
					testUserEmail = testAccCodefreshAccountUserAssociationGenerateUserEmail()
				},
				Config: testAccCodefreshAccountUserAssociationConfig(testUserEmail, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", testUserEmail),
					resource.TestCheckResourceAttr(resourceName, "admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "new"),
					func(s *terraform.State) error {
						newResourceId, err := testAccGetResourceId(s, resourceName)
						if err != nil {
							return err
						}
						if resourceId == newResourceId {
							return fmt.Errorf("expected email change on activated user to force a new resource")
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccCodefreshAccountUserAssociationConfig(email string, admin bool) string {
	return fmt.Sprintf(`
resource "codefresh_account_user_association" "test_user" {
	email = "%s"
	admin = %t
}
`, email, admin)
}
