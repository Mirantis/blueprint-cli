package distro

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
	"github.com/mirantiscontainers/boundless-cli/pkg/utils"

	"github.com/k0sproject/version"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K0s is the k0s provider
type K0s struct {
	name       string
	k0sConfig  string
	kubeConfig *k8s.KubeConfig
	client     *kubernetes.Clientset
}

// NewK0sProvider returns a new k0s provider
func NewK0sProvider(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) *K0s {
	provider := &K0s{
		name:       blueprint.Metadata.Name,
		kubeConfig: kubeConfig,
	}

	k0sConfig, err := CreateTempK0sConfig(blueprint)
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
	if err := WriteK0sKubeConfig(k.k0sConfig, k.kubeConfig); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}
	log.Trace().Msgf("kubeconfig file for k0s cluster: %s", kubeConfigPath)

	return nil
}

// Refresh reapplies k0sctl config with the current version of blueprint
func (k *K0s) Refresh() error {
	kubeConfigPath := k.kubeConfig.GetConfigPath()
	log.Debug().Msgf("Refreshing k0s cluster %q with kubeConfig at: %s", k.name, kubeConfigPath)

	if err := utils.ExecCommand(fmt.Sprintf("k0sctl apply --config %s --no-wait", k.k0sConfig)); err != nil {
		return fmt.Errorf("k0sctl apply failed: %w", err)
	}

	return nil
}

// Update updates k0s using k0sctl
func (k *K0s) Upgrade() error {
	if err := utils.ExecCommand(fmt.Sprintf("k0sctl apply --config %s --no-wait", k.k0sConfig)); err != nil {

		return fmt.Errorf("failed to update k0s: %w", err)
	}
	return nil
}

// SetupClient sets up the kubernets client for the distro
func (k *K0s) SetupClient() error {
	var err error
	k.client, err = k8s.GetClient(k.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}
	return k.WaitForNodes()
}

