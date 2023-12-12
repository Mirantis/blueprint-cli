package cmd

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/internal/color"
	"github.com/spf13/cobra"
)

var (
	version, commit, date = "0.0.0-dev", "", ""
)

func versionCmd() *cobra.Command {
	var short bool

	command := cobra.Command{
		Use:   "version",
		Short: "Print version/build info",
		Long:  "Print version/build information",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion(short)
		},
	}

	command.PersistentFlags().BoolVarP(&short, "short", "s", false, "Prints bctl version info in short format")

	return &command
}

func printVersion(short bool) {
	const fmat = "%-20s %s\n"

	var outputColor color.Paint

	if short {
		outputColor = -1
	} else {
		outputColor = color.Cyan
	}
	printTuple(fmat, "Version", version, outputColor)
	printTuple(fmat, "Commit", commit, outputColor)
	printTuple(fmat, "Date", date, outputColor)
}

func printTuple(fmat, section, value string, outputColor color.Paint) {
	if value != "" {
		if outputColor != -1 {
			fmt.Fprintf(out, fmat, color.Colorize(section+":", outputColor), value)
			return
		}
		fmt.Fprintf(out, fmat, section, value)
	}
}
