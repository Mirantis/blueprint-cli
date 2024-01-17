package distro

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/internal/utils"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"gopkg.in/yaml.v2"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
)

// K0s is the k0s provider
type K0s struct {
	name       string
	k0sConfig  string
	kubeConfig *k8s.KubeConfig
}

// NewK0sProvider returns a new k0s provider
func NewK0sProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) *K0s {
	provider := &K0s{
		name:       blueprint.Metadata.Name,
		kubeConfig: kubeConfig,
	}

	k0sConfig, err := createTempK0sConfig(blueprint)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get k0s config path")
	}
	provider.k0sConfig = k0sConfig

	return provider
}

// Install installs k0s using k0sctl
func (k *K0s) Install() error {
	kubeConfigPath := k.kubeConfig.GetConfigPath()
	log.Debug().Msgf("Creating k0s cluster %q with kubeConfig at: %s", k.name, kubeConfigPath)

	if err := utils.ExecCommand(fmt.Sprintf("k0sctl apply --config %s --no-wait", k.k0sConfig)); err != nil {
		return fmt.Errorf("failed to install k0s: %w", err)
	}

	// create kubeconfig
	if err := writeK0sKubeConfig(k.k0sConfig, k.kubeConfig); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}
	log.Trace().Msgf("kubeconfig file for k0s cluster: %s", kubeConfigPath)

	return nil
}

// Reset resets k0s using k0sctl
func (k *K0s) Reset() error {
	log.Debug().Msgf("Resetting k0s cluster: %s", k.name)

	if err := utils.ExecCommand(fmt.Sprintf("k0sctl reset --config %s", k.k0sConfig)); err != nil {
		return fmt.Errorf("failed to reset k0s: %w", err)
	}

	return nil
}

// GetKubeConfigContext returns the kubeconfig context for k0s
func (k *K0s) GetKubeConfigContext() string {
	return k.name
}

func writeK0sKubeConfig(k0sctlConfig string, kubeConfig *k8s.KubeConfig) error {
	c := exec.Command("k0sctl", "kubeconfig", "--config", k0sctlConfig)
	c.Stderr = os.Stderr

	buf := new(bytes.Buffer)
	c.Stdout = buf

	err := c.Run()
	if err != nil {
		return fmt.Errorf("failed to generate kubeconfig: %w", err)
	}

	configClient, err := clientcmd.NewClientConfigFromBytes(buf.Bytes())
	if err != nil {
		return err
	}

	rawConfig, err := configClient.RawConfig()
	if err != nil {
		return err
	}
	err = kubeConfig.MergeConfig(rawConfig)
	if err != nil {
		return err
	}

	return nil
}

// createTempK0sConfig creates a k0s config file from the blueprint in the tmp directory
func createTempK0sConfig(blueprint *types.Blueprint) (string, error) {
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

// writeToTempFile writes the k0sctl config to a tmp file and returns the path to it
func writeToTempFile(k0sctlConfig []byte) (string, error) {
	tmpfile, err := os.CreateTemp("", "k0sctl.yaml")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(k0sctlConfig)
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
