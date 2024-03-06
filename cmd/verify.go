package cmd

import (
	"github.com/mirantiscontainers/boundless-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func verifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verifies the blueprint is valid and can be applied to the cluster. Specifically checks helm chart addons",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Verify(&blueprint, kubeConfig)
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)

	return cmd
}
