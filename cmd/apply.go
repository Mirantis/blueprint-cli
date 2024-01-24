package cmd

import (
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/internal/boundless"
	"github.com/mirantiscontainers/boundless-cli/internal/distro"
	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/constants"

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

	exists, err := provider.Exists()
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}
	// If we are working with an unsupported provider, we need to make sure it exists
	// If we are working with a supported provider, we need to make sure it does not exist
	if provider.Type() != constants.ProviderExisting {
		if exists {
			return fmt.Errorf("cluster %q already exists", blueprint.Metadata.Name)
		}

		// Install the distro
		if err := provider.Install(); err != nil {
			return fmt.Errorf("failed to install cluster: %w", err)
		}
	} else {
		if !exists {
			return fmt.Errorf("cluster does not exist: %s", blueprint.Metadata.Name)
		}
	}

	if err = kubeConfig.TryLoad(); err != nil {
		return err
	}

	// Setup the client
	if err := provider.SetupClient(); err != nil {
		return fmt.Errorf("failed to setup client: %w", err)
	}

	// @todo: display the version of the operator
	log.Info().Msgf("Installing Boundless Operator")
	log.Trace().Msgf("Installing boundless operator using manifest file: %s", operatorUri)
	if err = k8s.ApplyYaml(kubeConfig, operatorUri); err != nil {
		return fmt.Errorf("failed to install Boundless Operator: %w", err)
	}

	// Wait for the pods to be ready
	if err := provider.WaitForPods(); err != nil {
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
