package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"boundless-cli/internal/boundless"
	"boundless-cli/internal/distro"
	"boundless-cli/internal/k0sctl"
	"boundless-cli/internal/k8s"
)

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply the blueprint to the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply()
		},
	}

	flags := cmd.Flags()
	addConfigFlags(flags)
	addCustomBOPFlag(flags)
	addKubeFlags(flags)
	cmd.Flags().MarkHidden("bop")

	return cmd
}

func runApply() error {
	var err error

	if blueprint.Spec.Kubernetes != nil {
		log.Info().Msgf("Installing Kubernetes distribution: %s", blueprint.Spec.Kubernetes.Provider)

		// TODO (ranyodh): Refactor the follow to use provider interface
		switch blueprint.Spec.Kubernetes.Provider {
		case distro.ProviderK0s:
			k0sctlConfigPath, err := k0sctl.GetConfigPath(blueprint)
			if err = distro.InstallK0s(k0sctlConfigPath, kubeConfig); err != nil {
				return err
			}
		case distro.ProviderKind:
			if err = distro.InstallKind(blueprint.Metadata.Name, kubeConfig); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid Kubernetes distribution provider: %s", blueprint.Spec.Kubernetes.Provider)
		}
	}

	if err = kubeConfig.TryLoad(); err != nil {
		return err
	}

	// TODO (ranyodh): The following should be moved to distro specific types
	// create the k8sClient
	k8sClient, err := k8s.GetClient(kubeConfig)
	if err := k8s.WaitForNodes(k8sClient); err != nil {
		return fmt.Errorf("failed to wait for nodes: %w", err)
	}

	log.Info().Msgf("Installing Boundless Operator")

	var file string
	// Check if the user provided custom bop manifest
	if customBOPFlag != "" {
		log.Debug().Msgf("Installing Boundless Operator using custom manifest file: %q", customBOPFlag)
		file = customBOPFlag
	} else {
		log.Debug().Msgf("Installing Boundless Operator using manifest file: %s", boundless.ManifestUrl)
		file = boundless.ManifestUrl
	}
	err = k8s.Apply(file, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to install Boundless Operator: %w", err)
	}

	if err := k8s.WaitForPods(k8sClient, boundless.NamespaceBoundless); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}

	// install components
	log.Info().Msgf("Applying Boundless Operator resource")
	err = boundless.ApplyBlueprint(kubeConfig, blueprint)
	if err != nil {
		return fmt.Errorf("failed to install components: %w", err)
	}

	log.Info().Msgf("Finished installing Boundless Operator")

	return nil
}
