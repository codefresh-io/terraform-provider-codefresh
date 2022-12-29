provider "codefresh" {
  api_url = var.api_url
  token   = var.token
}

module "pipeline" {
  source  = "../pipelines"
  api_url = var.api_url
  token   = var.token
}

resource "codefresh_pipeline_cron_trigger" "default" {
  pipeline_id = module.pipeline.id
  expression  = "0 0 2 ? * MON-FRI"
  message     = "test"
}
