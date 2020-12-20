package codefresh

import (
	"fmt"
	"log"
	"regexp"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func convertStringArr(ifaceArr []interface{}) []string {
	return convertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

func convertVariables(vars []cfClient.Variable) map[string]string {
	res := make(map[string]string, len(vars))
	for _, v := range vars {
		res[v.Key] = v.Value
	}
	return res
}

func flattenStringArr(sArr []string) []interface{} {
	iArr := []interface{}{}
	for _, s := range sArr {
		iArr = append(iArr, s)
	}
	return iArr
}

func stringIsYaml(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if _, err := normalizeYamlString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
	}

	return warnings, errors
}

func normalizeFieldName(fieldName string) string {
	reg, err := regexp.Compile("[^a-z0-9_]+")
	if err != nil {
		log.Printf("[DEBUG] Unable to compile regexp for field name normalization. Error = %v", err)
	}
	return reg.ReplaceAllString(fieldName, "")
}

func normalizeYamlString(yamlString interface{}) (string, error) {
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

func suppressEquivalentYamlDiffs(k, old, new string, d *schema.ResourceData) bool {
	normalizedOld, err := normalizeYamlString(old)

	if err != nil {
		log.Printf("[ERROR] Unable to normalize data body: %s", err)
		return false
	}

	normalizedNew, err := normalizeYamlString(new)

	if err != nil {
		log.Printf("[ERROR] Unable to normalize data body: %s", err)
		return false
	}

	return normalizedOld == normalizedNew
}
