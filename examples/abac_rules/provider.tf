provider "codefresh" {
  api_url =  var.api_url
  api_url_v2 =  var.api_url_v2
  token = var.token # If token isn't set the provider expects the $CODEFRESH_API_KEY env variable
}