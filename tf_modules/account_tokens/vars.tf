variable api_url {
  type = string
}

# 
variable admin_token {
  type = string
  default = ""
}


## Set of account names
variable accounts {
  type = set(string)
}