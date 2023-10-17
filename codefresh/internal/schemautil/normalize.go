package schemautil

import (
	"log"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	NormalizedFieldNameRegex string = `[^a-z0-9_]+`
)

// NormalizeFieldName normalizes a field name to be lowercase and contain only alphanumeric characters and dashes.
func NormalizeFieldName(fieldName string) (string, error) {
	reg, err := regexp.Compile(NormalizedFieldNameRegex)
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(fieldName, ""), nil
}

// MustNormalizeFieldName is the same as NormalizeFieldName, but will log an error (legacy logging) instead of returning it.
func MustNormalizeFieldName(fieldName string) string {
	normalizedFieldName, err := NormalizeFieldName(fieldName)
	if err != nil {
		log.Printf("[ERROR] Failed to normalize field name %q: %s", fieldName, err)
	}
	return normalizedFieldName
}

// NormalizeYAMLString normalizes a YAML string to a standardized order, format and indentation.
func NormalizeYamlString(yamlString interface{}) (string, error) {
	var j map[string]interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)
	err := yaml.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	bytes, _ := yaml.Marshal(j)
	return string(bytes[:]), nil
}

// MustNormalizeYamlString is the same as NormalizeYamlString, but will log an error (legacy logging) instead of returning it.
func MustNormalizeYamlString(yamlString interface{}) string {
	normalizedYamlString, err := NormalizeYamlString(yamlString)
	if err != nil {
		log.Printf("[ERROR] Failed to normalize YAML string %q: %s", yamlString, err)
	}
	return normalizedYamlString
}
