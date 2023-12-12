package distro

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/internal/types"
	"github.com/mirantiscontainers/boundless-cli/internal/utils"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
)

// InstallK0s installs k0s using k0sctl
func InstallK0s(k0sConfig string, kubeConfig *k8s.KubeConfig) error {
	log.Debug().Msgf("Installing k0s with config: %q", k0sConfig)
	if err := utils.ExecCommand("k0sctl", "apply", "--config", k0sConfig, "--no-wait"); err != nil {
		return fmt.Errorf("failed to install k0s: %w", err)
	}

	// create kubeconfig
	if err := writeK0sKubeConfig(k0sConfig, kubeConfig); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}
	log.Trace().Msgf("kubeconfig file for k0s cluster: %s", kubeConfig.GetConfigPath())

	return nil
}

// ResetK0s resets k0s using k0sctl
func ResetK0s(k0sConfig string) error {
	log.Debug().Msgf("Resetting k0s with config: %q", k0sConfig)
	if err := utils.ExecCommand("k0sctl", "reset", "--config", k0sConfig); err != nil {
		return fmt.Errorf("failed to reset k0s: %w", err)
	}

	return nil
}

func GetKubeConfigContextK0s(blueprint types.Blueprint) string {
	return blueprint.Metadata.Name
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
