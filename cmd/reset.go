package cmd

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/components"
	"github.com/mirantiscontainers/boundless-cli/pkg/distro"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
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
		PreRunE: actions(loadBlueprint, loadKubeConfig),
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
	log.Info().Msg("Resetting cluster")

	// Determine the distro
	provider, err := distro.GetProvider(&blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	// Uninstall components
	log.Info().Msgf("Reset Boundless Operator resources")
	err = components.RemoveComponents(provider.GetKubeConfig(), blueprint)
	if err != nil {
		return fmt.Errorf("failed to reset components: %w", err)
	}

	log.Info().Msgf("Uninstalling Boundless Operator")
	log.Trace().Msgf("Uninstalling boundless operator using manifest file: %s", operatorUri)
	if err = k8s.DeleteYamlObjects(kubeConfig, operatorUri); err != nil {
		return fmt.Errorf("failed to uninstall Boundless Operator: %w", err)
	}

	// Reset the cluster
	if err := provider.Reset(); err != nil {
		return fmt.Errorf("failed to reset cluster: %w", err)
	}

	return nil
}
