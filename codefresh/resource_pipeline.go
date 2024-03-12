package codefresh

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var terminationPolicyOnCreateBranchAttributes = []string{"branchName", "ignoreTrigger", "ignoreBranch"}

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Description: "The central component of the Codefresh Platform. Pipelines are workflows that contain individual steps. Each step is responsible for a specific action in the process.",
		Create:      resourcePipelineCreate,
		Read:        resourcePipelineRead,
		Update:      resourcePipelineUpdate,
		Delete:      resourcePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The display name for the pipeline.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"original_yaml_string": {
				Description: `
A string with original yaml pipeline.

For example:

<code>original_yaml_string = "version: \\"1.0\\"\nsteps:\n	test:\n	image: alpine:latest\n	commands:\n	- echo \\"ACC tests\\"</code>

Or: <code>original_yaml_string = file("/path/to/my/codefresh.yml")</code>
				`,
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Description: "The ID of the project that the pipeline belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_public": {
				Description: "Boolean that specifies if the build logs are publicly accessible (default: `false`).",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"revision": {
				Description: "The pipeline's revision. Should be added to the **lifecycle/ignore_changes** or incremented mannually each update.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"tags": {
				Description: "A list of tags to mark a project for easy management and access control.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spec": {
				Description: "The pipeline's specs.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Description: "Helps to organize the order of builds execution in case of reaching the concurrency limit (default: `0`).",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"concurrency": {
							Description: "The maximum amount of concurrent builds. Zero is unlimited (default: `0`).",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"branch_concurrency": {
							Description: "The maximum amount of concurrent builds that may run for each branch. Zero is unlimited (default: `0`).",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"trigger_concurrency": {
							Description: "The maximum amount of concurrent builds that may run for each trigger (default: `0`).",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"permit_restart_from_failed_steps": {
							Description: "Defines whether it is permitted to restart builds in this pipeline from failed step. Defaults to true",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"spec_template": {
							Description: "The pipeline's spec template.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Description: "The location of the spec template (default: `git`).",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "git",
									},
									"repo": {
										Description: "The repository of the spec template (owner/repo).",
										Type:        schema.TypeString,
										Required:    true,
									},
									"path": {
										Description: "The relative path to the Codefresh pipeline file.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"revision": {
										Description: "The git revision of the spec template. Possible values: '', *name of branch*. Use '' to autoselect a branch.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"context": {
										Description: "The Codefresh git context (default: `github`).",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "github",
									},
								},
							},
						},
						"variables": {
							Description: "The pipeline's variables.",
							Type:        schema.TypeMap,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"encrypted_variables": {
							Description: "Pipeline level encrypted variables. Please note that drift will not be detected for encrypted variables",
							Type:        schema.TypeMap,
							Optional:    true,
							Elem: &schema.Schema{
								Type:      schema.TypeString,
								Sensitive: true,
							},
						},
						"trigger": {
							Description: "The pipeline's triggers (currently the only nested trigger supported is git; for other trigger types, use the `codefresh_pipeline_*_trigger` resources).",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "The name of the trigger.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"description": {
										Description: "The description of the trigger.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"type": {
										Description: "The type of the trigger (default: `git`; see notes above).",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "git",
										ValidateDiagFunc: schemautil.StringMatchesRegExp(
											"git",
											schemautil.WithSummary("Invalid trigger type"),
											schemautil.WithDetailFormat("The trigger type %s is invalid. The only supported type is %s."),
										),
									},
									"repo": {
										Description: "The repository name, (owner/repo)",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"branch_regex": {
										Description:      " A regular expression and will only trigger for branches that match this naming pattern (default: `/.*/gi`).",
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "/.*/gi",
										ValidateDiagFunc: schemautil.StringIsValidRegExp(),
									},
									"branch_regex_input": {
										Description:  "Flag to manage how the `branch_regex` field is interpreted. Possible values: `multiselect-exclude`, `multiselect`, `regex` (default: `regex`).",
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "regex",
										ValidateFunc: validation.StringInSlice([]string{"multiselect-exclude", "multiselect", "regex"}, false),
									},
									"pull_request_target_branch_regex": {
										Description:      "A regular expression and will only trigger for pull requests to branches that match this naming pattern.",
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: schemautil.StringIsValidRegExp(),
									},
									"comment_regex": {
										Description:      " A regular expression and will only trigger for pull requests where a comment matches this naming pattern (default: `/.*/gi`).",
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "/.*/gi",
										ValidateDiagFunc: schemautil.StringIsValidRegExp(),
									},
									"modified_files_glob": {
										Description: "Allows to constrain the build and trigger it only if the modified files from the commit match this glob expression (default: `\"\"`).",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
									},
									"events": {
										Description: "A list of GitHub events for which a Pipeline is triggered.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"provider": {
										Description: "The git provider tied to the trigger.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "github",
									},
									"disabled": {
										Description: "Flag to disable the trigger.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"options": {
										Description: "The trigger's options.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"no_cache": {
													Description: "If true, docker layer cache is disabled",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"no_cf_cache": {
													Description: "If true, extra Codefresh caching is disabled.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"reset_volume": {
													Description: "If true, all files on volume will be deleted before each execution.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"enable_notifications": {
													Description: "If false the pipeline will not send notifications to Slack and status updates back to the Git provider.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
											},
										},
									},
									"pull_request_allow_fork_events": {
										Description: "If this trigger is also applicable to git forks.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"commit_status_title": {
										Description: "The commit status title pushed to the git provider.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"context": {
										Description: "The Codefresh git context.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "github",
									},
									"contexts": {
										Description: "A list of strings representing the contexts ([shared_configuration](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/)) to be loaded when the trigger is executed.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"runtime_environment": {
										Description: "The runtime environment for the trigger.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Description: "The name of the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"memory": {
													Description: "The memory allocated to the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"cpu": {
													Description: "The CPU allocated to the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"dind_storage": {
													Description: "The storage allocated to the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"required_available_storage": {
													Description: "Minimum disk space required for build filesystem ( unit Gi is required).",
													Type:        schema.TypeString,
													Optional:    true,
												},
											},
										},
									},
									"variables": {
										Description: "Trigger variables.",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"encrypted_variables": {
										Description: "Trigger level encrypted variables. Please note that drift will not be detected for encrypted variables",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem: &schema.Schema{
											Type:      schema.TypeString,
											Sensitive: true,
										},
									},
								},
							},
						},
						"cron_trigger": {
							Description: "The pipeline's cron triggers. Conflicts with the deprecated [codefresh_pipeline_cron_trigger](https://registry.terraform.io/providers/codefresh-io/codefresh/latest/docs/resources/pipeline_cron_trigger) resource.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "The name of the cron trigger.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"type": {
										Description: "The type of the trigger (default: `cron`; see notes above).",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "cron",
										ValidateDiagFunc: schemautil.StringMatchesRegExp(
											"cron",
											schemautil.WithSummary("Invalid cron trigger type"),
											schemautil.WithDetailFormat("The cron trigger type %s is invalid. The only supported type is %s."),
										),
									},
									"disabled": {
										Description: "Flag to disable the trigger.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"expression": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: schemautil.CronExpression(),
									},
									"message": {
										Type:     schema.TypeString,
										Required: true,
										ValidateDiagFunc: schemautil.StringMatchesRegExp(
											schemautil.ValidCronMessageRegex,
											schemautil.WithSeverity(diag.Error),
											schemautil.WithSummary("Invalid cron trigger message"),
											schemautil.WithDetailFormat("The message %q is invalid (must match %q)."),
										),
									},
									"git_trigger_id": {
										Description: "Related git-trigger id. Will by used to take all possible git information by branch.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"branch": {
										Description: "Branch that should be passed for build triggered by this cron trigger.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"options": {
										Description: "The trigger's options.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"no_cache": {
													Description: "If true, docker layer cache is disabled.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"no_cf_cache": {
													Description: "If true, extra Codefresh caching is disabled.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"reset_volume": {
													Description: "If true, all files on volume will be deleted before each execution.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"enable_notifications": {
													Description: "If false the pipeline will not send notifications to Slack and status updates back to the Git provider.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
											},
										},
									},
									"runtime_environment": {
										Description: "The runtime environment for the trigger.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Description: "The name of the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"memory": {
													Description: "The memory allocated to the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"cpu": {
													Description: "The CPU allocated to the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"dind_storage": {
													Description: "The storage allocated to the runtime environment.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"required_available_storage": {
													Description: "Minimum disk space required for build filesystem ( unit Gi is required).",
													Type:        schema.TypeString,
													Optional:    true,
												},
											},
										},
									},
									"variables": {
										Description: "Trigger variables.",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"encrypted_variables": {
										Description: "Trigger level encrypted variables. Please note that drift will not be detected for encrypted variables",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem: &schema.Schema{
											Type:      schema.TypeString,
											Sensitive: true,
										},
									},
								},
							},
						},
						"contexts": {
							Description: "A list of strings representing the contexts ([shared_configuration](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/)) to be configured for the pipeline.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"termination_policy": {
							Description: "The termination policy for the pipeline.",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"on_create_branch": {
										Description: `
The following table presents how to configure this block based on the options available in the UI:

| Option Description                                                            | Value Selected           | on_create_branch | branch_name | ignore_trigger | ignore_branch |
| ----------------------------------------------------------------------------- |:------------------------:|:----------------:|:-----------:|---------------:| -------------:|
| Once a build is created terminate previous builds from the same branch        | Disabled                 |        Omit      |     N/A     |       N/A      |      N/A      |
| Once a build is created terminate previous builds from the same branch        | From the SAME trigger    |       Defined    |     N/A     |      false     |      N/A      |
| Once a build is created terminate previous builds from the same branch        | From ANY trigger         |       Defined    |     N/A     |      true      |      N/A      |
| Once a build is created terminate previous builds only from a specific branch | Disabled                 |        Omit      |     N/A     |       N/A      |      N/A      |
| Once a build is created terminate previous builds only from a specific branch | From the SAME trigger    |       Defined    |    Regex    |      false     |      N/A      |
| Once a build is created terminate previous builds only from a specific branch | From ANY trigger         |       Defined    |    Regex    |      true      |      N/A      |
| Once a build is created, terminate all other running builds                   | Disabled                 |        Omit      |     N/A     |       N/A      |      N/A      |
| Once a build is created, terminate all other running builds                   | From the SAME trigger    |       Defined    |     N/A     |      false     |      true     |
| Once a build is created, terminate all other running builds                   | From ANY trigger         |       Defined    |     N/A     |      true      |      true     |
										`,
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"branch_name": {
													Description:      "A regular expression to filter the branches on with the termination policy applies.",
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: schemautil.StringIsValidRegExp(),
													ConflictsWith:    []string{"spec.0.termination_policy.0.on_create_branch.0.ignore_branch"},
												},
												"ignore_trigger": {
													Description: "Whether to ignore the trigger.",
													Optional:    true,
													Type:        schema.TypeBool,
												},
												"ignore_branch": {
													Description: "Whether to ignore the branch.",
													Optional:    true,
													Type:        schema.TypeBool,
												},
											},
										},
									},
									"on_terminate_annotation": {
										Description: "Enables the policy `Once a build is terminated, terminate all child builds initiated from it`.",
										Optional:    true,
										Type:        schema.TypeBool,
										Default:     false,
									},
								},
							},
						},
						"pack_id": {
							Description: "SAAS pack (`5cd1746617313f468d669013` for Small; `5cd1746717313f468d669014` for Medium; `5cd1746817313f468d669015` for Large; `5cd1746817313f468d669017` for XL; `5cd1746817313f468d669018` for XXL); `5cd1746817313f468d669020` for 4XL).",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"required_available_storage": {
							Description: " Minimum disk space required for build filesystem ( unit Gi is required).",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"runtime_environment": {
							Description: "The runtime environment for the pipeline.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "The name of the runtime environment.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"memory": {
										Description: "The memory allocated to the runtime environment.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"cpu": {
										Description: "The CPU allocated to the runtime environment.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"dind_storage": {
										Description: "The storage allocated to the runtime environment.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"required_available_storage": {
										Description: "Minimum disk space required for build filesystem ( unit Gi is required).",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"options": {
							Description: "The options for the pipeline.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"keep_pvcs_for_pending_approval": {
										Description: `
When build enters 'Pending Approval' state, volume should:
	* Default (attribute not specified): "Use Setting accounts"
	* true: "Remain (build remains active)"
	* false: "Be removed"
										`,
										Type:     schema.TypeBool,
										Optional: true,
									},
									"pending_approval_concurrency_applied": {
										Description: `
Pipeline concurrency policy: Builds on 'Pending Approval' state should be:
	* Default (attribute not specified): "Use Setting accounts"
	* true: "Included in concurrency"
	* false: "Not included in concurrency"
										`,
										Type:     schema.TypeBool,
										Optional: true,
									},
									"enable_notifications": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
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

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	pipeline, err := mapResourceToPipeline(d)
	if err != nil {
		return err
	}

	resp, err := client.CreatePipeline(pipeline)
	if err != nil {
		return err
	}

	d.SetId(resp.Metadata.ID)

	return resourcePipelineRead(d, meta)
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	pipelineID := d.Id()

	if pipelineID == "" {
		d.SetId("")
		return nil
	}

	pipeline, err := client.GetPipeline(pipelineID)
	if err != nil {
		return err
	}

	err = mapPipelineToResource(*pipeline, d)
	if err != nil {
		return err
	}

	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	pipeline, err := mapResourceToPipeline(d)
	if err != nil {
		return err
	}

	pipeline.Metadata.ID = d.Id()

	_, err = client.UpdatePipeline(pipeline)
	if err != nil {
		return err
	}

	return resourcePipelineRead(d, meta)
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfclient.Client)

	err := client.DeletePipeline(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapPipelineToResource(pipeline cfclient.Pipeline, d *schema.ResourceData) error {

	err := d.Set("name", pipeline.Metadata.Name)
	if err != nil {
		return err
	}

	err = d.Set("revision", pipeline.Metadata.Revision)
	if err != nil {
		return err
	}

	err = d.Set("project_id", pipeline.Metadata.ProjectId)
	if err != nil {
		return err
	}

	err = d.Set("is_public", pipeline.Metadata.IsPublic)
	if err != nil {
		return err
	}

	flattenedSpec := flattenSpec(pipeline.Spec)

	// Set encrypted variables from resource data, as otherwise they cause constant diff as the value is always returned as *****
	encryptedVariables, ok := flattenedSpec[0]["encrypted_variables"].(map[string]string)

	if ok {
		if len(encryptedVariables) > 0 {
			setEncryptedVariablesValuesFromResource(d, encryptedVariables, "spec.0.encrypted_variables")
		}
	}

	// Set trigger encrypted variables from resource data
	triggers, getTriggersOK := flattenedSpec[0]["trigger"]

	if getTriggersOK {
		for triggerIndex, triggerSpec := range triggers.([]map[string]interface{}) {

			triggerEncryptedVariables, ok := triggerSpec["encrypted_variables"].(map[string]string)

			if ok {
				if len(triggerEncryptedVariables) > 0 {
					setEncryptedVariablesValuesFromResource(d, triggerEncryptedVariables, fmt.Sprintf("spec.0.trigger.%d.encrypted_variables", triggerIndex))
				}
			}
		}
	}

	// Set cron trigger encrypted variables from resource data
	cronTriggers, getCronTriggersOK := flattenedSpec[0]["cron_trigger"]

	if getCronTriggersOK {
		for triggerIndex, triggerSpec := range cronTriggers.([]map[string]interface{}) {

			triggerEncryptedVariables, ok := triggerSpec["encrypted_variables"].(map[string]string)

			if ok {
				if len(triggerEncryptedVariables) > 0 {
					setEncryptedVariablesValuesFromResource(d, triggerEncryptedVariables, fmt.Sprintf("spec.0.cron_trigger.%d.encrypted_variables", triggerIndex))
				}
			}
		}
	}

	err = d.Set("spec", flattenedSpec)

	if err != nil {
		return err
	}

	err = d.Set("tags", pipeline.Metadata.Labels.Tags)
	if err != nil {
		return err
	}

	err = d.Set("original_yaml_string", pipeline.Metadata.OriginalYamlString)
	if err != nil {
		return err
	}

	return nil
}

func flattenSpec(spec cfclient.Spec) []map[string]interface{} {

	var res = make([]map[string]interface{}, 0)
	m := make(map[string]interface{})

	if len(spec.Triggers) > 0 {
		m["trigger"] = flattenTriggers(spec.Triggers)
	}

	if len(spec.CronTriggers) > 0 {
		m["cron_trigger"] = flattenCronTriggers(spec.CronTriggers)
	}

	if spec.SpecTemplate != nil {
		m["spec_template"] = flattenSpecTemplate(*spec.SpecTemplate)
	}

	if len(spec.Variables) != 0 {
		// Do not set encrypted variables because they cause constant diff
		m["variables"], m["encrypted_variables"] = datautil.ConvertVariables(spec.Variables)
	}

	if spec.RuntimeEnvironment != (cfclient.RuntimeEnvironment{}) {
		m["runtime_environment"] = flattenSpecRuntimeEnvironment(spec.RuntimeEnvironment)
	}

	if len(spec.TerminationPolicy) > 0 {
		m["termination_policy"] = flattenSpecTerminationPolicy(spec.TerminationPolicy)
	}

	if len(spec.Options) > 0 {
		var resOptions []map[string]bool
		options := map[string]bool{}
		for keyOption, valueOption := range spec.Options {
			switch {
			case keyOption == "keepPVCsForPendingApproval":
				options["keep_pvcs_for_pending_approval"] = valueOption
			case keyOption == "pendingApprovalConcurrencyApplied":
				options["pending_approval_concurrency_applied"] = valueOption
			case keyOption == "enableNotifications":
				options["enable_notifications"] = valueOption
			}
		}
		resOptions = append(resOptions, options)
		m["options"] = resOptions
	}

	m["pack_id"] = spec.PackId
	m["required_available_storage"] = spec.RequiredAvailableStorage

	m["concurrency"] = spec.Concurrency
	m["branch_concurrency"] = spec.BranchConcurrency
	m["trigger_concurrency"] = spec.TriggerConcurrency
	m["permit_restart_from_failed_steps"] = spec.PermitRestartFromFailedSteps

	m["priority"] = spec.Priority

	m["contexts"] = spec.Contexts

	res = append(res, m)
	return res
}

func flattenSpecTerminationPolicy(terminationPolicy []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	attribute := make(map[string]interface{})
	for _, policy := range terminationPolicy {
		eventName, _ := policy["event"]
		typeName, _ := policy["type"]
		attributeName := convertOnCreateBranchAttributeToPipelineFormat(eventName.(string) + "_" + typeName.(string))
		switch attributeName {
		case "on_create_branch":
			var valueList []map[string]interface{}
			attributeValues := make(map[string]interface{})
			for _, eventAttribute := range terminationPolicyOnCreateBranchAttributes {
				if item, ok := policy[eventAttribute]; ok {
					attributeValues[convertOnCreateBranchAttributeToPipelineFormat(eventAttribute)] = item
				}
			}
			attribute[attributeName] = append(valueList, attributeValues)
		case "on_terminate_annotation":
			if value, ok := policy["key"]; ok && value == "cf_predecessor" {
				attribute[attributeName] = true
			}
		default:
			log.Fatal("Unsupported event found in TerminationPolicy")
		}
	}
	res = append(res, attribute)
	return res
}

func flattenSpecTemplate(spec cfclient.SpecTemplate) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"location": spec.Location,
			"repo":     spec.Repo,
			"context":  spec.Context,
			"revision": spec.Revision,
			"path":     spec.Path,
		},
	}
}

func flattenSpecRuntimeEnvironment(spec cfclient.RuntimeEnvironment) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":                       spec.Name,
			"memory":                     spec.Memory,
			"cpu":                        spec.CPU,
			"dind_storage":               spec.DindStorage,
			"required_available_storage": spec.RequiredAvailableStorage,
		},
	}
}

