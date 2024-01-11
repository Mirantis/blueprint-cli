package cmd

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/internal/distro"
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
For a cluster with an external Kubernetes provider, this will remove Boundless Operator and all associated resources.
`,
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReset()
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)
	return cmd
}

func runReset() error {
	log.Info().Msgf("Resetting blueprint %s", blueprintFlag)

	// Determine the distro
	provider, err := distro.GetProvider(&blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	// Reset the cluster
	if err := provider.Reset(); err != nil {
		return fmt.Errorf("failed to reset cluster: %w", err)
	}

	return nil
}
