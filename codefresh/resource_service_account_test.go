package codefresh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var serviceUserNamePrefix = "TerraformAccTest_"

func TestAccCodefreshServiceUser_WithTeamAssignment(t *testing.T) {
	name := serviceUserNamePrefix + acctest.RandString(10)

	resourceName := "codefresh_service_account.test_serviceaccount"
	teamResourceName := "codefresh_team.test_serviceaccount"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshServiceUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshServiceUserTeam(name, name, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshServiceUserExists(resourceName),
					testAccCheckCodefreshServiceUserAssignedToTeam(resourceName, teamResourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "assign_admin_role", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCodefreshServiceUserExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		serviceUserID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetServiceUserByID(serviceUserID)

		if err != nil {
			return fmt.Errorf("error fetching serviceUser with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCodefreshServiceUserAssignedToTeam(serviceUserResource string, teamResource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		serviceUserState, ok := state.RootModule().Resources[serviceUserResource]

		if !ok {
			return fmt.Errorf("Not found: %s", serviceUserResource)
		}

		if serviceUserState.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		teamState, ok := state.RootModule().Resources[teamResource]

		if !ok {
			return fmt.Errorf("Not found: %s", teamResource)
		}

		if teamState.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for team")
		}

		serviceUserID := serviceUserState.Primary.ID
		teamID := teamState.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		serviceUser, err := apiClient.GetServiceUserByID(serviceUserID)

		if err != nil {
			return fmt.Errorf("error fetching serviceUser with resource %s. %s", serviceUserID, err)
		}

		isTeamAssigned := false

		for _, team := range serviceUser.Teams {
			if team.ID == teamID {
				isTeamAssigned = true
				break
			}
		}

		if !isTeamAssigned {
			return fmt.Errorf("service user %s is not assigned to team %s", serviceUserID, teamID)
		}

		return nil
	}
}

func testAccCheckCodefreshServiceUserDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfclient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_service_account" {
			continue
		}

		_, err := apiClient.GetServiceUserByID(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "does not exist"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

func testAccCodefreshServiceUserTeam(serviceUserName string, teamName string, assignAdminRole bool) string {
	return fmt.Sprintf(`
resource "codefresh_team" "test_serviceaccount" {
  name = "%s"
}

resource "codefresh_service_account" "test_serviceaccount" {
  name = "%s"
  assigned_teams = [codefresh_team.test_serviceaccount.id]
  assign_admin_role = %t

}
`, serviceUserName, teamName, assignAdminRole)
}
