package k8s

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"

	operatorv1alpha1 "github.com/mirantiscontainers/boundless-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Delete deletes a kubernetes object
func Delete(config *KubeConfig, obj client.Object) error {
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
		obj.SetFinalizers(existing.GetFinalizers())

		err = kubeClient.Delete(context.Background(), obj)
		if err != nil {
			return fmt.Errorf("failed to delete cluster object: %v", err)
		}
	}

	return nil
}
