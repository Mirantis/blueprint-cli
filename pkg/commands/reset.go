package commands

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/components"
	"github.com/mirantiscontainers/boundless-cli/pkg/distro"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/rs/zerolog/log"
)

// Reset resets the cluster
func Reset(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig, operatorUri string) error {
	log.Info().Msg("Resetting cluster")

	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
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
