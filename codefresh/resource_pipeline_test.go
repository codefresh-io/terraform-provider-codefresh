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

var pipelineNamePrefix = "TerraformAccTest_"

func TestAccCodefreshPipeline_basic(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.spec_template.0.revision", "master"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.spec_template.0.context", "git"),
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

func TestAccCodefreshPipeline_Concurrency(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigConcurrency(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "1", "2", "3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.concurrency", "1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.branch_concurrency", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger_concurrency", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfigConcurrency(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "4", "5", "6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.concurrency", "4"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.branch_concurrency", "5"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger_concurrency", "6"),
				),
			},
		},
	})
}

func TestAccCodefreshPipeline_Tags(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigTags(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "testTag1", "testTag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "testTag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "testTag1"),
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

func TestAccCodefreshPipeline_Variables(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigVariables(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "var1", "val1", "var2", "val2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.variables.var1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.variables.var2", "val2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfigVariables(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "var1", "val1_updated", "var2", "val2_updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.variables.var1", "val1_updated"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.variables.var2", "val2_updated"),
				),
			},
		},
	})
}

func TestAccCodefreshPipeline_RuntimeEnvironment(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	runtimeName := "system/default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigRuntimeEnvironment(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", runtimeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.runtime_environment.0.name", runtimeName),
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

func TestAccCodefreshPipeline_OriginalYamlString(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	originalYamlString := "version: \"1.0\"\nsteps:\n  test:\n    image: alpine:latest\n    commands:\n      - echo \"ACC tests\""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigOriginalYamlString(name, originalYamlString),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "original_yaml_string", originalYamlString),
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

func TestAccCodefreshPipeline_Triggers(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigTriggers(
					name,
					"codefresh-contrib/react-sample-app",
					"./codefresh.yml",
					"master",
					"git",
					"commits",
					"git",
					"push.heads",
					"codefresh-contrib/react-sample-app",
					"tags",
					"git",
					"push.tags",
					"codefresh-contrib/react-sample-app",
					"triggerTestVar",
					"triggerTestValue",
					"commitstatustitle",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.name", "commits"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.name", "tags"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfigTriggers(
					name,
					"codefresh-contrib/react-sample-app",
					"./codefresh.yml",
					"master",
					"git",
					"commits",
					"git",
					"push.heads",
					"codefresh-contrib/react-sample-app",
					"tags",
					"git",
					"push.tags",
					"codefresh-contrib/react-sample-app",
					"triggerTestVar",
					"triggerTestValue",
					"commitstatustitle",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.variables.triggerTestVar", "triggerTestValue"),
				),
			},
		},
	})
}

func TestAccCodefreshPipeline_Revision(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "revision", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "development", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "revision", "1"),
				),
			},
		},
	})
}

func testAccCheckCodefreshPipelineExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		pipelineID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfClient.Client)
		_, err := apiClient.GetPipeline(pipelineID)

		if err != nil {
			return fmt.Errorf("error fetching pipeline with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCodefreshPipelineDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfClient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_pipeline" {
			continue
		}

		_, err := apiClient.GetPipeline(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}

		notFoundErr := "PIPELINE_NOT_FOUND_ERROR"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

// CONFIGS
func testAccCodefreshPipelineBasicConfig(rName, repo, path, revision, context string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
	spec_template {
    	repo        = %q
    	path        = %q
    	revision    = %q
    	context     = %q
    }
  }
}
`, rName, repo, path, revision, context)
}

func testAccCodefreshPipelineBasicConfigTags(rName, repo, path, revision, context, tag1, tag2 string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
	spec_template {
    	repo        = %q
    	path        = %q
    	revision    = %q
    	context     = %q
    }
  }

  tags = [
	  %q,
	  %q
  ]
}
`, rName, repo, path, revision, context, tag1, tag2)
}

func testAccCodefreshPipelineBasicConfigVariables(rName, repo, path, revision, context, var1Name, var1Value, var2Name, var2Value string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
	spec_template {
    	repo        = %q
    	path        = %q
    	revision    = %q
    	context     = %q
	}

	variables = {
		%q = %q
		%q = %q
	}
  }
}
`, rName, repo, path, revision, context, var1Name, var1Value, var2Name, var2Value)
}

func testAccCodefreshPipelineBasicConfigContexts(rName, repo, path, revision, context, sharedContext1, sharedContext2 string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
	spec_template {
    	repo        = %q
    	path        = %q
    	revision    = %q
    	context     = %q
	}

	contexts = [
		%q,
		%q
	]

  }
}
`, rName, repo, path, revision, context, sharedContext1, sharedContext2)
}

func testAccCodefreshPipelineBasicConfigConcurrency(rName, repo, path, revision, context, concurrency, concurrencyBranch, concurrencyTrigger string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
	spec_template {
    	repo        = %q
    	path        = %q
    	revision    = %q
    	context     = %q
	}

	concurrency 	    = %q
	branch_concurrency  = %q
	trigger_concurrency = %q

  }
}
`, rName, repo, path, revision, context, concurrency, concurrencyBranch, concurrencyTrigger)
}

func testAccCodefreshPipelineBasicConfigTriggers(
	rName,
	repo,
	path,
	revision,
	context,
	trigger1Name,
	trigger1Context,
	trigger1Event,
	trigger1Repo,
	trigger2Name,
	trigger2Context,
	trigger2Event,
	trigger2Repo,
	trigger2VarName,
	trigger2VarValue,
	trigger2CommitStatusTitle string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
        spec_template {
            repo        = %q
            path        = %q
            revision    = %q
            context     = %q
        }

        trigger {
        name = %q
        branch_regex = "/.*/gi"
        context = %q
        description = ""
        disabled = false
        events = [
          %q
        ]
        modified_files_glob = ""
        provider = "github"
        repo = %q
        type = "git"
    }

    trigger {
        name = %q
        branch_regex = "/.*/gi"
        context = %q
        description = ""
        disabled = false
        events = [
          %q
        ]
        modified_files_glob = ""
        pull_request_allow_fork_events = true
        provider = "github"
        repo = %q
        type = "git"

        variables = {
            %q = %q
		}

		commit_status_title = "%s"
    }
  }
}
`,
		rName,
		repo,
		path,
		revision,
		context,
		trigger1Name,
		trigger1Context,
		trigger1Event,
		trigger1Repo,
		trigger2Name,
		trigger2Context,
		trigger2Event,
		trigger2Repo,
		trigger2VarName,
		trigger2VarValue,
		trigger2CommitStatusTitle)
}

func testAccCodefreshPipelineBasicConfigRuntimeEnvironment(rName, repo, path, revision, context, runtimeName string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  spec {
	spec_template {
		repo        = %q
		path        = %q
		revision    = %q
		context     = %q
	}

	runtime_environment {
		name = %q
	}
  }
}
`, rName, repo, path, revision, context, runtimeName)
}

func testAccCodefreshPipelineBasicConfigOriginalYamlString(rName, originalYamlString string) string {
	return fmt.Sprintf(`
resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name = "%s"

  original_yaml_string = %#v

  spec {}

}
`, rName, originalYamlString)
}

func TestAccCodefreshPipeline_Contexts(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigContexts(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "context1", "context2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.contexts.0", "context1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.contexts.1", "context2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfigContexts(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "context1_updated", "context2_updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec.0.contexts.0", "context1_updated"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.contexts.1", "context2_updated"),
				),
			},
		},
	})
}
