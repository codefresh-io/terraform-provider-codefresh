package codefresh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	funk "github.com/thoas/go-funk"
)

func TestAccCodefreshAbacRulesConfig(t *testing.T) {
	resourceName := "codefresh_abac_rules.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshAbacRulesConfig(
					"gitopsApplications",
					"LABEL",
					"KEY",
					"VALUE",
					[]string{"SYNC", "REFRESH"},
					[]string{"production", "*"},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAbacRulesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_type", "gitopsApplications"),
					resource.TestCheckResourceAttr(resourceName, "actions.0", "REFRESH"),
					resource.TestCheckResourceAttr(resourceName, "actions.1", "SYNC"),
					resource.TestCheckResourceAttr(resourceName, "attribute.0.name", "LABEL"),
					resource.TestCheckResourceAttr(resourceName, "attribute.0.key", "KEY"),
					resource.TestCheckResourceAttr(resourceName, "attribute.0.value", "VALUE"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "production"),
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

func testAccCheckCodefreshAbacRulesExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		abacRuleID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetAbacRuleByID(abacRuleID)

		if err != nil {
			return fmt.Errorf("error fetching abac rule with ID %s. %s", abacRuleID, err)
		}
		return nil
	}
}

// CONFIGS
func testAccCodefreshAbacRulesConfig(entityType, name, key, value string, actions, tags []string) string {
	escapeString := func(str string) string {
		if str == "null" {
			return str // null means Terraform should ignore this field
		}
		return fmt.Sprintf(`"%s"`, str)
	}
	tagsEscaped := funk.Map(tags, escapeString).([]string)
	actionsEscaped := funk.Map(actions, escapeString).([]string)

	attribute := ""
	if name != "" && value != "" {
		keyStr := ""
		if key != "" {
			keyStr = fmt.Sprintf(`key = %s`, escapeString(key))
		}
		attribute = fmt.Sprintf(`
		attribute {
					  name  = %s
					  %s
					  value = %s
				   }
		`, escapeString(name), keyStr, escapeString(value))
	}

	return fmt.Sprintf(`
	data "codefresh_team" "users" {
		name = "users"
	}

	resource "codefresh_abac_rules" "test" {
		teams            = [data.codefresh_team.users.id]
		entity_type      = %s
		actions          = [%s]
		%s
		tags             = [%s]
	}
`, escapeString(entityType), strings.Join(actionsEscaped[:], ","), attribute, strings.Join(tagsEscaped[:], ","))
}
