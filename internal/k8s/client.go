package k8s

import (
	"k8s.io/client-go/kubernetes"
)

// GetClient returns a handle to api server or die.
func GetClient(config *KubeConfig) (kubernetes.Interface, error) {
	cfg, err := config.RESTConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(cfg)
}
