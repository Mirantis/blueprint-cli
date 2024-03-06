package commands

import (
	"context"
	"fmt"

	"github.com/mirantiscontainers/boundless-cli/pkg/components"
	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/distro"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Apply installs the Boundless Operator and applies the components defined in the blueprint
func Apply(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig, operatorUri string) error {
	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	exists, err := provider.Exists()
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	// If we are working with an unsupported provider, we need to make sure it exists
	// For other supported providers, we check whether boundless is already installed
	if provider.Type() == constants.ProviderExisting {
		if !exists {
			return fmt.Errorf("cluster %q already exists", blueprint.Metadata.Name)
		}
	}
	if exists {
		log.Info().Msgf("Cluster %q already exists", blueprint.Metadata.Name)
	} else {
		if err := provider.Install(); err != nil {
			return fmt.Errorf("failed to install cluster: %w", err)
		}
	}

	if err = kubeConfig.TryLoad(); err != nil {
		return err
	}

	// Setup the client
	if err := provider.SetupClient(); err != nil {
		return fmt.Errorf("failed to setup client: %w", err)
	}
	k8sclient, err := k8s.GetClient(kubeConfig)
	if err != nil {
		panic(err)
	}

	// For existing clusters, determine whether boundless is currently installed
	installOperator := true
	if exists {
		_, err := k8sclient.AppsV1().Deployments(constants.NamespaceBoundless).Get(context.TODO(), constants.BoundlessOperatorDeployment, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Warn().Msgf("Could not determine existing Boundless Operator installation: %s", err)
			}
		} else {
			// @todo: determine operator version
			installOperator = false
		}
	}

	// @todo: display the version of the operator
	if installOperator {
		log.Info().Msgf("Installing Boundless Operator")
		log.Trace().Msgf("Installing Boundless Operator using manifest file: %s", operatorUri)
		if err = k8s.ApplyYaml(kubeConfig, operatorUri); err != nil {
			return fmt.Errorf("failed to install Boundless Operator: %w", err)
		}
	} else {
		log.Info().Msg("Boundless Operator already installed")
	}

	// Wait for the pods to be ready
	if err := provider.WaitForPods(); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}

	// install components
	log.Info().Msgf("Applying Boundless Operator resource")
	err = components.ApplyBlueprint(kubeConfig, blueprint)
	if err != nil {
		return fmt.Errorf("failed to install components: %w", err)
	}

	log.Info().Msgf("Finished installing Boundless Operator")

	return nil
}
