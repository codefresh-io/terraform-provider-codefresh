# Context Resource
A Context is an  authentication/configuration that is used by Codefresh system and engine.
There are multiple types of context available in Codefresh but they all have the following main components to define them:
* Name: A unique identifier for the context
* Type: A string representing the type of context
* Data: A data structure that provide the information related to the Context. This differs based on the type of context selected
For more details of the Context spec see in the [CLI official documentation](https://codefresh-io.github.io/cli/contexts/spec/)

## Supported types
Currently the provider support the following types of Context:
* config (Shared Config )
* secret (Shared Secret)
* yaml (YAML Configuration Context)
* secret-yaml (Secret YAML Configuration Context)

### Shared Configuration
A Shared Configuration is the entity in Codefresh that allow to create values in a central place that can then be consumed in pipelines to keep them DRY.
More details in the official [Shared Configuration documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/)

#### Example Usage of config (Shared Config)
```hcl
resource "codefresh_context" "test-config" {
    name = "my-shared-config"
    spec {
        config {
            data = {
                var1 = "value1"
                var2 = "value2"
            }
        }
    }
}
```

#### Example Usage of secret (Shared Secret)
```hcl
resource "codefresh_context" "test-secret" {
    name = "my-shared-secret"
    spec {
        secret {
            data = {
                var1 = "value1"
                var2 = "value2"
            }
        }
    }
}
```

#### Example Usage of yaml (YAML Configuration Context)
```hcl
resource "codefresh_context" "test-yaml" {
    name = "my-shared-yaml"
    spec {
        # NOTE: you can also load the yaml from a file with `yaml = file("PATH-TO-FILE.yaml")`
        yaml = <<YAML
test:
  nested_value: value1
  list:
    - test2
    - test3
another_element: value
YAML
    }
}
```

#### Example Usage of secret-yaml (Secret YAML Configuration Context)
```hcl
resource "codefresh_context" "test-secret-yaml" {
    name = "my-shared-secret-yaml"
    spec {
        # NOTE: The `-` from secret-yaml is stripped because the character is not allowed in Field name
        # File passed MUST be a valid YAML
        secretyaml = file("test.yaml")
    }
}
```


## Argument Reference

- `name` - (Required) The display name for the context.
- `spec` - (Required) A `spec` block as documented below.

---

`spec` supports the following (Note: only 1 of the below can be specified at any time):

- `config`      - (Optional) A `config` block as documented below. Shared Config [spec](https://codefresh-io.github.io/cli/contexts/spec/config/).
- `secret`      - (Optional) A `secret` block as documented below. Shared Secret [spec](https://codefresh-io.github.io/cli/contexts/spec/secret/).
- `yaml`        - (Optional) A `yaml` block as documented below. Yaml Configuration Context [spec](https://codefresh-io.github.io/cli/contexts/spec/yaml/).
- `secretyaml`  - (Optional) A `secretyaml` block as documented below. Secret Yaml Configuration Context[spec](https://codefresh-io.github.io/cli/contexts/spec/secret-yaml/).

---

`config` supports the following:

- `data` - (Required) Map of strings representing the variables to be defined in the Shared Config.

---

`secret` supports the following:

- `data` - (Required) Map of strings representing the variables to be defined in the Shared Config.

---

`yaml` supports the following:

- `data` - (Required) String representing a YAML file content

---

`secretyaml` supports the following:

- `data` - (Required) String representing a YAML file content

---