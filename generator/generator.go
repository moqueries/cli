// Package generator generates Moqueries mocks
package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"golang.org/x/tools/go/packages"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/metrics"
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
	WorkingDir     string   `json:"working-dir"`
}

// Generate generates a moq
func Generate(reqs ...GenerateRequest) error {
	m := metrics.NewMetrics(logs.IsDebug, logs.Debugf)
	cache := ast.NewCache(packages.Load, m)
	start := time.Now()
	for _, req := range reqs {
		err := generate(cache, req)
		if err != nil {
			return err
		}
	}
	m.TotalProcessingTimeInc(time.Since(start))
	m.Finalize()
	return nil
}

func generate(cache *ast.Cache, req GenerateRequest) error {
	newConverterFn := func(typ Type, export bool) Converterer {
		return NewConverter(typ, export, cache)
	}
	gen := New(cache, os.Getwd, newConverterFn)

	file, destPath, err := gen.Generate(req)
	if err != nil {
		return fmt.Errorf("error generating moqs: %w", err)
	}

	tempFile, err := ioutil.TempFile("", "*.go")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	defer func() {
		err = tempFile.Close()
		if err != nil {
			logs.Error("Error closing temp file", err)
		}
	}()

	destDir := filepath.Dir(destPath)
	if _, err = os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			logs.Errorf(
				"Error creating destination directory %s from working director %s: %v",
				destDir, req.WorkingDir, err)
		}
	}

	restorer := decorator.NewRestorerWithImports(destDir, gopackages.New(destDir))
	err = restorer.Fprint(tempFile, file)
	if err != nil {
		return fmt.Errorf("invalid moq: %w", err)
	}

	err = os.Rename(tempFile.Name(), destPath)
	if err != nil {
		logs.Debugf("Error removing destination file: %v", err)
	}
	logs.Debugf("Wrote file: %s", destPath)

	return nil
}
