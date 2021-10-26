# Data Source: codefresh_registry
This data source allows retrieving information on any existing registry

## Example Usage

```hcl
# some pre-existing registry
data "codefresh_registry" "dockerhub" {
  name = "dockerhub"
}

# example with using data reference to existing registry, not managed by terraform
# "dockerhub" registry will be used as fallback for "dockerhub1"
resource "codefresh_registry" "dockerhub1" {
  name              = "dockerhub1"
  primary           = !data.codefresh_registry.dockerhub.primary

  spec {
    dockerhub {
      username = "test"
      password = "test"
    }
  }
  fallback_registry = data.codefresh_registry.dockerhub.id
}
```

## Argument Reference

* `name` - (Required) Name of the registry to be retrieved

## Attributes Reference

* `domain` - String.
* `registry_provider` - String identifying the type of registry. E.g. `dockerhub, ecr, acr` and others
* `default` - Bool.
* `primary` - Bool.
* `fallback_registry` - String representing the id of the fallback registry.
* `repository prefix` - String representing the optional prefix for registry.
