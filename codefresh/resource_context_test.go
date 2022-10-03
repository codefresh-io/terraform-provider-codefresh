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

var contextNamePrefix = "TerraformAccTest_"

func TestAccCodefreshContextConfigWithCharactersToBeEscaped(t *testing.T) {
	name := contextNamePrefix + "cf ctx/test +?#@ special" + acctest.RandString(10)
	resourceName := "codefresh_context.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshContextConfig(name, "config1", "value1", "config2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshContextExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.config.0.data.config1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.config.0.data.config2", "value2"),
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

func TestAccCodefreshContextConfig(t *testing.T) {
	name := contextNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_context.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshContextConfig(name, "config1", "value1", "config2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshContextExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.config.0.data.config1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.config.0.data.config2", "value2"),
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

func TestAccCodefreshContextSecret(t *testing.T) {
	name := contextNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_context.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshContextSecret(name, "config1", "value1", "config2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshContextExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.secret.0.data.config1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.secret.0.data.config2", "value2"),
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

func TestAccCodefreshContextYaml(t *testing.T) {
	name := contextNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_context.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshContextYaml(name, "rootKey", "plainKey", "plainValue", "listKey", "listValue1", "listValue2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshContextExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.yaml.0.data", "rootKey:\n  listKey:\n  - listValue1\n  - listValue2\n  plainKey: plainValue\n"),
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
func TestAccCodefreshContextSecretYaml(t *testing.T) {
	name := contextNamePrefix + acctest.RandString(10)
	resourceName := "codefresh_context.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshContextSecretYaml(name, "rootKey", "plainKey", "plainValue", "listKey", "listValue1", "listValue2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshContextExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.secretyaml.0.data", "rootKey:\n  listKey:\n  - listValue1\n  - listValue2\n  plainKey: plainValue\n"),
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

func testAccCheckCodefreshContextExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		contextID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfClient.Client)
		_, err := apiClient.GetContext(contextID)

		if err != nil {
			return fmt.Errorf("error fetching context with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCodefreshContextDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfClient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_context" {
			continue
		}

		_, err := apiClient.GetContext(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}

		notFoundErr := "Context .* not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

// CONFIGS
func testAccCodefreshContextConfig(rName, dataKey1, dataValue1, dataKey2, dataValue2 string) string {

	return fmt.Sprintf(`
resource "codefresh_context" "test" {

  name = "%s"

  spec {
	config {
		data = { 
			%q = %q
			%q = %q
		}
	}
  }
}
`, rName, dataKey1, dataValue1, dataKey2, dataValue2)
}

func testAccCodefreshContextSecret(rName, dataKey1, dataValue1, dataKey2, dataValue2 string) string {

	return fmt.Sprintf(`
resource "codefresh_context" "test" {

  name = "%s"

  spec {
	secret {
		data = { 
			%q = %q
			%q = %q
		}
	}
  }
}
`, rName, dataKey1, dataValue1, dataKey2, dataValue2)
}

func testAccCodefreshContextYaml(rName, rootKey, plainKey, plainValue, listKey, listValue1, listValue2 string) string {

	return fmt.Sprintf(`
resource "codefresh_context" "test" {

  name = "%s"

  spec {
	yaml {
		data = "%s: \n %s: %s\n %s: \n  - %s\n  - %s"
	}
  }
}
`, rName, rootKey, plainKey, plainValue, listKey, listValue1, listValue2)
}

func testAccCodefreshContextSecretYaml(rName, rootKey, plainKey, plainValue, listKey, listValue1, listValue2 string) string {

	return fmt.Sprintf(`
resource "codefresh_context" "test" {

  name = "%s"

  spec {
	secretyaml {
		data = "%s: \n %s: %s\n %s: \n  - %s\n  - %s"
	}
  }
}
`, rName, rootKey, plainKey, plainValue, listKey, listValue1, listValue2)
}
