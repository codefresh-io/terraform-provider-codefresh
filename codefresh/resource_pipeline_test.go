package codefresh

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var pipelineNamePrefix = "TerraformAccTest_"

func TestAccCodefreshPipeline_basic(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigConcurrency(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "1", "2", "3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
	var pipeline cfclient.Pipeline

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigTags(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "testTag1", "testTag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigVariables(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "var1", "val1", "var2", "val2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
	runtimeName := "system/default/hybrid/k8s" // must be present in test environment
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigRuntimeEnvironment(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", runtimeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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

func TestAccCodefreshPipeline_OriginalYamlString_Steps(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	originalYamlString := `version: 1.0
steps:
  cc_firstStep:
    image: alpine
    commands:
      - echo Hello World First Step
  bb_secondStep:
    image: alpine
    commands:
      - echo Hello World Second jStep
  aa_secondStep:
    image: alpine
    commands:
      - echo Hello World Third Step`

	expectedSpecAttributes := &cfclient.Spec{
		Steps: &cfclient.Steps{
			Steps: `{"cc_firstStep":{"image":"alpine","commands":["echo Hello World First Step"]},"bb_secondStep":{"image":"alpine","commands":["echo Hello World Second jStep"]},"aa_secondStep":{"image":"alpine","commands":["echo Hello World Third Step"]}}`,
		},
		Stages: &cfclient.Stages{
			Stages: `[]`,
		},
	}

	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigOriginalYamlString(name, originalYamlString),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "original_yaml_string", originalYamlString),
					testAccCheckCodefreshPipelineOriginalYamlStringAttributePropagation(resourceName, expectedSpecAttributes),
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

func TestAccCodefreshPipeline_OriginalYamlString_All(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	originalYamlString := `version: 1.0
fail_fast: false
stages:
  - test
mode: parallel
hooks: 
  on_finish:
    steps:
      secondmycleanup:
        commands:
          - echo echo cleanup step
        image: alpine:3.9
      firstmynotification:
        commands:
          - echo Notify slack
        image: cloudposse/slack-notifier
  on_elected:
    exec:
      commands:
       - echo 'Creating an adhoc test environment'
      image: alpine:3.9
    annotations:
      set:
        - annotations:
            - my_annotation_example1: 10.45
            - my_string_annotation: Hello World
          entity_type: build
steps:
  zz_firstStep:
    stage: test
    image: alpine
    commands:
      - echo Hello World First Step
  aa_secondStep:
    stage: test
    image: alpine
    commands:
      - echo Hello World Second Step`

	expectedSpecAttributes := &cfclient.Spec{
		Steps: &cfclient.Steps{
			Steps: `{"zz_firstStep":{"stage":"test","image":"alpine","commands":["echo Hello World First Step"]},"aa_secondStep":{"stage":"test","image":"alpine","commands":["echo Hello World Second Step"]}}`,
		},
		Stages: &cfclient.Stages{
			Stages: `["test"]`,
		},
		Hooks: &cfclient.Hooks{
			Hooks: `{"on_finish":{"steps":{"secondmycleanup":{"commands":["echo echo cleanup step"],"image":"alpine:3.9"},"firstmynotification":{"commands":["echo Notify slack"],"image":"cloudposse/slack-notifier"}}},"on_elected":{"exec":{"commands":["echo 'Creating an adhoc test environment'"],"image":"alpine:3.9"},"annotations":{"set":[{"annotations":[{"my_annotation_example1":10.45},{"my_string_annotation":"Hello World"}],"entity_type":"build"}]}}}`,
		},
	}

	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigOriginalYamlString(name, originalYamlString),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "original_yaml_string", originalYamlString),
					testAccCheckCodefreshPipelineOriginalYamlStringAttributePropagation(resourceName, expectedSpecAttributes),
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
	var pipeline cfclient.Pipeline

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
					"/^(?!(master)$).*/gi",
					"multiselect",
					"/^(?!(master)$).*/gi",
					"/^PR comment$/gi",
					"shared_context1",
					"git",
					"push.heads",
					"codefresh-contrib/react-sample-app",
					"tags",
					"git",
					"shared_context2",
					true,
					true,
					true,
					true,
					"push.tags",
					"codefresh-contrib/react-sample-app",
					"triggerTestVar",
					"triggerTestValue",
					"commitstatustitle",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.branch_regex", "/^(?!(master)$).*/gi"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.branch_regex_input", "multiselect"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.pull_request_target_branch_regex", "/^(?!(master)$).*/gi"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.comment_regex", "/^PR comment$/gi"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.name", "commits"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.name", "tags"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.contexts.0", "shared_context2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.no_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.no_cf_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.reset_volume", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.enable_notifications", "true"),
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
					"/release/gi",
					"multiselect-exclude",
					"/release/gi",
					"/PR comment2/gi",
					"shared_context1_update",
					"git",
					"push.heads",
					"codefresh-contrib/react-sample-app",
					"tags",
					"git",
					"shared_context2_update",
					true,
					true,
					false,
					false,
					"push.tags",
					"codefresh-contrib/react-sample-app",
					"triggerTestVar",
					"triggerTestValue",
					"commitstatustitle",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.branch_regex", "/release/gi"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.branch_regex_input", "multiselect-exclude"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.pull_request_target_branch_regex", "/release/gi"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.0.comment_regex", "/PR comment2/gi"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.variables.triggerTestVar", "triggerTestValue"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.contexts.0", "shared_context2_update"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.no_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.no_cf_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.reset_volume", "false"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.trigger.1.options.0.enable_notifications", "false"),
				),
			},
		},
	})
}

