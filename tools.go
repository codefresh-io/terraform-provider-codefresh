//go:build tools
// +build tools

package main

import (
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlint"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
