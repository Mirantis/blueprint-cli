package commands

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/rs/zerolog/log"
)

// Upgrade upgrades the Boundless Operator
func Upgrade(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig, operatorUri string) error {
	log.Info().Msgf("Upgrading Boundless Operator using manifest file %q", operatorUri)
	if err := k8s.ApplyYaml(kubeConfig, operatorUri); err != nil {
		return fmt.Errorf("failed to upgrade boundless operator: %w", err)
	}

	log.Info().Msgf("Finished updating Boundless Operator")
	return nil
}
