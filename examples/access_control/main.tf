module "codefresh_access_control" {
  source = "../../tf_modules/access_control"
  api_url = var.api_url
  default_idps = var.default_idps
  accounts = var.accounts
  users = var.users

  default_acccount_limits = var.default_acccount_limits
}