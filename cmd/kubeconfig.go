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

	if provider.Type() == constants.ProviderK0s {
		k0sConfig, err := distro.CreateTempK0sConfig(&blueprint)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get k0s config path")
		}

		if err := utils.ExecCommand(fmt.Sprintf("k0sctl kubeconfig --config %s", k0sConfig)); err != nil {
			return fmt.Errorf("failed to get kubeconfig for k0s cluster %s : %w ", blueprint.Metadata.Name, err)
		}

	} else if provider.Type() == constants.ProviderKind {
		if err := utils.ExecCommand(fmt.Sprintf("kind get kubeconfig --name %s", blueprint.Metadata.Name)); err != nil {
			return fmt.Errorf("failed to get kubeconfig for kind cluster %s : %w ", blueprint.Metadata.Name, err)
		}
	} else if provider.Type() == constants.ProviderExisting {
		return fmt.Errorf("provider: %s not supported.", constants.ProviderExisting)
	}

	return nil
}
