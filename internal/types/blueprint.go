package types

import (
	"errors"
	"fmt"
	"os"

	"github.com/k0sproject/dig"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Blueprint struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   Metadata      `yaml:"metadata"`
	Spec       BlueprintSpec `yaml:"spec"`
}

// Validate validates the blueprint structure and its children
func (b *Blueprint) Validate() error {
	// TODO Check the APIVersion
	// TODO Check the cluster kind
	// TODO Check the metadata

	if err := b.Spec.Validate(); err != nil {
		return err
	}

	return nil
}

type BlueprintSpec struct {
	Kubernetes *Kubernetes `yaml:"kubernetes,omitempty"`
	Components Components  `yaml:"components"`
}

// Validate validates the blueprint spec structure and its children
func (bs *BlueprintSpec) Validate() error {
	if err := bs.Kubernetes.Validate(); err != nil {
		return err
	}

	// TODO Check the components

	return nil
}

type Infra struct {
	Hosts []Host `yaml:"hosts"`
}

type Kubernetes struct {
	Provider   string      `yaml:"provider"`
	Version    string      `yaml:"version,omitempty"`
	Config     dig.Mapping `yaml:"config,omitempty"`     // This is for defining the config within the blueprint
	ConfigPath string      `yaml:"configPath,omitempty"` // This is for passing the config as a separate file
	Infra      *Infra      `yaml:"infra,omitempty"`
}

// Validate validates the Kubernetes structure and its children
func (k *Kubernetes) Validate() error {
	// TODO Check the provider
	// TODO Check the version
	// TODO Check the config
	// TODO Check the infra

	if k.ConfigPath != "" && k.Config != nil {
		return fmt.Errorf("cannot specify both config and configPath")
	}

	// Only ConfigPath is specified
	if k.ConfigPath != "" && k.Config == nil {
		// Check that the file exists
		if _, err := os.Stat(k.ConfigPath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config file %q does not exist: %s", k.ConfigPath, err)
		}
	}

	// TODO Check the Config

	return nil
}

type Components struct {
	Core   *Core    `yaml:"core,omitempty"`
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

// Addons defines the desired state of Addon
type Addons struct {
	Name      string        `yaml:"name"`
	Kind      string        `yaml:"kind"`
	Enabled   bool          `yaml:"enabled"`
	Namespace string        `yaml:"namespace,omitempty"`
	Chart     *ChartInfo    `json:"chart,omitempty"`
	Manifest  *ManifestInfo `json:"manifest,omitempty"`
}

// ChartInfo defines the desired state of chart
type ChartInfo struct {
	Name    string                        `yaml:"name"`
	Repo    string                        `yaml:"repo"`
	Version string                        `yaml:"version"`
	Set     map[string]intstr.IntOrString `yaml:"set,omitempty"`
	Values  string                        `yaml:"values,omitempty"`
}

// ManifestInfo defines the desired state of manifest
type ManifestInfo struct {
	URL string `json:"url"`
}
