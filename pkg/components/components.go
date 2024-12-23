package components

import (
	"bytes"
	"fmt"
	"os"

	"github.com/k0sproject/dig"
	"github.com/rs/zerolog/log"
	yamlDecoder "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sYaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"

	"github.com/mirantiscontainers/blueprint-operator/api/v1alpha1"

	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
	"github.com/mirantiscontainers/blueprint-cli/pkg/k8s"
	"github.com/mirantiscontainers/blueprint-cli/pkg/types"
)

// ApplyBlueprint applies a Blueprint object to the cluster
func ApplyBlueprint(kubeConfig *k8s.KubeConfig, cluster *types.Blueprint) error {
	components := cluster.Spec.Components

	// Get the list of addons
	addons, err := getAddons(&components)
	if err != nil {
		return err
	}

	c := v1alpha1.Blueprint{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Metadata.Name,
			Namespace: v1.NamespaceDefault,
		},
		Spec: v1alpha1.BlueprintSpec{
			Components: v1alpha1.Component{
				Addons: addons,
			},
			Resources: getResources(cluster.Spec.Resources),
		},
	}

	log.Info().Msg("Applying Blueprint")
	if err := k8s.CreateOrUpdate(kubeConfig, &c); err != nil {
		return fmt.Errorf("failed to create/update Blueprint object: %v", err)
	}

	return nil
}

// RemoveComponents removes all components from the cluster
func RemoveComponents(kubeConfig *k8s.KubeConfig, cluster *types.Blueprint) error {
	components := cluster.Spec.Components

	// Get the list of addons
	addons, err := getAddons(&components)
	if err != nil {
		return err
	}

	c := v1alpha1.Blueprint{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Metadata.Name,
			Namespace: v1.NamespaceDefault,
		},
		Spec: v1alpha1.BlueprintSpec{
			Components: v1alpha1.Component{
				Addons: addons,
			},
			Resources: getResources(cluster.Spec.Resources),
		},
	}

	log.Info().Msg("Resetting Blueprint")
	if err := k8s.Delete(kubeConfig, &c); err != nil {
		return fmt.Errorf("failed to reset Blueprint object: %v", err)
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

func getAddons(components *types.Components) ([]v1alpha1.AddonSpec, error) {
	var addons []v1alpha1.AddonSpec

	for _, addon := range components.Addons {
		if addon.Kind == constants.AddonChart {
			addons = append(addons, v1alpha1.AddonSpec{
				Name:      addon.Name,
				Kind:      addon.Kind,
				Enabled:   addon.Enabled,
				DryRun:    addon.DryRun,
				Namespace: addon.Namespace,
				Chart: &v1alpha1.ChartInfo{
					Name:    addon.Chart.Name,
					Repo:    addon.Chart.Repo,
					Version: addon.Chart.Version,
					Set:     addon.Chart.Set,
					Values:  addon.Chart.Values,
				},
			})
		} else if addon.Kind == constants.AddonManifest {
			addons = append(addons, v1alpha1.AddonSpec{
				Name:      addon.Name,
				Kind:      addon.Kind,
				Enabled:   addon.Enabled,
				Namespace: addon.Namespace,
				Manifest: &v1alpha1.ManifestInfo{
					URL:           addon.Manifest.URL,
					FailurePolicy: addon.Manifest.FailurePolicy,
					Timeout:       addon.Manifest.Timeout,
					Values:        addon.Manifest.Values,
				},
			})
		} else {
			return nil, fmt.Errorf("unknown addon kind %q (valid values: %s|%s)", addon.Kind, constants.AddonChart, constants.AddonManifest)
		}
	}

	return addons, nil
}

func jsonValues(values dig.Mapping) (string, error) {
	valuesYaml := new(bytes.Buffer)

	encoder := yamlDecoder.NewEncoder(valuesYaml)
	err := encoder.Encode(&values)
	if err != nil {
		return "", err
	}

	json, err := k8sYaml.ToJSON(valuesYaml.Bytes())
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func Encode(blueprint types.Blueprint) error {
	encoder := yamlDecoder.NewEncoder(os.Stdout)
	return encoder.Encode(&blueprint)
}

var valueString = `service:
type: ClusterIP
`

var DefaultComponents = types.Components{
	Addons: []types.Addon{
		{
			Name:      "example-server",
			Kind:      constants.AddonChart,
			Enabled:   true,
			Namespace: "default",
			Chart: &types.ChartInfo{
				Name:    "nginx",
				Repo:    "https://charts.bitnami.com/bitnami",
				Version: "15.1.1",
				Values:  ConvertValues([]byte(valueString)),
			},
		},
	},
}

func getResources(resources *types.Resources) v1alpha1.Resources {
	if resources == nil {
		return v1alpha1.Resources{}
	}

	return v1alpha1.Resources{
		CertManagement: resources.CertManagement.CertManagement,
	}
}

// ConvertValues converts a byte slice to a JSON object
func ConvertValues(in []byte) *apiextensionsv1.JSON {
	var values *apiextensionsv1.JSON
	if in != nil {
		v, _ := yaml.YAMLToJSON(in)
		values = &apiextensionsv1.JSON{Raw: v}
	}

	return values
}
