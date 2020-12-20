# Data Source: codefresh_context
This data source allows to retrieve information on any defined context

## Example Usage

```hcl
# Assuming runtimes-list is a context of type "config" with the following values
# runtime_a: dev
# runtime_b: test
# runtime_c: prod

data "codefresh_context" "runtimes_list" {
  name = "runtimes-list"
}

resource "codefresh_project" "test" {
  name = "myproject"
}

resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name    = "${codefresh_project.test.name}/react-sample-app"

  runtime_environment {
    name = yamldecode(data.codefresh_context.runtimes_list.data).runtime_a
  }

  spec {

    spec_template {
      repo        = "codefresh-contrib/react-sample-app"
      path        = "./codefresh.yml"
      revision    = "master"
      context     = "git"
    }
  }
}
```

## Argument Reference

* `name` - (Required) Name of the context to be retrived

## Attributes Reference

* `type` - String identifying the type of extracted context. E.g. `config`, `secret`, `git.github-app`, etc.
* `data` - The yaml string representing the context. Use the `yamldecode` function to access the values belonging the context.
