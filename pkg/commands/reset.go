package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"

	"github.com/mirantiscontainers/blueprint-cli/pkg/components"
	"github.com/mirantiscontainers/blueprint-cli/pkg/distro"
	"github.com/mirantiscontainers/blueprint-cli/pkg/k8s"
	"github.com/mirantiscontainers/blueprint-cli/pkg/types"
)

// Reset resets the cluster
func Reset(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig, force bool) error {
	log.Info().Msg("Resetting cluster")

	if !force {
		color.Red("This will remove all resources and completely destroy the cluster. Are you sure? (N/y)")
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		if answer != "y\n" {
			return nil
		}
	}

	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	// Uninstall components
	log.Info().Msgf("Reset Blueprint Operator resources")
	err = components.RemoveComponents(provider.GetKubeConfig(), blueprint)
	if err != nil {
		return fmt.Errorf("failed to reset components: %w", err)
	}

	uri, err := determineOperatorUri(blueprint.Spec.Version)
	if err != nil {
		return fmt.Errorf("failed to determine operator URI: %w", err)
	}

	log.Info().Msgf("Uninstalling Blueprint Operator")
	log.Debug().Msgf("Uninstalling blueprint operator using manifest file: %s", uri)
	if err = k8s.DeleteYamlObjects(kubeConfig, uri); err != nil {
		return fmt.Errorf("failed to uninstall Blueprint Operator: %w", err)
	}

	// Reset the cluster
	if err := provider.Reset(); err != nil {
		return fmt.Errorf("failed to reset cluster: %w", err)
	}

	return nil
}
