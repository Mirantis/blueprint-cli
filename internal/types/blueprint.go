package types

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"

	"github.com/k0sproject/dig"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var blueprintKinds = []string{"Blueprint"}

type Blueprint struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   Metadata      `yaml:"metadata"`
	Spec       BlueprintSpec `yaml:"spec"`
}

// Validate checks the Blueprint structure and its children
func (b *Blueprint) Validate() error {
	// APIVersion checks
	if b.APIVersion == "" {
		return fmt.Errorf("apiVersion field cannot be left blank")
	}

	// Kind checks
	if b.Kind == "" {
		return fmt.Errorf("kind field cannot be left blank")
	}
	if !slices.Contains(blueprintKinds, b.Kind) {
		return fmt.Errorf("invalid cluster kind: %s", b.Kind)
	}

	// Metadata checks
	if err := b.Metadata.Validate(); err != nil {
		return err
	}

	// Spec checks
	if err := b.Spec.Validate(); err != nil {
		return err
	}

	return nil
}

type BlueprintSpec struct {
	Kubernetes *Kubernetes `yaml:"kubernetes,omitempty"`
	Components Components  `yaml:"components"`
}

// Validate checks the BlueprintSpec structure and its children
func (bs *BlueprintSpec) Validate() error {

	// Kubernetes checks
	if bs.Kubernetes != nil {
		if err := bs.Kubernetes.Validate(); err != nil {
			return err
		}
	}

	// Components checks
	if err := bs.Components.Validate(); err != nil {
		return err
	}

	return nil
}

type Infra struct {
	Hosts []Host `yaml:"hosts"`
}

// Validate checks the Infra structure and its children
func (i *Infra) Validate() error {

	// Host checks
	for _, host := range i.Hosts {
		if err := host.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Kubernetes struct {
	Provider string      `yaml:"provider"`
	Version  string      `yaml:"version,omitempty"`
	Config   dig.Mapping `yaml:"config,omitempty"`
	Infra    *Infra      `yaml:"infra,omitempty"`
}

var providerKinds = []string{"kind", "k0s"}

// Validate checks the Kubernetes structure and its children
func (k *Kubernetes) Validate() error {
	// Provider checks
	if k.Provider == "" {
		return fmt.Errorf("kubernetes.provider field cannot be left blank")
	}
	if !slices.Contains(providerKinds, k.Provider) {
		return fmt.Errorf("invalid kubernetes.provider: %s", k.Provider)
	}

	// Version checks
	// The version can be left empty, but if it's not, it must be a valid k0s semver
	if k.Version != "" {
		// This regex gives us semver with an optional "+k0s.0"
		re, _ := regexp.Compile(`/^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$/gm`)
		if !re.MatchString(k.Version) {
			return fmt.Errorf("invalid kubernetes.version: %s", k.Version)
		}
	}

	// Infra checks
	if k.Infra != nil {
		if err := k.Infra.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Components struct {
	Core   *Core   `yaml:"core,omitempty"`
	Addons []Addon `yaml:"addons,omitempty"`
}

// Validate checks the Components structure and its children
func (c *Components) Validate() error {
	// TODO Core components aren't checked because they will likely be removed/moved to MKE4

	// Addon checks
	for _, addon := range c.Addons {
		if err := addon.Validate(); err != nil {
			return err
		}
	}

	return nil
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

var addonKinds = []string{"chart", "Chart", "manifest", "Manifest"}

// Addon defines the desired state of an Addon
type Addon struct {
	Name      string        `yaml:"name"`
	Kind      string        `yaml:"kind"`
	Enabled   bool          `yaml:"enabled"`
	Namespace string        `yaml:"namespace,omitempty"`
	Chart     *ChartInfo    `json:"chart,omitempty"`
	Manifest  *ManifestInfo `json:"manifest,omitempty"`
}

// Validate checks the Addon structure and its children
func (a *Addon) Validate() error {

	// Name checks
	if a.Name == "" {
		return fmt.Errorf("addons.name field cannot be left blank")
	}

	// Kind checks
	if a.Kind == "" {
		return fmt.Errorf("addons.kind field cannot be left blank")
	}
	if !slices.Contains(addonKinds, a.Kind) {
		return fmt.Errorf("%s addons.kind field is an invalid kind: %s", a.Name, a.Kind)
	}

	// Chart checks
	if a.Chart != nil {
		if err := a.Chart.Validate(); err != nil {
			return err
		}
	}

	// Manifest checks
	if a.Manifest != nil {
		if err := a.Manifest.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ChartInfo defines the desired state of chart
type ChartInfo struct {
	Name    string                        `yaml:"name"`
	Repo    string                        `yaml:"repo"`
	Version string                        `yaml:"version"`
	Set     map[string]intstr.IntOrString `yaml:"set,omitempty"`
	Values  string                        `yaml:"values,omitempty"`
}

// Validate checks the ChartInfo structure and its children
func (ci *ChartInfo) Validate() error {
	// Name checks
	if ci.Name == "" {
		return fmt.Errorf("chart.name field cannot be left blank")
	}

	// Repo checks
	if ci.Repo == "" {
		return fmt.Errorf("chart.repo field cannot be left blank")
	}

	// Version checks
	if ci.Version == "" {
		return fmt.Errorf("chart.version field cannot be left blank")
	}

	return nil
}

// ManifestInfo defines the desired state of manifest
type ManifestInfo struct {
	URL string `json:"url"`
}

// Validate checks the ManifestInfo structure and its children
func (mi *ManifestInfo) Validate() error {
	// URL checks
	if mi.URL == "" {
		return fmt.Errorf("manifest.url field cannot be left blank")
	}
	if _, err := url.ParseRequestURI(mi.URL); err != nil {
		return fmt.Errorf("manifest.url field must be a valid url: %v", mi.URL)
	}

	return nil
}
