package boundlessclientset

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/mirantiscontainers/boundless-operator/api/v1alpha1"
)

// AddonInterface is an interface containing the operations that can be done on Addons
type AddonInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.AddonList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Addon, error)
	Create(addon *v1alpha1.Addon) (*v1alpha1.Addon, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type addonClient struct {
	restClient rest.Interface
	namespace  string
}

func (c *addonClient) List(opts metav1.ListOptions) (*v1alpha1.AddonList, error) {
	result := v1alpha1.AddonList{}
	err := c.restClient.
		Get().
		Namespace(c.namespace).
		Resource("addons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *addonClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Addon, error) {
	result := v1alpha1.Addon{}
	err := c.restClient.
		Get().
		Namespace(c.namespace).
		Resource("addons").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *addonClient) Create(addon *v1alpha1.Addon) (*v1alpha1.Addon, error) {
	result := v1alpha1.Addon{}
	err := c.restClient.
		Post().
		Namespace(c.namespace).
		Resource("addons").
		Body(addon).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *addonClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.namespace).
		Resource("addons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
