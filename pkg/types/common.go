package types

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
)

type Metadata struct {
	Name string `yaml:"name"`
}

// Validate checks the Metadata structure and its children
func (m *Metadata) Validate() error {
	// This is just a placeholder for now

	return nil
}

type Host struct {
	SSH          *SSHHost   `yaml:"ssh,omitempty"`
	LocalHost    *LocalHost `yaml:"localhost,omitempty"`
	Role         string     `yaml:"role"`
	InstallFlags []string   `yaml:"installFlags,omitempty"`
}

var nodeRoles = []string{"single", "controller", "worker", "controller+worker"}

// Validate checks the Host structure and its children
func (h *Host) Validate() error {

	// SSH checks
	if h.SSH != nil {
		if err := h.SSH.Validate(); err != nil {
			return err
		}
	}

	// Localhost checks
	if h.LocalHost != nil {
		if err := h.LocalHost.Validate(); err != nil {
			return err
		}
	}

	// Role checks
	if h.Role == "" {
		return fmt.Errorf("hosts.role field cannot be left blank")
	}
	if !slices.Contains(nodeRoles, h.Role) {
		return fmt.Errorf("invalid hosts.role: %s\nValid hosts.role values: %s", h.Role, nodeRoles)
	}

	return nil
}

type SSHHost struct {
	Address string `yaml:"address"`
	KeyPath string `yaml:"keyPath"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
}

// Validate checks the SSHHost structure and its children
func (sh *SSHHost) Validate() error {

	// Address checks
	if sh.Address == "" {
		return fmt.Errorf("hosts.ssh.address field cannot be left empty")
	}
	// This regex is for either valid hostnames or ip addresses
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if !re.MatchString(sh.Address) {
		return fmt.Errorf("invalid hosts.ssh.address: %s", sh.Address)
	}

	// KeyPath checks
	if sh.KeyPath == "" {
		return fmt.Errorf("hosts.ssh.keypath field cannot be left empty")
	}
	if _, err := os.Stat(sh.KeyPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("hosts.ssh.keypath does not exist: %s", sh.KeyPath)
	}

	// Port checks
	if sh.Port <= 0 || sh.Port > 65535 {
		return fmt.Errorf("hosts.ssh.port outside of valid range 0-65535")
	}

	// User checks
	if sh.User == "" {
		return fmt.Errorf("hosts.ssh.user cannot be left empty")
	}

	return nil
}

type LocalHost struct {
	Enabled bool `yaml:"enabled"`
}

// Validate checks the LocalHost structure and its children
func (l *LocalHost) Validate() error {
	// This is just a placeholder for now
	return nil
}
