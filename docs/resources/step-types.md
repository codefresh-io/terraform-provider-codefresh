# Step-types Resource

The Step-types resource allows to create your own typed step and manage all it's published versions.
The resource allows to handle the life-cycle of the version by allowing specifying multiple blocks `version` where the user provides a version number and the yaml file representing the plugin.
More about custom steps in the [official documentation](https://codefresh.io/docs/docs/codefresh-yaml/steps/#creating-a-typed-codefresh-plugin).

## Known limitations and disclaimers
### Version and name in yaml Metadata are ignored.
The version and name of the step declared in the yaml files are superseeded by the attributes specified at resource level:
- `name` : at top level
- `version_numer`: specified in the `version` block
The above are added/replaced at runtime time.

### Number of API requests
This resource makes a lot of additional API calls to validate the steps and retrieve all the version available.
Caution is recommended on the amount of versions maintained and the number of resources defined in a single project.


## Example Usage

```hcl

data "codefresh_current_account" "acc" {
}

resource "codefresh_step_types_versions" "my-custom-step" {
  name = "${data.codefresh_current_account.acc.name}/my-custom-step"

  version {
    version_number = "0.0.1"
    step_types_yaml = file("./templates/plugin-0.0.1.yaml")
  }
  version {
    version_number = "0.0.2"
    step_types_yaml = file("./templates/plugin-0.0.2.yaml")
  }
  ....
}
}
```

## Argument Reference
- `name` - (Required) The name for the step-type
- `version` - (At least 1 Required) A collection of `version` blocks as documented below.

---

`version` supports the following:
- `version_number` - (Required) String representing the semVer for the step
- `step_types_yaml` (Required) YAML String containing a valid definition of a typed plugin
