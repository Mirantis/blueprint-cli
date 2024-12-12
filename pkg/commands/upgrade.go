package commands

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mirantiscontainers/blueprint-cli/pkg/distro"
	"github.com/mirantiscontainers/blueprint-cli/pkg/k8s"
	"github.com/mirantiscontainers/blueprint-cli/pkg/types"
)

// Upgrade upgrades the Blueprint Operator
func Upgrade(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) error {
	uri, err := determineOperatorUri(blueprint.Spec.Version)
	if err != nil {
		return fmt.Errorf("failed to determine operator URI: %w", err)
	}

	log.Info().Msgf("Upgrading Blueprint Operator using manifest file %q", uri)
	if err := k8s.ApplyYaml(kubeConfig, uri); err != nil {
		return fmt.Errorf("failed to upgrade blueprint operator: %w", err)
	}

	log.Info().Msgf("Finished updating Blueprint Operator")

	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}
	if err := provider.SetupClient(); err != nil {
		return fmt.Errorf("failed to setup client: %w", err)
	}
	// Wait for the pods to be ready
	if err := provider.WaitForPods(); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}
	return nil
}
