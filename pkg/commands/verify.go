package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/mirantiscontainers/blueprint-cli/pkg/components"
	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
	"github.com/mirantiscontainers/blueprint-cli/pkg/distro"
	"github.com/mirantiscontainers/blueprint-cli/pkg/k8s"
	"github.com/mirantiscontainers/blueprint-cli/pkg/types"
	"github.com/mirantiscontainers/blueprint-cli/pkg/utils"

	"github.com/rs/zerolog/log"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// Verifies the contents of a blueprint against an existing cluster
func Verify(blueprint *types.Blueprint, kubeConfig *k8s.KubeConfig) error {
	ctx := context.Background()

	// Determine the distro
	provider, err := distro.GetProvider(blueprint, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to determine kubernetes provider: %w", err)
	}

	exists, err := provider.Exists()
	if err != nil {
		return fmt.Errorf("unable to check if provider exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("provider must be installed to verify addons, try applying a blueprint with no addons first")
	}

	// for now ignore manifests from the blueprint
	var helmAddons []types.Addon
	for _, addon := range blueprint.Spec.Components.Addons {
		if addon.Kind == constants.AddonChart {
			addon.DryRun = true
			helmAddons = append(helmAddons, addon)
		} else {
			log.Warn().Msgf("Manifest validation not available, manifest %s will not be verified", addon.Name)
		}
	}

	blueprint.Spec.Components.Addons = helmAddons

	defer func() {
		// dry running helm charts still creates the addon and chart CR , although helm chart contents are not created
		// clean up the dry addons before exiting

		blueprint.Spec.Components.Addons = nil
		err = components.ApplyBlueprint(kubeConfig, blueprint)
		if err != nil {
			log.Error().Msgf("failed to reset blueprint: %v", err)
		}

	}()

	err = components.ApplyBlueprint(kubeConfig, blueprint)
	if err != nil {
		return fmt.Errorf("failed to install components: %w", err)
	}

	k8sclient, err := k8s.GetClient(kubeConfig)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, addon := range blueprint.Spec.Components.Addons {
		// pods for failed dry-runs get deleted after running so run verification in parallel
		// this way we can get the detailed dry-run logs before the pod is deleted
		wg.Add(1)
		go verifyAddon(ctx, kubeConfig, addon, k8sclient, &wg)
	}

	wg.Wait()

	return nil
}

func verifyAddon(ctx context.Context, kubeConfig *k8s.KubeConfig, addon types.Addon, k8sclient *kubernetes.Clientset, wg *sync.WaitGroup) error {
	defer wg.Done()
	if addon.Kind == constants.AddonManifest {
		//TODO: BOP-309 Add Manifest Validation
		return nil
	}

	log.Info().Msgf("Verifying helmchart addon %s", addon.Chart.Name)

	var dryRunPod corev1.Pod

	dryRunPod, err := waitForJobReady(ctx, addon, k8sclient)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Warn().Msgf("Verification timed out for helmchart %s", addon.Chart.Name)
		}
		return err
	}

	// start streaming the job pod logs
	dryRunOutput, err := getPodLogs(kubeConfig, dryRunPod)
	if err != nil {
		log.Warn().Msgf("Verification failed for helmchart %s: %v", addon.Chart.Name, err)
		return err
	}

	// write dry-run output to a file and let user know where they can get detailed logs
	fileName, err := utils.WriteTempFile([]byte(dryRunOutput), fmt.Sprintf("%s-dry-run.log", addon.Chart.Name))
	if err != nil {
		return err
	}

	err = getJobResult(ctx, addon, k8sclient, err, fileName)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Warn().Msgf("Verification timed out for helmchart %s", addon.Chart.Name)
		}
		return err
	}

	return nil
}

func waitForJobReady(ctx context.Context, addon types.Addon, k8sclient *kubernetes.Clientset) (corev1.Pod, error) {
	var dryRunPod corev1.Pod

	// wait for the job that runs the helm install to start
	err := wait.PollUntilContextTimeout(ctx, 1*time.Second, constants.DryRunTimeout, true, func(ctx context.Context) (bool, error) {
		pods, err := k8sclient.CoreV1().Pods(addon.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("batch.kubernetes.io/job-name=helm-install-%s", addon.Chart.Name),
			Limit:         1,
		})
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debug().Msgf("Pod for helmchart addon %s not found, retrying", addon.Chart.Name)
				time.Sleep(constants.DryRunWaitInterval)
				return false, nil
			}
			log.Warn().Msgf("Failed to get pod for helmchart addon %s: %v", addon.Chart.Name, err)
			return false, err
		}

		if len(pods.Items) == 0 {
			log.Debug().Msgf("Pod for helmchart addon %s not found, retrying", addon.Chart.Name)
			time.Sleep(constants.DryRunWaitInterval)
			return false, nil
		}

		dryRunPod = pods.Items[0]

		// check if pod is ready
		if dryRunPod.Status.Phase == corev1.PodPending {
			log.Debug().Msgf("Pod for helmchart addon %s is still pending, retrying", addon.Chart.Name)
			time.Sleep(constants.DryRunWaitInterval)
			return false, nil
		}

		return true, nil
	})

	return dryRunPod, err
}

// getJobResult waits for the job to complete and checks the job result
func getJobResult(ctx context.Context, addon types.Addon, k8sclient *kubernetes.Clientset, err error, fileName string) error {
	var dryRunJob *batchv1.Job

	return wait.PollUntilContextTimeout(ctx, 5*time.Second, constants.DryRunTimeout, true, func(ctx context.Context) (bool, error) {
		dryRunJob, err = k8sclient.BatchV1().Jobs(addon.Namespace).Get(context.TODO(), fmt.Sprintf("helm-install-%s", addon.Chart.Name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debug().Msgf("Job for helmchart addon %s not found, retrying", addon.Chart.Name)
				return false, nil
			}
			return false, fmt.Errorf("failed to get job for helmchart addon %s: %v", addon.Chart.Name, err)
		}

		if len(dryRunJob.Status.Conditions) == 0 {
			log.Debug().Msgf("waiting on job %s conditions", addon.Chart.Name)
			return false, nil
		}

		for _, condition := range dryRunJob.Status.Conditions {
			log.Debug().Msgf("Condition: %s, Status: %s", condition.Type, condition.Status)
			if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
				log.Error().Msgf("Verification failed for helmchart %s, detailed dry-run for helmchart %s written to %s", addon.Chart.Name, addon.Chart.Name, fileName)
				return true, fmt.Errorf("verification failed for helmchart %s", addon.Chart.Name)
			}
			if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
				log.Info().Msgf("Verification completed for helmchart %s, detailed dry-run for helmchart %s written to %s", addon.Chart.Name, addon.Chart.Name, fileName)
				return true, nil
			}
		}

		return false, nil
	})

}

func getPodLogs(kubeConfig *k8s.KubeConfig, pod corev1.Pod) (string, error) {
	podLogOpts := corev1.PodLogOptions{Follow: true}

	clientset, err := k8s.GetClient(kubeConfig)
	if err != nil {
		return "unable to get k8s client", err
	}

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "error in opening stream", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf", err
	}
	str := buf.String()

	return str, nil
}
