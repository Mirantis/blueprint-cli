package cmd

import (
	"github.com/mirantiscontainers/boundless-cli/pkg/commands"
	"github.com/spf13/cobra"
)

// updateCmd represents the apply command
func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the cluster according to the blueprint",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Update(&blueprint, kubeConfig)
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}
