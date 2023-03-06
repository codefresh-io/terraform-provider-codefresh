## Development

We are currently using [Terraform Plugin SDK v2](https://github.com/hashicorp/terraform-plugin-sdk).

It is possible that we will switch to the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) sometime in the future.

### Prerequisites (other than Terraform)

- GNU Make
- [Go](https://golang.org/doc/install) `1.18.x` (minimum supported Go version required to build the provider).

### Building and Running a Local Build of the Provider

```bash
make install
```

Set the [developer overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) to point Terraform at the locally-built binary:

```terraform
# `~/.terraformrc (Windows: %APPDATA%/.terraformrc)
provider_installation {
    dev_overrides {
        "codefresh.io/codefresh" = "[REPLACE WITH GOPATH]/bin"
    }
    direct {}
}
```

Note that if developer overrides are set, Terraform will ignore the version pinned in `versions.tf`, so you do not need to remove the version pin when testing. You can keep it.

### Debugging with Delve

[Reference guide](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html#support-for-debuggable-provider-binaries)  

[SDK code](https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.0.0-rc.2/plugin/debug.go#L97)  

Run the provider with `CODEFRESH_PLUGIN_DEBUG=true` in Delve debugger.

For vscode, set `launch.json` as follows:

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.1",
    "configurations": [
        {
            "name": "terraform-provider-codefresh",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "port": 2345,
            "host": "127.0.0.1",
            "env": {"CODEFRESH_PLUGIN_DEBUG": "true"},
            "program": "/home/mitchellh/go/src/github.com/codefresh-io/terraform-provider-codefresh/main.go",
            "showLog": true,
            "trace": "verbose"
        }
    ]
}
```

Then, copy the value of `TF_REATTACH_PROVIDERS` from the output of debug console and set it for terraform exec:

```bash
export TF_REATTACH_PROVIDERS='{"registry.terraform.io/-/codefresh":{"Protocol":"grpc","Pid":614875,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin369955425"}}}'

terraform apply
```
