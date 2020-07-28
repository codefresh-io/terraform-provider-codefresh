package codefresh

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

var projectNamePrefix = "TerraformAccTest_"

func TestAccCodefreshProject_basic(t *testing.T) {
	name := projectNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshProjectBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshProjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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

func TestAccCodefreshProject_Tags(t *testing.T) {
	name := projectNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshProjectBasicConfigTags(name, "testTag1", "testTag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshProjectExists(resourceName),
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

func TestAccCodefreshProject_Variables(t *testing.T) {
	name := projectNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshProjectBasicConfigVariables(name, "var1", "val1", "var2", "val2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshProjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "variables.var1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "variables.var2", "val2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodefreshProjectBasicConfigVariables(name, "var1", "val1_updated", "var2", "val2_updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshProjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "variables.var1", "val1_updated"),
					resource.TestCheckResourceAttr(resourceName, "variables.var2", "val2_updated"),
				//   resource.TestCheckResourceAttr(resourceName, "variables.", name),
				),
			},
		},
	})
}

func testAccCheckCodefreshProjectExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		projectID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfClient.Client)
		_, err := apiClient.GetProjectByID(projectID)

		if err != nil {
			return fmt.Errorf("error fetching project with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCodefreshProjectDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfClient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_project" {
			continue
		}

		_, err := apiClient.GetProjectByID(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

// CONFIGS
func testAccCodefreshProjectBasicConfig(rName string) string {
	return fmt.Sprintf(` 
resource "codefresh_project" "test" { 
  name = "%s" 
} 
`, rName)
}

func testAccCodefreshProjectBasicConfigTags(rName, tag1, tag2 string) string {
	return fmt.Sprintf(`
resource "codefresh_project" "test" { 
  name = "%s" 
  tags = [
	%q,
    %q,
  ]
} 
`, rName, tag1, tag2)
}

func testAccCodefreshProjectBasicConfigVariables(rName, var1Name, var1Value, var2Name, var2Value string) string {
	return fmt.Sprintf(`
resource "codefresh_project" "test" {
  name = "%s"
  variables = {
	%q = %q
	%q = %q
  }
}
`, rName, var1Name, var1Value, var2Name, var2Value)
}
