# Project Resource

The top-level concept in Codefresh. You can create projects to group pipelines that are related. In most cases a single project will be a single application (that itself contains many micro-services). You are free to use projects as you see fit. For example, you could create a project for a specific Kubernetes cluster or a specific team/department.
More about pipeline concepts see in the [official documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/pipelines/#pipeline-concepts).

## Example Usage

```hcl
resource "codefresh_project" "test" {
    name = "myproject"

    tags = [
      "production",
      "docker",
    ]

    variables = {
      go_version = "1.13"
   }
}
```

## Argument Reference

- `name` (Required) The display name for the project.
- `tags` (Optional) A list of tags to mark a project for easy management and access control.
- `variables` (Optional) project variables.

## Attributes Reference

- `id` - The Project ID

## Import

```sh
terraform import codefresh_project.test xxxxxxxxxxxxxxxxxxx
```
