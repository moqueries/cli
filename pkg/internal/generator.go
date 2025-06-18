// Package internal is the internal implementation for generate mocks for an
// entire package or module
package internal

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dave/dst"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/generator"
	"moqueries.org/cli/metrics"
)

// ErrSkipTooManyPackageDirs is returned by Generate when skipPkgDirs requests
// that more directories should be skipped than directories observed
var ErrSkipTooManyPackageDirs = errors.New("skipping too many package dirs")

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

type PackageGenerateRequest struct {
	DestinationDir      string
	SkipPkgDirs         int
	PkgPatterns         []string
	ExcludePkgPathRegex string
}

// Generate generates mocks for several packages at once
func Generate(
	cache TypeCache,
	mProcessor metrics.Metrics,
	genFn GenerateWithTypeCacheFn,
	req PackageGenerateRequest,
) error {
	start := time.Now()

	var reg *regexp.Regexp
	if req.ExcludePkgPathRegex != "" {
		var err error
		reg, err = regexp.Compile(req.ExcludePkgPathRegex)
		if err != nil {
			return fmt.Errorf("%w: could not compile exclude package regex \"%s\"",
				err, req.ExcludePkgPathRegex)
		}
	}

	for _, pkgPattern := range req.PkgPatterns {
		err := cache.LoadPackage(pkgPattern)
		if err != nil {
			return err
		}
	}

	typs := cache.MockableTypes(true)
	logs.Debugf("Mocking %d types", len(typs))

	for _, id := range typs {
		if reg != nil && reg.MatchString(id.Path) {
			logs.Warnf("Skipping %s because of package exclusion %s",
				id.String(), req.ExcludePkgPathRegex)
			continue
		}

		pkgDestDir, err := skipDirs(id.Path, req.SkipPkgDirs)
		if err != nil {
			return err
		}
		pkgDestDir = filepath.Join(req.DestinationDir, pkgDestDir)
		logs.Debugf("Package generating,"+
			" destination-dir: %s,"+
			" type: %s",
			pkgDestDir, id.String())

		err = genFn(cache, generator.GenerateRequest{
			Types:              []string{id.Name},
			Export:             true,
			DestinationDir:     pkgDestDir,
			Import:             id.Path,
			ExcludeNonExported: true,
		})
		if err != nil {
			if errors.Is(err, generator.ErrNonExported) || errors.Is(err, ast.ErrMixedRecvTypes) {
				logs.Debugf("Skipping generation of mock for %s, %s",
					id.String(), err.Error())
				continue
			}
			return err
		}
	}

	mProcessor.TotalProcessingTimeInc(time.Since(start))
	mProcessor.Finalize()

	return nil
}

func skipDirs(pkgPath string, skipDirs int) (string, error) {
	orig := pkgPath
	for range skipDirs {
		if pkgPath == "." {
			return "", fmt.Errorf("%w: skipping %d directories on package %s",
				ErrSkipTooManyPackageDirs, skipDirs, orig)
		}
		idx := strings.Index(pkgPath, string(filepath.Separator))
		if idx == -1 {
			pkgPath = "."
		} else {
			pkgPath = pkgPath[idx+1:]
		}
	}
	return pkgPath, nil
}
