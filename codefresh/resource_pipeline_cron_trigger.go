package codefresh

import (
	"log"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePipelineCronTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCronTriggerCreate,
		Read:   resourcePipelineCronTriggerRead,
		Update: resourcePipelineCronTriggerUpdate,
		Delete: resourcePipelineCronTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_id": {
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

	err = client.CreateHermesTriggerByEventAndPipeline(eventString, hermesTrigger.Pipeline)
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
	return resourcePipelineCronTriggerCreate(d, meta)
}

func resourcePipelineCronTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	hermesTrigger := *mapResourceToPipelineCronTrigger(d)

	err := client.DeleteHermesTriggerByEventAndPipeline(hermesTrigger.Event, hermesTrigger.Pipeline)
	if err != nil {
		log.Printf("Failed to delete cron trigger: %s", err)
	}

	return nil
}

func mapPipelineCronTriggerToResource(hermesTrigger *cfClient.HermesTrigger, d *schema.ResourceData) error {

	d.SetId(hermesTrigger.Event)
	d.Set("pipeline", hermesTrigger.Pipeline)

	return nil
}

func mapResourceToPipelineCronTrigger(d *schema.ResourceData) *cfClient.HermesTrigger {

	triggerId := d.Id()
	hermesTrigger := &cfClient.HermesTrigger{
		Event:    triggerId,
		Pipeline: d.Get("pipeline_id").(string),
	}

	return hermesTrigger
}
