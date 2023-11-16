package k8s

import (
	"context"
	"fmt"

	"github.com/mirantis/boundless-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// CreateOrUpdate creates or updates a kubernetes object
func CreateOrUpdate(config *KubeConfig, obj client.Object) error {
	// TODO (ranyodh): This is a temporary hack to use the kube client
	// Should be replaced with our own logger
	logf.SetLogger(zap.New())

	scheme := runtime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)

	restConfig, err := config.RESTConfig()
	if err != nil {
		return err
	}
	kubeClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	existing := &v1alpha1.Blueprint{}
	err = kubeClient.Get(context.Background(), client.ObjectKeyFromObject(obj), existing)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return fmt.Errorf("failed to get cluster object: %v", err)
		}
	}
	if existing.Name != "" {
		obj.SetResourceVersion(existing.GetResourceVersion())
		err = kubeClient.Update(context.Background(), obj)
		if err != nil {
			return fmt.Errorf("failed to update cluster object: %v", err)
		}
	} else {
		if err := kubeClient.Create(context.Background(), obj); err != nil {
			return fmt.Errorf("failed to create cluster object: %v", err)
		}
	}

	return nil
}
