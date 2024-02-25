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

func TestAccCodefreshPipelineCronTriggerValidExpressions(t *testing.T) {
	pipelineName := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline_cron_trigger.test"
	var pipelineCronTrigger cfclient.HermesTrigger

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineCronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				// https://crontab.guru/between-certain-hours
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 9-17 * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 9-17 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
			},
			{
				// https://crontab.guru/daily
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 0 * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 0 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
			},
			{
				// documented example: https://codefresh.io/docs/docs/configure-ci-cd-pipeline/triggers/cron-triggers
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 */20 * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 */20 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
			},
			{
				// Test resource import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccCodefreshPipelineCronTriggerImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodefreshPipelineCronTriggerInvalidExpressions(t *testing.T) {
	pipelineName := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline_cron_trigger.test"
	var pipelineCronTrigger cfclient.HermesTrigger
	expectedError := regexp.MustCompile("The cron expression .* is invalid: .*")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineCronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				// invalid cron field
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "*/1 * * * * test", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * test"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
				ExpectError: expectedError,
			},
			{
				// empty expression
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", ""),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
				ExpectError: expectedError,
			},
			{
				// too few cron fields (Codefresh Cron Triggers have 6 fields, the following only has 5)
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "* * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
				ExpectError: expectedError,
			},
			{
				// too many cron fields (Codefresh Cron Triggers have 6 fields, the following has 7)
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "* * * * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
				ExpectError: expectedError,
			},
			{
				// not a cron expression
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "foo", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "foo"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccCodefreshPipelineCronTriggerValidMessages(t *testing.T) {
	pipelineName := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline_cron_trigger.test"
	var pipelineCronTrigger cfclient.HermesTrigger

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineCronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				// Valid message
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 * * * * *", "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test"),
				),
			},
		},
	})
}

func TestAccCodefreshPipelineCronTriggerInvalidMessages(t *testing.T) {
	pipelineName := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline_cron_trigger.test"
	var pipelineCronTrigger cfclient.HermesTrigger
	expectedError := regexp.MustCompile("The message .* is invalid.*")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineCronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				// Empty message
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 * * * * *", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", ""),
				),
				ExpectError: expectedError,
			},
			{
				// Message contains !
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 * * * * *", "Triggered by cron!"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "Triggered by cron!"),
				),
				ExpectError: expectedError,
			},
			{
				// Message contains ,
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 * * * * *", "Triggered, by cron"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "Triggered, by cron"),
				),
				ExpectError: expectedError,
			},
			{
				// Message with URI
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 * * * * *", "Check out https://go.dev/play/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "https://go.dev/play/"),
				),
				ExpectError: expectedError,
			},
		},
	})
}

// Ensure that after a pipeline cron trigger is updated, the new cron trigger is created and the old one is deleted.
// That is, that this process only results in one cron trigger, and not two. The Hermes API does not have a proper
// update method, so we have to delete and recreate the cron trigger, and if we are not careful, we can end up with
// an additional cron trigger (if we forget to delete the old one).
func TestAccCodefreshPipelineCronTriggerUpdateNoDuplicates(t *testing.T) {
	pipelineName := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline_cron_trigger.test"
	var pipelineCronTrigger cfclient.HermesTrigger

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineCronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 * * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 * * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
			},
			{
				Config: testAccCodefreshPipelineCronTriggerBasicConfig(pipelineName, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "0 1 * * * *", "test message"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineCronTriggerExists(resourceName, &pipelineCronTrigger),
					resource.TestCheckResourceAttr(resourceName, "expression", "0 1 * * * *"),
					resource.TestCheckResourceAttr(resourceName, "message", "test message"),
				),
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
	apiClient := testAccProvider.Meta().(*cfclient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_pipline_cron_trigger" {
			continue
		}

		_, err := apiClient.GetHermesTriggerByEventAndPipeline(rs.Primary.ID, rs.Primary.Attributes["pipeline_id"])

		if err == nil {
			return fmt.Errorf("pipeline Cron Trigger still exists")
		}

		notFoundErr := "PIPELINE_NOT_FOUND_ERROR"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

func testAccCheckCodefreshPipelineCronTriggerExists(resource string, pipelineCronTrigger *cfclient.HermesTrigger) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		apiClient := testAccProvider.Meta().(*cfclient.Client)
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
