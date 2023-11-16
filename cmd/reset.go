package cmd

import (
	"github.com/spf13/cobra"

	"boundless-cli/internal/distro"
	"boundless-cli/internal/k0sctl"
)

// resetCmd represents the apply command
func resetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reset",
		Short:   "Reset the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resetFunc()
		},
	}

	flags := cmd.Flags()
	addConfigFlags(flags)
	addKubeFlags(flags)
	return cmd
}

func resetFunc() error {
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
