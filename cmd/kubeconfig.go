package cmd

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/distro"
	"github.com/mirantiscontainers/boundless-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func kubeConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kubeconfig",
		Short:   "Generate kubeconfig file for the Blueprint",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKubeConfig()
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)

	return cmd
}
func runKubeConfig() error {
	log.Info().Msgf("Generating kubeconfig for blueprint %s", blueprintFlag)

	// Determine the distro
	provider, err := distro.GetProvider(&blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	// Check if the cluster exists
	exists, err := provider.Exists()
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("cluster doesn't exist: %s", blueprint.Metadata.Name)
	}
	log.Info().Msgf("Cluster %q exists", blueprint.Metadata.Name)

	// Create kubeconfig
	if provider.Type() == constants.ProviderK0s {

	} else if provider.Type() == constants.ProviderKind {
		if err := utils.ExecCommand(fmt.Sprintf("kind get kubeconfig --name %s", blueprint.Metadata.Name)); err != nil {
			return fmt.Errorf("failed to get kubeconfig for cluster %s : %w ", blueprint.Metadata.Name, err)
		}
	} else if provider.Type() == constants.ProviderExisting {

	}

	// Set the context

	log.Info().Msgf("Finished generating kubeconfig")

	return nil
}
