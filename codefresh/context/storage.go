package context

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func convertStorageContext(context []interface{}, auth map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	data["auth"] = auth
	return data
}

func ConvertJsonConfigStorageContext(context []interface{}) map[string]interface{} {
	contextData := context[0].(map[string]interface{})
	contextAuth := contextData["auth"].([]interface{})[0].(map[string]interface{})
	auth := make(map[string]interface{})
	auth["type"] = contextAuth["type"]
	auth["jsonConfig"] = contextAuth["json_config"]
	return convertStorageContext(context, auth)
}

func ConvertAzureStorageContext(context []interface{}) map[string]interface{} {
	contextData := context[0].(map[string]interface{})
	contextAuth := contextData["auth"].([]interface{})[0].(map[string]interface{})
	auth := make(map[string]interface{})
	auth["type"] = contextAuth["type"]
	auth["accountName"] = contextAuth["account_name"]
	auth["accountKey"] = contextAuth["account_key"]
	return convertStorageContext(context, auth)
}

func flattenStorageContextConfig(spec cfclient.ContextSpec, auth map[string]interface{}) []interface{} {

	var res = make([]interface{}, 0)
	m := make(map[string]interface{})

	dataList := make([]interface{}, 0)
	data := make(map[string]interface{})

	authList := make([]interface{}, 0)
	authList = append(authList, auth)

	data["auth"] = authList

	dataList = append(dataList, data)

	m["data"] = dataList
	res = append(res, m)
	return res

}

func FlattenJsonConfigStorageContextConfig(spec cfclient.ContextSpec) []interface{} {
	auth := make(map[string]interface{})
	auth["json_config"] = spec.Data["auth"].(map[string]interface{})["jsonConfig"]
	auth["type"] = spec.Data["type"]
	return flattenStorageContextConfig(spec, auth)
}

func FlattenAzureStorageContextConfig(spec cfclient.ContextSpec) []interface{} {
	auth := make(map[string]interface{})
	authParams := spec.Data["auth"].(map[string]interface{})
	auth["account_name"] = authParams["accountName"]
	auth["account_key"] = authParams["accountKey"]
	auth["type"] = spec.Data["type"]
	return flattenStorageContextConfig(spec, auth)
}

func storageSchema(authSchema *schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"data": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"auth": authSchema,
						},
					},
				},
			},
		},
	}
}

func GcsSchema() *schema.Schema {
	sch := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"json_config": {
					Type:     schema.TypeMap,
					Required: true,
				},
			},
		},
	}
	return storageSchema(sch)
}

func S3Schema() *schema.Schema {
	sch := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"json_config": {
					Type:     schema.TypeMap,
					Required: true,
				},
			},
		},
	}
	return storageSchema(sch)
}

func AzureStorage() *schema.Schema {
	sch := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"account_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"account_key": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
	return storageSchema(sch)
}
