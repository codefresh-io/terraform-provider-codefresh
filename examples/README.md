# Example

In the example the Codefresh Provider is configured to authenticate with Codefresh API, and new project and pipeline are created.
Pipeline includes link to the original __codefresh.yml__ spec and two git triggres.

Run `terraform plan` or `terraform apply` as usual. Note this will modify the actual Codefresh configuration.

```hcl
provider "codefresh" {
  api_url = "https://my.onpremcodefresh.com/api"
  token = "xxxxxxxxxxxxxxx.xxxxxxxxxxxxxx"
}

resource "codefresh_project" "test" {
  name = "myproject"

  tags = [
    "docker",
  ]

  variables {
    go_version = "1.13"
  }
}

resource "codefresh_pipeline" "test" {
  name    = "${codefresh_project.test.name}/react-sample-app"

  tags = [
    "production",
    "docker",
  ]

  spec {
    concurrency = 1
    priority    = 5

    spec_template {
      repo        = "codefresh-contrib/react-sample-app"
      path        = "./codefresh.yml"
      revision    = "master"
      context     = "git"
    }

    trigger {
      branch_regex  = "/.*/gi"
      context       = "git"
      description   = "Trigger for commits"
      disabled      = false
      events        = [
        "push.heads"
      ]
      modified_files_glob = ""
      name                = "commits"
      provider            = "github"
      repo                = "codefresh-contrib/react-sample-app"
      type                = "git"
    }

    trigger {
      branch_regex  = "/.*/gi"
      context       = "git"
      description   = "Trigger for tags"
      disabled      = false
      events        = [
        "push.tags"
      ]
      modified_files_glob = ""
      name                = "tags"
      provider            = "github"
      repo                = "codefresh-contrib/react-sample-app"
      type                = "git"
    }

    variables = {
      MY_PIP_VAR      = "value"
      ANOTHER_PIP_VAR = "another_value"
    }
  }
}
```
