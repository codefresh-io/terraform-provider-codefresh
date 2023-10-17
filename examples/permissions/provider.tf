provider "codefresh" {
  api_url = var.api_url
  token   = var.token # If token isn't set the provider expects the $CODEFRESH_API_KEY env variable
}