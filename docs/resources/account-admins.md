# Account Admins resource

Use this resource to set a list of admins for any account.

## Example usage

#### Example 1

```hcl
resource "codefresh_account_admins" "test" {

  account_id = <ACCOUNT ID>

  users = [
    <USER ID>,
  ]
}
```

#### Example 2

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

data "codefresh_user" "admin" {
  email = "<EXISTING USER EMAIL>"
}

resource "codefresh_account_admins" "test" {

  account_id = codefresh_account.test.id

  users = [
    data.codefresh_user.admin.user_id
  ]
}
```

## Argument Reference

- `account_id` - (Required) The account id where to set up a list of admins.
- `users` - (Required) A list of users to set up as account admins.

## Attributes Reference

- `id` - The Account ID.

## Import

```sh
terraform import codefresh_account_admins.test xxxxxxxxxxxxxxxxxxx
```
