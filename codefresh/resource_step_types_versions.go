package codefresh

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/Masterminds/semver"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStepTypesVersions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStepTypesVersionCreate,
		ReadContext:   resourceStepTypesVersionRead,
		UpdateContext: resourceStepTypesVersionUpdate,
		DeleteContext: resourceStepTypesVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"version": {
				Type:       schema.TypeSet,
				Required:   true,
				MinItems:   1,
				Set:        resourceStepTypesVersionsConfigHash,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_number": {
							Type:     schema.TypeString,
							Required: true,
						},
						"step_types_yaml": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     stringIsYaml,
							DiffSuppressFunc: suppressEquivalentYamlDiffs,
							StateFunc: func(v interface{}) string {
								template, _ := normalizeYamlStringStepTypes(v)
								return template
							},
						},
					},
				},
			},
		},
	}
}

func normalizeYamlStringStepTypes(yamlString interface{}) (string, error) {
	var j map[string]interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)
	err := yaml.Unmarshal([]byte(s), &j)
	metadataMap := j["metadata"].(map[string]interface{})
	//Removing "latest" attribute from metadata since it's transient based on the version
	delete(metadataMap, "latest")
	if err != nil {
		return s, err
	}

	bytes, _ := yaml.Marshal(j)
	return string(bytes[:]), nil
}

func resourceStepTypesVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*cfClient.Client)
	stepTypes := *mapResourceToStepTypesVersions(d)

	name := d.Get("name").(string)
	d.SetId(name)

	// Extract all the versions so that we can order the set based on semantic versioning
	mapVersion := make(map[string]cfClient.StepTypes)
	var versions []string
	for _, version := range stepTypes.Versions {
		version.StepTypes.Metadata["name"] = name
		version.StepTypes.Metadata["version"] = version.VersionNumber
		log.Printf("[DEBUG] Length: %q, %v", versions, len(stepTypes.Versions))
		versions = append(versions, version.VersionNumber)

		mapVersion[version.VersionNumber] = version.StepTypes

	}

	// Create the versions in order based on semver
	orderedVersions := sortVersions(versions)
	for _, version := range orderedVersions {
		step := mapVersion[version.String()]
		log.Printf("[DEBUG] Version for create: %q", version)
		_, err := client.CreateStepTypes(&step)
		if err != nil {
			return diag.Errorf("[DEBUG] Error while creating step types OnCreate. Error = %v", err)
		}
	}

	return resourceStepTypesVersionRead(ctx, d, meta)
}

func resourceStepTypesVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cfClient.Client)

	stepTypesIdentifier := d.Id()
	if stepTypesIdentifier == "" {

		d.SetId("")
		return nil
	}

	//Extracting the step just based on the name to validate it exists
	stepTypes, err := client.GetStepTypes(stepTypesIdentifier)
	if err != nil {
		log.Printf("[DEBUG] Step Not found %v. Error = %v", stepTypesIdentifier, err)
		d.SetId("")
		return nil
	}

	var stepVersions cfClient.StepTypesVersions
	name := stepTypes.Metadata["name"].(string)
	stepVersions.Name = name
	versions := d.Get("version").(*schema.Set)
	// Try to retrieve defined versions and add to the list if it exists
	for _, step := range versions.List() {
		version := step.(map[string]interface{})["version_number"].(string)
		log.Printf("[DEBUG] Get step version FromList %v", version)
		if version != "" {
			stepTypes, err := client.GetStepTypes(stepTypesIdentifier + ":" + version)
			log.Printf("[DEBUG] Get step version %v", version)
			if err != nil {
				log.Printf("[DEBUG] StepVersion not found %v. Error = %v", stepTypesIdentifier+":"+version, err)
			} else {
				cleanUpStepFromTransientValues(stepTypes, name, version)
				stepVersion := cfClient.StepTypesVersion{
					VersionNumber: version,
					StepTypes:     *stepTypes,
				}
				stepVersions.Versions = append(stepVersions.Versions, stepVersion)

			}
		}
	}

	err = mapStepTypesVersionsToResource(stepVersions, d)

	if err != nil {
		return diag.Errorf("[DEBUG] Error while mapping stepTypes to resource for READ. Error = %v", err)
	}

	return nil
}

func resourceStepTypesVersionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*cfClient.Client)
	name := d.Get("name").(string)
	stepTypesVersions := mapResourceToStepTypesVersions(d)
	mapVersionToCreate := make(map[string]cfClient.StepTypes)
	versionsPreviouslyDefined := make(map[string]string)
	versionsDefined := make(map[string]string)
	// Name is set to ForceNew so if we reach this function "version" is changed. Skipping check on HasChange
	// Retrieving old version of the resource to enable comparsion with new and determine which versions should be removed
	old, _ := d.GetChange("version")

	for _, oldStep := range old.(*schema.Set).List() {
		oldVersion := oldStep.(map[string]interface{})["version_number"].(string)
		versionsPreviouslyDefined[oldVersion] = oldVersion
	}

	// Parse current set: new versions that need to be created are added to a data structure
	// that will be sorted later for the creation
	// Updates are performed immediately
	for _, version := range stepTypesVersions.Versions {
		versionNumber := version.VersionNumber
		versionsDefined[versionNumber] = versionNumber

		_, err := client.GetStepTypes(name + ":" + versionNumber)
		cleanUpStepFromTransientValues(&version.StepTypes, name, versionNumber)
		if err != nil {
			// If an error occured during Get, we assume step doesn't exist
			log.Printf("[DEBUG] Recording for creation: %q", versionNumber)
			mapVersionToCreate[versionNumber] = version.StepTypes
		} else {
			log.Printf("[DEBUG] Update Version step: %q", versionNumber)
			_, err := client.UpdateStepTypes(&version.StepTypes)
			if err != nil {
				return diag.Errorf("[DEBUG] Error while updating stepTypes. Error = %v", err)

			}
		}
	}

	// Order versions for creation
	createVersions := make([]string, len(mapVersionToCreate))
	i := 0
	for k := range mapVersionToCreate {
		createVersions[i] = k
		i++
	}
	orderedVersions := sortVersions(createVersions)
	for _, version := range orderedVersions {
		step := mapVersionToCreate[version.String()]
		log.Printf("[DEBUG] Creating version %s for step types: %s", step.Metadata["version"], step.Metadata["name"])
		_, err := client.CreateStepTypes(&step)
		if err != nil {
			return diag.Errorf("[DEBUG] Error while creating step types OnUpdate function. Error = %v", err)
		}
	}

	// If a version is not listed in versionsDefined we can remove it from the system
	for version := range versionsPreviouslyDefined {
		if _, ok := versionsDefined[version]; !ok {
			log.Printf("[DEBUG] Deleting version: %s", version)
			// If not defined we remove from the system
			err := client.DeleteStepTypes(d.Id() + ":" + version)
			if err != nil {
				return diag.Errorf("[DEBUG] Error while deleting step_types_versions. Error = %v", err)
			}
		}
	}

	return resourceStepTypesVersionRead(ctx, d, meta)
}

func resourceStepTypesVersionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*cfClient.Client)
	log.Printf("[DEBUG] Deleting step type: %s", d.Id())
	err := client.DeleteStepTypes(d.Id())
	if err != nil {
		return diag.Errorf("[DEBUG] Error while deleting step_types %s. Error = %v", d.Id(), err)
	}

	return nil
}

