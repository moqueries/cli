package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

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
	public, err := cmd.PersistentFlags().GetBool(publicFlag)
	if err != nil {
		logs.Panic("Error getting public flag", err)
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
		logs.Debugf("moqueries invoked, debug: %t, public: %t, destination: %s,"+
			" package: %s, import: %s, types: %s, current working directory: %s",
			debug, public, dest, pkg, imp, typs, cwd)
	}

	converter := generator.NewConverter()
	gen := generator.New(public, pkg, dest, ast.LoadTypes, converter)

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
