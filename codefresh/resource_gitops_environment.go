package codefresh

import (
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/datautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGitopsEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "An environment in Codefresh GitOps is a logical grouping of one or more Kubernetes clusters and namespaces, representing a deployment context for your Argo CD applications. See [official documentation](https://codefresh.io/docs/gitops/environments/environments-overview/) for more information.",
		Create:      resourceGitopsEnvironmentCreate,
		Read:        resourceGitopsEnvironmentRead,
		Update:      resourceGitopsEnvironmentUpdate,
		Delete:      resourceGitopsEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Environment ID",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the environment. Must be unique per account",
			},
			"kind": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of environment. Possible values: NON_PROD, PROD",
				ValidateFunc: validation.StringInSlice([]string{"NON_PROD", "PROD"}, false),
			},
			"cluster": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target cluster name",
						},
						"runtime_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Runtime name where the target cluster is registered",
						},
						"namespaces": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "List of namespaces in the target cluster",
						},
					},
				},
			},
			"label_pairs": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of labels and values in the format label=value that can be used to assign applications to the environment. Example: ['codefresh.io/environment=prod']",
			},
		},
	}
}

// func resourceGitopsEnvironmentCreate(d *schema.ResourceData, m interface{}) error {
func resourceGitopsEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	environment := mapResourceToGitopsEnvironment(d)
	newEnvironment, err := client.CreateGitopsEnvironment(environment)

	if err != nil {
		return err
	}

	d.SetId(newEnvironment.ID)

	return resourceGitopsEnvironmentRead(d, meta)
}

func resourceGitopsEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	environment := mapResourceToGitopsEnvironment(d)
	_, err := client.UpdateGitopsEnvironment(environment)

	if err != nil {
		return err
	}

	return resourceGitopsEnvironmentRead(d, meta)
}

func resourceGitopsEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	id := d.Id()
	if id == "" {
		d.SetId("")
		return nil
	}

	_, err := client.DeleteGitopsEnvironment(id)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGitopsEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	id := d.Id()
	if id == "" {
		d.SetId("")
		return nil
	}

	environment, err := client.GetGitopsEnvironmentById(id)

	if err != nil {
		return err
	}

	if environment == nil {
		d.SetId("")
		return nil
	}

	return mapGitopsEnvironmentToResource(d, environment)
}

func mapResourceToGitopsEnvironment(d *schema.ResourceData) *cfclient.GitopsEnvironment {

	clusters := expandClusters(d.Get("cluster").([]interface{}))

	labelPairs := []string{}

	if len(d.Get("label_pairs").([]interface{})) > 0 {
		labelPairs = datautil.ConvertStringArr(d.Get("label_pairs").([]interface{}))
	}

	return &cfclient.GitopsEnvironment{
		ID:         d.Get("id").(string),
		Name:       d.Get("name").(string),
		Kind:       d.Get("kind").(string),
		Clusters:   clusters,
		LabelPairs: labelPairs,
	}
}

func mapGitopsEnvironmentToResource(d *schema.ResourceData, environment *cfclient.GitopsEnvironment) error {
	if err := d.Set("id", environment.ID); err != nil {
		return err
	}

	if err := d.Set("name", environment.Name); err != nil {
		return err
	}

	if err := d.Set("kind", environment.Kind); err != nil {
		return err
	}

	if err := d.Set("cluster", flattenClusters(environment.Clusters)); err != nil {
		return err
	}

	if err := d.Set("label_pairs", environment.LabelPairs); err != nil {
		return err
	}
	return nil
}

func flattenClusters(clusters []cfclient.GitopsEnvironmentCluster) []map[string]interface{} {

	var res = make([]map[string]interface{}, 0)

	for _, cluster := range clusters {
		m := make(map[string]interface{})
		m["name"] = cluster.Name
		m["runtime_name"] = cluster.RuntimeName
		m["namespaces"] = cluster.Namespaces
		res = append(res, m)
	}

	return res
}

func expandClusters(list []interface{}) []cfclient.GitopsEnvironmentCluster {
	var clusters = make([]cfclient.GitopsEnvironmentCluster, 0)

	for _, item := range list {
		clusterMap := item.(map[string]interface{})
		cluster := cfclient.GitopsEnvironmentCluster{
			Name:        clusterMap["name"].(string),
			RuntimeName: clusterMap["runtime_name"].(string),
			Namespaces:  datautil.ConvertStringArr(clusterMap["namespaces"].([]interface{})),
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}
