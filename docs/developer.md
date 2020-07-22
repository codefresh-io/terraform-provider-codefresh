## Developer guide

We are using [Terraform Plugin SDK v2](https://github.com/hashicorp/terraform-plugin-sdk/tree/v2.0.0-rc.2)  

### Run with Delve Debugger
[Reference guide](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html#support-for-debuggable-provider-binaries)  
[sdk code](https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.0.0-rc.2/plugin/debug.go#L97)  

run pluging with set env CODEFRESH_PLUGIN_DEBUG=true in delve  
- for vscode set launch.json like this:
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "terraform-provider-codefresh",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "port": 2345,
            "host": "127.0.0.1",
            "env": {"CODEFRESH_PLUGIN_DEBUG": "true"},
            "program": "/d1/home/kosta/devel/go/src/github.com/codefresh-io/terraform-provider-codefresh/main.go",
            "showLog": true,
            "trace": "verbose"
        }
    ]
}
```
- copy value of TF_REATTACH_PROVIDERS from the output of debug console and set it for terraform exec:
```
export TF_REATTACH_PROVIDERS='{"registry.terraform.io/-/codefresh":{"Protocol":"grpc","Pid":614875,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin369955425"}}}'

terraform apply
```





