// Package pkg is used to generate mocks for an entire package or module
package pkg

import (
	"os"

	"golang.org/x/tools/go/packages"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/generator"
	"moqueries.org/cli/metrics"
	"moqueries.org/cli/pkg/internal"
)

type PackageGenerateRequest = internal.PackageGenerateRequest

// Generate generates mocks for several packages at once
func Generate(req PackageGenerateRequest) error {
	m := metrics.NewMetrics(logs.IsDebug, logs.Debugf)
	cache := ast.NewCache(packages.Load, os.Stat, os.ReadFile, m)

	err := internal.Generate(cache, m, generator.GenerateWithTypeCache, req)
	if err != nil {
		return err
	}

	return nil
}
