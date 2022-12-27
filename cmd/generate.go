package cmd

import (
	"github.com/spf13/cobra"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/bulk"
	"moqueries.org/cli/generator"
)

const (
	exportFlag = "export"

	destinationFlag    = "destination"
	destinationDirFlag = "destination-dir"

	packageFlag    = "package"
	importFlag     = "import"
	testImportFlag = "test-import"
)

func addGenerateFlags() {
	rootCmd.Flags().Bool(exportFlag, false,
		"If true, generated mocks will be exported and accessible from other packages")

	rootCmd.Flags().String(destinationFlag, "",
		"The file path where mocks are generated. Relative paths are "+
			"allowed (relative to the directory containing the generate "+
			"directive or relative to the current directory) (defaults to "+
			"./moq_<type>.go when exported or ./moq_<type>_test.go when not "+
			"exported)")
	addDestinationDirFlag(rootCmd)

	rootCmd.Flags().String(packageFlag, "",
		"The package to generate code into (defaults to the test package of "+
			"the destination directory when --export=false or the package of "+
			"the destination directory when --export=true)")
	rootCmd.Flags().String(importFlag, ".",
		"The package containing the type (interface or function type) to be "+
			"mocked (defaults to the directory containing generate directive)")
	rootCmd.Flags().Bool(testImportFlag, false,
		"Look for the types to be mocked in the test package")
}

func addDestinationDirFlag(cmd *cobra.Command) {
	cmd.Flags().String(destinationDirFlag, "",
		"The file directory where mocks are generated relative. Relative "+
			"paths are allowed (relative to the directory containing the "+
			"generate directive (or relative to the current directory) "+
			"(defaults to .)")
}

// generate gathers details on the environment and calls the generator
func generate(cmd *cobra.Command, typs []string) {
	root := rootSetup(cmd)

	export, err := cmd.Flags().GetBool(exportFlag)
	if err != nil {
		logs.Panic("Error getting export flag", err)
	}
	dest, err := cmd.Flags().GetString(destinationFlag)
	if err != nil {
		logs.Panic("Error getting destination flag", err)
	}
	destDir, err := cmd.Flags().GetString(destinationDirFlag)
	if err != nil {
		logs.Panic("Error getting destination dir flag", err)
	}
	pkg, err := cmd.Flags().GetString(packageFlag)
	if err != nil {
		logs.Panic("Error getting package flag", err)
	}
	imp, err := cmd.Flags().GetString(importFlag)
	if err != nil {
		logs.Panic("Error getting import flag", err)
	}
	testImp, err := cmd.Flags().GetBool(testImportFlag)
	if err != nil {
		logs.Panic("Error getting test-import flag", err)
	}

	logs.Debugf("Moqueries generate invoked,"+
		" export: %t,"+
		" destination: %s,"+
		" destination-dir: %s,"+
		" package: %s,"+
		" import: %s,"+
		" types: %s",
		export, dest, destDir, pkg, imp, typs)

	req := generator.GenerateRequest{
		Types:          typs,
		Export:         export,
		Destination:    dest,
		DestinationDir: destDir,
		Package:        pkg,
		Import:         imp,
		TestImport:     testImp,
		WorkingDir:     root.workingDir,
	}
	if root.stateFile == "" {
		err = generator.Generate(req)
		if err != nil {
			logs.Panicf("Error generating mock for %s in %s: %#v", typs, imp, err)
		}
	} else {
		err = bulk.Append(root.stateFile, req)
		if err != nil {
			logs.Panicf("Error appending mock request for %s in %s: %#v", typs, imp, err)
		}
	}
}
