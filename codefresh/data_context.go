package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContext() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve information on any defined context.",
		Read:        dataSourceContextRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceContextRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	var context *cfclient.Context
	var err error

	if name, nameOk := d.GetOk("name"); nameOk {
		context, err = client.GetContext(name.(string))
	} else {
		return fmt.Errorf("data.codefresh_context - must specify name")
	}
	if err != nil {
		return err
	}

	if context == nil {
		return fmt.Errorf("data.codefresh_context - cannot find context")
	}

	return mapDataContextToResource(context, d)
}

func mapDataContextToResource(context *cfclient.Context, d *schema.ResourceData) error {

	if context == nil || context.Metadata.Name == "" {
		return fmt.Errorf("data.codefresh_context - failed to mapDataContextToResource")
	}
	d.SetId(context.Metadata.Name)

	err := d.Set("name", context.Metadata.Name)

	if err != nil {
		return err
	}

	err = d.Set("type", context.Spec.Type)

	if err != nil {
		return err
	}

	data, err := yaml.Marshal(context.Spec.Data)
	if err != nil {
		return err
	}

	err = d.Set("data", string(data))

	if err != nil {
		return err
	}

	return nil
}
