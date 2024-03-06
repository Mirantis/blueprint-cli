package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mirantiscontainers/boundless-cli/pkg/commands"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Get the status of the blueprint",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return commands.AddonSpecificStatus(kubeConfig, args[0])
			}

			return commands.Status(kubeConfig)
		},
	}

	flags := cmd.Flags()
	addOperatorUriFlag(flags)
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}
