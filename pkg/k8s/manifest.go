package k8s

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mirantiscontainers/boundless-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// readYamlManifest reads a Kubernetes YAML manifest file containing multiple objects and returns the contents
// as array of unstructured objects. The order of the objects in the returned slice is the same as in the file.
// The uri argument can be a file path or a URL.
func readYamlManifest(uri string) ([]unstructured.Unstructured, error) {
	log.Debug().Msgf("Reading YAML manifest from %q", uri)
	b, err := utils.ReadURI(uri)
	if err != nil {
		return nil, err
	}

	return decodeObjects(b)
}

func decodeObjects(data []byte) ([]unstructured.Unstructured, error) {
	var objs []unstructured.Unstructured
	decoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(data))

	var o unstructured.Unstructured
	for {
		if err := decoder.Decode(&o); err != nil {
			if err != io.EOF {
				return objs, fmt.Errorf("error decoding yaml manifest file: %s", err)
			}
			break
		}
		objs = append(objs, o)

	}
	return objs, nil
}
