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
  event = "registry:dockerhub:korenyoni:codefresh-web-app:push:47e5d8141593"
  pipeline = module.pipeline.id
}