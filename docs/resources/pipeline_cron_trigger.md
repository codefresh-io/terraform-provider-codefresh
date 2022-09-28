# Pipeline Cron Trigger resource

Pipeline Cron Trigger is used to create cron-based triggers for pipeilnes.
See the [documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/triggers/cron-triggers/).

## Example usage

```hcl
resource "codefresh_project" "test" {
  name = "myproject"
}

resource "codefresh_pipeline" "test" {

  name    = "${codefresh_project.test.name}/react-sample-app"

  ...
}

resource "codefresh_pipeline_cron_trigger" "test" {
	pipeline_id =  codefresh_pipeline.test.id
	expression  = "*/1 * * * *"
	message     = "Example Cron Trigger"
}
```

## Argument Reference

- `pipeline_id` - (Required) The pipeline to which this trigger belongs.
- `expression` - (Required) The cron expression. Visit [this page](https://github.com/codefresh-io/cronus/blob/master/docs/expression.md) to learn about the supported cron expression format and aliases.
- `message` - (Required) The message which will be passed to the pipeline upon each trigger.

## Attributes Reference

Along with all arguments above, the following attributes are exported:

- `event` - The Event ID assigned to this trigger.