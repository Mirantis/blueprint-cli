package types

import (
	"github.com/k0sproject/dig"
)

type K0sCluster struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   Metadata       `yaml:"metadata"`
	Spec       K0sClusterSpec `yaml:"spec"`
}

type K0sClusterSpec struct {
	Hosts []Host `yaml:"hosts"`
	K0S   K0s    `yaml:"k0s"`
}

type K0s struct {
	Version       string      `yaml:"version"`
	DynamicConfig bool        `yaml:"dynamicConfig"`
	Config        dig.Mapping `yaml:"config,omitempty"`
	//Spec          K0sSpec     `yaml:"spec"`
}

//type K0sSpec struct {
//	Network K0sNetwork `yaml:"network"`
//}
//
//type K0sNetwork struct {
//	Provider string `yaml:"provider"`
//	Calico   Calico `yaml:"calico"`
//}