func TestAccCodefreshPipeline_CronTriggers(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigCronTriggers(
					name,
					"codefresh-contrib/react-sample-app",
					"./codefresh.yml",
					"master",
					"git",
					"cT1",
					"first",
					"0/1 * 1/1 * *",
					"64abd1550f02a62699b10df7",
					"runtime1",
					"100mb",
					"1cpu",
					"1gb",
					"1gb",
					"cT2",
					"second",
					"0/1 * 1/1 * *",
					"64abd1550f02a62699b10df7",
					true,
					true,
					true,
					true,
					"MY_VAR",
					"test",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.message", "first"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.expression", "0/1 * 1/1 * *"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.name", "runtime1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.memory", "100mb"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.cpu", "1cpu"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.dind_storage", "1gb"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.required_available_storage", "1gb"),

					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.message", "second"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.expression", "0/1 * 1/1 * *"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.git_trigger_id", "64abd1550f02a62699b10df7"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.no_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.no_cf_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.reset_volume", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.enable_notifications", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfigCronTriggers(
					name,
					"codefresh-contrib/react-sample-app",
					"./codefresh.yml",
					"master",
					"git",
					"cT1",
					"first-1",
					"0/1 * 1/1 * *",
					"00abd1550f02a62699b10df7",
					"runtime2",
					"500mb",
					"2cpu",
					"2gb",
					"3gb",
					"cT2",
					"second",
					"1/1 * 1/1 * *",
					"00abd1550f02a62699b10df7",
					true,
					true,
					false,
					false,
					"MY_VAR",
					"test",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.name", "cT1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.message", "first-1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.expression", "0/1 * 1/1 * *"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.name", "runtime2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.memory", "500mb"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.cpu", "2cpu"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.dind_storage", "2gb"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.required_available_storage", "3gb"),

					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.name", "cT2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.message", "second"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.expression", "1/1 * 1/1 * *"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.git_trigger_id", "00abd1550f02a62699b10df7"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.no_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.no_cf_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.reset_volume", "false"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.enable_notifications", "false"),
				),
			},
		},
	})
}

// Same config as TestAccCodefreshPipeline_CronTriggers but with invalid cron expression (too many fields)
func TestAccCodefreshPipeline_CronTriggersInvalid(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigCronTriggers(
					name,
					"codefresh-contrib/react-sample-app",
					"./codefresh.yml",
					"master",
					"git",
					"cT1",
					"first",
					"0 0/1 * 1/1 * *",
					"64abd1550f02a62699b10df7",
					"runtime1",
					"100mb",
					"1cpu",
					"1gb",
					"1gb",
					"cT2",
					"second",
					"0 0/1 * 1/1 * *",
					"64abd1550f02a62699b10df7",
					true,
					true,
					true,
					true,
					"MY_VAR",
					"test",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.name", "cT1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.message", "first"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.expression", "0 0/1 * 1/1 * *"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.name", "runtime1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.memory", "100mb"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.cpu", "1cpu"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.dind_storage", "1gb"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.0.runtime_environment.0.required_available_storage", "1gb"),

					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.name", "cT2"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.message", "second"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.expression", "0 0/1 * 1/1 * *"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.git_trigger_id", "64abd1550f02a62699b10df7"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.no_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.no_cf_cache", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.reset_volume", "true"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.cron_trigger.1.options.0.enable_notifications", "true"),
				),
				ExpectError: regexp.MustCompile("The cron expression .* is invalid: Expected exactly 5 fields.*"),
			},
		},
	})
}

