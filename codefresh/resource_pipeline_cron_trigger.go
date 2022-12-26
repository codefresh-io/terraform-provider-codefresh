package codefresh

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codefresh-io/cronus/pkg/cronexp"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePipelineCronTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCronTriggerCreate,
		Read:   resourcePipelineCronTriggerRead,
		Update: resourcePipelineCronTriggerUpdate,
		Delete: resourcePipelineCronTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")

				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected EVENT,PIPELINE_ID", d.Id())
				}

				event := idParts[0]
				pipelineID := idParts[1]
				d.SetId(event)
				d.Set("pipeline_id", pipelineID)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(v interface{}, path cty.Path) (diags diag.Diagnostics) {
					expression := v.(string)
					var cronExp cronexp.CronExpression

					if _, err := cronExp.DescribeCronExpression(expression); err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid cron expression.",
							Detail:   fmt.Sprintf("The cron expression %q is invalid: %s", expression, err),
						})
					}

					return
				},
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePipelineCronTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	eventString, err := client.CreateHermesTriggerEvent(&cfClient.HermesTriggerEvent{
		Type:   "cron",
		Kind:   "codefresh",
		Secret: "!generate",
		Values: map[string]string{
			"expression": d.Get("expression").(string),
			"message":    d.Get("message").(string),
		},
	})
	if err != nil {
		return err
	}

	hermesTrigger := *mapResourceToPipelineCronTrigger(d)

	err = client.CreateHermesTriggerByEventAndPipeline(eventString, hermesTrigger.PipelineID)
	if err != nil {
		return err
	}

	d.SetId(eventString)

	return nil
}

func resourcePipelineCronTriggerRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	event := d.Id()
	pipeline := d.Get("pipeline_id").(string)

	hermesTrigger, err := client.GetHermesTriggerByEventAndPipeline(event, pipeline)
	if err != nil {
		return err
	}

	err = mapPipelineCronTriggerToResource(hermesTrigger, d)
	if err != nil {
		return err
	}

	return nil
}

func resourcePipelineCronTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	err := resourcePipelineCronTriggerCreate(d, meta)
	if err != nil {
		return err
	}
	// Delete old trigger after creating one in its place — the Hermes API does not support updating the trigger.
	// We delete the old trigger after creating the new one, in case the new trigger fails to create.
	return resourcePipelineCronTriggerDelete(d, meta)
}

func resourcePipelineCronTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	hermesTrigger := *mapResourceToPipelineCronTrigger(d)

	err := client.DeleteHermesTriggerByEventAndPipeline(hermesTrigger.Event, hermesTrigger.PipelineID)
	if err != nil {
		return fmt.Errorf("failed to delete cron trigger: %v", err)
	}

	return nil
}

func mapPipelineCronTriggerToResource(hermesTrigger *cfClient.HermesTrigger, d *schema.ResourceData) error {

	d.SetId(hermesTrigger.Event)
	d.Set("pipeline_id", hermesTrigger.PipelineID)

	if hermesTrigger.Event != "" {
		r := regexp.MustCompile("[^:]+:[^:]+:[^:]+:[^:]+")
		eventStringAttributes := strings.Split(hermesTrigger.Event, ":")
		if !r.MatchString(hermesTrigger.Event) {
			return fmt.Errorf("event string must be in format 'cron:codefresh:[expression]:[message]:[uid]': %s", hermesTrigger.Event)
		}
		d.Set("expression", eventStringAttributes[2])
		d.Set("message", eventStringAttributes[3])
	}

	return nil
}

func mapResourceToPipelineCronTrigger(d *schema.ResourceData) *cfClient.HermesTrigger {

	triggerId := d.Id()
	hermesTrigger := &cfClient.HermesTrigger{
		Event:      triggerId,
		PipelineID: d.Get("pipeline_id").(string),
	}

	return hermesTrigger
}
