# IDP Admins resource

The resource adds the list of provided account IDs to the IDP.  
Because of the current Codefresh API limitation it's impossible to remove account from IDP, only adding is supporting.

## Example usage

```hcl
resource "codefresh_account" "test" {
  name = "<MY ACCOUNT NAME>"
}

resource "codefresh_idp_accounts" "test" {

  idp = "azure"

  accounts = [
    codefresh_account.test.id,
    "<ANY ACCOUNT ID>"
  ]
}
```

## Argument Reference

- `idp` - (Required) The IDP client name.
- `accounts` - (Required) A list of account IDs.

## Import

```sh
terraform import codefresh_idp_accounts.test xxxxxxxxxxxxxxxxxxx
```
