package distro

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
)

const (
	// ProviderK0s is the name of the k0s provider
	ProviderK0s = "k0s"
	// ProviderKind is the name of the kind provider
	ProviderKind = "kind"
	// ProviderExisting is the name of the existing provider
	ProviderExisting = "existing"
)

// Provider is the interface for a distro provider
type Provider interface {
	Install() error
	Reset() error
	GetKubeConfigContext() string
}

// GetProvider returns a new provider
func GetProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) (Provider, error) {
	switch blueprint.Spec.Kubernetes.Provider {
	case ProviderK0s:
		return NewK0sProvider(blueprint, kubeConfig), nil
	case ProviderKind:
		return NewKindProvider(blueprint, kubeConfig), nil
	default:
		return nil, fmt.Errorf("invalid kubernetes distribution provider: %s", blueprint.Spec.Kubernetes.Provider)
	}
}