func flattenTriggerOptions(options cfclient.TriggerOptions) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"no_cache":             options.NoCache,
			"no_cf_cache":          options.NoCfCache,
			"reset_volume":         options.ResetVolume,
			"enable_notifications": options.EnableNotifications,
		},
	}
}

func flattenTriggers(triggers []cfclient.Trigger) []map[string]interface{} {
	var res = make([]map[string]interface{}, len(triggers))
	for i, trigger := range triggers {
		m := make(map[string]interface{})
		m["name"] = trigger.Name
		m["description"] = trigger.Description
		m["context"] = trigger.Context
		m["contexts"] = trigger.Contexts
		m["repo"] = trigger.Repo
		m["branch_regex"] = trigger.BranchRegex
		m["branch_regex_input"] = trigger.BranchRegexInput
		m["pull_request_target_branch_regex"] = trigger.PullRequestTargetBranchRegex
		m["comment_regex"] = trigger.CommentRegex
		m["modified_files_glob"] = trigger.ModifiedFilesGlob
		m["disabled"] = trigger.Disabled
		if trigger.Options != nil {
			m["options"] = flattenTriggerOptions(*trigger.Options)
		}
		m["pull_request_allow_fork_events"] = trigger.PullRequestAllowForkEvents
		m["commit_status_title"] = trigger.CommitStatusTitle
		m["provider"] = trigger.Provider
		m["type"] = trigger.Type
		m["events"] = trigger.Events
		m["variables"], m["encrypted_variables"] = datautil.ConvertVariables(trigger.Variables)
		if trigger.RuntimeEnvironment != nil {
			m["runtime_environment"] = flattenSpecRuntimeEnvironment(*trigger.RuntimeEnvironment)
		}
		res[i] = m
	}
	return res
}

