package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/mirantiscontainers/boundless-cli/boundlessclientset"
	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/mirantiscontainers/boundless-cli/pkg/k8s"
	"github.com/mirantiscontainers/boundless-cli/pkg/utils"
	"github.com/mirantiscontainers/boundless-operator/api/v1alpha1"
)

const (
	helmControllerDeployment     = "helm-controller"
	boundlessOperatorDeployment  = "boundless-operator-controller-manager"
	kubernetesManagedByLabel     = "app.kubernetes.io/managed-by"
	kubernetesManagedByHelmValue = "Helm"
	kubernetesInstanceLabel      = "app.kubernetes.io/instance"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Get the status of the blueprint",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			//log.Info().Msgf("Args are %v", args)
			if len(args) > 0 {
				return runAddonSpecificStatus(args[0])
			}

			return runStatus()
		},
	}

	flags := cmd.Flags()
	addOperatorUriFlag(flags)
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}

func runStatus() error {
	k8sclient, err := k8s.GetClient(kubeConfig)
	if err != nil {
		panic(err)
	}

	operatorDeployment, err := k8sclient.AppsV1().Deployments(constants.NamespaceBoundless).Get(context.TODO(), boundlessOperatorDeployment, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("No boundless operator installation detected")
		} else {
			panic(err)
		}
	} else {
		utils.PrintDeploymentStatus(*operatorDeployment)
	}

	helmController, err := k8sclient.AppsV1().Deployments(constants.NamespaceBoundless).Get(context.TODO(), helmControllerDeployment, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("No helm controller detected - Chart addons may not function")
		} else {
			panic(err)
		}
	} else {
		utils.PrintDeploymentStatus(*helmController)
	}

	fmt.Println("-------------------------------------------------------")

	addonList, err := getAddons()
	if err != nil {
		panic(err)
	}

	if len(addonList.Items) == 0 {
		fmt.Println("No addons installed")
		return nil
	}

	fmt.Printf("%-20s %-10s %-10s\n", "NAME", "KIND", "STATUS")
	for _, addon := range addonList.Items {
		fmt.Printf("%-20s %-10s %-10s\n", addon.Name, addon.Spec.Kind, addon.Status.Type)
	}

	return nil
}

func getAddons() (*v1alpha1.AddonList, error) {
	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := getBoundlessClientSet()
	if err != nil {
		return nil, err
	}

	addonList, err := clientSet.Addons(constants.NamespaceBoundless).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return addonList, nil
}

func getAddon(addonName string) (*v1alpha1.Addon, error) {
	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := getBoundlessClientSet()
	if err != nil {
		return nil, err
	}

	addon, err := clientSet.Addons(constants.NamespaceBoundless).Get(addonName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return addon, nil
}

func getBoundlessClientSet() (*boundlessclientset.BoundlessV1Alpha1Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig.GetConfigPath())
	if err != nil {
		return nil, err
	}

	clientSet, err := boundlessclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func runAddonSpecificStatus(providedAddonName string) error {
	providedAddon, err := getAddon(providedAddonName)
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("invalid input %s, no addon named %s exists", providedAddonName, providedAddonName)
		}
		return err
	}

	fmt.Printf("%-20s %-10s %-10s\n", "NAME", "KIND", "STATUS")
	fmt.Printf("%-20s %-10s %-10s\n\n", providedAddon.Name, providedAddon.Spec.Kind, providedAddon.Status.Type)

	fmt.Printf("Status Reason: %s\n", providedAddon.Status.Reason)
	fmt.Printf("Detailed Status Message: %s\n\n", providedAddon.Status.Message)

	k8sclient, err := k8s.GetClient(kubeConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("-------------------------------------------------------")
	fmt.Println("ADDON RESOURCES")
	if strings.EqualFold(providedAddon.Spec.Kind, "chart") {
		printHelmchartResources(k8sclient, *providedAddon)

	} else {
		printManifestResources(*providedAddon, k8sclient)

	}
	fmt.Println("-------------------------------------------------------")

	// lastly show any events created by boundless
	// kubernetes events are relatively short-lived, so we can't rely on them always being here

	var eventMsgs []string

	eventList, err := k8sclient.EventsV1().Events(constants.NamespaceBoundless).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, event := range eventList.Items {
		if event.ObjectMeta.Annotations["Addon"] == providedAddonName {
			eventMsgs = append(eventMsgs, event.Note)
		}
	}

	if len(eventMsgs) > 0 {
		fmt.Println("\nBOUNDLESS SYSTEM EVENTS")
		for _, msg := range eventMsgs {
			fmt.Printf("%s\n", msg)
		}
	} else {
		fmt.Printf("No boundless system events for addon %s\n", providedAddonName)
	}

	return nil
}

func printManifestResources(providedAddon v1alpha1.Addon, k8sclient *kubernetes.Clientset) {
	clientSet, err := getBoundlessClientSet()
	if err != nil {
		panic(err)
	}

	manifest, err := clientSet.Manifests(constants.NamespaceBoundless).Get(providedAddon.Spec.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	for _, obj := range manifest.Spec.Objects {
		if obj.Kind == "DaemonSet" {
			ds, err := k8sclient.AppsV1().DaemonSets(obj.Namespace).Get(context.TODO(), obj.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Printf("Unable to get Daemonset %s\n", ds.Name)
				continue
			}
			utils.PrintDaemonsetStatus(*ds)
		}

		if obj.Kind == "Deployment" {
			deployment, err := k8sclient.AppsV1().Deployments(obj.Namespace).Get(context.TODO(), obj.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Printf("Unable to get Deployment %s\n", deployment.Name)
				continue
			}
			utils.PrintDeploymentStatus(*deployment)

		}
	}
}

func printHelmchartResources(k8sclient *kubernetes.Clientset, providedAddon v1alpha1.Addon) {
	// show resources related to the helm chart
	// limited to pods,services,daemonsets,deployments - similar to how `kubectl get all` only shows those resources

	deploymentList, err := k8sclient.AppsV1().Deployments(providedAddon.Spec.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err == nil && len(deploymentList.Items) > 0 {
		for _, deployment := range deploymentList.Items {
			if len(deployment.Labels) > 0 && deployment.Labels[kubernetesManagedByLabel] == kubernetesManagedByHelmValue && deployment.Labels[kubernetesInstanceLabel] == providedAddon.Spec.Chart.Name {
				utils.PrintDeploymentStatus(deployment)
			}
		}
	}

	daemonsetList, err := k8sclient.AppsV1().DaemonSets(providedAddon.Spec.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err == nil && len(daemonsetList.Items) > 0 {
		for _, ds := range daemonsetList.Items {
			if len(ds.Labels) > 0 && ds.Labels[kubernetesManagedByLabel] == kubernetesManagedByHelmValue && ds.Labels[kubernetesInstanceLabel] == providedAddon.Spec.Chart.Name {
				utils.PrintDaemonsetStatus(ds)
			}

		}
	}

	statefulSetList, err := k8sclient.AppsV1().StatefulSets(providedAddon.Spec.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err == nil && len(statefulSetList.Items) > 0 {
		for _, ss := range statefulSetList.Items {
			if len(ss.Labels) > 0 && ss.Labels[kubernetesManagedByLabel] == kubernetesManagedByHelmValue && ss.Labels[kubernetesInstanceLabel] == providedAddon.Spec.Chart.Name {
				utils.PrintStatefulsetStatus(ss)
			}

		}
	}
}
