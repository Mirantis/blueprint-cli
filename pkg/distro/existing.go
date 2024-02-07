package distro

import (
	"fmt"
	"strings"

	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/mirantiscontainers/boundless-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
)

// Existing is the existing provider
type Existing struct {
	kubeConfig *k8s.KubeConfig
	client     *kubernetes.Clientset
}

// NewExistingProvider returns a new existing provider
func NewExistingProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) *Existing {
	return &Existing{
		kubeConfig: kubeConfig,
	}
}

// SetupClient sets up the kubernets client for the distro
func (e *Existing) SetupClient() error {
	var err error
	e.client, err = k8s.GetClient(e.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}
	return e.WaitForNodes()
}

// WaitForNodes waits for nodes to be ready
func (e *Existing) WaitForNodes() error {
	if err := k8s.WaitForNodes(e.client); err != nil {
		return fmt.Errorf("failed to wait for nodes: %w", err)
	}

	return nil
}

// WaitForPods waits for pods to be ready
func (e *Existing) WaitForPods() error {
	if err := k8s.WaitForPods(e.client, constants.NamespaceBoundless); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}

	return nil
}

// Install installs the existing cluster
func (e *Existing) Install() error {
	log.Debug().Msgf("Nothing done to install an unsupported existing cluster")
	return nil
}

// Update updates the existing cluster
func (e *Existing) Upgrade() error {
	log.Debug().Msgf("Nothing done to update an unsupported existing cluster")
	return nil
}

// Exists checks if the cluster exists
func (e *Existing) Exists() (bool, error) {
	config, err := e.kubeConfig.RESTConfig()
	if err != nil {
		return false, err
	}

	// This checks if the cluster exists but doesn't use authentication
	err = utils.ExecCommandQuietly("bash", "-c", fmt.Sprintf("curl -k %s/livesz/verbose", config.Host))
	// Exists but we have no authentication
	if err != nil && strings.Contains(err.Error(), "exit status 6") {
		return true, nil
	}
	// Can't be reached/doesn't exist
	if err != nil && strings.Contains(err.Error(), "exit status 7") {
		return false, nil
	}
	// Some other error
	if err != nil {
		return false, err
	}

	return true, nil
}

// Reset resets the existing cluster
func (e *Existing) Reset() error {
	log.Debug().Msgf("Nothing done to reset an unsupported existing cluster")
	return nil
}

// GetKubeConfigContext returns the kubeconfig context
func (e *Existing) GetKubeConfigContext() string {
	return ""
}

// Type returns the type of the provider
func (e *Existing) Type() string {
	return constants.ProviderExisting
}

// GetKubeConfig returns the kubeconfig
func (e *Existing) GetKubeConfig() *k8s.KubeConfig {
	return e.kubeConfig
}

// NeedsUpgrade returns false for existing cluster
func (e *Existing) NeedsUpgrade(blueprint *types.Blueprint) (bool, error) {
	return false, nil
}

// ValidateProviderUpgrade returns nil for existing cluster
func (e *Existing) ValidateProviderUpgrade(blueprint *types.Blueprint) error {
	return nil
}
