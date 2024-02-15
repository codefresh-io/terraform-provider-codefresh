package codefresh

import (
	"fmt"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestIDPCodefreshProject_AccountSpecific(t *testing.T) {
	uniqueId := acctest.RandString(10)
	resourceName := "codefresh_idp.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckCodefreshProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testIDPCodefreshProjectAccountSpecificConfig("onelogin", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshIDPAccountSpecficExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-onelogin-%s", uniqueId)),
				),
			},
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckCodefreshIDPAccountSpecficExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		idpID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetAccountIdpByID(idpID)

		if err != nil {
			return fmt.Errorf("error fetching project with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testIDPCodefreshProjectAccountSpecificConfig(idpType string, uniqueId string) string {

	idpResource := ""

	if idpType == "onelogin" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_idp" "test" { 
			display_name = "tf-test-onelogin-%s"

			onelogin {
				client_id = "onelogin-%s"
				client_secret = "myoneloginsecret1"
				domain = "myonelogindomain"
				app_id = "myappid"
				api_client_id = "myonelogindomain"
				api_client_secret = "myapiclientsecret1"
			}
		}`, uniqueId, uniqueId)
	}

	return idpResource
}