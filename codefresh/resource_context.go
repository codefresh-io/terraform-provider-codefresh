package codefresh

import (
	"log"

	storageContext "github.com/codefresh-io/terraform-provider-codefresh/codefresh/context"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/schemautil"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	contextConfig        = "config"
	contextSecret        = "secret"
	contextYaml          = "yaml"
	contextSecretYaml    = "secret-yaml"
	contextGoogleStorage = "storage.gc"
	contextS3Storage     = "storage.s3"
	contextAzureStorage  = "storage.azuref"
)

var supportedContextType = []string{
	contextConfig,
	contextSecret,
	contextYaml,
	contextSecretYaml,
}

var encryptedContextTypes = []string{
	contextSecret,
	contextSecretYaml,
	contextS3Storage,
	contextAzureStorage,
}

func getConflictingContexts(context string) []string {
	var conflictingTypes []string
	normalizedContext := schemautil.MustNormalizeFieldName(context)
	for _, value := range supportedContextType {
		normlizedValue := schemautil.MustNormalizeFieldName(value)
		if normlizedValue != normalizedContext {
			conflictingTypes = append(conflictingTypes, "spec.0."+normlizedValue)
		}
	}
	return conflictingTypes
}

func resourceContext() *schema.Resource {
	return &schema.Resource{
		Description: "A Context is an authentication/configuration resource used by the Codefresh system and engine.",
		Create:      resourceContextCreate,
		Read:        resourceContextRead,
		Update:      resourceContextUpdate,
		Delete:      resourceContextDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The display name for the context.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"spec": {
				Description: "The context's specs.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						schemautil.MustNormalizeFieldName(contextConfig): {
							Type:          schema.TypeList,
							ForceNew:      true,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingContexts(contextConfig),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Description: "The map of variables representing the shared config.",
										Type:        schema.TypeMap,
										Required:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						schemautil.MustNormalizeFieldName(contextSecret): {
							Type:          schema.TypeList,
							Optional:      true,
							ForceNew:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingContexts(contextSecret),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Description: "The map of variables representing the shared config (secret).",
										Type:        schema.TypeMap,
										Required:    true,
										Sensitive:   true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						schemautil.MustNormalizeFieldName(contextYaml): {
							Type:          schema.TypeList,
							Optional:      true,
							ForceNew:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingContexts(contextYaml),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Description:      "The YAML string representing the shared config.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: schemautil.StringIsValidYaml(),
										DiffSuppressFunc: schemautil.SuppressEquivalentYamlDiffs(),
										StateFunc: func(v interface{}) string {
											return schemautil.MustNormalizeYamlString(v)
										},
									},
								},
							},
						},
						schemautil.MustNormalizeFieldName(contextSecretYaml): {
							Type:          schema.TypeList,
							Optional:      true,
							ForceNew:      true,
							MaxItems:      1,
							ConflictsWith: getConflictingContexts(contextSecretYaml),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Description:      "The YAML string representing the shared config (secret).",
										Type:             schema.TypeString,
										Required:         true,
										Sensitive:        true,
										ValidateDiagFunc: schemautil.StringIsValidYaml(),
										DiffSuppressFunc: schemautil.SuppressEquivalentYamlDiffs(),
										StateFunc: func(v interface{}) string {
											return schemautil.MustNormalizeYamlString(v)
										},
									},
								},
							},
						},
						schemautil.MustNormalizeFieldName(contextGoogleStorage): storageContext.GcsSchema(),
						schemautil.MustNormalizeFieldName(contextS3Storage):     storageContext.S3Schema(),
						schemautil.MustNormalizeFieldName(contextAzureStorage):  storageContext.AzureStorage(),
					},
				},
			},
		},
	}
}

func resourceContextCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)
	resp, err := client.CreateContext(mapResourceToContext(d))
	if err != nil {
		log.Printf("[DEBUG] Error while creating context. Error = %v", err)
		return err
	}

	d.SetId(resp.Metadata.Name)
	return resourceContextRead(d, meta)
}

func resourceContextRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	contextName := d.Id()

	if contextName == "" {
		d.SetId("")
		return nil
	}

	context, err := client.GetContext(contextName)

	if err != nil {
		log.Printf("[DEBUG] Error while getting context. Error = %v", contextName)
		return err
	}

	err = mapContextToResource(*context, d)

	if err != nil {
		log.Printf("[DEBUG] Error while mapping context to resource. Error = %v", err)
		return err
	}

	return nil
}

func resourceContextUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	context := *mapResourceToContext(d)
	context.Metadata.Name = d.Id()

	_, err := client.UpdateContext(&context)
	if err != nil {
		log.Printf("[DEBUG] Error while updating context. Error = %v", err)
		return err
	}

	return resourceContextRead(d, meta)
}

func resourceContextDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	err := client.DeleteContext(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapContextToResource(context cfclient.Context, d *schema.ResourceData) error {

	err := d.Set("name", context.Metadata.Name)

	if err != nil {
		return err
	}

	// Read spec from API if context is not encrypted or forbitDecrypt is not set
	if !context.IsEncrypred {

		err = d.Set("spec", flattenContextSpec(context.Spec))

		if err != nil {
			log.Printf("[DEBUG] Failed to flatten Context spec = %v", context.Spec)
			return err
		}
	}

	return nil
}

func flattenContextSpec(spec cfclient.ContextSpec) []interface{} {

	var res = make([]interface{}, 0)
	m := make(map[string]interface{})

	switch currentContextType := spec.Type; currentContextType {
	case contextConfig, contextSecret:
		m[schemautil.MustNormalizeFieldName(currentContextType)] = flattenContextConfig(spec)
	case contextYaml, contextSecretYaml:
		m[schemautil.MustNormalizeFieldName(currentContextType)] = flattenContextYaml(spec)
	case contextGoogleStorage, contextS3Storage:
		m[schemautil.MustNormalizeFieldName(currentContextType)] = storageContext.FlattenJsonConfigStorageContextConfig(spec)
	case contextAzureStorage:
		m[schemautil.MustNormalizeFieldName(currentContextType)] = storageContext.FlattenAzureStorageContextConfig(spec)
	default:
		return nil
	}

	res = append(res, m)
	return res
}

func flattenContextConfig(spec cfclient.ContextSpec) []interface{} {
	var res = make([]interface{}, 0)
	m := make(map[string]interface{})
	m["data"] = spec.Data
	res = append(res, m)
	return res
}

func flattenContextYaml(spec cfclient.ContextSpec) []interface{} {
	var res = make([]interface{}, 0)
	m := make(map[string]interface{})
	data, err := yaml.Marshal(spec.Data)
	if err != nil {
		return nil
	}
	m["data"] = string(data)
	res = append(res, m)
	return res
}

func mapResourceToContext(d *schema.ResourceData) *cfclient.Context {

	var normalizedContextType string
	var normalizedContextData map[string]interface{}

	if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextConfig) + ".0.data"); ok {
		normalizedContextType = contextConfig
		normalizedContextData = data.(map[string]interface{})
	} else if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextSecret) + ".0.data"); ok {
		normalizedContextType = contextSecret
		normalizedContextData = data.(map[string]interface{})
	} else if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextYaml) + ".0.data"); ok {
		normalizedContextType = contextYaml
		_ = yaml.Unmarshal([]byte(data.(string)), &normalizedContextData)
	} else if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextSecretYaml) + ".0.data"); ok {
		normalizedContextType = contextSecretYaml
		_ = yaml.Unmarshal([]byte(data.(string)), &normalizedContextData)
	} else if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextGoogleStorage) + ".0.data"); ok {
		normalizedContextType = contextGoogleStorage
		normalizedContextData = storageContext.ConvertJsonConfigStorageContext(data.([]interface{}))
	} else if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextS3Storage) + ".0.data"); ok {
		normalizedContextType = contextS3Storage
		normalizedContextData = storageContext.ConvertJsonConfigStorageContext(data.([]interface{}))
	} else if data, ok := d.GetOk("spec.0." + schemautil.MustNormalizeFieldName(contextAzureStorage) + ".0.data"); ok {
		normalizedContextType = contextAzureStorage
		normalizedContextData = storageContext.ConvertAzureStorageContext(data.([]interface{}))
	}

	return &cfclient.Context{
		Metadata: cfclient.ContextMetadata{
			Name: d.Get("name").(string),
		},
		Spec: cfclient.ContextSpec{
			Type: normalizedContextType,
			Data: normalizedContextData,
		},
	}
}
