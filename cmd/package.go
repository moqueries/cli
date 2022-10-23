package cmd

import (
	"github.com/spf13/cobra"

	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/pkg"
)

const skipPkgDirsFlag = "skip-pkg-dirs"

func packageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package [packages]",
		Short: "Generate mocks for a package",
		Run:   pkgGen,
	}

	addDestinationDirFlag(cmd)

	cmd.PersistentFlags().Int(skipPkgDirsFlag, 0,
		"Skips specified number of directories in the package path when "+
			"determining the destination directory (defaults to 0)")

	return cmd
}

func pkgGen(cmd *cobra.Command, pkgPatterns []string) {
	destDir, err := cmd.Flags().GetString(destinationDirFlag)
	if err != nil {
		logs.Panic("Error getting destination dir flag", err)
	}
	skipPkgDirs, err := cmd.Flags().GetInt(skipPkgDirsFlag)
	if err != nil {
		logs.Panic("Error getting the skip package dirs flag", err)
	}

	err = pkg.Generate(destDir, skipPkgDirs, pkgPatterns...)
	if err != nil {
		logs.Panicf("Error generating mocks for %s packages: %#v", pkgPatterns, err)
	}
}
