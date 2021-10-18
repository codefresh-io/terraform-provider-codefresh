package context

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ConvertStorageContext(context []interface{}) map[string]interface{} {
	contextData := context[0].(map[string]interface{})
	contextAuth := contextData["auth"].([]interface{})[0].(map[string]interface{})
	data := make(map[string]interface{})
	auth := make(map[string]interface{})
	auth["type"] = contextAuth["type"]
	auth["jsonConfig"] = contextAuth["json_config"]
	data["auth"] = auth
	return data
}

func FlattenStorageContextConfig(spec cfClient.ContextSpec) []interface{} {

	var res = make([]interface{}, 0)
	m := make(map[string]interface{})

	dataList := make([]interface{}, 0)
	data := make(map[string]interface{})

	auth := make(map[string]interface{})
	auth["json_config"] = spec.Data["auth"].(map[string]interface{})["jsonConfig"]
	auth["type"] = spec.Data["type"]

	authList := make([]interface{}, 0)
	authList = append(authList, auth)

	data["auth"] = authList

	dataList = append(dataList, data)

	m["data"] = dataList
	res = append(res, m)
	return res

}

func storageSchema() *schema.Schema {
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
							"auth": {
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
							},
						},
					},
				},
			},
		},
	}
}

func GcsSchema() *schema.Schema {
	return storageSchema()
}

func S3Schema() *schema.Schema {
	return storageSchema()
}
