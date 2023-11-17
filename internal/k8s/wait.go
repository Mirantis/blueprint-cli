package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// WaitForNodes waits for all nodes to be ready
func WaitForNodes(client kubernetes.Interface) error {
	log.Info().Msgf("Waiting for nodes to be ready")
	return waitForNodes(context.Background(), client)
}

// WaitForPods waits for all pods in the given namespace to be running
func WaitForPods(client kubernetes.Interface, namespace string) error {
	log.Info().Msgf("Waiting for all pods to be ready")
	return waitForPods(context.Background(), client, namespace)
}

func waitForPods(ctx context.Context, clientset kubernetes.Interface, namepsace string) error {
	// wait for all pods
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, 5*time.Minute)
	defer cancelFunc()

	return wait.PollUntilContextCancel(timeoutCtx, 5*time.Second, true, func(ctx context.Context) (bool, error) {
		pods, err := clientset.CoreV1().Pods(namepsace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, fmt.Errorf("failed to list pods: %v", err)
		}

		if len(pods.Items) == 0 {
			return false, nil
		}

		allRunning := true
		for _, pod := range pods.Items {
			log.Debug().Msgf("pod %s is %s", pod.Name, pod.Status.Phase)
			if !podInPhase([]v1.PodPhase{v1.PodRunning, v1.PodSucceeded}, pod.Status.Phase) {
				allRunning = false
				break
			}
		}

		return allRunning, nil
	})
}

func podInPhase(strings []v1.PodPhase, s v1.PodPhase) bool {
	for _, str := range strings {
		if str == s {
			return true
		}
	}
	return false
}

func waitForNodes(ctx context.Context, clientset kubernetes.Interface) error {
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, 5*time.Minute)
	defer cancelFunc()

	return wait.PollUntilContextCancel(timeoutCtx, 5*time.Second, true, func(ctx context.Context) (bool, error) {
		// wait for node to be ready
		nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, fmt.Errorf("failed to list nodes: %v", err)
		}

		if len(nodes.Items) == 0 {
			return false, nil
		}

		allReady := true
		for _, node := range nodes.Items {
			ready := false
			for _, condition := range node.Status.Conditions {
				if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
					ready = true
					break
				}
			}

			allReady = allReady && ready
		}
		return allReady, nil
	})
}
