---
layout: "codefresh"
page_title: "Provider: Codefresh"
sidebar_current: "docs-codefresh-index"
description: |-
  The Codefresh provider is used to manage Codefresh resources.
---

# Codefresh Provider

The Codefresh Provider can be used to configure [Codefresh](https://codefresh.io/) resources - pipelines, projects, accounts, etc using the [Codefresh API](https://codefresh.io/docs/docs/integrations/codefresh-api/).

## Authenticating to Codefresh

The Codefresh API requires the [authentication key](https://codefresh.io/docs/docs/integrations/codefresh-api/#authentication-instructions) to authenticate.
The key can be passed either as provider's attribute or as environment variable - `CODEFRESH_API_KEY`.

## Example Usage

```hcl
provider "codefresh" {
    token = "xxxxxxxxx.xxxxxxxxxx"
}

resource "codefresh_project" "project" {
    name = "myproject"

    tags = [
      "production",
      "docker",
    ]

    variables = {
      myProjectVar = "value"
   }
}

resource "codefresh_pipeline" "pipeline" {
    name  = "${codefresh_project.project.name}/mypipeline"

    spec {

        spec_template {
            repo        = "my-github-account/my-repository"
            path        = "./codefresh.yml"
            revision    = "master"
            context     = "github"
        }

        variables = {
            goVersion = "1.13"
            release = "true"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

- `token` - (Optional) The client API token. This can also be sourced from the `CODEFRESH_API_KEY` environment variable.
- `api_url` -(Optional) Default value - https://g.codefresh.io/api.