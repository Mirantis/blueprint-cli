package cmd

import (
	"github.com/mirantiscontainers/boundless-cli/pkg/commands"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// resetCmd represents the apply command
func resetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the cluster to a clean state",
		Long: `
Reset the cluster to a clean state.

For a cluster with k0s, this will remove traces of k0s from all the hosts.
For a cluster with kind, it will delete the cluster (same as 'kind delete cluster <CLUSTER NAME>').
For a cluster with an external Kubernetes provider, this will remove Blueprint Operator and all associated resources.
`,
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msgf("Resetting blueprint at %s", blueprintFlag)
			return commands.Reset(&blueprint, kubeConfig, force)
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)
	addForceFlag(flags)
	return cmd
}
