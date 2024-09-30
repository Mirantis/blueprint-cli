package commands

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/components"
	"github.com/mirantiscontainers/boundless-cli/pkg/distro"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/rs/zerolog/log"
)

// Update updates the Blueprint Operator and applies the components defined in the blueprint
func Update(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) error {
	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	needsUpgrade, err := provider.NeedsUpgrade(blueprint)
	if err != nil {
		return err
	}

	if needsUpgrade {
		if err := provider.ValidateProviderUpgrade(blueprint); err != nil {
			return fmt.Errorf("provider failed pre-upgrade validation and may require manual changes: %w", err)
		}

		log.Info().Msgf("Updating provider")
		if err := provider.Upgrade(); err != nil {
			return fmt.Errorf("failed to update provider: %w", err)
		}
	}

	log.Info().Msgf("Applying Blueprint Operator resources")
	if err := components.ApplyBlueprint(kubeConfig, blueprint); err != nil {
		return fmt.Errorf("failed to update components: %w", err)
	}

	log.Info().Msgf("Finished updating Blueprint Operator")
	return nil
}
