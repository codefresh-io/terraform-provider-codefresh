package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
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