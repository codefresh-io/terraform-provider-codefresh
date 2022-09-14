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
			"event": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePipelineCronTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cfClient.Client)

	hermesTrigger := *mapResourceToPipelineCronTrigger(d)

	_, err := client.CreateHermesTriggerByEventAndPipeline(hermesTrigger.Event, hermesTrigger.Pipeline)
	if err != nil {
		return err
	}

	return nil
}

func resourcePipelineCronTriggerRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	event := d.Id()
	if event == "" {
		d.SetId("")
		return nil
	}

	hermesTrigger, err := client.GetHermesTriggerByEvent(event)
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
		log.Printf("Failed to delete hermes trigger: %s", err)
	}

	return nil
}

func mapPipelineCronTriggerToResource(hermesTrigger *cfClient.HermesTrigger, d *schema.ResourceData) error {

	d.SetId(hermesTrigger.Event)
	return nil
}

func mapResourceToPipelineCronTrigger(d *schema.ResourceData) *cfClient.HermesTrigger {

	hermesTrigger := &cfClient.HermesTrigger{
		Event:    d.Get("event").(string),
		Pipeline: d.Get("pipeline").(string),
	}

	return hermesTrigger
}
