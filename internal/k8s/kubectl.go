package k8s

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// Apply applies a kubernetes manifest
// TODO (ranyodh): use client-go instead of kubectl and remove kubectl dependency
func Apply(path string, kc *KubeConfig) error {

	contextName, err := kc.CurrentContextName()
	if err != nil {
		return fmt.Errorf("failed to get current context name: %v", err)
	}

	log.Debug().Msgf("kubeconfig file: %q with context : %q", kc.GetConfigPath(), contextName)
	cmd := exec.Command("kubectl", "apply", "-f", path, "--kubeconfig", kc.GetConfigPath(), "--context", contextName)
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
