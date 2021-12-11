# Terraform provider for Codefresh

This provider was initialized by [LightStep](https://lightstep.com/) and will be maintained as the official Terraform provider for Codefresh.

The provider is still under development, and can be used as a terraform [third-party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) only.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+ ;
- [Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin).

## Download Provider
Download and extract terraform-provider-codefresh from [releases](https://github.com/codefresh-io/terraform-provider-codefresh/releases)

## Building the Provider

```sh
go build -o terraform-provider-codefresh
```

## Using the Provider

Compile or take from the [Releases](https://github.com/codefresh-contrib/terraform-provider-codefresh/releases) `terraform-provider-codefresh` binary and place it locally according the Terraform plugins [documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

### for terraform 0.12:
- _~/.terraform.d/plugins/linux\_amd64_
- _./terraform.d/plugins/linux\_amd64_. The relative path in your Terraform project.

### for terraform 0.13 follow [required providers](https://www.terraform.io/docs/configuration/provider-requirements.html):
```bash
# OS  is linux|windows|darwin, ARCH is amd64|arm|x86
PLUGIN_DIR=~/.terraform.d/plugins/codefresh.io/app/codefresh/0.1.0/<OS_ARCH>
mkdir -p ${PLUGIN_DIR}
cp terraform-provider-codefresh  ${PLUGIN_DIR}/
```

add [required_providers block](https://www.terraform.io/docs/configuration/provider-requirements.html#requiring-providers)
```terraform
terraform {

  required_providers {
    codefresh = {
      version = "0.1.0"
      source = "codefresh.io/app/codefresh"
    }
  }
}
```


## [Documentations](./docs)

## [Examples](./examples)

## To configure Codefresh provider:

```hcl
provider "codefresh" {
  api_url = "<MY API URL>" # Default value - https://g.codefresh.io/api
  token = "<MY API TOKEN>" # If token isn't set the provider expects the $CODEFRESH_API_KEY env variable
}
```

Get an API key from [Codefresh](https://g.codefresh.io/user/settings) and set the following scopes:

- Environments-V2
- Pipeline
- Project
- Repos
- Step-Type
- Step-Types
- View

```bash
export CODEFRESH_API_KEY='xyz'
```

## Testing the Provider

## License

Copyright 2020 Codefresh.

The Codefresh Provider is available under [MPL2.0 license](./LICENSE).

