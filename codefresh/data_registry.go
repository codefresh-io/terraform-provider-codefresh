package codefresh

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRegistry() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows retrieving information on any existing registry.",
		Read:        dataSourceRegistryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"fallback_registry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRegistryRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	var registry *cfclient.Registry
	var err error

	if name, nameOk := d.GetOk("name"); nameOk {
		registry, err = client.GetRegistry(name.(string))
	} else {
		return fmt.Errorf("data.codefresh_registry - must specify name")
	}
	if err != nil {
		return err
	}

	if registry == nil {
		return fmt.Errorf("data.codefresh_registry - cannot find registry")
	}

	return mapDataRegistryToResource(registry, d)
}

func mapDataRegistryToResource(registry *cfclient.Registry, d *schema.ResourceData) error {

	if registry == nil || registry.Name == "" {
		return fmt.Errorf("data.codefresh_registry - failed to mapDataRegistryToResource")
	}
	d.SetId(registry.Id)

	d.Set("name", registry.Name)
	d.Set("registry_provider", registry.Provider)
	d.Set("kind", registry.Kind)
	d.Set("domain", registry.Domain)
	d.Set("primary", registry.Primary)
	d.Set("default", registry.Default)
	d.Set("fallback_registry", registry.FallbackRegistry)
	d.Set("repository_prefix", registry.RepositoryPrefix)

	return nil
}
