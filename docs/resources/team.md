# Team resource

Team is used as a part of access control and allow to define what teams have access to which clusters and pipelines.
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
```

## Argument Reference

- `name` - (Required) The display name for the team.
- `type` - (Optional) The type of the team. Possible values:
  - __default__
  - __admin__
- `tags` - (Optional) A list of tags to mark a team for easy management.
- `users` - (Optional) A list of user IDs that should be in the team.

## Attributes Reference

- `id` - The Team ID.
- `account_id` - The relevant Account ID.