func cleanUpStepFromTransientValues(stepTypes *cfClient.StepTypes, name, version string) {
	if stepTypes != nil {
		// Remove transient attributes from metadata
		for _, attribute := range []string{"created_at", "accountId", "id", "updated_at"} {
			if _, ok := stepTypes.Metadata[attribute]; ok {
				delete(stepTypes.Metadata, attribute)
			}
		}
		// Forcing latest to false
		// This is needed because in some cases (e.g. adding an old version) the latest attribute is set to `null` by Codefresh
		// Having `null` as value causes subsequent calls to fail validation against this attribute
		stepTypes.Metadata["latest"] = false

		// If name of version are empty strings we remove them from the data structure
		// The use case is for the calculation of the Hash of the Set item, where we don't have access to this information.
		// Since that is coming from the other attribute of the resource there's no point to actually consider it for hashing
		if name != "" {
			stepTypes.Metadata["name"] = name
		} else {
			delete(stepTypes.Metadata, "name")
		}
		if version != "" {
			stepTypes.Metadata["version"] = version
		} else {
			delete(stepTypes.Metadata, "version")
		}

	}

}

func sortVersions(versions []string) []*semver.Version {
	log.Printf("[DEBUG] Sorting: %q", versions)
	var vs []*semver.Version
	for _, version := range versions {
		v, err := semver.NewVersion(version)
		if err != nil {
			diag.Errorf("Error parsing version: %s", err)
		}
		vs = append(vs, v)
	}

	sort.Sort(semver.Collection(vs))
	return vs
}

func mapStepTypesVersionsToResource(stepTypesVersions cfClient.StepTypesVersions, d *schema.ResourceData) error {

	err := d.Set("name", stepTypesVersions.Name)
	if err != nil {
		return err
	}

	err = d.Set("version", flattenVersions(stepTypesVersions.Name, stepTypesVersions.Versions))
	return err

}

func resourceStepTypesVersionsConfigHash(v interface{}) int {

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s", m["version_number"].(string)))
	var stepTypes cfClient.StepTypes
	stepTypesYaml := m["step_types_yaml"].(string)
	yaml.Unmarshal([]byte(stepTypesYaml), &stepTypes)
	// Remove runtime attributes, name and version to avoid discrepancies when comparing hashes
	cleanUpStepFromTransientValues(&stepTypes, "", "")
	stepTypesYamlByteArray, _ := yaml.Marshal(stepTypes)
	buf.WriteString(fmt.Sprintf("%s", string(stepTypesYamlByteArray)))
	hash := hashcode.String(buf.String())
	return hash
}

func flattenVersions(name string, versions []cfClient.StepTypesVersion) *schema.Set {

	stepVersions := make([]interface{}, 0)
	for _, version := range versions {
		m := make(map[string]interface{})
		m["version_number"] = version.VersionNumber
		cleanUpStepFromTransientValues(&version.StepTypes, name, version.VersionNumber)
		stepTypesYaml, _ := yaml.Marshal(version.StepTypes)
		m["step_types_yaml"] = string(stepTypesYaml)
		stepVersions = append(stepVersions, m)
	}

	return schema.NewSet(resourceStepTypesVersionsConfigHash, stepVersions)
}

func mapResourceToStepTypesVersions(d *schema.ResourceData) *cfClient.StepTypesVersions {
	var stepTypesVersions cfClient.StepTypesVersions
	stepTypesVersions.Name = d.Get("name").(string)
	versions := d.Get("version").(*schema.Set)

	for _, step := range versions.List() {
		version := step.(map[string]interface{})["version_number"].(string)
		if version != "" {
			var stepTypes cfClient.StepTypes
			stepTypesYaml := step.(map[string]interface{})["step_types_yaml"].(string)
			yaml.Unmarshal([]byte(stepTypesYaml), &stepTypes)
			cleanUpStepFromTransientValues(&stepTypes, stepTypesVersions.Name, version)
			stepVersion := cfClient.StepTypesVersion{
				VersionNumber: version,
				StepTypes:     stepTypes,
			}

			stepTypesVersions.Versions = append(stepTypesVersions.Versions, stepVersion)
		}
	}

	return &stepTypesVersions
}