func flattenCronTriggers(cronTriggers []cfclient.CronTrigger) []map[string]interface{} {
	var res = make([]map[string]interface{}, len(cronTriggers))
	for i, trigger := range cronTriggers {
		m := make(map[string]interface{})
		m["name"] = trigger.Name
		m["type"] = trigger.Type
		m["expression"] = trigger.Expression
		m["message"] = trigger.Message
		m["disabled"] = trigger.Disabled
		m["git_trigger_id"] = trigger.GitTriggerId
		m["branch"] = trigger.Branch
		m["variables"], m["encrypted_variables"] = datautil.ConvertVariables(trigger.Variables)
		if trigger.Options != nil {
			m["options"] = flattenTriggerOptions(*trigger.Options)
		}
		if trigger.RuntimeEnvironment != nil {
			m["runtime_environment"] = flattenSpecRuntimeEnvironment(*trigger.RuntimeEnvironment)
		}
		res[i] = m
	}
	return res
}

func mapResourceToPipeline(d *schema.ResourceData) (*cfclient.Pipeline, error) {

	tags := d.Get("tags").(*schema.Set).List()

	originalYamlString := strings.Replace(
		d.Get("original_yaml_string").(string),
		"\n",
		"\n",
		-1)
	pipeline := &cfclient.Pipeline{
		Metadata: cfclient.Metadata{
			Name:      d.Get("name").(string),
			Revision:  d.Get("revision").(int),
			ProjectId: d.Get("project_id").(string),
			IsPublic:  d.Get("is_public").(bool),
			Labels: cfclient.Labels{
				Tags: datautil.ConvertStringArr(tags),
			},
			OriginalYamlString: originalYamlString,
		},
		Spec: cfclient.Spec{
			PackId:                       d.Get("spec.0.pack_id").(string),
			RequiredAvailableStorage:     d.Get("spec.0.required_available_storage").(string),
			Priority:                     d.Get("spec.0.priority").(int),
			Concurrency:                  d.Get("spec.0.concurrency").(int),
			BranchConcurrency:            d.Get("spec.0.branch_concurrency").(int),
			TriggerConcurrency:           d.Get("spec.0.trigger_concurrency").(int),
			PermitRestartFromFailedSteps: d.Get("spec.0.permit_restart_from_failed_steps").(bool),
		},
	}

	if _, ok := d.GetOk("spec.0.spec_template"); ok {
		pipeline.Spec.SpecTemplate = &cfclient.SpecTemplate{
			Location: d.Get("spec.0.spec_template.0.location").(string),
			Repo:     d.Get("spec.0.spec_template.0.repo").(string),
			Path:     d.Get("spec.0.spec_template.0.path").(string),
			Revision: d.Get("spec.0.spec_template.0.revision").(string),
			Context:  d.Get("spec.0.spec_template.0.context").(string),
		}
	} else {
		err := extractSpecAttributesFromOriginalYamlString(originalYamlString, pipeline)
		if err != nil {
			return nil, err
		}
	}

	if _, ok := d.GetOk("spec.0.runtime_environment"); ok {
		pipeline.Spec.RuntimeEnvironment = cfclient.RuntimeEnvironment{
			Name:                     d.Get("spec.0.runtime_environment.0.name").(string),
			Memory:                   d.Get("spec.0.runtime_environment.0.memory").(string),
			CPU:                      d.Get("spec.0.runtime_environment.0.cpu").(string),
			DindStorage:              d.Get("spec.0.runtime_environment.0.dind_storage").(string),
			RequiredAvailableStorage: d.Get("spec.0.runtime_environment.0.required_available_storage").(string),
		}
	}

	if contexts, ok := d.GetOk("spec.0.contexts"); ok {
		pipeline.Spec.Contexts = contexts.([]interface{})
	}

	if variables, ok := d.GetOk("spec.0.variables"); ok {
		pipeline.SetVariables(variables.(map[string]interface{}), false)
	}

	if encryptedVariables, ok := d.GetOk("spec.0.encrypted_variables"); ok {
		pipeline.SetVariables(encryptedVariables.(map[string]interface{}), true)
	}

	if triggers, ok := d.GetOk("spec.0.trigger"); ok {
		for idx := range triggers.([]interface{}) {
			events := d.Get(fmt.Sprintf("spec.0.trigger.%v.events", idx)).([]interface{})
			contexts := d.Get(fmt.Sprintf("spec.0.trigger.%v.contexts", idx)).([]interface{})
			codefreshTrigger := cfclient.Trigger{
				Name:                         d.Get(fmt.Sprintf("spec.0.trigger.%v.name", idx)).(string),
				Description:                  d.Get(fmt.Sprintf("spec.0.trigger.%v.description", idx)).(string),
				Type:                         d.Get(fmt.Sprintf("spec.0.trigger.%v.type", idx)).(string),
				Repo:                         d.Get(fmt.Sprintf("spec.0.trigger.%v.repo", idx)).(string),
				BranchRegex:                  d.Get(fmt.Sprintf("spec.0.trigger.%v.branch_regex", idx)).(string),
				BranchRegexInput:             d.Get(fmt.Sprintf("spec.0.trigger.%v.branch_regex_input", idx)).(string),
				PullRequestTargetBranchRegex: d.Get(fmt.Sprintf("spec.0.trigger.%v.pull_request_target_branch_regex", idx)).(string),
				CommentRegex:                 d.Get(fmt.Sprintf("spec.0.trigger.%v.comment_regex", idx)).(string),
				ModifiedFilesGlob:            d.Get(fmt.Sprintf("spec.0.trigger.%v.modified_files_glob", idx)).(string),
				Provider:                     d.Get(fmt.Sprintf("spec.0.trigger.%v.provider", idx)).(string),
				Disabled:                     d.Get(fmt.Sprintf("spec.0.trigger.%v.disabled", idx)).(bool),
				PullRequestAllowForkEvents:   d.Get(fmt.Sprintf("spec.0.trigger.%v.pull_request_allow_fork_events", idx)).(bool),
				CommitStatusTitle:            d.Get(fmt.Sprintf("spec.0.trigger.%v.commit_status_title", idx)).(string),
				Context:                      d.Get(fmt.Sprintf("spec.0.trigger.%v.context", idx)).(string),
				Contexts:                     datautil.ConvertStringArr(contexts),
				Events:                       datautil.ConvertStringArr(events),
			}
			variables := d.Get(fmt.Sprintf("spec.0.trigger.%v.variables", idx)).(map[string]interface{})
			codefreshTrigger.SetVariables(variables, false)

			encryptedVariables := d.Get(fmt.Sprintf("spec.0.trigger.%v.encrypted_variables", idx)).(map[string]interface{})
			codefreshTrigger.SetVariables(encryptedVariables, true)

			if _, ok := d.GetOk(fmt.Sprintf("spec.0.trigger.%v.options", idx)); ok {
				options := cfclient.TriggerOptions{
					NoCache:             d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.no_cache", idx)).(bool),
					NoCfCache:           d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.no_cf_cache", idx)).(bool),
					ResetVolume:         d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.reset_volume", idx)).(bool),
					EnableNotifications: d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.enable_notifications", idx)).(bool),
				}
				codefreshTrigger.Options = &options
			}
			if _, ok := d.GetOk(fmt.Sprintf("spec.0.trigger.%v.runtime_environment", idx)); ok {
				triggerRuntime := cfclient.RuntimeEnvironment{
					Name:                     d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.name", idx)).(string),
					Memory:                   d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.memory", idx)).(string),
					CPU:                      d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.cpu", idx)).(string),
					DindStorage:              d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.dind_storage", idx)).(string),
					RequiredAvailableStorage: d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.required_available_storage", idx)).(string),
				}
				codefreshTrigger.RuntimeEnvironment = &triggerRuntime
			}
			pipeline.Spec.Triggers = append(pipeline.Spec.Triggers, codefreshTrigger)
		}
	}

	if cronTriggers, ok := d.GetOk("spec.0.cron_trigger"); ok {
		for idx := range cronTriggers.([]interface{}) {
			codefreshCronTrigger := cfclient.CronTrigger{
				Name:         d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.name", idx)).(string),
				Type:         d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.type", idx)).(string),
				Expression:   d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.expression", idx)).(string),
				Message:      d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.message", idx)).(string),
				Disabled:     d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.disabled", idx)).(bool),
				GitTriggerId: d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.git_trigger_id", idx)).(string),
				Branch:       d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.branch", idx)).(string),
			}
			variables := d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.variables", idx)).(map[string]interface{})
			codefreshCronTrigger.SetVariables(variables, false)
			encryptedVariables := d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.encrypted_variables", idx)).(map[string]interface{})
			codefreshCronTrigger.SetVariables(encryptedVariables, true)

			if _, ok := d.GetOk(fmt.Sprintf("spec.0.cron_trigger.%v.options", idx)); ok {
				options := cfclient.TriggerOptions{
					NoCache:             d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.options.0.no_cache", idx)).(bool),
					NoCfCache:           d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.options.0.no_cf_cache", idx)).(bool),
					ResetVolume:         d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.options.0.reset_volume", idx)).(bool),
					EnableNotifications: d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.options.0.enable_notifications", idx)).(bool),
				}
				codefreshCronTrigger.Options = &options
			}
			if _, ok := d.GetOk(fmt.Sprintf("spec.0.cron_trigger.%v.runtime_environment", idx)); ok {
				triggerRuntime := cfclient.RuntimeEnvironment{
					Name:                     d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.runtime_environment.0.name", idx)).(string),
					Memory:                   d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.runtime_environment.0.memory", idx)).(string),
					CPU:                      d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.runtime_environment.0.cpu", idx)).(string),
					DindStorage:              d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.runtime_environment.0.dind_storage", idx)).(string),
					RequiredAvailableStorage: d.Get(fmt.Sprintf("spec.0.cron_trigger.%v.runtime_environment.0.required_available_storage", idx)).(string),
				}
				codefreshCronTrigger.RuntimeEnvironment = &triggerRuntime
			}
			pipeline.Spec.CronTriggers = append(pipeline.Spec.CronTriggers, codefreshCronTrigger)
		}
	}

	var codefreshTerminationPolicy []map[string]interface{}

	if _, ok := d.GetOk("spec.0.termination_policy.0.on_create_branch"); ok {
		var onCreatBranchPolicy = make(map[string]interface{})
		onCreatBranchPolicy = getSupportedTerminationPolicyAttributes("on_create_branch")
		for _, attribute := range terminationPolicyOnCreateBranchAttributes {
			if attributeValue, ok := d.GetOk(fmt.Sprintf("spec.0.termination_policy.0.on_create_branch.0.%s", convertOnCreateBranchAttributeToPipelineFormat(attribute))); ok {
				onCreatBranchPolicy[attribute] = attributeValue
			}
		}
		codefreshTerminationPolicy = append(codefreshTerminationPolicy, onCreatBranchPolicy)
	}
	if _, ok := d.GetOk("spec.0.termination_policy.0.on_terminate_annotation"); ok {
		var onTerminateAnnotationPolicy = make(map[string]interface{})
		onTerminateAnnotationPolicy = getSupportedTerminationPolicyAttributes("on_terminate_annotation")
		onTerminateAnnotationPolicy["key"] = "cf_predecessor"
		codefreshTerminationPolicy = append(codefreshTerminationPolicy, onTerminateAnnotationPolicy)
	}
	if _, ok := d.GetOk("spec.0.options"); ok {
		pipelineSpecOption := make(map[string]bool)
		if keepPVCs, ok := d.GetOkExists("spec.0.options.0.keep_pvcs_for_pending_approval"); ok {
			pipelineSpecOption["keepPVCsForPendingApproval"] = keepPVCs.(bool)
		}
		if pendingApprovalConcurrencyApplied, ok := d.GetOkExists("spec.0.options.0.pending_approval_concurrency_applied"); ok {
			pipelineSpecOption["pendingApprovalConcurrencyApplied"] = pendingApprovalConcurrencyApplied.(bool)
		}
		if enableNotifications, ok := d.GetOkExists("spec.0.options.0.enable_notifications"); ok {
			pipelineSpecOption["enableNotifications"] = enableNotifications.(bool)
		}
		pipeline.Spec.Options = pipelineSpecOption
	} else {
		pipeline.Spec.Options = nil
	}

	pipeline.Spec.TerminationPolicy = codefreshTerminationPolicy

	return pipeline, nil
}

