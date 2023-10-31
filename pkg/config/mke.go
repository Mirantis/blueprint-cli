package config

import (
	"github.com/k0sproject/dig"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Blueprint struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       ClusterSpec `yaml:"spec"`
}

type ClusterSpec struct {
	Infra      Infra      `yaml:"infra,omitempty"`
	Kubernetes Kubernetes `yaml:"kubernetes,omitempty"`
	Components Components `yaml:"components"`
}

type Infra struct {
	Hosts []Host `yaml:"hosts"`
}

type Kubernetes struct {
	Provider string      `yaml:"provider"`
	Version  string      `yaml:"version,omitempty"`
	Config   dig.Mapping `yaml:"config,omitempty"`
}

type Components struct {
	Core   Core     `yaml:"core,omitempty"`
	Addons []Addons `yaml:"addons,omitempty"`
}

type Core struct {
	Cni        *CoreComponent `yaml:"cni,omitempty"`
	Ingress    *CoreComponent `yaml:"ingress,omitempty"`
	DNS        *CoreComponent `yaml:"dns,omitempty"`
	Logging    *CoreComponent `yaml:"logging,omitempty"`
	Monitoring *CoreComponent `yaml:"monitoring,omitempty"`
}

type CoreComponent struct {
	Enabled  bool        `yaml:"enabled"`
	Provider string      `yaml:"provider"`
	Config   dig.Mapping `yaml:"config,omitempty"`
}

type Addons struct {
	Name      string `yaml:"name"`
	Kind      string `yaml:"kind"`
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace,omitempty"`
	Chart     Chart  `yaml:"chart"`
}

type Chart struct {
	Name    string                        `yaml:"name"`
	Repo    string                        `yaml:"repo"`
	Version string                        `yaml:"version"`
	Set     map[string]intstr.IntOrString `yaml:"set,omitempty"`
	Values  string                        `yaml:"values,omitempty"`
}
