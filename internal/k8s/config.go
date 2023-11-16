package k8s

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// UsePersistentConfig caches client config to avoid reloads.
	UsePersistentConfig = true
)

// KubeConfig tracks a kubernetes configuration.
type KubeConfig struct {
	flags *genericclioptions.ConfigFlags
}

// NewConfig returns a new k8s config or an error if the flags are invalid.
func NewConfig(f *genericclioptions.ConfigFlags) *KubeConfig {
	return &KubeConfig{
		flags: f,
	}
}

func (c *KubeConfig) RESTConfig() (*restclient.Config, error) {
	return c.clientConfig().ClientConfig()
}

func (c *KubeConfig) clientConfig() clientcmd.ClientConfig {
	return c.flags.ToRawKubeConfigLoader()
}

// MergeConfig merges a new config into the existing config and writes it back to disk.
func (c *KubeConfig) MergeConfig(newConfig clientcmdapi.Config) error {
	cfg, err := c.clientConfig().RawConfig()
	if err != nil {
		return err
	}

	merge(&cfg, &newConfig)
	cfg.CurrentContext = newConfig.CurrentContext

	acc := c.ConfigAccess()
	return clientcmd.ModifyConfig(acc, cfg, true)
}

// DelContext remove a given context from the configuration.
func (c *KubeConfig) DelContext(n string) error {
	cfg, err := c.clientConfig().RawConfig()
	if err != nil {
		return err
	}
	delete(cfg.Contexts, n)

	acc := c.ConfigAccess()
	return clientcmd.ModifyConfig(acc, cfg, true)
}

// merge kind config into an existing config
func merge(existing, cfg *clientcmdapi.Config) {
	// insert or append cluster entry
	for name, cluster := range cfg.Clusters {
		existing.Clusters[name] = cluster
	}

	// insert or append user entry
	for name, info := range cfg.AuthInfos {
		existing.AuthInfos[name] = info
	}

	// insert or append context entry
	for name, context := range cfg.Contexts {
		existing.Contexts[name] = context
	}
}

// CurrentContextName returns the currently active config context.
func (c *KubeConfig) CurrentContextName() (string, error) {
	if isSet(c.flags.Context) {
		return *c.flags.Context, nil
	}
	cfg, err := c.clientConfig().RawConfig()
	if err != nil {
		return "", err
	}

	return cfg.CurrentContext, nil
}

// GetContext fetch a given context or error if it does not exists.
func (c *KubeConfig) GetContext(n string) (*clientcmdapi.Context, error) {
	cfg, err := c.clientConfig().RawConfig()
	if err != nil {
		return nil, err
	}
	if c, ok := cfg.Contexts[n]; ok {
		return c, nil
	}

	return nil, fmt.Errorf("invalid context `%s specified", n)
}

// ConfigAccess return the current kubeconfig api server access configuration.
func (c *KubeConfig) ConfigAccess() clientcmd.ConfigAccess {
	return c.clientConfig().ConfigAccess()
}

func (c *KubeConfig) GetConfigPath() string {
	return c.ConfigAccess().GetDefaultFilename()
}

func isSet(s *string) bool {
	return s != nil && len(*s) != 0
}
