package datautil

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/sclevine/yj/convert"
	"gopkg.in/op/go-logging.v1"
)

// Yq gets a value from a YAML string using yq
func Yq(yamlString string, expression string) (string, error) {
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

// YamlToJson converts a YAML string to JSON
//
// This function preserves the order of map keys (courtesy of yj package).
// If this were to use yaml.Unmarshal() and json.Marshal() instead, the order of map keys would be lost.
func YamlToJson(yamlString string) (string, error) {
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
