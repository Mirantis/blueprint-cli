package types

import (
	"github.com/k0sproject/dig"
)

type K0sCluster struct {
	APIVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       string         `yaml:"kind" json:"kind"`
	Metadata   Metadata       `yaml:"metadata" json:"metadata"`
	Spec       K0sClusterSpec `yaml:"spec" json:"spec"`
}

type K0sClusterSpec struct {
	Hosts []Host `yaml:"hosts" json:"hosts"`
	K0S   K0s    `yaml:"k0s" json:"k0S"`
}

type K0s struct {
	Version       string      `yaml:"version" json:"version"`
	DynamicConfig bool        `yaml:"dynamicConfig" json:"dynamicConfig"`
	Config        dig.Mapping `yaml:"config,omitempty" json:"config"`
}
