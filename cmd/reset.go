package cmd

import (
	"github.com/mirantiscontainers/boundless-cli/internal/distro"
	"github.com/mirantiscontainers/boundless-cli/internal/k0sctl"

	"github.com/spf13/cobra"
)

// resetCmd represents the apply command
func resetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the cluster to a clean state",
		Long: `
Reset the cluster to a clean state.

For a cluster with k0s, this will remove traces of k0s from all the hosts.
For a cluster with kind, it will delete the cluster (same as 'kind delete cluster <CLUSTER NAME>').
For a cluster with an external Kubernetes provider, this will remove Boundless Operator and all associated resources.
`,
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReset()
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)
	return cmd
}

func runReset() error {
	switch blueprint.Spec.Kubernetes.Provider {
	case "k0s":
		path, err := k0sctl.GetConfigPath(blueprint)
		if err != nil {
			return err
		}
		return distro.ResetK0s(path)
	case "kind":
		return distro.ResetKind(blueprint.Metadata.Name)
	}
	return nil
}
