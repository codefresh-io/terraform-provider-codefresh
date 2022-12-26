package codefresh

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// can be used instead the yaml.MapSlice (as an argument to yaml.Unmarshal) in order to preserve the keys order when converting to JSON later
//
//	// Usage example:
//	ms := OrderedMapSlice{}
//	yaml.Unmarshal([]byte(originalYamlString), &ms)
//	orderedJson, _ := json.Marshal(ms)
//
// implements json.Marshaler interface
type OrderedMapSlice []yaml.MapItem

func (ms OrderedMapSlice) MarshalJSON() ([]byte, error) {
	//	keep the order of keys while converting to json with json.Marshal(ms)

	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, mi := range ms {
		b, err := json.Marshal(&mi.Value)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprintf("%v", mi.Key)))
		buf.Write(b)
		if i < len(ms)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}
