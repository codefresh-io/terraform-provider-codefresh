package datautil

import (
	"io"
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
	yqEncoder := yqlib.NewYamlEncoder(yqlib.YamlPreferences{Indent: 0, ColorsEnabled: false})

	if outputformat == YQ_OUTPUT_FORMAT_JSON {
		yqEncoder = yqlib.NewJSONEncoder(yqlib.JsonPreferences{Indent: 0, ColorsEnabled: false, UnwrapScalar: false})
	}
	yqDecoder := yqlib.NewYamlDecoder(yqlib.NewDefaultYamlPreferences())
	yqEvaluator := yqlib.NewStringEvaluator()

	// Disable yq logging
	yqLogBackend := logging.AddModuleLevel(logging.NewLogBackend(io.Discard, "", 0))
	yqlib.GetLogger().SetBackend(yqLogBackend)

	yamlString, err := yqEvaluator.Evaluate(yamlString, expression, yqEncoder, yqDecoder)
	yamlString = strings.TrimSpace(yamlString)

	if yamlString == "null" { // yq's Evaluate() returns "null" if the expression does not match anything
		return "", err
	}
	return yamlString, err
}
