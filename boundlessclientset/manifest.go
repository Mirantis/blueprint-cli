package boundlessclientset

import (
	"context"
	"github.com/mirantis/boundless-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ManifestInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ManifestList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Manifest, error)
	Create(addon *v1alpha1.Manifest) (*v1alpha1.Manifest, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type manifestClient struct {
	restClient rest.Interface
	ns         string
}

func (c *manifestClient) List(opts metav1.ListOptions) (*v1alpha1.ManifestList, error) {
	result := v1alpha1.ManifestList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("manifests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *manifestClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Manifest, error) {
	result := v1alpha1.Manifest{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("manifests").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *manifestClient) Create(addon *v1alpha1.Manifest) (*v1alpha1.Manifest, error) {
	result := v1alpha1.Manifest{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("manifests").
		Body(addon).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *manifestClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("manifests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
