# Example

```yaml
resource "codefresh_project" "test" {
    name = "test"
}
```

Run `terraform plan` or `terraform apply` as usual. Note this will modify the actual Codefresh configuration.

# Syntax Examples

## Project

```yaml
resource "codefresh_project" "docker" {
  name = "docker"
}
```

## Pipeline

```yaml
resource "codefresh_pipeline" "docker_monorepo" {
  name    = "docker/docker-monorepo"
  project = "docker"

  spec = {
    repo        = "abcinc/monorepo"
    path        = "./codefresh/docker/docker-monorepo.yaml"
    revision    = "master"
    concurrency = 1
    priority    = 5
  }

  tags = [
    "docker",
  ]

  variables {
    TAG = "master"
  }
}
```

## Cron Trigger

```yaml
resource "codefresh_cron_event" "docker_monorepo_cron" {
  expression = "40 0 * * *"
  message    = "build monorepo docker"
}

resource "codefresh_cron_trigger" "docker_monorepo_cron" {
  pipeline = "${codefresh_pipeline.docker_monorepo.id}"
  event    = "${codefresh_cron_event.docker_monorepo_cron.id}"
}
```

## Environment

```yaml
resource "codefresh_environment" "staging" {
  account_id = "<redacted>"
  name       = "staging"
  namespace  = "staging"
  cluster    = "abcinc-staging"
}
```

## User
```yaml
resource "codefresh_user" "john_doe" {
  email = "jdoe@abcinc.com"
}
```
