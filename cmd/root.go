package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/logs"
)

const (
	debugFlag   = "debug"
	debugEnvVar = "MOQ_DEBUG"
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
}

// Execute generates one or more mocks
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logs.Panic("Error executing command", err)
	}
}

// rootSetup is called by all subcommands for general setup
func rootSetup(cmd *cobra.Command) {
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
}
