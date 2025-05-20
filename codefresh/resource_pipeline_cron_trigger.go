package codefresh

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robfig/cron"
)

func resourcePipelineCronTrigger() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated and will be removed in a future version of the Codefresh Terraform provider. Please use the cron_triggers attribute of the codefresh_pipeline resource instead.",
		Description:        "This resource is used to create cron-based triggers for pipeilnes.",
		Create:             resourcePipelineCronTriggerCreate,
		Read:               resourcePipelineCronTriggerRead,
		Update:             resourcePipelineCronTriggerUpdate,
		Delete:             resourcePipelineCronTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")

				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected EVENT,PIPELINE_ID", d.Id())
				}

				event := idParts[0]
				pipelineID := idParts[1]
				d.SetId(event)

				err := d.Set("pipeline_id", pipelineID)

				if err != nil {
					return nil, err
				}

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
				ValidateDiagFunc: schemautil.CronExpression(
					// Legacy cron parser, still used by standalone cron triggers
					schemautil.WithCronParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)),
				),
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
		},
		// Force new resource if any field changes. This is because the Codefresh API does not support updating cron triggers.
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("pipeline_id", func(ctx context.Context, old, new, meta interface{}) bool {
				return true
			}),
			customdiff.ForceNewIfChange("expression", func(ctx context.Context, old, new, meta interface{}) bool {
				return true
			}),
			customdiff.ForceNewIfChange("message", func(ctx context.Context, old, new, meta interface{}) bool {
				return true
			}),
		),
	}
}

func resourcePipelineCronTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	eventString, err := client.CreateHermesTriggerEvent(&cfclient.HermesTriggerEvent{
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

	client := meta.(*cfclient.Client)

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
	// see notes in resourcePipelineCronTrigger()
	return fmt.Errorf("cron triggers cannot be updated")
}

func resourcePipelineCronTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfclient.Client)

	hermesTrigger := *mapResourceToPipelineCronTrigger(d)

	err := client.DeleteHermesTriggerByEventAndPipeline(hermesTrigger.Event, hermesTrigger.PipelineID)
	if err != nil {
		return fmt.Errorf("failed to delete cron trigger: %v", err)
	}

	return nil
}

func mapPipelineCronTriggerToResource(hermesTrigger *cfclient.HermesTrigger, d *schema.ResourceData) error {

	d.SetId(hermesTrigger.Event)
	err := d.Set("pipeline_id", hermesTrigger.PipelineID)

	if err != nil {
		return err
	}

	if hermesTrigger.Event != "" {
		r := regexp.MustCompile("[^:]+:[^:]+:[^:]+:[^:]+")
		eventStringAttributes := strings.Split(hermesTrigger.Event, ":")
		if !r.MatchString(hermesTrigger.Event) {
			return fmt.Errorf("event string must be in format 'cron:codefresh:[expression]:[message]:[uid]': %s", hermesTrigger.Event)
		}
		err = d.Set("expression", eventStringAttributes[2])

		if err != nil {
			return err
		}

		err = d.Set("message", eventStringAttributes[3])

		if err != nil {
			return err
		}
	}

	return nil
}

func mapResourceToPipelineCronTrigger(d *schema.ResourceData) *cfclient.HermesTrigger {

	triggerId := d.Id()
	hermesTrigger := &cfclient.HermesTrigger{
		Event:      triggerId,
		PipelineID: d.Get("pipeline_id").(string),
	}

	return hermesTrigger
}
