package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/k0sproject/dig"
	"github.com/mirantiscontainers/boundless-cli/internal/boundless"
	"github.com/mirantiscontainers/boundless-cli/internal/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var isKind bool

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a blueprint file template",
		RunE:  runInit,
	}

	// @TODO This should be a subcommand of init
	// bctl init kind
	// bctl init k0s (default)
	cmd.Flags().BoolVar(&isKind, "kind", false, "Install a kind cluster")
	cmd.Flags()

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	if isKind {
		return encode(types.ConvertToClusterWithKind("boundless-cluster", defaultComponents))
	}

	// @TODO Include pFlags for k0sctl init
	cmd2 := exec.Command("k0sctl", "init")
	cmd2.Stdin = os.Stdin
	cmd2.Stderr = os.Stderr

	buf := new(bytes.Buffer)
	cmd2.Stdout = buf
	err := cmd2.Run()
	if err != nil {
		return err
	}

	k0sConfig, err := types.ParseK0sCluster(buf.Bytes())
	if err != nil {
		return err
	}

	return encode(types.ConvertToClusterWithK0s(k0sConfig, defaultComponents))
}

func encode(blueprint types.Blueprint) error {
	encoder := yaml.NewEncoder(os.Stdout)
	return encoder.Encode(&blueprint)
}

var defaultComponents = types.Components{
	Core: &types.Core{
		Ingress: &types.CoreComponent{
			Enabled:  true,
			Provider: "ingress-nginx",
			Config: dig.Mapping{
				"controller": dig.Mapping{
					"service": dig.Mapping{
						"type": "NodePort",
						"nodePorts": dig.Mapping{
							"http":  30000,
							"https": 30001,
						},
					},
				},
			},
		},
	},
	Addons: []types.Addon{
		{
			Name:      "example-server",
			Kind:      boundless.AddonKindChart,
			Enabled:   true,
			Namespace: "default",
			Chart: &types.ChartInfo{
				Name:    "nginx",
				Repo:    "https://charts.bitnami.com/bitnami",
				Version: "15.1.1",
				Values: `"service":
  "type": "ClusterIP"
`,
			},
		},
	},
}
