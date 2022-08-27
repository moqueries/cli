// Package internal is the internal implementation for generate mocks for an
// entire package or module
package internal

import (
	"path/filepath"
	"time"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/metrics"
)

//go:generate moqueries TypeCache

// TypeCache defines the interface to the Cache type
type TypeCache interface {
	LoadPackage(pkgPattern string) error
	MockableTypes(onlyExported bool) []dst.Ident
	generator.TypeCache
}

//go:generate moqueries GenerateWithTypeCacheFn

// GenerateWithTypeCacheFn is the function type for generator.GenerateWithTypeCache
type GenerateWithTypeCacheFn func(cache generator.TypeCache, req generator.GenerateRequest) error

// Generate generates mocks for several packages at once
func Generate(
	cache TypeCache,
	mProcessor metrics.Metrics,
	genFn GenerateWithTypeCacheFn,
	destDir string,
	pkgPatterns []string,
) error {
	start := time.Now()
	for _, pkgPattern := range pkgPatterns {
		err := cache.LoadPackage(pkgPattern)
		if err != nil {
			return err
		}
	}

	typs := cache.MockableTypes(true)
	logs.Debugf("Mocking %d types", len(typs))

	for _, id := range typs {
		pkgDestDir := filepath.Join(destDir, id.Path)
		logs.Debugf("Package generating,"+
			" destination-dir: %s,"+
			" type: %s",
			pkgDestDir, id.String())

		err := genFn(cache, generator.GenerateRequest{
			Types:          []string{id.Name},
			Export:         true,
			DestinationDir: pkgDestDir,
			Import:         id.Path,
		})
		if err != nil {
			return err
		}
	}

	mProcessor.TotalProcessingTimeInc(time.Since(start))
	mProcessor.Finalize()

	return nil
}
