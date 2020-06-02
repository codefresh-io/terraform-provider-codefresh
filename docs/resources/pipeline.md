# Pipeline Resource

The central component of the Codefresh Platform. Pipelines are workflows that contain individual steps. Each step is responsible for a specific action in the process.
See the [documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/introduction-to-codefresh-pipelines/) for the details.

## Example Usage

```hcl
resource "codefresh_project" "test" {
  name = "myproject"
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

## Argument Reference

- `name` - (Required) The display name for the pipeline.
- `tags` - (Optional) A list of tags to mark a project for easy management and access control.
- `spec` - (Required) A collection of `spec` blocks as documented below.
- `original_yaml_string` - (Optional) A string with original yaml pipeline.
  - `original_yaml_string` = "version: \"1.0\"\nsteps:\n  test:\n    image: alpine:latest\n    commands:\n      - echo \"ACC tests\""
  - or `original_yaml_string` = file("/path/to/my/codefresh.yml")

---

`spec` supports the following:

- `concurrency` - (Optional) The maximum amount of concurrent builds.
- `priority` - (optional) Helps to organize the order of builds execution in case of reaching the concurrency limit.
- `variables` - (Optional) Pipeline variables.
- `trigger` - (Optional) A collection of `trigger` blocks as documented below. Triggers [documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/triggers/git-triggers/).
- `spec_template` - (Optional) A collection of `spec_template` blocks as documented below.
- `runtime_environment` - (Optional) A collection of `runtime_environment` blocks as documented below.

---

`spec_template` supports the following:

- `location` - (Optional) Default value - **git**.
- `repo` - (Required) The GitHub `account/repo_name`.
- `path` - (Required) The relative path to the Codefresh pipeline file.
- `revison` - (Required) The git revision.
- `context` - (Optional) The Codefresh Git [context](https://codefresh.io/docs/docs/integrations/git-providers/).

---

`trigger` supports the following:

- `name` - (Optional) The display name for the pipeline.
- `description` - (Optional) The trigger description.
- `type` - (Optional) The trigger type. Default value - **git**.
- `repo` - (Optional) The GitHub `account/repo_name`.
- `branch_regex` - (Optional) A regular expression and will only trigger for branches that match this naming pattern.
- `modified_files_glob` - (Optional) Allows to constrain the build and trigger it only if the modified files from the commit match this glob expression.
- `events` - (Optional) A list of GitHub events for which a Pipeline is triggered. Default value - **push.heads**.
- `provider` - (Optional) Default value - **github**.
- `context` - (Optional) Codefresh Git context.
- `variables` - (Optional) Trigger variables.

---

`runtime_environment` supports the following:

- `name` - (Required) A name of runtime.
- `cpu` - (Optional) A required amount of CPU.
- `memory` - (Optional) A required amount of memory.
- `dind_storage` - (Optional) A pipeline shared storage.

## Attributes Reference

- `id` - The Pipeline ID.

## Import

```sh
terraform import codefresh_pipeline.test xxxxxxxxxxxxxxxxxxx
```
