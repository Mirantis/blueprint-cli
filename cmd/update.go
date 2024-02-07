package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mirantiscontainers/boundless-cli/pkg/components"
	"github.com/mirantiscontainers/boundless-cli/pkg/distro"
)

// updateCmd represents the apply command
func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the cluster according to the blueprint",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd)
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}

func runUpdate(cmd *cobra.Command) error {
	// Determine the distro
	provider, err := distro.GetProvider(&blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	needsUpgrade, err := provider.NeedsUpgrade(&blueprint)
	if err != nil {
		return err
	}

	if needsUpgrade {
		if err := provider.ValidateProviderUpgrade(&blueprint); err != nil {
			return fmt.Errorf("provider failed pre-upgrade validation and may require manual changes: %w", err)
		}

		log.Info().Msgf("Updating provider")
		if err := provider.Upgrade(); err != nil {
			return fmt.Errorf("failed to update provider: %w", err)
		}
	}

	log.Info().Msgf("Applying Boundless Operator resources")
	if err := components.ApplyBlueprint(kubeConfig, blueprint); err != nil {
		return fmt.Errorf("failed to update components: %w", err)
	}

	log.Info().Msgf("Finished updating Boundless Operator")
	return nil
}
