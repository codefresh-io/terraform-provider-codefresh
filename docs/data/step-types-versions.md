# Data Source: codefresh_step_types_versions
This data source allows to retrieve the list of published versions of a step-types

## Example Usage

```hcl
data "codefresh_step_types_versions" "freestyle" {
    name = "freestyle"
}

output "versions" {
    value = data.codefresh_step_types_versions.freestyle.versions
}

```

## Argument Reference

* `name` - (Required) Name of the step-types to be retrieved

## Attributes Reference

* `versions` - List of versions available for the custom plugin (step-types).
