package k0sctl

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/mirantiscontainers/boundless-cli/internal/types"
)

// GetConfigPath writes the k0sctl config file to a temporary file and returns the path to it
func GetConfigPath(blueprint *types.Blueprint) (string, error) {
	k0sctlConfig := types.ConvertToK0s(blueprint)

	data, err := yaml.Marshal(k0sctlConfig)
	if err != nil {
		return "", err
	}

	k0sctlConfigFile, err := writeToTempFile(data)
	if err != nil {
		return "", err
	}

	return k0sctlConfigFile, nil
}

func writeToTempFile(data []byte) (string, error) {
	// create tmp file for k0sctl config file
	tmpfile, err := os.CreateTemp("", "k0sctl.yaml")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(data)
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
