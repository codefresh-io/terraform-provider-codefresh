# API Key resource

By default **terraform-provider-codefresh** uses API key, passed as provider's attribute, but it's possible to generate a new one.  
Codefresh API allows to operate only with entities in the current account, which owns the provided API Key.  
To be able to operate with entities in different accounts - you should create a new key in the relevant account and use providers [alias](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-instances).

## Example usage

```hcl
provider "codefresh" {
  api_url = "my API URL"
  token = "my init API token"
}

resource "codefresh_account" "test" {
  name = "my new account"
}

resource "random_string" "random" {
  length = 16
  special = false
}

resource "codefresh_api_key" "new" {
  account_id = codefresh_account.test.id
  user_id = data.codefresh_account.test_account_user.user_id
  name = "tfkey_${random_string.random.result}"

  scopes = [
    "agent",
    "agents",
    "audit",
    "build",
    "cluster",
    "clusters",
    "environments-v2",
    "github-action",
    "helm",
    "kubernetes",
    "pipeline",
    "project",
    "repos",
    "runner-installation",
    "step-type",
    "step-types",
    "view",
    "workflow",
  ]
}

provider "codefresh" {
  alias = "new_account"
  api_url = "my API URL"
  token = codefresh_api_key.new.token
}


resource "codefresh_team" "team_1" {

  provider = codefresh.new_account

  name = "team name"
}
```

## Argument Reference

- `name` - (Required) The display name for the API key.
- `account_id` - (Required) The ID of account in which the API key will be created.
- `user_id` - (Required) The ID of a user within the referenced `account_id` that will own the API key. 
- `scopes` - (Optional) A list of access scopes, that can be targeted. The possible values:
  - `agent`
  - `agents`
  - `audit`
  - `build`
  - `cluster`
  - `clusters`
  - `environments-v2`
  - `github-action`
  - `helm`
  - `kubernetes`
  - `pipeline`
  - `project`
  - `repos`
  - `runner-installation`
  - `step-type`
  - `step-types`
  - `view`
  - `workflow`

## Attributes Reference

- `id` - The Key ID.
- `token` - The Token, that should used as a new provider's token attribute.
