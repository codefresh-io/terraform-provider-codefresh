package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"regexp"
	"testing"
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
					resource.TestCheckResourceAttr(resourceName, "tags.3247715412", "testTag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.3938019223", "testTag1"),
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
