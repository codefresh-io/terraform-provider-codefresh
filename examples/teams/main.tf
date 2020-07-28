variable api_url {
  type = string
}

variable token {
  type = string
  default = ""
}
provider "codefresh" {
  api_url = var.api_url
  token = var.token
}

variable teams {
  type = map(any)
}

module "teams" {
    source = "../../tf_modules/teams"
    teams = var.teams
}
