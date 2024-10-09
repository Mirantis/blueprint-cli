package k8s

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"

	operatorv1alpha1 "github.com/mirantiscontainers/blueprint-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdate creates or updates a kubernetes object
func CreateOrUpdate(config *KubeConfig, obj client.Object) error {
	// TODO (ranyodh): This is currently using in-cluster client. We should switch to:
	// - either a dynamic client,
	// - or generate a client in the `blueprint-operator` to be used here
	scheme := runtime.NewScheme()
	_ = operatorv1alpha1.AddToScheme(scheme)

	restConfig, err := config.RESTConfig()
	if err != nil {
		return err
	}
	restConfig.WarningHandler = rest.NoWarnings{}
	kubeClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	existing := &operatorv1alpha1.Blueprint{}
	err = kubeClient.Get(context.Background(), client.ObjectKeyFromObject(obj), existing)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return fmt.Errorf("failed to get existing blueprints: %v", err)
		}
	}
	if existing.Name != "" {
		obj.SetResourceVersion(existing.GetResourceVersion())
		// @TODO add support for .Patch() or merge exisiting and obj.
		obj.SetFinalizers(existing.GetFinalizers())

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
