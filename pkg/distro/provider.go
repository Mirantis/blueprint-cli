package distro

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
)

// Provider is the interface for a distro provider
type Provider interface {
	Install() error
	Upgrade() error
	SetupClient() error
	Exists() (bool, error)
	Reset() error
	GetKubeConfigContext() string
	Type() string
	GetKubeConfig() *k8s.KubeConfig
	WaitForNodes() error
	WaitForPods() error
	NeedsUpgrade(blueprint *types.Blueprint) (bool, error)
	ValidateProviderUpgrade(blueprint *types.Blueprint) error
}

// GetProvider returns a new provider
func GetProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) (Provider, error) {
	switch blueprint.Spec.Kubernetes.Provider {
	case constants.ProviderK0s:
		return NewK0sProvider(blueprint, kubeConfig), nil
	case constants.ProviderKind:
		return NewKindProvider(blueprint, kubeConfig), nil
	case constants.ProviderExisting:
		return NewExistingProvider(blueprint, kubeConfig), nil
	default:
		return nil, fmt.Errorf("invalid kubernetes distribution provider: %s", blueprint.Spec.Kubernetes.Provider)
	}
}