// This function is used to extract the spec attributes from the original_yaml_string attribute.
// Typically, unmarshalling the YAML string is problematic because the order of the attributes is not preserved.
// Namely, we care a lot about the order of the steps and stages attributes.
// Luckily, the yj package introduces a MapSlice type that preserves the order Map items (see utils.go).
func extractSpecAttributesFromOriginalYamlString(originalYamlString string, pipeline *cfclient.Pipeline) error {
	for _, attribute := range []string{"stages", "steps", "hooks"} {
		yamlString, err := datautil.Yq(fmt.Sprintf(".%s", attribute), originalYamlString)
		if err != nil {
			return fmt.Errorf("error while extracting '%s' from original YAML string: %v", attribute, err)
		} else if yamlString == "" {
			continue
		}

		attributeJson, err := datautil.YamlToJson(yamlString)
		if err != nil {
			return fmt.Errorf("error while converting '%s' YAML to JSON: %v", attribute, err)
		}

		switch attribute {
		case "stages":
			pipeline.Spec.Stages = &cfclient.Stages{
				Stages: attributeJson,
			}
		case "steps":
			pipeline.Spec.Steps = &cfclient.Steps{
				Steps: attributeJson,
			}
		case "hooks":
			pipeline.Spec.Hooks = &cfclient.Hooks{
				Hooks: attributeJson,
			}
		}
	}

	mode, err := datautil.Yq(".mode", originalYamlString)
	if err != nil {
		return fmt.Errorf("error while extracting 'mode' from original YAML string: %v", err)
	} else if mode != "" {
		pipeline.Spec.Mode = mode
	}

	ff, err := datautil.Yq(".fail_fast", originalYamlString)
	if err != nil {
		return fmt.Errorf("error while extracting 'mode' from original YAML string: %v", err)
	} else if ff != "" {
		ff_b, err := strconv.ParseBool(strings.TrimSpace(ff))
		if err != nil {
			return fmt.Errorf("error while parsing 'fail_fast' as boolean: %v", err)
		}
		pipeline.Spec.FailFast = &ff_b
	}
	return nil
}

