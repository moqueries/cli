package cmd

import (
	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/pkg/logs"
)

const (
	debugFlag       = "debug"
	exportFlag      = "export"
	destinationFlag = "destination"
	packageFlag     = "package"
	importFlag      = "import"
	testImportFlag  = "test-import"
)

var rootCmd = &cobra.Command{
	Use:   "moqueries [interfaces and/or function types to mock]",
	Short: "Moqueries generates simple but thread-safe mocks",
	Args:  cobra.MinimumNArgs(1),
	Run:   generate,
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().Bool(debugFlag, false,
		"If true, debugging output will be logged")
	rootCmd.PersistentFlags().Bool(exportFlag, false,
		"If true, mocks will be exported and accessible from other packages")

	rootCmd.PersistentFlags().String(destinationFlag, "",
		"File path where mocks are generated relative to directory containing generate directive")
	err := rootCmd.MarkPersistentFlagRequired(destinationFlag)
	if err != nil {
		logs.Panic("Error configuring required flag", err)
	}

	rootCmd.PersistentFlags().String(packageFlag, "",
		"Package generated code will be created in (defaults to <destination dir>_test)")
	rootCmd.PersistentFlags().String(importFlag, ".",
		"Package containing interface to be mocked (defaults to directory containing generate directive)")
	rootCmd.PersistentFlags().Bool(testImportFlag, false,
		"Look for the types to be mocked in the test package")
}

// Execute generates one or more mocks
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logs.Panic("Error executing command", err)
	}
}
