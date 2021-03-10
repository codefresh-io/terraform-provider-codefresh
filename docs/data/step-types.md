# Data Source: codefresh_step_types
This data source allows to retrieve the published versions of a step-types

## Example Usage

```hcl
data "codefresh_step_types" "freestyle" {
    name = "freestyle"
}

local {
  freestyle_map = { for step_definition in data.codefresh_step_types.freestyle.version: step_definition.version_number => step_definition }
}

output "test" {
  # Value is return as YAML
  value = local.freestyle_map[keys(local.freestyle_map)[0]].version_number
}

```

## Argument Reference

* `name` - (Required) Name of the step-types to be retrieved

## Attributes Reference

- `version` -  A Set of `version` blocks as documented below.

---

`version` provides the following:
- `version_number` - String representing the semVer for the step
- `step_types_yaml` - YAML String containing the definition of a typed plugin
