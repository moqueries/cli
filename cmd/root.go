// Package cmd centralizes the command line interface
package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"moqueries.org/runtime/logs"
)

const (
	debugFlag = "debug"

	stateFileEnvVar = "MOQ_BULK_STATE_FILE"
	debugEnvVar     = "MOQ_DEBUG"
)

var rootCmd = &cobra.Command{
	Use:   "moqueries [interfaces and/or function types to mock]",
	Short: "Moqueries generates lock-free mocks",
	Args:  cobra.MinimumNArgs(1),
	Run:   generate,
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().Bool(debugFlag, false,
		"If true, debugging output will be logged")

	addGenerateFlags()

	rootCmd.AddCommand(summarizeMetricsCmd)
	rootCmd.AddCommand(initializeCmd)
	rootCmd.AddCommand(finalizeCmd)
	rootCmd.AddCommand(packageCmd())
}

// Execute generates one or more mocks
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logs.Panic("Error executing command", err)
	}
}

type rootInfo struct {
	workingDir string
	stateFile  string
}

// rootSetup is called by all subcommands for general setup
func rootSetup(cmd *cobra.Command) rootInfo {
	debug, err := cmd.Root().PersistentFlags().GetBool(debugFlag)
	if err != nil {
		logs.Panic("Error getting debug flag", err)
	}

	debugStr, ok := os.LookupEnv(debugEnvVar)
	if ok {
		envVar, err := strconv.ParseBool(debugStr)
		if err != nil {
			logs.Panic("Error parsing "+debugEnvVar+" environment variable", err)
		}
		debug = debug || envVar
	}

	logs.Init(debug)

	workingDir, err := os.Getwd()
	if err != nil {
		logs.Panic("Could not get working directory", err)
	}

	stateFile := os.Getenv(stateFileEnvVar)

	logs.Debugf("Moqueries root info,"+
		" bulk processing state file: %s,"+
		" working directory: %s",
		stateFile, workingDir)

	return rootInfo{
		workingDir: workingDir,
		stateFile:  stateFile,
	}
}
