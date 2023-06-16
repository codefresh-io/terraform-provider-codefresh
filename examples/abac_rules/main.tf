data "codefresh_team" "admins" {
  name = "admins"
}

data "codefresh_team" "users" {
  name = "users"
}

resource "codefresh_abac_rules" "dev_pipeline" {
  for_each = toset(["run", "create", "update", "delete", "read"])
  team     = data.codefresh_team.users.id
  action   = each.value
  resource = "pipeline"
  tags     = ["dev", "untagged"]
}

resource "codefresh_permission" "admin_pipeline" {
  for_each = toset(["run", "create", "update", "delete", "read", "approve"])
  team     = data.codefresh_team.admins.id
  action   = each.value
  resource = "pipeline"
  tags     = ["production", "*"]
}

resource "codefresh_permission" "admin_pipeline_related_resource" {
  for_each         = toset(["run", "create", "update", "delete", "read", "approve"])
  team             = data.codefresh_team.admins.id
  action           = each.value
  resource         = "pipeline"
  related_resource = "project"
  tags             = ["production", "*"]
}
