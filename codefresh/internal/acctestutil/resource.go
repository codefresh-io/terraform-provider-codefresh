package acctestutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccGetResourceId returns the ID of the resource with the given name,
// when provided with the Terraform state.
//
// This is useful for acceptance tests, in order to verify that a resource has
// been recreated and hence its ID has changed.
func GetResourceId(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource %s not found", resourceName)
	}
	return rs.Primary.ID, nil
}
