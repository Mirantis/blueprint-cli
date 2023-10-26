package cmd

import (
	"bytes"
	"fmt"

	"github.com/k0sproject/dig"
	log "github.com/sirupsen/logrus"
	yamlDecoder "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/mirantis/boundless-operator/api/v1alpha1"

	"boundless-cli/pkg/kube"

	"boundless-cli/pkg/config"
)

func installComponents(components config.Components) error {
	ingressConfig, err := yamlValues(components.Core.Ingress.Config)
	if err != nil {
		return fmt.Errorf("failed to convert ingress config to yaml: %w", err)
	}

	var addons []v1alpha1.AddonSpec
	for _, addon := range components.Addons {
		addons = append(addons, v1alpha1.AddonSpec{
			Name:      addon.Name,
			Kind:      addon.Kind,
			Enabled:   addon.Enabled,
			Namespace: addon.Namespace,
			Chart: v1alpha1.Chart{
				Name:    addon.Chart.Name,
				Repo:    addon.Chart.Repo,
				Version: addon.Chart.Version,
				Set:     addon.Chart.Set,
				Values:  addon.Chart.Values,
			},
		})
	}

	c := v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mkecluster",
			Namespace: v1.NamespaceDefault,
		},
		Spec: v1alpha1.ClusterSpec{
			Components: v1alpha1.Component{
				Core: v1alpha1.Core{
					Ingress: v1alpha1.IngressSpec{
						Enabled:  components.Core.Ingress.Enabled,
						Provider: components.Core.Ingress.Provider,
						Config:   ingressConfig,
					},
				},
				Addons: addons,
			},
		},
	}

	log.Info("Creating/Updating MKE cluster")
	if err := kube.CreateOrUpdate(&c); err != nil {
		return fmt.Errorf("failed to create/update mke cluster object: %v", err)
	}

	return nil
}

func yamlValues(values dig.Mapping) (string, error) {
	valuesYaml := new(bytes.Buffer)

	encoder := yamlDecoder.NewEncoder(valuesYaml)
	err := encoder.Encode(&values)
	if err != nil {
		return "", err
	}
	return valuesYaml.String(), nil
}

func jsonValues(values dig.Mapping) (string, error) {
	valuesYaml := new(bytes.Buffer)

	encoder := yamlDecoder.NewEncoder(valuesYaml)
	err := encoder.Encode(&values)
	if err != nil {
		return "", err
	}

	json, err := yaml.ToJSON(valuesYaml.Bytes())
	if err != nil {
		return "", err
	}
	return string(json), nil
}