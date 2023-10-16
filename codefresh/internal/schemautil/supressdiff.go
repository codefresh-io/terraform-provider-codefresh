// Package schemautil provides utilities for working with Terraform resource schemas.
//
// Note that this package uses legacy logging because the provider context is not available
package schemautil

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	normalizationFailedErrorFormat = "[ERROR] Unable to normalize data body: %s"
)

// SuppressEquivalentYamlDiffs returns SchemaDiffSuppressFunc that suppresses diffs between
// equivalent YAML strings.
func SuppressEquivalentYamlDiffs() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		normalizedOld, err := NormalizeYamlString(old)
		if err != nil {
			log.Printf(normalizationFailedErrorFormat, err)
			return false
		}

		normalizedNew, err := NormalizeYamlString(new)
		if err != nil {
			log.Printf(normalizationFailedErrorFormat, err)
			return false
		}

		return normalizedOld == normalizedNew
	}
}