func TestAccCodefreshPipeline_Revision(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "revision", "1"),
				),
			},
		},
	})
}

func TestAccCodefreshPipeline_IsPublic(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineIsPublic(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "development", "git", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
				),
			},
		},
	})
}

func TestAccCodefreshPipelineOnCreateBranchIgnoreTrigger(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineOnCreateBranchIgnoreTrigger(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.termination_policy.0.on_create_branch.0.ignore_trigger", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckNoResourceAttr(resourceName, "spec.0.termination_policy.0.on_create_branch"),
				),
			},
		},
	})
}

func TestAccCodefreshPipelineOptions(t *testing.T) {
	name := pipelineNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_pipeline.test"
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineOptions(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "spec.0.options.*", map[string]string{
						"keep_pvcs_for_pending_approval":       "true",
						"pending_approval_concurrency_applied": "false",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshPipelineBasicConfig(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckNoResourceAttr(resourceName, "spec.0.options.#"),
				),
			},
		},
	})
}

func testAccCheckCodefreshPipelineExists(resource string, pipeline *cfclient.Pipeline) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		pipelineID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		retrievedPipeline, err := apiClient.GetPipeline(pipelineID)

		if err != nil {
			return fmt.Errorf("error fetching pipeline with resource %s. %s", resource, err)
		}

		*pipeline = *retrievedPipeline

		return nil
	}
}

func testAccCheckCodefreshPipelineDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfclient.Client)

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

func testAccCheckCodefreshPipelineOriginalYamlStringAttributePropagation(resource string, spec *cfclient.Spec) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		pipelineID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		pipeline, err := apiClient.GetPipeline(pipelineID)

		if !reflect.DeepEqual(pipeline.Spec.Steps, spec.Steps) {
			return fmt.Errorf("Expected Step %v. Got %v", spec.Steps, pipeline.Spec.Steps)
		}
		if !reflect.DeepEqual(pipeline.Spec.Stages, spec.Stages) {
			return fmt.Errorf("Expected Stages %v. Got %v", spec.Stages, pipeline.Spec.Stages)
		}
		if !reflect.DeepEqual(pipeline.Spec.Hooks, spec.Hooks) {
			return fmt.Errorf("Expected Hooks %v. Got %v", spec.Hooks, pipeline.Spec.Hooks)
		}
		if err != nil {
			return fmt.Errorf("error fetching pipeline with resource %s. %s", resource, err)
		}
		return nil
	}
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
	trigger1Regex,
	trigger1RegexInput,
	trigger1PrTargetBranchRegex,
	trigger1CommentRegex,
	trigger1Context,
	trigger1Contexts,
	trigger1Event,
	trigger1Repo,
	trigger2Name,
	trigger2Context,
	trigger2Contexts string,
	trigger2NoCache,
	trigger2NoCfCache,
	trigger2ResetVolume,
	trigger2EnableNotifications bool,
	trigger2Event,
	trigger2Repo,
	trigger2VarName,
	trigger2VarValue,
	trigger2CommitStatusTitle string,
) string {
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
		branch_regex = %q
		branch_regex_input = %q
		pull_request_target_branch_regex = %q
		comment_regex = %q

		context = %q
		contexts = [
			%q
		]
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
        contexts = [
            %q
        ]
        description = ""
        disabled = false
        options {
            no_cache             = %t
            no_cf_cache          = %t
            reset_volume         = %t
            enable_notifications = %t
        }
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
		trigger1Regex,
		trigger1RegexInput,
		trigger1PrTargetBranchRegex,
		trigger1CommentRegex,
		trigger1Context,
		trigger1Contexts,
		trigger1Event,
		trigger1Repo,
		trigger2Name,
		trigger2Context,
		trigger2Contexts,
		trigger2NoCache,
		trigger2NoCfCache,
		trigger2ResetVolume,
		trigger2EnableNotifications,
		trigger2Event,
		trigger2Repo,
		trigger2VarName,
		trigger2VarValue,
		trigger2CommitStatusTitle)
}

