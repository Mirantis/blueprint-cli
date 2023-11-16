package k8s

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// Apply applies a kubernetes manifest
// TODO (ranyodh): use client-go instead of kubectl and remove kubectl dependency
func Apply(path string, kc *KubeConfig) error {
	log.Debug().Msgf("kubeconfig file: %s", kc.GetConfigPath())
	cmd := exec.Command("kubectl", "apply", "-f", path, "--kubeconfig", kc.GetConfigPath())
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
