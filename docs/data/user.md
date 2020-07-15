# User Data Source

Use this data source to get the User from existing users for use in other resources.

## Example usage

```hcl
data "codefresh_user" "admin" {
  email = "admin@codefresh.io"
}

resource "codefresh_team" "admins" {

  name = "testsuperteam123"

  users = [
    data.codefresh_user.admin.user_id,
    "<ANY USER ID>",
  ]
}
```

## Argument Reference

- `email` - (Required) The email of user to filter.

## Attributes Reference

- `user_name`.
- `email`.
- `user_id`.
- `personal`. A collection of `personal` blocks as documented below.
- `short_profile`. A collection of `short_profile` blocks as documented below.
- `roles`. A list of roles.
- `status`. User status - `new`, `pending`, etc.
- `logins`. A collection of `short_profile` blocks as documented below.

---

`personal` includes the following:
- `first_name`.
- `last_name`.
- `company_name`.
- `phone_number`.
- `country`.

---

`short_profile` includes the following:
- `user_name`.

---

`logins` includes the following:
- `credentials`
    - `permissions`
- `idp`
    - `id`
    - `client_type