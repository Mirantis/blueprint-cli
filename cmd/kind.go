package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func installKindCluster(name string, kubeconfig string) error {
	log.Debugf("creating kind cluster %q with kubeconfig at: %s", name, kubeconfig)
	if err := cmdRun("kind", "create", "cluster", "-n", name, "--kubeconfig", kubeconfig); err != nil {
		return fmt.Errorf("failed to create kind cluster %w", err)
	}

	log.Debugf("kubeconfig file for kind cluster: %s", kubeconfig)
	return nil
}
