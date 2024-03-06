package cmd

import (
	"github.com/mirantiscontainers/boundless-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply the blueprint to the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Apply(&blueprint, kubeConfig, operatorUri)
		},
	}

	flags := cmd.Flags()
	addOperatorUriFlag(flags)
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}
