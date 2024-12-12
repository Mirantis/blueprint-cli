package commands

import (
	"context"
	"fmt"

	"github.com/mirantiscontainers/blueprint-cli/pkg/components"
	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
	"github.com/mirantiscontainers/blueprint-cli/pkg/distro"
	"github.com/mirantiscontainers/blueprint-cli/pkg/k8s"
	"github.com/mirantiscontainers/blueprint-cli/pkg/types"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Apply installs the Blueprint Operator and applies the components defined in the blueprint
func Apply(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig, providerInstallOnly bool) error {
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
	// For other supported providers, we check whether blueprint is already installed
	if provider.Type() == constants.ProviderExisting {
		if !exists {
			return fmt.Errorf("cluster %q already exists", blueprint.Metadata.Name)
		}
	}
	if !exists {
		if err := provider.Install(); err != nil {
			return fmt.Errorf("failed to install cluster: %w", err)

		}
	} else {
		log.Info().Msgf("Cluster %q already exists", blueprint.Metadata.Name)
		if err = provider.Refresh(); err != nil {
			return fmt.Errorf("failed to refresh cluster: %w", err)
		}
	}

	if providerInstallOnly {
		return nil
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

	// For existing clusters, determine whether blueprint is currently installed
	installOperator := true
	if exists {
		_, err := k8sclient.AppsV1().Deployments(constants.NamespaceBlueprint).Get(context.TODO(), constants.BlueprintOperatorDeployment, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Warn().Msgf("Could not determine existing Blueprint Operator installation: %s", err)
			}
		} else {
			// @todo: determine operator version
			installOperator = false
		}
	}

	uri, err := determineOperatorUri(blueprint.Spec.Version)
	if err != nil {
		return fmt.Errorf("failed to determine operator URI: %w", err)
	}

	// @todo: display the version of the operator
	if installOperator {
		log.Info().Msgf("Installing Blueprint Operator")
		log.Debug().Msgf("Installing Blueprint Operator using manifest file: %s", blueprint.Spec.Version)
		if err = k8s.ApplyYaml(kubeConfig, uri); err != nil {
			return fmt.Errorf("failed to install Blueprint Operator: %w", err)
		}
	} else {
		log.Info().Msg("Blueprint Operator already installed")
	}

	// Wait for the pods to be ready
	if err := provider.WaitForPods(); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}

	// install components
	log.Info().Msgf("Applying Blueprint Operator resource")
	err = components.ApplyBlueprint(kubeConfig, blueprint)
	if err != nil {
		return fmt.Errorf("failed to install components: %w", err)
	}

	log.Info().Msgf("Finished installing Blueprint Operator")

	return nil
}
