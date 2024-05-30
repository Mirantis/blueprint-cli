package distro

import (
	"fmt"
	"strings"

	"github.com/k0sproject/dig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/mirantiscontainers/boundless-cli/pkg/utils"
)

// Kind is the kind provider
type Kind struct {
	name       string
	kindConfig dig.Mapping
	kubeConfig *k8s.KubeConfig
	client     *kubernetes.Clientset
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

func (k *Kind) Upgrade() error {
	return nil
}

// SetupClient sets up the kubernets client for the distro
func (k *Kind) SetupClient() error {
	var err error
	k.client, err = k8s.GetClient(k.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}
	return k.WaitForNodes()
}

// Exists checks if kind exists
func (k *Kind) Exists() (bool, error) {
	err := utils.ExecCommandQuietly("bash", "-c", fmt.Sprintf("kind get clusters -q | grep -x %s", k.name))
	if err != nil && strings.Contains(err.Error(), "exit status 1") {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Reset deletes the kind cluster
func (k *Kind) Reset(force bool) error {
	log.Debug().Msgf("Resetting kind cluster %q", k.name)

	if err := utils.ExecCommand(fmt.Sprintf("kind delete clusters %s", k.name)); err != nil {
		return fmt.Errorf("failed to delete kind cluster: %w", err)
	}

	return nil
}

// GetKubeConfigContext returns the kubeconfig context
func (k *Kind) GetKubeConfigContext() string {
	return "kind-" + k.name
}

// Type returns the type of the provider
func (k *Kind) Type() string {
	return constants.ProviderKind
}

// GetKubeConfig returns the kubeconfig
func (k *Kind) GetKubeConfig() *k8s.KubeConfig {
	return k.kubeConfig
}

// WaitForPods waits for pods to be ready
func (k *Kind) WaitForPods() error {
	if err := k8s.WaitForPods(k.client, constants.NamespaceBlueprint); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}

	return nil
}

// WaitForNodes waits for nodes to be ready
func (k *Kind) WaitForNodes() error {
	if err := k8s.WaitForNodes(k.client); err != nil {
		return fmt.Errorf("failed to wait for nodes: %w", err)
	}

	return nil
}

// NeedsUpgrade returns false for Kind
func (k *Kind) NeedsUpgrade(blueprint *types.Blueprint) (bool, error) {
	return false, nil
}

// ValidateProviderUpgrade returns nil for Kind
func (k *Kind) ValidateProviderUpgrade(blueprint *types.Blueprint) error {
	return nil
}
