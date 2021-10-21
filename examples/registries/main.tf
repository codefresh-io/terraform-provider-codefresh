
variable api_url {
  type = string
}

variable token {
  type    = string
  default = ""
}

variable test_password {
  type    = string
  default = ""
}

provider "codefresh" {
  api_url = var.api_url
  token   = var.token
}

resource "codefresh_registry" "acr" {
  name = "acr"
  spec {
    acr {
      domain            = "acr.io"
      client_id         = "test"
      client_secret     = "test"
      repository_prefix = "test"
    }
  }
}

resource "codefresh_registry" "gcr" {
  name = "gcr"
  spec {
    gcr {
      domain            = "gcr.io"
      keyfile           = "test"
      repository_prefix = "codefresh-inc"
    }
  }
}

resource "codefresh_registry" "gar" {
  name = "gar"
  spec {
    gar {
      domain            = "asia"
      keyfile           = "test"
      repository_prefix = "codefresh-inc"
    }
  }
}

resource "codefresh_registry" "dockerhub" {
  name    = "dockerhub1"
  primary = false
  spec {
    dockerhub {
      username = "test"
      password = "test"
    }
  }
}

resource "codefresh_registry" "bintray" {
  name = "bintray"
  spec {
    bintray {
      domain   = "bintray.io"
      username = "test"
      token    = "test"
    }
  }
}

resource "codefresh_registry" "other" {
  name    = "other"
  primary = true
  spec {
    other {
      domain   = "other.io"
      username = "test"
      password = "test"
    }
  }
}

resource "codefresh_registry" "other1" {
  name    = "other1"
  primary = false
  spec {
    other {
      domain   = "other.io"
      username = "test"
      password = "test"
    }
  }
}

resource "codefresh_registry" "other2" {
  name    = "other2"
  primary = false
  spec {
    other {
      domain   = "other.io"
      username = "test"
      password = "test"
    }
  }
}
