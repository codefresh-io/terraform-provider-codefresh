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

// ConvertVariables converts an array of cfclient. Variables to 2 maps of key/value pairs - first one for un-encrypted variables second one for encrypted variables.
func ConvertVariables(vars []cfclient.Variable) (map[string]string,map[string]string) {
	
	numberOfEncryptedVars := 0
	
	for _, v := range vars {
		if v.Encrypted {
			numberOfEncryptedVars++
		}
	}

	resUnencrptedVars := make(map[string]string, len(vars) - numberOfEncryptedVars)
	resEncryptedVars := make(map[string]string, numberOfEncryptedVars)

	for _, v := range vars {
		if v.Encrypted {
			resEncryptedVars[v.Key] = v.Value
		} else {
			resUnencrptedVars[v.Key] = v.Value
		}
	}

	return resUnencrptedVars,resEncryptedVars
}

// FlattenStringArr flattens an array of strings.
func FlattenStringArr(sArr []string) []interface{} {
	iArr := []interface{}{}
	for _, s := range sArr {
		iArr = append(iArr, s)
	}
	return iArr
}
