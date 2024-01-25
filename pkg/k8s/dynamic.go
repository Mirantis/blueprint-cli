package k8s

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// ApplyYaml applies a yaml manifest to the cluster from the URI. The URI can be a file path or a URL
// It creates CRDs first and then other objects
// @TODO: Make this function testable by passing a "uri reader" and kubernetes clients
func ApplyYaml(kc *KubeConfig, uri string) error {
	var err error
	var client kubernetes.Interface
	var dynamicClient dynamic.Interface

	if client, err = GetClient(kc); err != nil {
		return fmt.Errorf("failed to get kubernetes client: %q", err)
	}
	if dynamicClient, err = GetDynamicClient(kc); err != nil {
		return fmt.Errorf("failed to get kubernetes dynamic client: %q", err)
	}

	objs, err := readYamlManifest(uri)
	if err != nil {
		return fmt.Errorf("failed to read manifest from %q: %w", uri, err)
	}

	// separate out the CRDs and other objects
	// CRDs need to be created first
	crds, others := splitCrdAndOthers(objs)
	log.Trace().Msgf("Found %d CRDs and %d other objects", len(crds), len(others))

	ctx := context.Background()
	for _, o := range crds {
		if err = createOrUpdateObject(ctx, client, dynamicClient, &o); err != nil {
			return fmt.Errorf("failed to apply crds resources from manifest at %q: %w", uri, err)
		}
	}

	// create other objects
	for _, o := range others {
		if err = createOrUpdateObject(ctx, client, dynamicClient, &o); err != nil {
			return fmt.Errorf("failed to apply resources from manifest at %q: %w", uri, err)
		}
	}
	return nil
}

// DeleteYamlObjects deletes all objects in the cluster that are specified in the yaml
func DeleteYamlObjects(kc *KubeConfig, uri string) error {
	var err error
	var client kubernetes.Interface
	var dynamicClient dynamic.Interface

	if client, err = GetClient(kc); err != nil {
		return fmt.Errorf("failed to get kubernetes client: %q", err)
	}
	if dynamicClient, err = GetDynamicClient(kc); err != nil {
		return fmt.Errorf("failed to get kubernetes dynamic client: %q", err)
	}

	objs, err := readYamlManifest(uri)
	if err != nil {
		return fmt.Errorf("failed to read manifest from %q: %w", uri, err)
	}

	log.Info().Msgf("Deleting %d objects", len(objs))
	ctx := context.Background()
	for _, o := range objs {
		if err = deleteObject(ctx, client, dynamicClient, &o); err != nil {
			return fmt.Errorf("failed to reset obj resources from manifest at %q: %w", uri, err)
		}
	}

	return nil
}

func splitCrdAndOthers(objs []unstructured.Unstructured) ([]unstructured.Unstructured, []unstructured.Unstructured) {
	var crds []unstructured.Unstructured
	var others []unstructured.Unstructured
	for _, o := range objs {
		if o.GetKind() == "CustomResourceDefinition" {
			crds = append(crds, o)
		} else {
			others = append(others, o)
		}
	}
	return crds, others
}

func createOrUpdateObject(ctx context.Context, client kubernetes.Interface, dynamicClient dynamic.Interface, obj *unstructured.Unstructured) error {
	gvr, _ := getResource(client, obj)
	namespace := obj.GetNamespace()
	objName := obj.GetName()

	existing, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, objName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Trace().Msgf("Creating %q of kind %q", objName, obj.GetKind())
		_, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create resource: %q", err)
		}
		log.Trace().Msgf("Created %q of kind %q", objName, obj.GetKind())
	} else {
		log.Trace().Msgf("Updating %q of kind %q", objName, obj.GetKind())
		obj.SetResourceVersion(existing.GetResourceVersion())
		_, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(ctx, obj, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update resource: %q", err)
		}
		log.Trace().Msgf("Updated %q of kind %q", objName, obj.GetKind())
	}

	return nil
}

func deleteObject(ctx context.Context, client kubernetes.Interface, dynamicClient dynamic.Interface, obj *unstructured.Unstructured) error {
	gvr, _ := getResource(client, obj)
	namespace := obj.GetNamespace()
	objName := obj.GetName()

	_, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, objName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Trace().Msgf("%q was not found. No changes made", objName)
		return nil
	}

	log.Trace().Msgf("Deleting %q of kind %q", objName, obj.GetKind())
	err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete resource: %q", err)
	}
	log.Trace().Msgf("Deleted %q of kind %q", objName, obj.GetKind())

	return nil
}

// getResource returns the GroupVersionResource for a given object
// Especially, it discovers the resource name for the given object
func getResource(client kubernetes.Interface, object *unstructured.Unstructured) (schema.GroupVersionResource, error) {
	gvk := object.GroupVersionKind()
	apiResourceList, err := client.Discovery().ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	var resource *metav1.APIResource
	for _, r := range apiResourceList.APIResources {
		if r.Kind == gvk.Kind {
			resource = &r
			break
		}
	}

	log.Trace().Msgf("Found resource: %q", resource.Name)
	return schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: resource.Name}, nil
}
