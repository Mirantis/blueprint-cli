package utils

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
)

func PrintDeploymentStatus(deployment v1.Deployment) {
	detailedStatus := fmt.Sprintf("Desired: %d, Ready: %d/%d, Available: %d/%d",
		deployment.Status.Replicas,
		deployment.Status.ReadyReplicas,
		deployment.Status.Replicas,
		deployment.Status.AvailableReplicas,
		deployment.Status.Replicas,
	)
	fmt.Printf("%-30s %-30s %-30s\n", "Deployment", deployment.Name, detailedStatus)
}

func PrintDaemonsetStatus(ds v1.DaemonSet) {
	detailedStatus := fmt.Sprintf("Desired: %d, Ready: %d/%d, Available: %d/%d",
		ds.Status.DesiredNumberScheduled,
		ds.Status.NumberReady,
		ds.Status.DesiredNumberScheduled,
		ds.Status.NumberAvailable,
		ds.Status.DesiredNumberScheduled,
	)
	fmt.Printf("%-30s %-30s %-30s\n", "Daemonset", ds.Name, detailedStatus)
}

func PrintStatefulsetStatus(ss v1.StatefulSet) {
	detailedStatus := fmt.Sprintf("Desired: %d, Ready: %d/%d, Available: %d/%d",
		ss.Status.Replicas,
		ss.Status.ReadyReplicas,
		ss.Status.Replicas,
		ss.Status.AvailableReplicas,
		ss.Status.Replicas,
	)
	fmt.Printf("%-30s %-30s %-30s\n", "StatefulSet", ss.Name, detailedStatus)
}
