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

		pkgDestDir := filepath.Join(req.DestinationDir, id.Path)
		pkgDestDir, err := skipDirs(pkgDestDir, req.SkipPkgDirs)
		if err != nil {
			return err
		}
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
			if errors.Is(err, generator.ErrNonExported) {
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

func skipDirs(dir string, skipDirs int) (string, error) {
	orig := dir
	for n := 0; n < skipDirs; n++ {
		if dir == "." {
			return "", fmt.Errorf("%w: skipping %d directories on %s path",
				ErrSkipTooManyPackageDirs, skipDirs, orig)
		}
		idx := strings.Index(dir, string(filepath.Separator))
		if idx == -1 {
			dir = "."
		} else {
			dir = dir[idx+1:]
		}
	}
	return dir, nil
}
