variable "teams" {
  type = map(any)
}

module "teams" {
  source = "../.modules/teams"
  teams  = var.teams
}
