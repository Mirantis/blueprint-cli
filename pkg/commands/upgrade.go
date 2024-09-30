package commands

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/rs/zerolog/log"
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
	return nil
}
