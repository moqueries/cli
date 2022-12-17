// Package pkg is used to generate mocks for an entire package or module
package pkg

import (
	"golang.org/x/tools/go/packages"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/generator"
	"moqueries.org/cli/logs"
	"moqueries.org/cli/metrics"
	"moqueries.org/cli/pkg/internal"
)

// Generate generates mocks for several packages at once
func Generate(destinationDir string, skipPkgDirs int, pkgPatterns ...string) error {
	m := metrics.NewMetrics(logs.IsDebug, logs.Debugf)
	cache := ast.NewCache(packages.Load, m)

	err := internal.Generate(cache, m, generator.GenerateWithTypeCache, destinationDir, skipPkgDirs, pkgPatterns)
	if err != nil {
		return err
	}

	return nil
}
