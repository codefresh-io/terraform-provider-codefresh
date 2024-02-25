package datautil

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
)

// ConvertStringArr converts an array of interfaces to an array of strings.
func ConvertStringArr(ifaceArr []interface{}) []string {
	return ConvertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

// ConvertAndMapStringArr converts an array of interfaces to an array of strings,
// applying the supplied function to each element.
func ConvertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

// ConvertVariables converts an array of cfclient.Variables to a map of key/value pairs.
func ConvertVariables(vars []cfclient.Variable) map[string]string {
	res := make(map[string]string, len(vars))
	for _, v := range vars {
		res[v.Key] = v.Value
	}
	return res
}

// FlattenStringArr flattens an array of strings.
func FlattenStringArr(sArr []string) []interface{} {
	iArr := []interface{}{}
	for _, s := range sArr {
		iArr = append(iArr, s)
	}
	return iArr
}
