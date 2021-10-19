package codefresh

import (
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// todo: 1) decide about password
// todo: 2) add another registry types
// todo: 3) add data definitions
// todo: 4) add more fields

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"registry_provider": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"primary": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"behind_firewall": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"deny_composite_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceRegistryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)
	resp, err := client.CreateRegistry(mapResourceToRegistry(d))
	if err != nil {
		log.Printf("[DEBUG] Error while creating registry. Error = %v", err)
		return err
	}

	d.SetId(resp.Id)
	return resourceRegistryRead(d, meta)
}

func resourceRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	registryId := d.Id()

	if registryId == "" {
		d.SetId("")
		return nil
	}

	registry, err := client.GetRegistry(registryId)
	if err != nil {
		log.Printf("[DEBUG] Error while getting registry. Error = %v", err)
		return err
	}

	err = mapRegistryToResource(*registry, d)
	if err != nil {
		log.Printf("[DEBUG] Error while mapping registry to resource. Error = %v", err)
		return err
	}

	return nil
}

func resourceRegistryUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	registry := *mapResourceToRegistry(d)
	registry.Id = d.Id()

	_, err := client.UpdateRegistry(&registry)
	if err != nil {
		log.Printf("[DEBUG] Error while updating registry. Error = %v", err)
		return err
	}

	return resourceRegistryRead(d, meta)
}

func resourceRegistryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	err := client.DeleteRegistry(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToRegistry(d *schema.ResourceData) *cfClient.Registry {
	registry := &cfClient.Registry{
		Id:       d.Id(),
		Name:     d.Get("name").(string),
		Kind:     d.Get("kind").(string),
		Provider: d.Get("registry_provider").(string),
		Domain:   d.Get("domain").(string),
	}

	if data, ok := d.GetOk("username"); ok {
		registry.Username = data.(string)
	}
	if data, ok := d.GetOk("password"); ok {
		registry.Password = data.(string)
	}

	return registry
}

func mapRegistryToResource(registry cfClient.Registry, d *schema.ResourceData) error {
	d.SetId(registry.Id)
	d.Set("name", registry.Name)
	d.Set("kind", registry.Kind)
	d.Set("registry_provider", registry.Provider)
	d.Set("domain", registry.Domain)
	d.Set("username", registry.Username)
	d.Set("password", registry.Password)

	return nil
}
