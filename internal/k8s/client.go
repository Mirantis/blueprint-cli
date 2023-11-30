package k8s

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// GetClient returns a handle to api server.
func GetClient(config *KubeConfig) (*kubernetes.Clientset, error) {
	cfg, err := config.RESTConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get REST config for kube client: %v", err)
	}

	return kubernetes.NewForConfig(cfg)
}

func GetDynamicClient(config *KubeConfig) (*dynamic.DynamicClient, error) {
	cfg, err := config.RESTConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get REST config for dynaminc kube client: %v", err)
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate dynamic client: %v", err)
	}

	return dynamicClient, err
}
