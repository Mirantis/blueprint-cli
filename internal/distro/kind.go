package distro

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/internal/types"
	"github.com/mirantiscontainers/boundless-cli/internal/utils"
)

func InstallKind(name string, k8scfg *k8s.KubeConfig, configPath string) error {
	kubeconfig := k8scfg.GetConfigPath()
	log.Debug().Msgf("Creating kind cluster %q with kubeconfig at: %s", name, kubeconfig)

	if configPath != "" {
		if err := utils.ExecCommand("kind", "create", "cluster",
			"-n", name,
			"--kubeconfig", kubeconfig,
			fmt.Sprintf("--config=%s", configPath),
		); err != nil {
			return fmt.Errorf("failed to create kind cluster %w", err)
		}
	} else {
		if err := utils.ExecCommand("kind", "create", "cluster", "-n", name, "--kubeconfig", kubeconfig); err != nil {
			return fmt.Errorf("failed to create kind cluster %w", err)
		}
	}

	log.Debug().Msgf("kubeconfig file for kind cluster: %s", kubeconfig)
	return nil
}

func ResetKind(name string) error {
	log.Debug().Msgf("Deleting kind cluster %q", name)
	if err := utils.ExecCommand("kind", "delete", "clusters", name); err != nil {
		return fmt.Errorf("failed to delete kind cluster %w", err)
	}

	return nil
}

func GetKubeConfigContextKind(blueprint types.Blueprint) string {
	return "kind-" + blueprint.Metadata.Name
}
