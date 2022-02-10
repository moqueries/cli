package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/logs"
)

// GenerateRequest contains all the parameters needed to call Generate
type GenerateRequest struct {
	Types       []string
	Export      bool
	Destination string
	Package     string
	Import      string
	TestImport  bool
}

// Generate generates a moq
func Generate(reqs ...GenerateRequest) error {
	for _, req := range reqs {
		err := generate(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func generate(req GenerateRequest) error {
	if req.Export && strings.HasSuffix(req.Destination, "_test.go") {
		logs.Warn("Exported moq in a test file will not be accessible in" +
			" other packages. Remove --export option or set the --destination" +
			" to a non-test file.")
	}

	if req.Destination == "" {
		dest := "moq_"
		for n, typ := range req.Types {
			dest += strings.ToLower(typ)
			if n+1 < len(req.Types) {
				dest += "_"
			}
		}
		if !req.Export {
			dest += "_test"
		}
		dest += ".go"
		req.Destination = dest
	}

	cache := ast.NewCache(ast.LoadTypes)
	converter := NewConverter(req.Export, cache)
	gen := New(req.Export, req.Package, req.Destination, ast.FindPackage, cache, converter)

	_, file, err := gen.Generate(req.Types, req.Import, req.TestImport)
	if err != nil {
		return fmt.Errorf("error generating moqs: %w", err)
	}

	tempFile, err := ioutil.TempFile("", "*.go")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	logs.Debugf("Temp file created: %s", tempFile.Name())

	defer func() {
		err = tempFile.Close()
		if err != nil {
			logs.Error("Error closing temp file", err)
		}
	}()

	destDir := filepath.Dir(req.Destination)
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

	err = os.Rename(tempFile.Name(), req.Destination)
	if err != nil {
		logs.Debugf("Error removing destination file: %v", err)
	}

	return nil
}
