# User resource

Use this resource to create a new user.

## Example usage

```hcl
resource "codefresh_account" "test" {

  name = "mynewaccount"

  limits {
    collaborators = 25
    data_retention_weeks = 5
  }

  build {
    parallel = 2
  }
}

resource "codefresh_user" "new" {
  email = "<EMAIL>"
  user_name = "<USER>"

  activate = true

  roles = [
    "Admin",
    "User"
  ]

  login {
      idp_id = data.codefresh_idps.idp_azure.id
      sso = true
  }
  
  login  {
      idp_id = data.codefresh_idps.local.id
      //sso = false
  }


  personal {
    first_name = "John"
    last_name = "Smith"
  }

  accounts = [
    codefresh_account.test.id,
    "59009117c102763beda7ce71",
  ]
}
```

## Argument Reference

- `email` - (Required) A new user email.
- `user_name` - (Required) The new user name.
- `activate` - (Optional) Boolean. Activate the new use or not. If a new user is not activate, it'll be left pending.
- `accounts` - (Optional) A list of accounts to add to the user.
- `personal` - (Optional) A collection of `personal` blocks as documented below.
- `accounts` - (Optional) A list of user roles. Possible values - `Admin`, `User`.
- `login` - (Optional) A collection of `login` blocks as documented below.

---

`personal` supports the following:

- `first_name` - (Optional).
- `last_name` - (Optional).
- `company_name` - (Optional).
- `phone_number` - (Optional).
- `country` - (Optional).

---

`login` supports the following:
- `credentials`
  - `permissions` - (Optional) A list of permissions.
- `idp`
  - `idp_id` - (Optional) The id of IDP to the user to.
  - `client_type` -  (Optional) IDP type. ex. - `github`, `azure`, etc.

## Attributes Reference

- `id` - The User ID.
- `short_profile`
  - `user_name`
- `status`. Current status of the user. ex - `new`, `pengind`.



## Import

```sh
terraform import codefresh_user.new xxxxxxxxxxxxxxxxxxx
```

