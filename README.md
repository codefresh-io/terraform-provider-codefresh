# Terraform Provider for Codefresh

This is the official Terraform Provider for Codefresh.

Terraform Registry: [registry.terraform.io/providers/codefresh-io/codefresh](https://registry.terraform.io/providers/codefresh-io/codefresh/latest)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) `1.x.x` ;

## Download the Provider

Download and extract terraform-provider-codefresh from [releases](https://github.com/codefresh-io/terraform-provider-codefresh/releases)

## Building the Provider

```sh
go build -o terraform-provider-codefresh
```

## Using the Provider

In `versions.tf`:

```terraform
terraform {
  required_providers {
    codefresh = {
      version = "x.y.z" # Optional but recommended; replace with latest semantic version
      source = "codefresh.io/codefresh"
    }
  }
}
```


## [Documentation](./docs)

The documentation is generated using [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).

See: [CONTRIBUTING.md](./CONTRIBUTING.md#documentation)

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

```bash
make testacc
```

## Contributors

[![All Contributors](https://img.shields.io/github/all-contributors/codefresh-io/terraform-provider-aws?color=ee8449&style=flat-square)](#contributors)


<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## Acknowledgements

_This provider was initialized by [LightStep](https://lightstep.com/)_.

## License

Copyright 2023 Codefresh.

The Codefresh Provider is available under [MPL2.0 license](./LICENSE).
