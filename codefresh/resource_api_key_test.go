package codefresh

import (
	"fmt"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var apiKeyNamePrefix = "TerraformAccTest_"

func TestAccCodefreshAPIKey_ServiceUser(t *testing.T) {
	name := apiKeyNamePrefix + acctest.RandString(10)

	resourceName := "codefresh_api_key.test_apikey"
	serviceAccountResourceName := "codefresh_service_account.test_apikey"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshServiceUserAndAPIKeyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshAPIKeyServiceAccount(name, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshServiceUserAPIKeyExists(resourceName, serviceAccountResourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "scopes.0", "agent"),
				),
			},
			{
				ResourceName: resourceName,
				RefreshState: true,
			},
		},
	})
}

func testAccCheckCodefreshServiceUserAPIKeyExists(apiKeyResource string, serviceUserResource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		serviceUserState, ok := state.RootModule().Resources[serviceUserResource]

		if !ok {
			return fmt.Errorf("Not found: %s", serviceUserResource)
		}

		if serviceUserState.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		apiKeyState, ok := state.RootModule().Resources[apiKeyResource]

		if !ok {
			return fmt.Errorf("Not found: %s", apiKeyResource)
		}

		if apiKeyState.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for team")
		}

		serviceUserID := serviceUserState.Primary.ID
		apiKeyID := apiKeyState.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetAPIKeyServiceUser(apiKeyID, serviceUserID)

		if err != nil {
			return fmt.Errorf("error fetching service user api key for resource %s. %s", apiKeyID, err)
		}

		return nil
	}
}

func testAccCheckCodefreshServiceUserAndAPIKeyDestroyed(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfclient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_service_account" && rs.Type != "codefresh_api_key" {
			continue
		}

		var (
			serviceAccountId string
			apiKeyId         string
		)

		if rs.Type == "codefresh_service_account" {
			serviceAccountId = rs.Primary.ID
			_, err := apiClient.GetServiceUserByID(serviceAccountId)

			if err == nil {
				return fmt.Errorf("Alert service account still exists")
			}
		}

		if rs.Type == "codefresh_api_key" {
			apiKeyId = rs.Primary.ID
			_, err := apiClient.GetAPIKeyServiceUser(apiKeyId, serviceAccountId)

			if err == nil {
				return fmt.Errorf("Alert api key still exists")
			}
		}
	}

	return nil
}

func testAccCodefreshAPIKeyServiceAccount(apiKeyName string, serviceUserName string) string {
	return fmt.Sprintf(`
resource "codefresh_service_account" "test_apikey" {
  name = "%s"
}

resource "codefresh_api_key" "test_apikey" {
  service_account_id = codefresh_service_account.test_apikey.id
  name = "%s"
  scopes = [
    "agent",
    "agents",
    "audit",
    "api-keys"
  ]
}


`, serviceUserName, apiKeyName)
}
