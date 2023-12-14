package boundlessclientset

import (
	"github.com/mirantis/boundless-operator/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type BoundlessV1Alpha1Interface interface {
	Addons(namespace string) v1alpha1.Addon
}

type BoundlessV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*BoundlessV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &BoundlessV1Alpha1Client{restClient: client}, nil
}

func (c *BoundlessV1Alpha1Client) Addons(namespace string) AddonInterface {
	return &addonClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *BoundlessV1Alpha1Client) Manifests(namespace string) ManifestInterface {
	return &manifestClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
