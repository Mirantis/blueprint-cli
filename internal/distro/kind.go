package distro

import (
	"fmt"

	"github.com/k0sproject/dig"
	"github.com/rs/zerolog/log"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/internal/utils"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
)

// Kind is the kind provider
type Kind struct {
	name       string
	kindConfig dig.Mapping
	kubeConfig *k8s.KubeConfig
}

// NewKindProvider returns a new kind provider
func NewKindProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) *Kind {
	provider := &Kind{
		name:       blueprint.Metadata.Name,
		kubeConfig: kubeConfig,
		kindConfig: blueprint.Spec.Kubernetes.Config,
	}

	return provider
}

// Install creates a new kind cluster
func (k *Kind) Install() error {
	kubeConfigPath := k.kubeConfig.GetConfigPath()
	log.Debug().Msgf("Creating kind cluster %q with kubeConfig at: %s", k.name, kubeConfigPath)

	// Setup the kind create command
	command := fmt.Sprintf("kind create cluster -n %s", k.name)
	if k.kindConfig != nil {
		// Create the tmp kind config
		kindConfigYaml, err := yaml.Marshal(k.kindConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal kind config: %w", err)
		}

		// /dev/stdin is used to pass the config to kind without creating a file
		command = fmt.Sprintf("echo '%s' | %s --config /dev/stdin", kindConfigYaml, command)
	}

	if err := utils.ExecCommand(command); err != nil {
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}

	return nil
}

// Reset deletes the kind cluster
func (k *Kind) Reset() error {
	log.Debug().Msgf("Resetting kind cluster %q", k.name)

	if err := utils.ExecCommand(fmt.Sprintf("kind delete clusters %s", k.name)); err != nil {
		return fmt.Errorf("failed to delete kind cluster: %w", err)
	}

	return nil
}

func (k *Kind) GetKubeConfigContext() string {
	return "kind-" + k.name
}
