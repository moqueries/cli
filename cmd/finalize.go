package cmd

import (
	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/bulk"
	"github.com/myshkin5/moqueries/logs"
)

var finalizeCmd = &cobra.Command{
	Use:   "bulk-finalize",
	Short: "Finalize bulk processing and generate mocks",
	Args:  cobra.NoArgs,
	Run:   finalize,
}

func finalize(cmd *cobra.Command, _ []string) {
	rootInfo := rootSetup(cmd)

	if rootInfo.stateFile == "" {
		logs.Panic(stateFileEnvVar+" environment variable is required"+
			" when finalizing bulk processing", nil)
	}

	logs.Debugf("Moqueries finalize invoked")

	err := bulk.Finalize(rootInfo.stateFile, rootInfo.workingDir)
	if err != nil {
		logs.Panicf("Error finalizing bulk processing for %s: %#v",
			rootInfo.stateFile, err)
	}
}
