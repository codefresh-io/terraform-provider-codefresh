package codefresh

import (
	"fmt"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/gitops"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCodefreshAccountGitopsSettings_basic(t *testing.T) {
	resourceName := "codefresh_account_gitops_settings.test"

	expectedDefaultApiURLGithub, _ := gitops.GetDefaultAPIUrlForProvider(gitops.GitProviderGitHub)
	//expectedDefaultApiURLGithub := "https://bnlah.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckCodefreshProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccountGitopsSettingsGithubDefaultApiUrl("https://github.com/codefresh-io/terraform-provider-isc-test.git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitopsSettings(resourceName, gitops.GitProviderGitHub, *expectedDefaultApiURLGithub, "https://github.com/codefresh-io/terraform-provider-isc-test.git"),
					resource.TestCheckResourceAttr(resourceName, "git_provider", gitops.GitProviderGitHub),
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

func testAccCheckGitopsSettings(resource string, gitProvider string, gitProviderApiUrl string, sharedConfigRepository string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		apiClient := testAccProvider.Meta().(*cfclient.Client)

		accGitopsInfo, err := apiClient.GetActiveGitopsAccountInfo()

		if err != nil {
			return fmt.Errorf("failed getting gitops settings with error %s", err)
		}

		if accGitopsInfo.GitApiUrl != gitProviderApiUrl {
			return fmt.Errorf("expecting APIUrl to be %s but got %s", gitProviderApiUrl, accGitopsInfo.GitApiUrl)
		}

		if accGitopsInfo.GitProvider != gitProvider {
			return fmt.Errorf("expecting provider to be %s but got %s", gitProvider, accGitopsInfo.GitProvider)
		}

		if accGitopsInfo.SharedConfigRepo != sharedConfigRepository {
			return fmt.Errorf("expecting shared config repository to be %s but got %s", sharedConfigRepository, accGitopsInfo.SharedConfigRepo)
		}

		return nil
	}
}

// CONFIGS
func testAccountGitopsSettingsGithubDefaultApiUrl(sharedConfigRepository string) string {
	return fmt.Sprintf(`
	resource "codefresh_account_gitops_settings" "test" {
		git_provider = "GITHUB"
		shared_config_repository = "%s"
	  }`, sharedConfigRepository)
}
