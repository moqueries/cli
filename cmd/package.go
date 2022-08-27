package cmd

import (
	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/pkg"
)

func packageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package [packages]",
		Short: "Generate mocks for a package",
		Run:   pkgGen,
	}

	addDestinationDirFlag(cmd)

	return cmd
}

func pkgGen(cmd *cobra.Command, pkgPatterns []string) {
	destDir, err := cmd.Flags().GetString(destinationDirFlag)
	if err != nil {
		logs.Panic("Error getting destination dir flag", err)
	}

	err = pkg.Generate(destDir, pkgPatterns...)
	if err != nil {
		logs.Panicf("Error generating mocks for %s packages: %#v", pkgPatterns, err)
	}
}
