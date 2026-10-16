data "codefresh_team" "admins" {
  name = "admins"
}

data "codefresh_team" "users" {
  name = "users"
}

resource "codefresh_abac_rules" "app_rule" {
  entity_type = "gitopsApplications"
  teams       = [data.codefresh_team.users.id]
  actions     = ["REFRESH", "SYNC", "TERMINATE_SYNC", "VIEW_POD_LOGS", "APP_ROLLBACK"]

  attribute {
    name = "LABEL"
    key = "KEY"
    value = "VALUE"
  }

  tags        = ["dev", "untagged"]
}

resource "codefresh_abac_rules" "app_view_rule" {
  entity_type = "gitopsApplications"
  teams       = [data.codefresh_team.users.id]
  actions     = ["VIEW"]

  attribute {
    name  = "RUNTIME"
    value = "production-runtime"
  }
}