// Exists checks if k0s exists using k0sctl
func (k *K0s) Exists() (bool, error) {
	err := utils.ExecCommandQuietly("bash", "-c", fmt.Sprintf("k0sctl kubeconfig -c %s", k.k0sConfig))
	if err != nil && strings.Contains(err.Error(), "exit status 1") {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Reset resets k0s using k0sctl
func (k *K0s) Reset(force bool) error {
	log.Debug().Msgf("Resetting k0s cluster: %s", k.name)

	resetCmd := fmt.Sprintf("k0sctl reset --config %s", k.k0sConfig)
	if force {
		resetCmd = fmt.Sprintf("%s --force", resetCmd)
	}

	if err := utils.ExecCommand(resetCmd); err != nil {
		return fmt.Errorf("failed to reset k0s: %w", err)
	}

	return nil
}

// GetKubeConfigContext returns the kubeconfig context for k0s
func (k *K0s) GetKubeConfigContext() string {
	return k.name
}

// Type returns the type of the provider
func (k *K0s) Type() string {
	return constants.ProviderK0s
}

// GetKubeConfig returns the kubeconfig
func (k *K0s) GetKubeConfig() *k8s.KubeConfig {
	return k.kubeConfig
}

// WaitForNodes waits for nodes to be ready
func (k *K0s) WaitForNodes() error {
	if err := k8s.WaitForNodes(k.client); err != nil {
		return fmt.Errorf("failed to wait for nodes: %w", err)
	}

	return nil
}

// WaitForPods waits for pods to be ready
func (k *K0s) WaitForPods() error {
	if err := k8s.WaitForPods(k.client, constants.NamespaceBlueprint); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}

	return nil
}

func WriteK0sKubeConfig(k0sctlConfig string, kubeConfig *k8s.KubeConfig) error {
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

// CreateTempK0sConfig creates a k0s config file from the blueprint in the tmp directory
func CreateTempK0sConfig(blueprint *types.Blueprint) (string, error) {
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

// getInstalledVersion returns version of k0s
// for local k0s it will get the k0s version of the local machine
// otherwise it will get k0s version on the first controller node that does not throw an error
func (k *K0s) getInstalledVersion(blueprint *types.Blueprint) (string, error) {

	if k.isLocalK0s(blueprint) {
		out, err := utils.ExecCommandWithReturn("sudo k0s version")
		if err != nil {
			return "", fmt.Errorf("unable to get k0s version on local host : %w", err)
		}

		return out, nil
	}

	controllers := k.getControllerHosts(blueprint)

	for _, controller := range controllers {
		key, err := utils.ReadFile(controller.SSH.KeyPath)
		if err != nil {
			return "", err
		}

		// k0sctl has no apparent way to get version of k0s previously installed so get the k0s version directly on the first controller node
		stdout, stderr, err := utils.RemoteCommand(controller.SSH.User, controller.SSH.Address, string(key), "sudo k0s version")
		if err != nil {
			log.Warn().Msgf("unable to get k0s version on host %s : %s, %s", controller.SSH.Address, stderr, err)

			// try to get version from another controller
			continue
		}

		return stdout, nil

	}

	// if we got here all controllers have errors when getting version
	return "", fmt.Errorf("unable to get k0s version of cluster")
}

func (k *K0s) getControllerHosts(blueprint *types.Blueprint) []types.Host {
	var hosts []types.Host

	for _, host := range blueprint.Spec.Kubernetes.Infra.Hosts {
		// match both controller and controller+worker role
		if strings.Contains(host.Role, "controller") {
			hosts = append(hosts, host)
			break
		}
	}
	return hosts
}

func (k *K0s) isLocalK0s(blueprint *types.Blueprint) bool {
	// if running localhost there should just be 1 host
	return len(blueprint.Spec.Kubernetes.Infra.Hosts) > 0 && blueprint.Spec.Kubernetes.Infra.Hosts[0].LocalHost.Enabled
}

// NeedsUpgrade checks if an upgrade of the provider is required
// return true if the providedVersion is greater than the installed Version
// return false if the versions are equal
// throw an error if the providedVersion is lower than the installed Version (don't support downgrade)
func (k *K0s) NeedsUpgrade(blueprint *types.Blueprint) (bool, error) {
	installedVersion, err := k.getInstalledVersion(blueprint)
	if err != nil {
		return false, fmt.Errorf("failed to get installed k0s version: %w", err)
	}

	installed := version.MustParse(installedVersion)
	provided := version.MustParse(blueprint.Spec.Kubernetes.Version)

	if provided.GreaterThan(installed) {
		return true, nil
	}
	if installed.GreaterThan(provided) {
		return false, fmt.Errorf("downgrade version detected - cannot downgrade provider versions")
	}

	return false, nil
}

// ValidateProviderUpgrade does some validation that the controller nodes will be able to run the new version of k0s proposed in the blueprint
// First download new version of k0s binary and place in tmp folder
// Use new binary to run k0s config validate which validates config will work on new version
// In some k0s upgrade scenarios (such as previously existing config fields that have been removed in newer version) it requires user to update the node configs
func (k *K0s) ValidateProviderUpgrade(blueprint *types.Blueprint) error {
	controllers := k.getControllerHosts(blueprint)

	defer func() {
		// cleanup the temp k0s binaries used to validate each controller
		for _, controller := range controllers {
			key, err := utils.ReadFile(controller.SSH.KeyPath)
			if err != nil {
				log.Warn().Msgf("failed to read ssh key during cleanup for host %s", controller.SSH.Address)
			}

			_, cleanupErr, err := utils.RemoteCommand(controller.SSH.User, controller.SSH.Address, string(key), "sudo rm -f /tmp/k0s")
			if err != nil {
				if !errors.IsNotFound(err) {
					log.Warn().Msgf("failed to clean up temp k0s binary for host %s : %s", controller.SSH.Address, cleanupErr)
				}
			}
		}
	}()

	for _, controller := range controllers {
		key, err := utils.ReadFile(controller.SSH.KeyPath)
		if err != nil {
			return err
		}

		log.Info().Msg("Downloading new version of k0s binary")
		downloadCmd := fmt.Sprintf("curl -sSLf https://get.k0s.sh | sed -e 's;k0sInstallPath=/usr/local/bin;k0sInstallPath=/tmp;' | sudo K0S_VERSION=v%s sh", blueprint.Spec.Kubernetes.Version)

		_, downloadErr, err := utils.RemoteCommand(controller.SSH.User, controller.SSH.Address, string(key), downloadCmd)
		if err != nil {
			log.Error().Msgf("failed to install new version of k0s binary on host %s : %s", controller.SSH.Address, downloadErr)
			return err
		}

		log.Info().Msg("Validating existing config with new version of k0s binary")
		validateCmd := fmt.Sprintf("sudo /tmp/k0s config validate --config /etc/k0s/k0s.yaml")
		_, validateErr, err := utils.RemoteCommand(controller.SSH.User, controller.SSH.Address, string(key), validateCmd)
		if err != nil {
			log.Error().Msgf("validation of new provider version failed on host %s : %s", controller.SSH.Address, validateErr)
			return err
		}
	}

	log.Info().Msg("New provider version successfully validated")
	return nil

}
