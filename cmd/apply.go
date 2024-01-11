package cmd

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/internal/boundless"
	"github.com/mirantiscontainers/boundless-cli/internal/distro"
	"github.com/mirantiscontainers/boundless-cli/internal/k8s"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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
	addOperatorUriFlag(flags)
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}

func runApply() error {
	log.Info().Msgf("Applying blueprint %s", blueprintFlag)

	// Determine the distro
	provider, err := distro.GetProvider(&blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	// Install the distro
	if err := provider.Install(); err != nil {
		return fmt.Errorf("failed to install cluster: %w", err)
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

	// @todo: display the version of the operator
	log.Info().Msgf("Installing Boundless Operator")
	log.Trace().Msgf("Installing boundless operator using manifest file: %s", operatorUri)
	if err = k8s.ApplyYaml(kubeConfig, operatorUri); err != nil {
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
