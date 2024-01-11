package distro

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/internal/types"
	"github.com/mirantiscontainers/boundless-cli/internal/utils"
)

// Kind is the kind provider
type Kind struct {
	name       string
	kindConfig string
	kubeConfig *k8s.KubeConfig
}

// NewKindProvider returns a new kind provider
func NewKindProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) *Kind {
	provider := &Kind{
		name:       blueprint.Metadata.Name,
		kubeConfig: kubeConfig,
	}

	return provider
}

// Install creates a new kind cluster
func (k *Kind) Install() error {
	kubeConfigPath := k.kubeConfig.GetConfigPath()
	log.Debug().Msgf("Creating kind cluster %q with kubeConfig at: %s", k.name, kubeConfigPath)

	if err := utils.ExecCommand("kind", "create", "cluster", "-n", k.name, "--kubeconfig", kubeConfigPath); err != nil {
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}

	return nil
}

// Reset deletes the kind cluster
func (k *Kind) Reset() error {
	log.Debug().Msgf("Resetting kind cluster %q", k.name)

	if err := utils.ExecCommand("kind", "delete", "clusters", k.name); err != nil {
		return fmt.Errorf("failed to delete kind cluster: %w", err)
	}

	return nil
}

func (k *Kind) GetKubeConfigContext() string {
	return "kind-" + k.name
}
