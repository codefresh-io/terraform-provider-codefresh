package context

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GcsSchema() *schema.Schema {
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
