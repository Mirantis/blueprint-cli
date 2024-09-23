package types

import (
	"fmt"
	"github.com/k0sproject/dig"
	"github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1"
	v1betacluster "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/version"
	v1 "github.com/k3s-io/helm-controller/pkg/apis/helm.cattle.io/v1"
	"sigs.k8s.io/yaml"
)

const apiVersion = "blueprint.mirantis.com/v1alpha1"
const apiVersionK0s = "k0sctl.k0sproject.io/v1beta1"

func ParseK0sCluster(data []byte) (v1beta1.Cluster, error) {
	var cluster v1beta1.Cluster
	err := yaml.Unmarshal(data, &cluster)
	if err != nil {
		return v1beta1.Cluster{}, err
	}
	return cluster, nil
}

func ParseBoundlessCluster(data []byte) (Blueprint, error) {
	var cluster Blueprint
	err := yaml.Unmarshal(data, &cluster)
	if err != nil {
		return Blueprint{}, err
	}

	return cluster, nil
}

func ParseCoreComponentManifests(data []byte) (v1.HelmChart, error) {
	var helmChart v1.HelmChart
	err := yaml.Unmarshal(data, &helmChart)
	if err != nil {
		return v1.HelmChart{}, err
	}

	return helmChart, nil
}

func ConvertToK0s(cluster *Blueprint) (v1beta1.Cluster, error) {

	var convertedK0sHosts []*v1betacluster.Host
	for _, host := range cluster.Spec.Kubernetes.Infra.Hosts {
		k0sHost := v1betacluster.Host{
			Connection: rig.Connection{
				SSH: &rig.SSH{
					Address: host.SSH.Address,
					User:    host.SSH.User,
					Port:    host.SSH.Port,
					KeyPath: &host.SSH.KeyPath,
				},
				Localhost: &rig.Localhost{Enabled: host.LocalHost.Enabled},
				OpenSSH:   nil,
				OSVersion: nil,
			},
			Role:         host.Role,
			InstallFlags: host.InstallFlags,
		}
		convertedK0sHosts = append(convertedK0sHosts, &k0sHost)
	}

	k0sVersion, err := version.NewVersion(cluster.Spec.Kubernetes.Version)
	if err != nil {
		return v1beta1.Cluster{}, fmt.Errorf("unable to parse provided version as valid k0s version: %w", err)
	}

	k0sCluster := v1beta1.Cluster{
		APIVersion: apiVersionK0s,
		Kind:       "Cluster",
		Metadata: &v1beta1.ClusterMetadata{
			Name: cluster.Metadata.Name,
		},
		Spec: &v1betacluster.Spec{
			Hosts: convertedK0sHosts,
			K0s: &v1betacluster.K0s{
				Version:       k0sVersion,
				DynamicConfig: digBool(cluster.Spec.Kubernetes.Config, "dynamicConfig"),
				Config:        cluster.Spec.Kubernetes.Config,
			},
		},
	}

	return k0sCluster, nil
}

func ConvertToClusterWithK0s(k0s v1beta1.Cluster, components Components) Blueprint {

	var boundlessHosts []Host
	for _, k0sHost := range k0s.Spec.Hosts {
		boundlessHost := Host{
			SSH: &SSHHost{
				Address: k0sHost.SSH.Address,
				KeyPath: *k0sHost.SSH.KeyPath,
				Port:    k0sHost.SSH.Port,
				User:    k0sHost.SSH.User,
			},
			LocalHost:    &LocalHost{Enabled: k0sHost.Localhost.Enabled},
			Role:         k0sHost.Role,
			InstallFlags: k0sHost.InstallFlags,
		}
		boundlessHosts = append(boundlessHosts, boundlessHost)
	}

	return Blueprint{
		APIVersion: apiVersion,
		Kind:       "Blueprint",
		Metadata: Metadata{
			Name: k0s.Metadata.Name,
		},
		Spec: BlueprintSpec{
			Version: "latest",
			Kubernetes: &Kubernetes{
				Provider: "k0s",
				Version:  k0s.Spec.K0s.Version.String(),
				Config:   k0s.Spec.K0s.Config,
				Infra: &Infra{
					Hosts: boundlessHosts,
				},
			},
			Components: components,
		},
	}
}

func ConvertToClusterWithKind(name string, components Components) Blueprint {
	return Blueprint{
		APIVersion: apiVersion,
		Kind:       "Blueprint",
		Metadata: Metadata{
			Name: name,
		},
		Spec: BlueprintSpec{
			Version: "latest",
			Kubernetes: &Kubernetes{
				Provider: "kind",
			},
			Components: components,
		},
	}
}

func DigToString(m dig.Mapping, keys ...string) string {
	val := m.Dig(keys...)
	if val == nil {
		return ""
	}
	return val.(string)
}

func digBool(m dig.Mapping, keys ...string) bool {
	val := m.Dig(keys...)
	if val == nil {
		return false
	}
	return val.(bool)
}
