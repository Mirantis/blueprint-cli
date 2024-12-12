package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version, commit, date = "", "", "" // These are always injected at build time

	// verbose flag
	verbose bool
)

// versionCmd creates the version command
func versionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Print version/build info",
		Long:  "Print version/build information",
		RunE:  runVersion,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "Print more detailed version information")

	return &cmd
}

// runVersion prints the version information
func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("Version: %s\n", version)
	if verbose {
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
	}

	return nil
}
