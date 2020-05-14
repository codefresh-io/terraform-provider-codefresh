# Terraform provider for Codefresh

This provider was initialized by [LightStep](https://lightstep.com/) and will be maintained as the official Terraform provider for Codefresh.  

The provider is still under development, and can be used as a terraform [third-party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) only.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.11+ ;
- [Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin).

## Build

```sh
go build
```

## Usage

Compile or take from the [Releases](https://github.com/codefresh-contrib/terraform-provider-codefresh/releases) `terraform-provider-codefresh` binary and place it locally according the Terraform plugins [documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

For Linux OS it can be:

- _~/.terraform.d/plugins/linux\_amd64_
- _./terraform.d/plugins/linux\_amd64_. The relative path in your Terraform project.

## Configuration

There are two environment variables to configure Codefresh provider:

- `CODEFRESH_API_URL`. Default value - https://g.codefresh.io/api ;
- `CODEFRESH_API_KEY`.

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

## Examples

See the [examples](examples/).