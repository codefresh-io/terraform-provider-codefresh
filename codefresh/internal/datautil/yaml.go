package datautil

import (
	"io/ioutil"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

const (
	YQ_OUTPUT_FORMAT_JSON = "json"
	YQ_OUTPUT_FORMAT_YAML = "yaml"
)

// Yq gets a value from a YAML string using yq
func Yq(yamlString string, expression string, outputformat string) (string, error) {
	yqEncoder := yqlib.NewYamlEncoder(0, false, yqlib.NewDefaultYamlPreferences())

	if outputformat == YQ_OUTPUT_FORMAT_JSON {
		yqEncoder = yqlib.NewJSONEncoder(0, false, false)
	}
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
