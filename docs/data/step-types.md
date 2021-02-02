# Data Source: codefresh_step_types
This data source allows to retrieve the latest published version of a step-types

## Example Usage

```hcl
data "codefresh_step_types" "freestyle" {
    name = "freestyle"
}

output "test" {
  # Value is return as YAML
  value = yamldecode(data.codefresh_step_types.freestyle.step_types_yaml).metadata.updated_at
}

```

## Argument Reference

* `name` - (Required) Name of the step-types to be retrieved
* `version` - (Optional) Version to be retrieved. If not specified, the latest published will be returned

## Attributes Reference

* `step_types_yaml` - The yaml string representing the custom plugin (step-types).
