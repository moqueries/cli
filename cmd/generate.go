package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/logs"
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

	err = generator.Generate(generator.GenerateRequest{
		Types:       typs,
		Export:      export,
		Destination: dest,
		Package:     pkg,
		Import:      imp,
		TestImport:  testImp,
	})
	if err != nil {
		logs.Panic("Error generating mock", err)
	}
}
