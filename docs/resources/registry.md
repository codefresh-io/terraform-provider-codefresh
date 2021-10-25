# Registry Resource

Registry is the configuration that Codefresh uses to push/pull docker images.
For more details see the [Codefresh Docker Registries](https://codefresh.io/docs/docs/integrations/docker-registries/)


## Concurrency Limitation 

Codefresh Registry API was not designed initially to handle concurrent modifications on `registry` entity.
Thus, you need to take one of the following approaches to avoid **errors** and **non-expected behavior**:

1) run terraform write operations with `-parallelism=1` option
```shell
terraform apply -parallelism=1
terraform destroy -parallelism=1
```

2) make each registry resource `depend_on` each other - so the CRUD operations will be performed for each registry **sequentially**
```hcl
resource "codefresh_registry" "dockerhub" {
    name = "dockerhub"
  
    spec {
        dockerhub {
          # some specific fields here
        }
    }
}

# this registry will depend on the "dockerhub" registry
resource "codefresh_registry" "gcr" {
    name = "gcr"
    
    depends_on = [codefresh_registry.dockerhub]    
    spec {
        gcr {
          # some specific fields here
        }
    }
}
```

## Supported Registry Providers

Currently, Codefresh supports the following registry providers:
* dockerhub - [Docker Hub](https://codefresh.io/docs/docs/integrations/docker-registries/docker-hub/)
* acr - [Azure Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/azure-docker-registry)
* gcr - [Google Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/google-container-registry)
* gar - [Google Artifact Registry](https://codefresh.io/docs/docs/integrations/docker-registries/google-artifact-registry)
* ecr - [Amazon EC2 Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/amazon-ec2-container-registry)
* bintray - [Bintray / Artifactory](https://codefresh.io/docs/docs/integrations/docker-registries/bintray-io)
* other - any other provider including [Quay](https://codefresh.io/docs/docs/integrations/docker-registries/quay-io) and [Github Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/github-container-registry). See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/other-registries).

### Resource Spec

Each registry resource have some common fields and specific ones stored under the `spec`. Here is the template:

```hcl
resource "codefresh_registry" "some_registry" {
    name = "some_name"
    default = false
    primary = true
    fallback_registry = codefresh_registry.some_other_registry.id
  
    spec {
        <provider_name> {
          # some specific fields here
        }
    }
}
```


### Common fields

* `name` - `(Required)` some unique name for registry
* `default` - `(Optional, Default = false)` see the [Default Registry](https://codefresh.io/docs/docs/integrations/docker-registries/#the-default-registry)
* `primary` - `(Optional, Default = true)` see the [Multiple Registries](https://codefresh.io/docs/docs/ci-cd-guides/working-with-docker-registries/#working-with-multiple-registries-with-the-same-domain)
* `fallback_registry` - `(Optional)` see the [Fallback Registry](https://codefresh.io/docs/docs/integrations/docker-registries/#fallback-registry)


### Default registry usage

If you want to manage default registry by Codefresh terraform provider correctly, 
you need to mark only one registry as `default = true`

```hcl
resource "codefresh_registry" "dockerhub" {
    name = "dockerhub"
  
    spec {
        dockerhub {
          # some specific fields here
        }
    }
}

# this registry will be default
resource "codefresh_registry" "gcr" {
    name = "gcr"
    default = true
    fallback_registry = codefresh_registry.some_other_registry.id
    
    spec {
        gcr {
          # some specific fields here
        }
    }
}
```

### Primary registry usage

If you are using [Multiple Registries](https://codefresh.io/docs/docs/ci-cd-guides/working-with-docker-registries/#working-with-multiple-registries-with-the-same-domain) feature
you need to manually mark each registry of the same domain as non-primary and only one as primary

```hcl
# this registry will be primary
resource "codefresh_registry" "dockerhub" {
    name = "dockerhub"
    primary = true
  
    spec {
        dockerhub {
          # some specific fields here
        }
    }
}

resource "codefresh_registry" "dockerhub1" {
    name = "dockerhub1"
    primary = false

    spec {
        dockerhub {
          # some specific fields here
        }
    }
}

resource "codefresh_registry" "dockerhub2" {
    name = "dockerhub2"
    primary = false

    spec {
        dockerhub {
          # some specific fields here
        }
    }
}
```

### Fallback registry usage

If you want to use one of your registries as fallback you need to specify its id 
for `fallback_registry` field of another registry

```hcl
resource "codefresh_registry" "dockerhub" {
    name = "dockerhub"
  
    spec {
        dockerhub {
          # some specific fields here
        }
    }
}

resource "codefresh_registry" "gcr" {
    name = "gcr"
  
    # here we take the id of "dockerhub" registry
    fallback_registry = codefresh_registry.dockerhub.id
    
    spec {
        gcr {
          # some specific fields here
        }
    }
}
```

## Argument Reference

- `name` - _(Required)_ some unique name for registry
- `default` - _(Optional, Default = false)_ default registry
- `primary` - _(Optional, Default = true)_ primary registry
- `fallback_registry` - _(Optional)_ fallback registry
- `spec` - _(Required)_ A `spec` block as documented below.

---

`spec` supports the following (Note: only 1 of the below can be specified at any time):

- dockerhub - _(Optional)_ A `dockerhub` block as documented below ([Docker Hub Registry](https://codefresh.io/docs/docs/integrations/docker-registries/docker-hub/))
- acr - _(Optional)_ An `acr` block as documented below ([Azure Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/azure-docker-registry))
- gcr - _(Optional)_ A `gcr` block as documented below ([Google Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/google-container-registry))
- gar - _(Optional)_ A `gar` block as documented below ([Google Artifact Registry](https://codefresh.io/docs/docs/integrations/docker-registries/google-artifact-registry))
- ecr - _(Optional)_ An `ecr` block as documented below ([Amazon EC2 Container Registry](https://codefresh.io/docs/docs/integrations/docker-registries/amazon-ec2-container-registry))
- bintray - _(Optional)_ A `bintray` block as documented below ([Bintray / Artifactory](https://codefresh.io/docs/docs/integrations/docker-registries/bintray-io))
- other - _(Optional)_ `other` provider block described below ([Other Providers](https://codefresh.io/docs/docs/integrations/docker-registries/other-registries)).


---

`dockerhub` supports the following:

- `username` - _(Required)_ String.
- `password` - _(Required, Sensitive)_ String.

---

`acr` supports the following:

- `domain` - _(Required)_ String representing your acr registry domain.
- `client_id` - _(Required)_ String representing client id.
- `client_secret` - _(Required)_ String representing client secret.
- `repository_prefix` - _(Optional)_ String. See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).

---

`gcr` supports the following:

- `domain` - _(Required)_ String representing one of the Google's gcr domains
- `keyfile` - _(Required)_ String representing service account json file contents
- `repository_prefix` - _(Optional)_ String. See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).

---

`gar` supports the following:

- `location` - _(Required)_ String representing one of the Google's gar locations
- `keyfile` - _(Required)_ String representing service account json file contents
- `repository_prefix` - _(Optional)_ String. See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).

---

`ecr` supports the following:

- `region` - _(Required)_ String representing one of the Amazon regions
- `access_key_id` - _(Required)_ String representing access key id
- `secret_access_key` - _(Required)_ String representing secret access key
- `repository_prefix` - _(Optional)_ String. See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).

---

`bintray` supports the following:

- `domain` - _(Required)_ String representing the bintray domain
- `username` - _(Required)_ String representing the username
- `token` - _(Required)_ String representing token
- `repository_prefix` - _(Optional)_ String. See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).

---

`other` supports the following:

- `domain` - _(Required)_ String representing the bintray domain
- `username` - _(Required)_ String representing the username
- `password` - _(Required)_ String representing token
- `repository_prefix` - _(Optional)_ String. See the [docs](https://codefresh.io/docs/docs/integrations/docker-registries/#using-an-optional-repository-prefix).
- `behind_firewall` - _(Optional, Default = false)_ Bool. See the [docs](https://codefresh.io/docs/docs/administration/behind-the-firewall/#accessing-an-internal-docker-registry).

---