variable api_url {
  type = string
}

# 
variable token {
  type = string
  default = ""
}

## Set of account names
variable accounts {
  type = set(string)
}

module "account_tokens" {
    source = "../../tf_modules/account_tokens"
    api_url = var.api_url
    accounts = var.accounts
}

output "account_tokens" {
    value = module.account_tokens.tokens
}