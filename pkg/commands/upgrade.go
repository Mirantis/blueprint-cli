package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
	"github.com/mirantiscontainers/blueprint-cli/pkg/distro"
	"github.com/mirantiscontainers/blueprint-cli/pkg/k8s"
	"github.com/mirantiscontainers/blueprint-cli/pkg/types"
)

// Upgrade upgrades the Blueprint Operator
func Upgrade(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig, imageRegistry string) error {
	var client kubernetes.Interface
	var err error

	if client, err = k8s.GetClient(kubeConfig); err != nil {
		return fmt.Errorf("failed to get kubernetes client: %q", err)
	}

	bopDeployment, err := client.AppsV1().Deployments(constants.NamespaceBlueprint).Get(context.TODO(), constants.BlueprintOperatorDeployment, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get existing Blueprint Operator deployment: %w", err)
	}

	deployedRegistry, err := detectDeployedRegistry(bopDeployment.Spec.Template.Spec.Containers)
	if err != nil {
		return fmt.Errorf("failed to detect image registry of the deployed Blueprint Operator: %w", err)
	}

	if imageRegistry == "" {
		imageRegistry = deployedRegistry
	} else if deployedRegistry != imageRegistry {
		return fmt.Errorf(
			"image registry %s does not match the deployed blueprint operator image registry %s; "+
				"use --image-registry flag to upgrade with the same registry, "+
				"or run `bctl apply --image-registry` to change the image registry of the deployed BOP before upgrading",
			imageRegistry, deployedRegistry,
		)
	}

	uri, err := determineOperatorUri(blueprint.Spec.Version)
	if err != nil {
		return fmt.Errorf("failed to determine operator URI: %w", err)
	}

	var needCleanup bool
	uri, needCleanup, err = setImageRegistry(uri, imageRegistry)
	if err != nil {
		return fmt.Errorf("failed to set image registry in BOP manifest: %w", err)
	}
	if needCleanup {
		defer os.Remove(strings.TrimPrefix(uri, "file://"))
	}

	var dynamicClient dynamic.Interface

	if dynamicClient, err = k8s.GetDynamicClient(kubeConfig); err != nil {
		return fmt.Errorf("failed to get kubernetes dynamic client: %q", err)
	}

	log.Info().Msgf("Upgrading Blueprint Operator using manifest file %q", uri)
	if err := k8s.ApplyYaml(client, dynamicClient, uri); err != nil {
		return fmt.Errorf("failed to upgrade blueprint operator: %w", err)
	}

	log.Info().Msgf("Finished updating Blueprint Operator")

	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}
	if err := provider.SetupClient(); err != nil {
		return fmt.Errorf("failed to setup client: %w", err)
	}
	// Wait for the pods to be ready
	if err := provider.WaitForPods(); err != nil {
		return fmt.Errorf("failed to wait for pods: %w", err)
	}
	return nil
}
