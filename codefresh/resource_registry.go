package codefresh

import (
	"fmt"
	"log"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
		Description: "Registry is the configuration that Codefresh uses to push/pull container images.",
		Create:      resourceRegistryCreate,
		Read:        resourceRegistryRead,
		Update:      resourceRegistryUpdate,
		Delete:      resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The display name for the registry.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kind": {
				Description: "The kind of registry.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default": {
				Description: "Whether this registry is the default registry (default: `false`).",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"primary": {
				Description: "Whether this registry is the primary registry (default: `true`).",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"fallback_registry": {
				Description: "The name of the fallback registry.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"spec": {
				Description: "The registry's specs.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						providerAcr: {
							Description:   "An `acr` block as documented below ([Azure Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/azure-docker-registry)).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerAcr),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Description: "The Client ID.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"client_secret": {
										Description: "The Client Secret.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"domain": {
										Description: "The ACR registry domain.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"repository_prefix": {
										Description: "See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
						providerGcr: {
							Description:   "[Google Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/google-container-registry).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerGcr),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Description: "The GCR registry domain.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"keyfile": {
										Description: "The serviceaccount json file contents.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"repository_prefix": {
										Description: "See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
						providerGar: {
							Description:   "A `gar` block as documented below ([Google Artifact Registry](https://codefresh.io/docs/docs/integrations/docker-registries/google-artifact-registry)).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerGar),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Description: "The GAR location.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"keyfile": {
										Description: "The serviceaccount json file contents.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"repository_prefix": {
										Description: "See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
						providerEcr: {
							Description:   "An `ecr` block as documented below ([Amazon EC2 Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/amazon-ec2-container-registry)).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerEcr),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Description: "The AWS region.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"access_key_id": {
										Description: "The AWS access key ID.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"secret_access_key": {
										Description: "The AWS secret access key.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"repository_prefix": {
										Description: "See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
						providerBintray: {
							Description:   "A `bintray` block as documented below ([Bintray / Artifactory](https://codefresh.io/docs/docs/integrations/docker-registries/bintray-io)).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerBintray),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Description: "The Bintray username.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"token": {
										Description: "The Bintray token.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"domain": {
										Description: "The Bintray domain.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"repository_prefix": {
										Description: "See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
						providerDockerhub: {
							Description:   "A `dockerhub` block as documented below ([Docker Hub Registry](https://codefresh.io/docs/docs/integrations/docker-registries/docker-hub/)).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingProviders(providers, providerDockerhub),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Description: "The DockerHub username.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"password": {
										Description: "The DockerHub password.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
								},
							},
						},
						providerOther: {
							Description:   "`other` provider block described below ([Other Providers](https://codefresh.io/docs/docs/integrations/docker-registries/other-registries)).",
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							ConflictsWith: getConflictingProviders(providers, providerOther),
							MaxItems:      1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Description: "The username.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"password": {
										Description: "The password.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"domain": {
										Description: "The domain.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"repository_prefix": {
										Description: "See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
									"behind_firewall": {
										Description: "See the [docs](https://codefresh.io/docs/docs/administration/behind-the-firewall/#accessing-an-internal-docker-registry).",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										ForceNew:    true,
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

	client := meta.(*cfclient.Client)
	resp, err := client.CreateRegistry(mapResourceToRegistry(d))
	if err != nil {
		log.Printf("[DEBUG] Error while creating registry. Error = %v", err)
		return err
	}

	d.SetId(resp.Id)
	return resourceRegistryRead(d, meta)
}

func resourceRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

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

	client := meta.(*cfclient.Client)

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

	client := meta.(*cfclient.Client)

	err := client.DeleteRegistry(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapResourceToRegistry(d *schema.ResourceData) *cfclient.Registry {
	registry := &cfclient.Registry{
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
		registry.Domain = d.Get(providerKey + ".0.location").(string)
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

func mapRegistryToResource(registry cfclient.Registry, d *schema.ResourceData) error {
	d.SetId(registry.Id)
	err := d.Set("name", registry.Name)

	if err != nil {
		return err
	}

	err = d.Set("kind", registry.Kind)

	if err != nil {
		return err
	}

	err = d.Set("default", registry.Default)

	if err != nil {
		return err
	}

	err = d.Set("primary", registry.Primary)

	if err != nil {
		return err
	}

	err = d.Set("fallback_registry", registry.FallbackRegistry)

	if err != nil {
		return err
	}

	if registry.Provider == providerAcr {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerAcr), map[string]interface{}{
			"domain":            registry.Domain,
			"client_id":         registry.ClientId,
			"client_secret":     registry.ClientSecret,
			"repository_prefix": registry.RepositoryPrefix,
		})

		if err != nil {
			return err
		}
	}

	if registry.Provider == providerEcr {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerEcr), map[string]interface{}{
			"region":            registry.Domain,
			"access_key_id":     registry.AccessKeyId,
			"secret_access_key": registry.SecretAccessKey,
			"repository_prefix": registry.RepositoryPrefix,
		})

		if err != nil {
			return err
		}
	}

	if registry.Provider == providerGcr {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerGcr), map[string]interface{}{
			"domain":            registry.Domain,
			"keyfile":           registry.Keyfile,
			"repository_prefix": registry.RepositoryPrefix,
		})

		if err != nil {
			return err
		}
	}

	if registry.Provider == providerGar {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerGcr), map[string]interface{}{
			"location":          registry.Domain,
			"keyfile":           registry.Keyfile,
			"repository_prefix": registry.RepositoryPrefix,
		})

		if err != nil {
			return err
		}
	}

	if registry.Provider == providerBintray {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerBintray), map[string]interface{}{
			"domain":            registry.Domain,
			"username":          registry.Username,
			"token":             registry.Token,
			"repository_prefix": registry.RepositoryPrefix,
		})

		if err != nil {
			return err
		}
	}

	if registry.Provider == providerDockerhub {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerDockerhub), map[string]interface{}{
			"username": registry.Username,
			"password": registry.Password,
		})

		if err != nil {
			return err
		}
	}

	if registry.Provider == providerOther {
		err = d.Set(fmt.Sprintf("spec.0.%v.0", providerOther), map[string]interface{}{
			"domain":            registry.Domain,
			"username":          registry.Username,
			"password":          registry.Password,
			"behind_firewall":   registry.BehindFirewall,
			"repository_prefix": registry.RepositoryPrefix,
		})

		if err != nil {
			return err
		}
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
