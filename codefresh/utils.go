package codefresh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/sclevine/yj/convert"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	logging "gopkg.in/op/go-logging.v1"
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

func convertVariables(vars []cfclient.Variable) map[string]string {
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

// Get a value from a YAML string using yq
func yq(yamlString string, expression string) (string, error) {
	yqEncoder := yqlib.NewYamlEncoder(0, false, yqlib.NewDefaultYamlPreferences())
	yqDecoder := yqlib.NewYamlDecoder(yqlib.NewDefaultYamlPreferences())
	yqEvaluator := yqlib.NewStringEvaluator()

	// Disable yq logging
	yqLogBackend := logging.AddModuleLevel(logging.NewLogBackend(ioutil.Discard, "", 0))
	yqlib.GetLogger().SetBackend(yqLogBackend)

	yamlString, err := yqEvaluator.Evaluate(yamlString, expression, yqEncoder, yqDecoder)
	yamlString = strings.TrimSpace(yamlString)

	if yamlString == "null" { // yq's Evaluate() returns "null" if the expression does not match anything
		return "", err
	}
	return yamlString, err
}

// Convert a YAML string to JSON while preserving the order of map keys (courtesy of yj package).
// If this were to use yaml.Unmarshal() and json.Marshal() instead, the order of map keys would be lost.
func yamlToJson(yamlString string) (string, error) {
	yamlConverter := convert.YAML{}
	jsonConverter := convert.JSON{}

	yamlDecoded, err := yamlConverter.Decode(strings.NewReader(yamlString))
	if err != nil {
		return "", err
	}

	jsonBuffer := new(bytes.Buffer)
	err = jsonConverter.Encode(jsonBuffer, yamlDecoded)
	if err != nil {
		return "", err
	}

	return jsonBuffer.String(), nil
}

func testAccGetResourceId(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource %s not found", resourceName)
	}
	return rs.Primary.ID, nil
}