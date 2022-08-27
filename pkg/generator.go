// Package pkg is used to generate mocks for an entire package or module
package pkg

import (
	"golang.org/x/tools/go/packages"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/metrics"
	"github.com/myshkin5/moqueries/pkg/internal"
)

// Generate generates mocks for several packages at once
func Generate(destinationDir string, pkgPatterns ...string) error {
	m := metrics.NewMetrics(logs.IsDebug, logs.Debugf)
	cache := ast.NewCache(packages.Load, m)

	err := internal.Generate(cache, m, generator.GenerateWithTypeCache, destinationDir, pkgPatterns)
	if err != nil {
		return err
	}

	return nil
}
