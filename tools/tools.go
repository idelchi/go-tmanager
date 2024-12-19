//go:build tools

// Package tools provides the tools used by the go generate directives in the project.
// Run `go mod download` to download the tools.
package tools

import (
	_ "github.com/dmarkham/enumer"
)

// https://go.dev/wiki/GoGenerateTools
