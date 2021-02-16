# resource codefresh_permission
Permission are used to setup access control and allow to define which teams have access to which clusters and pipelines based on tags
See the [documentation](https://codefresh.io/docs/docs/administration/access-control/).

## Example usage

```hcl
resource "codefresh_team" "developers" {
  name = "developers"

  users = [
      "5efc3cb6355c6647041b6e49",
      "59009221c102763beda7cf04"
    ]
}

resource "codefresh_permission" "developers" {

  team = codefresh_team.developers.id
  resource = "pipeline"
  action = "run"
  tags = [
      "demo",
      "test"
    ]
}
```

## Argument Reference

- `action` - (Required) Action to be allowed. Possible values:
  - __create__
  - __read__
  - __update__
  - __delete__
  - __run__ (Only valid for `pipeline` resource)
  - __approve__ (Only valid for `pipeline` resource)
  - __debug__ (Only valid for `pipeline` resource)
- `resource` - (Required) The type of resource the permission applies to. Possible values:
  - __pipeline__
  - __cluster__
- `team` - (Required) The Id of the team the permissions apply to.
- `tags` - (Optional) The effective tags to apply the permission. It supports 2 custom tags:
  - __untagged__ is a “tag” which refers to all clusters that don’t have any tag.
  - __*__ (the star character) means all tags.

## Attributes Reference

- `id` - The permission ID.
