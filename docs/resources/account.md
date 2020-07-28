# Account resource

By creating different accounts for different teams within the same company a customer can achieve complete segregation of assets between the teams.
See the [documentation](https://codefresh.io/docs/docs/administration/ent-account-mng/).

## Example usage

```hcl
resource "codefresh_account" "test" {
  name = "my_account_name"

  limits {
    collaborators = 25
    data_retention_weeks = 5
  }

  build {
    parallel = 27
  }

  features = {
    OfflineLogging = true,
    ssoManagement = true,
    teamsManagement = true,
    abac = true,
    customKubernetesCluster = true,
    launchDarklyManagement = false,
  }
}
```

## Argument Reference

- `name` - (Required) The display name for the account.
- `limits` - (Optional) A collection of `limits` blocks as documented below.
- `build` -  (Optional) A collection of `build` blocks as documented below.
- `features` - (Optional) map of supported features toggles 
---

`limits` supports the following:
- `collaborators` - (Optional) Max account's collaborators number.
- `data_retention_weeks` -(Optional) How long in weeks will the builds be stored.

---

`build` supports the following:
- `parallel` - (Optional) How many pipelines can be run in parallel.
` `node` - (Optional) Number of nodes.

## Attributes Reference

- `id` - The Account ID.

## Import

```sh
terraform import codefresh_account.test xxxxxxxxxxxxxxxxxxx
```
