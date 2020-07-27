# data codefresh_team

*Note*: Teams resources should be called with account specific access token  

```
data "codefresh_team" "admin" {
  provider = codefresh.acc1
  name = "users"
}

resource "codefresh_permission" "permission2" {
  provider = codefresh.acc1
  team = data.codefresh_team.admin.id
  action = "create"
  resource = "pipeline"
  tags = ["frontend"]
}

```

