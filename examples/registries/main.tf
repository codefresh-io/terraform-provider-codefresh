
variable api_url {
  type = string
}

variable token {
  type = string
  default = ""
}

variable test_password {
  type = string
  default = ""
}

provider "codefresh" {
  api_url = var.api_url
  token = var.token
}

resource "codefresh_registry" "test" {
  name = "test"
  kind = "standard"
  registry_provider = "other"
  domain = "test1.io"
  username = "test"
  password = var.test_password
}