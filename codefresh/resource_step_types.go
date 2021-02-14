package codefresh

import (
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStepTypes() *schema.Resource {
	return &schema.Resource{
		Create: resourceStepTypesCreate,
		Read:   resourceStepTypesRead,
		Update: resourceStepTypesUpdate,
		Delete: resourceStepTypesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"step_types_yaml": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     stringIsYaml,
				DiffSuppressFunc: suppressEquivalentYamlDiffs,
				StateFunc: func(v interface{}) string {
					template, _ := normalizeYamlString(v)
					return template
				},
			},
		},
	}
}

func resourceStepTypesCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)
	stepTypes := *mapResourceToStepTypes(d)
	resp, err := client.CreateStepTypes(&stepTypes)
	if err != nil {
		log.Printf("[DEBUG] Error while creating step types for resource_step_types. Error = %v", err)
		return err
	}

	d.SetId(resp.GetID())
	return resourceStepTypesRead(d, meta)
}

func resourceStepTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	stepTypesIdentifier := d.Id()

	if stepTypesIdentifier == "" {
		d.SetId("")
		return nil
	}
	var stepTypesGetVersion cfClient.StepTypes
	stepTypesYaml := d.Get("step_types_yaml")
	yaml.Unmarshal([]byte(stepTypesYaml.(string)), &stepTypesGetVersion)
	version := stepTypesGetVersion.Metadata["version"].(string)

	stepTypes, err := client.GetStepTypes(stepTypesIdentifier + ":" + version)
	// Remove transient attributes from metadata
	for _, attribute := range []string{"created_at", "accountId", "id", "updated_at", "latest"} {
		if _, ok := stepTypes.Metadata[attribute]; ok {
			delete(stepTypes.Metadata, attribute)
		}
	}
	if err != nil {
		log.Printf("[DEBUG] Error while getting stepTypes. Error = %v", stepTypesIdentifier)
		return err
	}

	err = mapStepTypesToResource(*stepTypes, d)
	if err != nil {
		log.Printf("[DEBUG] Error while mapping stepTypes to resource. Error = %v", err)
		return err
	}

	return nil
}

func resourceStepTypesUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	stepTypes := *mapResourceToStepTypes(d)
	newVersion := stepTypes.Metadata["version"].(string)
	existingVersions, err := client.GetStepTypesVersions(stepTypes.Metadata["name"].(string))
	if err == nil {
		for _, version := range existingVersions {
			if version == newVersion {
				log.Printf("[DEBUG] Version %s already exists. Updating...", newVersion)
				_, err := client.UpdateStepTypes(&stepTypes)
				if err != nil {
					log.Printf("[DEBUG] Error while updating stepTypes. Error = %v", err)
					return err
				}
				return resourceStepTypesRead(d, meta)
			}
		}
	}
	log.Printf("[DEBUG] Creating new version %s", newVersion)
	_, err = client.CreateStepTypes(&stepTypes)
	if err != nil {
		log.Printf("[DEBUG] Error while Creating stepTypes. Error = %v", err)
		return err
	}

	return resourceStepTypesRead(d, meta)
}

func resourceStepTypesDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	err := client.DeleteStepTypes(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapStepTypesToResource(stepTypes cfClient.StepTypes, d *schema.ResourceData) error {

	stepTypesYaml, err := yaml.Marshal(stepTypes)
	log.Printf("[DEBUG] Marshalled Step Types yaml = %v", string(stepTypesYaml))
	if err != nil {
		log.Printf("[DEBUG] Failed to Marshal Step Types yaml = %v", stepTypes)
		return err
	}
	err = d.Set("step_types_yaml", string(stepTypesYaml))

	if err != nil {
		return err
	}

	return nil
}

func mapResourceToStepTypes(d *schema.ResourceData) *cfClient.StepTypes {

	var stepTypes cfClient.StepTypes
	stepTypesYaml := d.Get("step_types_yaml")
	yaml.Unmarshal([]byte(stepTypesYaml.(string)), &stepTypes)
	if stepTypes.Spec.Steps != nil {
		stepTypes.Spec.Steps = extractSteps(stepTypesYaml.(string))
	}

	return &stepTypes
}
