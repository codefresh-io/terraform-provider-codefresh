package codefresh

import (
	"fmt"
	"regexp"
	"testing"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCodefreshPipelineCronTrigger_basic(t *testing.T) {
	pipelineName := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline_cron_trigger.test"
	var pipelineCronTrigger cfClient.HermesTrigger

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineCronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "*/1 * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "*/1 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccCodefreshPipelineCronTriggerImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCodefreshPipelineCronTriggerBasicConfig(rName, repo, path, revision, context, expression, message string) string {
	return testAccCodefreshPipelineBasicConfig(rName, repo, path, revision, context) + fmt.Sprintf(`
resource "codefresh_pipeline_cron_trigger" "test" {
	pipeline_id =  codefresh_pipeline.test.id
	expression = "%s"
	message  = "%s"
  }
`, expression, message)
}

func testAccCheckCodefreshPipelineCronTriggerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfClient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_pipline_cron_trigger" {
			continue
		}

		_, err := apiClient.GetHermesTriggerByEventAndPipeline(rs.Primary.ID, rs.Primary.Attributes["pipeline_id"])

		if err == nil {
			return fmt.Errorf("Pipeline Cron Trigger still exists")
		}

		notFoundErr := "PIPELINE_NOT_FOUND_ERROR"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

func testAccCheckCodefreshPipelineCronTriggerExists(resource string, pipelineCronTrigger *cfClient.HermesTrigger) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		apiClient := testAccProvider.Meta().(*cfClient.Client)
		retrievedHermesTrigger, err := apiClient.GetHermesTriggerByEventAndPipeline(rs.Primary.ID, rs.Primary.Attributes["pipeline_id"])

		if err != nil {
			return fmt.Errorf("error fetching pipeline cron trigger with resource %s. %s", resource, err)
		}

		*pipelineCronTrigger = *retrievedHermesTrigger

		return nil
	}
}

func testAccCodefreshPipelineCronTriggerImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["pipeline_id"]), nil
	}
}