func getSupportedTerminationPolicyAttributes(policy string) map[string]interface{} {
	switch policy {
	case "on_create_branch":
		return map[string]interface{}{"event": "onCreate", "type": "branch"}
	case "on_terminate_annotation":
		return map[string]interface{}{"event": "onTerminate", "type": "annotation"}
	default:
		log.Fatal("Invalid termination policy selected: ", policy)
	}
	return nil
}

func convertOnCreateBranchAttributeToResourceFormat(src string) string {
	var re = regexp.MustCompile(`_[a-z]`)
	return re.ReplaceAllStringFunc(src, func(w string) string {
		return strings.ToUpper(strings.ReplaceAll(w, "_", ""))
	})
}

func convertOnCreateBranchAttributeToPipelineFormat(src string) string {
	var re = regexp.MustCompile(`[A-Z]`)
	return re.ReplaceAllStringFunc(src, func(w string) string {
		return "_" + strings.ToLower(w)
	})
}

func setEncryptedVariablesValuesFromResource(d *schema.ResourceData, flattenedVariables map[string]string, schemaPath string) error {

	if len(flattenedVariables) > 0 {
		// Iterate over variables and set the value from resource data
		for k := range flattenedVariables {
			flattenedVariables[k] = d.Get(fmt.Sprintf("%s.%s", schemaPath, k)).(string)
		}
	}

	return nil
}
