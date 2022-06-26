package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"golang.org/x/tools/go/packages"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/metrics"
)

// ErrInvalidConfig is returned when configuration values are invalid or
// conflict with each other
var ErrInvalidConfig = errors.New("invalid configuration")

// GenerateRequest contains all the parameters needed to call Generate
type GenerateRequest struct {
	Types          []string
	Export         bool
	Destination    string
	DestinationDir string
	Package        string
	Import         string
	TestImport     bool
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
	if req.Export && strings.HasSuffix(req.Destination, "_test.go") {
		logs.Warn("Exported moq in a test file will not be accessible in" +
			" other packages. Remove --export option or set the --destination" +
			" to a non-test file.")
	}

	if req.Destination != "" && req.DestinationDir != "" {
		return fmt.Errorf("both destination and destination dir flags"+
			"must not be present together: %w", ErrInvalidConfig)
	}

	if req.Destination == "" {
		dest := "moq_"
		for n, typ := range req.Types {
			dest += strings.ToLower(typ)
			if n+1 < len(req.Types) {
				dest += "_"
			}
		}
		if !req.Export || (req.Package != "" && strings.HasSuffix(req.Package, testPkgSuffix)) {
			dest += testPkgSuffix
		}
		dest += ".go"
		req.Destination = dest
	}

	newConverterFn := func(typ Type, export bool) Converterer {
		return NewConverter(typ, export, cache)
	}
	gen := New(cache, newConverterFn)

	_, file, err := gen.Generate(req)
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

	destDir := req.DestinationDir
	if destDir == "" {
		destDir = filepath.Dir(req.Destination)
	}
	if _, err = os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			wd, _ := os.Getwd()
			logs.Errorf(
				"Error creating destination directory %s from working director %s: %v",
				destDir, wd, err)
		}
	}

	restorer := decorator.NewRestorerWithImports(destDir, gopackages.New(destDir))
	err = restorer.Fprint(tempFile, file)
	if err != nil {
		return fmt.Errorf("invalid moq: %w", err)
	}

	out := path.Join(destDir, filepath.Base(req.Destination))
	err = os.Rename(tempFile.Name(), out)
	if err != nil {
		logs.Debugf("Error removing destination file: %v", err)
	}
	logs.Debugf("Wrote file: %s", out)

	return nil
}
