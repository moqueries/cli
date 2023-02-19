package cmd

import (
	"github.com/spf13/cobra"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/pkg"
)

const (
	skipPkgDirsFlag         = "skip-pkg-dirs"
	excludePkgPathRegexFlag = "exclude-pkg-path-regex"
)

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
	cmd.PersistentFlags().String(excludePkgPathRegexFlag, "",
		"Specifies a regular expression used to exclude package paths")

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
	excludePkgPathRegex, err := cmd.Flags().GetString(excludePkgPathRegexFlag)
	if err != nil {
		logs.Panic("Error getting the skip package dirs flag", err)
	}

	err = pkg.Generate(pkg.PackageGenerateRequest{
		DestinationDir:      destDir,
		SkipPkgDirs:         skipPkgDirs,
		PkgPatterns:         pkgPatterns,
		ExcludePkgPathRegex: excludePkgPathRegex,
	})
	if err != nil {
		logs.Panicf("Error generating mocks for %s packages: %#v", pkgPatterns, err)
	}
}
