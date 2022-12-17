package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"moqueries.org/cli/logs"
	"moqueries.org/cli/metrics"
)

var summarizeMetricsCmd = &cobra.Command{
	Use:   "summarize-metrics",
	Short: "Summarize metrics logging from multiple runs",
	Args:  cobra.MaximumNArgs(1),
	Run:   summarizeMetrics,
}

func summarizeMetrics(cmd *cobra.Command, files []string) {
	rootSetup(cmd)

	if len(files) > 1 {
		logs.Panicf("Expected one file argument or none (reads stdin)")
	}

	var f *os.File
	var err error
	if len(files) == 0 || files[0] == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(files[0])
		if err != nil {
			logs.Panicf("Error opening file %s: %#v", files[0], err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				logs.Error("error closing log file", err)
			}
		}()
	}

	err = metrics.Summarize(f, logs.Infof)
	if err != nil {
		logs.Panicf("Error summarizing file %s: %#v", files[0], err)
	}
}
