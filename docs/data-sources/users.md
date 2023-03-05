---
page_title: "codefresh_users Data Source - terraform-provider-codefresh"
subcategory: ""
description: |-
  This data source retrieves all users in the system.
---

# codefresh_users (Data Source)

This data source retrieves all users in the system.

## Example usage

```hcl
data "codefresh_users" "users" {}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `users` (List of Object) (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `email` (String)
- `logins` (List of Object) (see [below for nested schema](#nestedobjatt--users--logins))
- `personal` (List of Object) (see [below for nested schema](#nestedobjatt--users--personal))
- `roles` (Set of String)
- `short_profile` (List of Object) (see [below for nested schema](#nestedobjatt--users--short_profile))
- `status` (String)
- `user_id` (String)
- `user_name` (String)

<a id="nestedobjatt--users--logins"></a>
### Nested Schema for `users.logins`

Read-Only:

- `credentials` (List of Object) (see [below for nested schema](#nestedobjatt--users--logins--credentials))
- `idp` (List of Object) (see [below for nested schema](#nestedobjatt--users--logins--idp))

<a id="nestedobjatt--users--logins--credentials"></a>
### Nested Schema for `users.logins.credentials`

Read-Only:

- `permissions` (Set of String)


<a id="nestedobjatt--users--logins--idp"></a>
### Nested Schema for `users.logins.idp`

Read-Only:

- `client_type` (String)
- `id` (String)



<a id="nestedobjatt--users--personal"></a>
### Nested Schema for `users.personal`

Read-Only:

- `company_name` (String)
- `country` (String)
- `first_name` (String)
- `last_name` (String)
- `phone_number` (String)


<a id="nestedobjatt--users--short_profile"></a>
### Nested Schema for `users.short_profile`

Read-Only:

- `user_name` (String)