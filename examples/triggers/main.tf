provider "codefresh" {
  api_url = var.api_url
  token = var.token
}

module "pipeline" {
  source = "../pipelines"
  api_url = var.api_url
  token   = var.token
}

resource "codefresh_pipeline_cron_trigger" default {
  pipeline = module.pipeline.id
  expression = "0 0 2 ? * MON-FRI"
  message  = "Triggered by cron"
}