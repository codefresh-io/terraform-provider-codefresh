package codefresh

import (
	"fmt"
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// todo: 1) add data definitions
// todo: 2) add tests

const (
	providerOther     = "other"
	providerGcr       = "gcr"
	providerGar       = "gar"
	providerEcr       = "ecr"
	providerAcr       = "acr"
	providerDockerhub = "dockerhub"
	providerBintray   = "bintray"
)

var providers = []string{
	providerOther,
	providerGcr,
	providerGar,
	providerEcr,
	providerAcr,
	providerDockerhub,
	providerBintray,
}

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
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"primary": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"fallback_registry": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						providerAcr: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerAcr),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"client_secret": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
									"repository_prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						providerGcr: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerGcr),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
									"keyfile": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"repository_prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						providerGar: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerGar),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
									"keyfile": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"repository_prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						providerEcr: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerEcr),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Required: true,
									},
									"access_key_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"secret_access_key": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"repository_prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						providerBintray: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerBintray),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"token": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
									"repository_prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						providerDockerhub: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerDockerhub),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
						providerOther: {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							ConflictsWith: getConflictingProviders(providers, providerOther),
							MaxItems:      1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
									"repository_prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"behind_firewall": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
										ForceNew: true,
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
		Id:               d.Id(),
		Name:             d.Get("name").(string),
		Kind:             d.Get("kind").(string),
		Default:          d.Get("default").(bool),
		Primary:          d.Get("primary").(bool),
		FallbackRegistry: d.Get("fallback_registry").(string),
	}

	var providerKey string
	providerKey = "spec.0." + providerAcr
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerAcr
		registry.Domain = d.Get(providerKey + ".0.domain").(string)
		registry.ClientId = d.Get(providerKey + ".0.client_id").(string)
		registry.ClientSecret = d.Get(providerKey + ".0.client_secret").(string)
		registry.RepositoryPrefix = d.Get(providerKey + ".0.repository_prefix").(string)
		return registry
	}

	providerKey = "spec.0." + providerEcr
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerEcr
		registry.Region = d.Get(providerKey + ".0.region").(string)
		registry.AccessKeyId = d.Get(providerKey + ".0.access_key_id").(string)
		registry.SecretAccessKey = d.Get(providerKey + ".0.secret_access_key").(string)
		registry.RepositoryPrefix = d.Get(providerKey + ".0.repository_prefix").(string)
		return registry
	}

	providerKey = "spec.0." + providerGcr
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerGcr
		registry.Domain = d.Get(providerKey + ".0.domain").(string)
		registry.Keyfile = d.Get(providerKey + ".0.keyfile").(string)
		registry.RepositoryPrefix = d.Get(providerKey + ".0.repository_prefix").(string)
		return registry
	}

	providerKey = "spec.0." + providerGar
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerGar
		registry.Domain = d.Get(providerKey + ".0.domain").(string)
		registry.Keyfile = d.Get(providerKey + ".0.keyfile").(string)
		registry.RepositoryPrefix = d.Get(providerKey + ".0.repository_prefix").(string)
		return registry
	}

	providerKey = "spec.0." + providerDockerhub
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerDockerhub
		registry.Username = d.Get(providerKey + ".0.username").(string)
		registry.Password = d.Get(providerKey + ".0.password").(string)
		return registry
	}

	providerKey = "spec.0." + providerBintray
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerBintray
		registry.Domain = d.Get(providerKey + ".0.domain").(string)
		registry.Username = d.Get(providerKey + ".0.username").(string)
		registry.Token = d.Get(providerKey + ".0.token").(string)
		registry.RepositoryPrefix = d.Get(providerKey + ".0.repository_prefix").(string)
		return registry
	}

	providerKey = "spec.0." + providerOther
	if _, ok := d.GetOk(providerKey); ok {
		registry.Provider = providerOther
		registry.Domain = d.Get(providerKey + ".0.domain").(string)
		registry.Username = d.Get(providerKey + ".0.username").(string)
		registry.Password = d.Get(providerKey + ".0.password").(string)
		registry.BehindFirewall = d.Get(providerKey + ".0.behind_firewall").(bool)
		registry.RepositoryPrefix = d.Get(providerKey + ".0.repository_prefix").(string)
		return registry
	}

	return registry
}

func mapRegistryToResource(registry cfClient.Registry, d *schema.ResourceData) error {
	d.SetId(registry.Id)
	d.Set("name", registry.Name)
	d.Set("kind", registry.Kind)
	d.Set("default", registry.Default)
	d.Set("primary", registry.Primary)
	d.Set("fallback_registry", registry.FallbackRegistry)

	if registry.Provider == providerAcr {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerAcr), map[string]interface{}{
			"domain":            registry.Domain,
			"client_id":         registry.ClientId,
			"client_secret":     registry.ClientSecret,
			"repository_prefix": registry.RepositoryPrefix,
		})
	}

	if registry.Provider == providerEcr {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerEcr), map[string]interface{}{
			"region":            registry.Domain,
			"access_key_id":     registry.AccessKeyId,
			"secret_access_key": registry.SecretAccessKey,
			"repository_prefix": registry.RepositoryPrefix,
		})
	}

	if registry.Provider == providerGcr {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerGcr), map[string]interface{}{
			"domain":            registry.Domain,
			"keyfile":           registry.Keyfile,
			"repository_prefix": registry.RepositoryPrefix,
		})
	}

	if registry.Provider == providerGar {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerGcr), map[string]interface{}{
			"domain":            registry.Domain,
			"keyfile":           registry.Keyfile,
			"repository_prefix": registry.RepositoryPrefix,
		})
	}

	if registry.Provider == providerBintray {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerBintray), map[string]interface{}{
			"domain":            registry.Domain,
			"username":          registry.Username,
			"token":             registry.Token,
			"repository_prefix": registry.RepositoryPrefix,
		})
	}

	if registry.Provider == providerDockerhub {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerDockerhub), map[string]interface{}{
			"username": registry.Username,
			"password": registry.Password,
		})
	}

	if registry.Provider == providerOther {
		d.Set(fmt.Sprintf("spec.0.%v.0", providerOther), map[string]interface{}{
			"domain":            registry.Domain,
			"username":          registry.Username,
			"password":          registry.Password,
			"behind_firewall":   registry.BehindFirewall,
			"repository_prefix": registry.RepositoryPrefix,
		})
	}

	return nil
}

func getConflictingProviders(arr []string, exclude string) []string {
	filtered := make([]string, 0)
	for _, provider := range arr {
		if provider != exclude {
			filtered = append(filtered, "spec.0."+provider)
		}
	}
	return filtered
}
