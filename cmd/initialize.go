package cmd

import (
	"github.com/spf13/cobra"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/bulk"
)

var initializeCmd = &cobra.Command{
	Use:   "bulk-initialize",
	Short: "Initialize state for bulk processing",
	Args:  cobra.NoArgs,
	Run:   initialize,
}

func initialize(cmd *cobra.Command, _ []string) {
	rootInfo := rootSetup(cmd)

	if rootInfo.stateFile == "" {
		logs.Panic(stateFileEnvVar+" environment variable is required"+
			" when initializing bulk processing", nil)
	}

	logs.Debugf("Moqueries initialize invoked")

	err := bulk.Initialize(rootInfo.stateFile, rootInfo.workingDir)
	if err != nil {
		logs.Panicf("Error initializing state for bulk processing for %s: %#v",
			rootInfo.stateFile, err)
	}
}
