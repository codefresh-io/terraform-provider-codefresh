# Step-type Resource

The Step-type resource allows to create your own typed step.
More about custom steps in the [official documentation](https://codefresh.io/docs/docs/codefresh-yaml/steps/#creating-a-typed-codefresh-plugin).

## Known limitations and disclaimers
### Differences during plan phase
When executing `terraform plan` the diff presented will be the comparison between the latest published version and the version configured in the `step_types_yaml`.
At this stage the Read function doesn't have the reference to the new version in order to be able to retrieve the exact version for comparison.

### Deletion of resource
When executing `terraform destroy` the step-stype is completely removed (including all the existing version) 

## Example Usage

```hcl
resource "codefresh_step_types" "custom_step" {
   
  # NOTE: you can also load the yaml from a file with `step_types_yaml = file("PATH-TO-FILE.yaml")`
  # Example has been cut down for simplicity. Yaml schema must be compliant with the what specified in the documentation for typed plugins
  step_types_yaml = <<YAML
version: '1.0'
kind: step-type
metadata:
  name: <ACCOUNT_NAME>/custom-step
  ...
spec:
  arguments: |-
     {
       ....
     }
delimiters:
    left: '[['
    right: ']]'
  stepsTemplate: |-
    print_info_message:
      name: Test step
      ...
YAML
}
```

## Argument Reference

- `step_types_yaml` (Required) YAML String containing a valid definition of a typed plugin


