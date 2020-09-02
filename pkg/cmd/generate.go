package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/pkg/ast"
	"github.com/myshkin5/moqueries/pkg/generator"
	"github.com/myshkin5/moqueries/pkg/logs"
)

// generate gathers details on the environment and calls the generator
func generate(cmd *cobra.Command, typs []string) {
	debug, err := cmd.PersistentFlags().GetBool(debugFlag)
	if err != nil {
		logs.Panic("Error getting debug flag", err)
	}
	export, err := cmd.PersistentFlags().GetBool(exportFlag)
	if err != nil {
		logs.Panic("Error getting export flag", err)
	}
	dest, err := cmd.PersistentFlags().GetString(destinationFlag)
	if err != nil {
		logs.Panic("Error getting destination flag", err)
	}
	pkg, err := cmd.PersistentFlags().GetString(packageFlag)
	if err != nil {
		logs.Panic("Error getting package flag", err)
	}
	imp, err := cmd.PersistentFlags().GetString(importFlag)
	if err != nil {
		logs.Panic("Error getting import flag", err)
	}
	testImp, err := cmd.PersistentFlags().GetBool(testImportFlag)
	if err != nil {
		logs.Panic("Error getting test-import flag", err)
	}

	logs.Init(debug)
	if debug {
		cwd, _ := os.Getwd()
		logs.Debugf("moqueries invoked, debug: %t, export: %t, destination: %s,"+
			" package: %s, import: %s, types: %s, current working directory: %s",
			debug, export, dest, pkg, imp, typs, cwd)
	}

	if export && strings.HasSuffix(dest, "_test.go") {
		logs.Warn("Exported mock in a test file will not be accessible in" +
			" other packages. Remove --export option or set the --destination" +
			" to a non-test file.")
	}

	converter := generator.NewConverter(export)
	gen := generator.New(export, pkg, dest, ast.LoadTypes, converter)

	_, file, err := gen.Generate(typs, imp, testImp)
	if err != nil {
		logs.Panic("Error generating mocks", err)
	}

	tempFile, err := ioutil.TempFile("", "*.go")
	if err != nil {
		logs.Panic("Error creating temp file", err)
	}
	logs.Debugf("Temp file created: %s", tempFile.Name())

	defer func() {
		err = tempFile.Close()
		if err != nil {
			logs.Error("Error closing temp file", err)
		}
	}()

	destDir := filepath.Dir(dest)
	if _, err = os.Stat(destDir); os.IsNotExist(err) {
		err = os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			logs.Error("Error creating destination directory", err)
		}
	}

	restorer := decorator.NewRestorerWithImports(destDir, gopackages.New(destDir))
	err = restorer.Fprint(tempFile, file)
	if err != nil {
		logs.Panic("Invalid mock", err)
	}

	err = os.Rename(tempFile.Name(), dest)
	if err != nil {
		logs.Debugf("Error removing destination file: %v", err)
	}
}
