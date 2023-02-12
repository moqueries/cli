// Package generator generates Moqueries mocks
package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"golang.org/x/tools/go/packages"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/metrics"
)

// GenerateRequest contains all the parameters needed to call Generate
type GenerateRequest struct {
	Types          []string `json:"types"`
	Export         bool     `json:"export"`
	Destination    string   `json:"destination"`
	DestinationDir string   `json:"destination-dir"`
	Package        string   `json:"package"`
	Import         string   `json:"import"`
	TestImport     bool     `json:"test-import"`
	// WorkingDir is the current working directory. Optional, in which case
	// os.Getwd is used. Useful in cases where a request is serialized then
	// rerun in bulk processing from a different working directory. WorkingDir
	// is used for relative-path imports and relative path destination
	// files/directories.
	WorkingDir string `json:"working-dir"`
	// ExcludeNonExported causes the generator to exclude non-exported types
	// from the generated mock. This includes possibly returning ErrNonExported
	// if after exclusion, nothing is left to be written as a mock.
	ExcludeNonExported bool `json:"exclude-non-exported"`
}

// ErrNonExported is returned by Generate when ExcludeNonExported is set to
// true resulting in an empty mock.
var ErrNonExported = errors.New("non-exported types")

// Generate generates a moq
func Generate(reqs ...GenerateRequest) error {
	m := metrics.NewMetrics(logs.IsDebug, logs.Debugf)
	cache := ast.NewCache(packages.Load, os.Stat, os.ReadFile, m)
	start := time.Now()
	for _, req := range reqs {
		err := GenerateWithTypeCache(cache, req)
		if err != nil {
			return err
		}
	}
	m.TotalProcessingTimeInc(time.Since(start))
	m.Finalize()
	return nil
}

//go:generate moqueries TypeCache

// TypeCache defines the interface to the Cache type
type TypeCache interface {
	Type(id dst.Ident, contextPkg string, testImport bool) (ast.TypeInfo, error)
	IsComparable(expr dst.Expr) (bool, error)
	IsDefaultComparable(expr dst.Expr) (bool, error)
	FindPackage(dir string) (string, error)
}

// GenerateWithTypeCache generates a single moq using the provided type cache.
// This function is exposed for use in bulk operations that have already loaded
// a type.
func GenerateWithTypeCache(cache TypeCache, req GenerateRequest) error {
	newConverterFn := func(typ Type, export bool) Converterer {
		return NewConverter(typ, export, cache)
	}
	gen := New(cache, os.Getwd, newConverterFn)

	resp, err := gen.Generate(req)
	if err != nil {
		return err
	}

	destDir := filepath.Dir(resp.DestPath)
	if _, err = os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			logs.Errorf(
				"Error creating destination directory %s from working director %s: %v",
				destDir, req.WorkingDir, err)
		}
	}

	tempFile, err := ioutil.TempFile(destDir, "*.go-gen")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	defer func() {
		err = tempFile.Close()
		if err != nil {
			logs.Error("Error closing temp file", err)
		}
	}()

	restorer := decorator.NewRestorerWithImports(resp.OutPkgPath, gopackages.New(destDir))
	err = restorer.Fprint(tempFile, resp.File)
	if err != nil {
		return fmt.Errorf("invalid moq: %w", err)
	}

	err = os.Rename(tempFile.Name(), resp.DestPath)
	if err != nil {
		logs.Debugf("Error removing destination file: %v", err)
	}
	logs.Debugf("Wrote file: %s", resp.DestPath)

	return nil
}