func testAccCodefreshPipelineBasicConfigCronTriggers(
	rName,
	repo,
	path,
	revision,
	context,
	cronTrigger1Name string,
	cronTrigger1Message string,
	cronTrigger1Expression string,
	cronTrigger1GitTriggerId string,
	cronTrigger1REName string,
	cronTrigger1REMemory string,
	cronTrigger1RECpu string,
	cronTrigger1REDindStorage string,
	cronTrigger1RERequiredAvailableStorage string,
	cronTrigger2Name string,
	cronTrigger2Message string,
	cronTrigger2Expression string,
	cronTrigger2GitTriggerId string,
	cronTrigger2NoCache bool,
	cronTrigger2NoCfCache bool,
	cronTrigger2ResetVolume bool,
	cronTrigger2EnableNotifications bool,
	cronTrigger2VarName string,
	cronTrigger2VarValue string,
) string {
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

    cron_trigger {
        name = %q
        type = "cron"
        branch = "main"
        message = %q
        expression = %q
        git_trigger_id = %q
        disabled = true
        runtime_environment {
			name = %q
			memory = %q
			cpu = %q
			dind_storage = %q
			required_available_storage = %q
		}
    }

    cron_trigger {
        name = %q
        type = "cron"
        branch = "master"
        message = %q
        expression = %q
        git_trigger_id = %q
        disabled = false
        options {
            no_cache             = %t
            no_cf_cache          = %t
            reset_volume         = %t
            enable_notifications = %t
        }
        variables = {
            %q = %q
		}
    }
  }
}
`,
		rName,
		repo,
		path,
		revision,
		context,
		cronTrigger1Name,
		cronTrigger1Message,
		cronTrigger1Expression,
		cronTrigger1GitTriggerId,
		cronTrigger1REName,
		cronTrigger1REMemory,
		cronTrigger1RECpu,
		cronTrigger1REDindStorage,
		cronTrigger1RERequiredAvailableStorage,
		cronTrigger2Name,
		cronTrigger2Message,
		cronTrigger2Expression,
		cronTrigger2GitTriggerId,
		cronTrigger2NoCache,
		cronTrigger2NoCfCache,
		cronTrigger2ResetVolume,
		cronTrigger2EnableNotifications,
		cronTrigger2VarName,
		cronTrigger2VarValue,
	)
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
	var pipeline cfclient.Pipeline

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPipelineBasicConfigContexts(name, "codefresh-contrib/react-sample-app", "./codefresh.yml", "master", "git", "context1", "context2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
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
					testAccCheckCodefreshPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "spec.0.contexts.0", "context1_updated"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.contexts.1", "context2_updated"),
				),
			},
		},
	})
}

func testAccCodefreshPipelineOptions(rName, repo, path, revision, context string, keepPVCsForPendingApproval, pendingApprovalConcurrencyApplied bool) string {
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
	options {
		keep_pvcs_for_pending_approval = %t
		pending_approval_concurrency_applied = %t
	}
  }
}
`, rName, repo, path, revision, context, keepPVCsForPendingApproval, pendingApprovalConcurrencyApplied)
}

func testAccCodefreshPipelineOnCreateBranchIgnoreTrigger(rName, repo, path, revision, context string, ignoreTrigger bool) string {
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
	termination_policy {
		on_create_branch {
			ignore_trigger = %t
		}
	}
  }
}
`, rName, repo, path, revision, context, ignoreTrigger)
}

func testAccCodefreshPipelineOnCreateBranchIgnoreTriggerWithBranchName(rName, repo, path, revision, context, branchName string, ignoreTrigger bool) string {
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
	termination_policy {
		on_create_branch {
			branch_nane = %q
			ignore_trigger = %t
		}
	}
  }
}
`, rName, repo, path, revision, context, branchName, ignoreTrigger)
}

func testAccCodefreshPipelineIsPublic(rName, repo, path, revision, context string, isPublic bool) string {
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

  is_public = %t

}
`, rName, repo, path, revision, context, isPublic)
}
