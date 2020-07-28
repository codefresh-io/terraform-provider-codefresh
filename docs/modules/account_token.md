# modules account_token and account_tokens

To operate with Teams and Permission we should use token generated for Codefresh Account user (not adminCF token)  


[account_token](../../tf_modules/account_token) - creates and outputs token of single account, for usage in aliased providers

[account_tokens](../../tf_modules/account_tokens) - creates and outputs token for multiple accounts, for usage in other per-account configurations

### Example - account_token
```hcl
module "account_token" "acc1_token" {
  source = "../../tf_modules/account_token"
  account_name = "acc1"
}

provider "codefresh" {
  alias = "acc1"
  api_url = var.api_url
  token = module.change_account.acc1_token.token
}

resource "codefresh_team" "developers" {
  provider = codefresh.acc1
  name = "developers"
  account_id = data.codefresh_account.acc.id

  users = [
      data.codefresh_user.user.id
    ]
}

resource "codefresh_permission" "permission" {
  for_each = toset(["run", "create", "update", "delete", "read", "approve"])
  provider = codefresh.acc1
  team = codefresh_team.developers.id
  action = each.value
  resource = "pipeline"
  tags = [ "*", "untagged"]
}

```

### [Example account-tokens](../../examples/account_tokens)